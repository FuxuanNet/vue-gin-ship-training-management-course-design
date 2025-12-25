<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPendingEvaluations, submitGrading } from '../../api/teacher'

const loading = ref(true)
const courseItems = ref([])
const totalCount = ref(0)
const pendingCount = ref(0)
const statusFilter = ref('pending') // pending 或 all

// 获取待评分学员列表
const fetchPendingEvaluations = async () => {
  loading.value = true
  try {
    const response = await getPendingEvaluations({
      status: statusFilter.value
    })
    if (response.code === 200) {
      courseItems.value = response.data.courseItems || []
      totalCount.value = response.data.totalCount || 0
      pendingCount.value = response.data.pendingCount || 0
    }
  } catch (error) {
    ElMessage.error('获取学员列表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const selectedCourse = ref(null)
const dialogVisible = ref(false)
const currentStudent = ref(null)
const submitting = ref(false)
const gradingForm = ref({
  teacherScore: null,
  teacherComment: '',
  scoreRatio: 0.7
})

// 打开评分对话框
const openGradingDialog = (courseItem, student) => {
  selectedCourse.value = courseItem
  currentStudent.value = student
  
  gradingForm.value = {
    teacherScore: student.teacherScore || null,
    teacherComment: student.teacherComment || '',
    scoreRatio: student.scoreRatio || 0.7
  }
  
  dialogVisible.value = true
}

// 提交评分
const handleSubmitGrading = async () => {
  if (gradingForm.value.teacherScore === null) {
    ElMessage.warning('请输入评分')
    return
  }
  
  if (gradingForm.value.teacherScore < 0 || gradingForm.value.teacherScore > 100) {
    ElMessage.warning('评分应在0-100之间')
    return
  }
  
  submitting.value = true
  try {
    const response = await submitGrading({
      itemId: selectedCourse.value.itemId,
      personId: currentStudent.value.personId,
      teacherScore: gradingForm.value.teacherScore,
      teacherComment: gradingForm.value.teacherComment,
      scoreRatio: gradingForm.value.scoreRatio
    })
    
    if (response.code === 200) {
      ElMessage.success('评分提交成功')
      dialogVisible.value = false
      // 重新获取列表
      await fetchPendingEvaluations()
    }
  } catch (error) {
    ElMessage.error('提交评分失败：' + (error.message || '未知错误'))
  } finally {
    submitting.value = false
  }
}

// 使用AI评分
const useAIGrading = async () => {
  if (!gradingForm.value.teacherComment.trim()) {
    ElMessage.warning('请先填写评语，AI将根据评语生成分数')
    return
  }
  
  try {
    // 调用后端AI评分接口
    const response = await submitGrading({
      itemId: selectedCourse.value.itemId,
      personId: currentStudent.value.personId,
      teacherScore: 0, // AI将根据评语生成
      teacherComment: gradingForm.value.teacherComment,
      scoreRatio: gradingForm.value.scoreRatio
    })
    
    if (response.code === 200) {
      gradingForm.value.teacherScore = response.data.teacherScore
      ElMessage.success(`AI已根据评语生成评分：${response.data.teacherScore}分`)
      dialogVisible.value = false
      await fetchPendingEvaluations()
    }
  } catch (error) {
    ElMessage.error('AI评分失败')
  }
}

onMounted(() => {
  fetchPendingEvaluations()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>学员评分</h1>
      <p class="subtitle">对学员的学习表现进行评价和打分</p>
    </div>
    
    <div class="content-wrapper" v-loading="loading">
      <div class="filter-bar">
        <el-radio-group v-model="statusFilter" @change="fetchPendingEvaluations" size="large">
          <el-radio-button label="pending">待评分</el-radio-button>
          <el-radio-button label="all">全部</el-radio-button>
        </el-radio-group>
        <div class="stats-info">
          <el-tag type="warning" size="large">待评分：{{ pendingCount }}</el-tag>
          <el-tag type="info" size="large">总数：{{ totalCount }}</el-tag>
        </div>
      </div>

      <el-collapse accordion>
        <el-collapse-item 
          v-for="courseItem in courseItems" 
          :key="courseItem.itemId"
          :name="courseItem.itemId"
        >
          <template #title>
            <div class="course-title-item">
              <el-tag type="primary">{{ courseItem.classDate }}</el-tag>
              <span class="course-name">{{ courseItem.courseName }}</span>
              <el-tag size="small">{{ courseItem.students?.length || 0 }}人</el-tag>
            </div>
          </template>
          
          <el-table :data="courseItem.students" border>
            <el-table-column prop="personName" label="学员姓名" width="120" />
            <el-table-column prop="selfScore" label="自评分" width="100" align="center">
              <template #default="{ row }">
                <span class="score-value">{{ row.selfScore || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="selfComment" label="自评内容" min-width="200" show-overflow-tooltip />
            <el-table-column prop="teacherScore" label="我的评分" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.teacherScore" type="success">{{ row.teacherScore }}</el-tag>
                <span v-else style="color: #909399;">未评分</span>
              </template>
            </el-table-column>
            <el-table-column label="综合得分" width="100" align="center">
              <template #default="{ row }">
                <el-tag 
                  v-if="row.teacherScore && row.selfScore"
                  effect="dark"
                >
                  {{ (row.selfScore * (1 - row.scoreRatio) + row.teacherScore * row.scoreRatio).toFixed(1) }}
                </el-tag>
                <span v-else style="color: #909399;">-</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button 
                  type="primary" 
                  size="small"
                  @click="openGradingDialog(courseItem, row)"
                >
                  {{ row.teacherScore ? '修改' : '评分' }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-collapse-item>
      </el-collapse>
      
      <el-empty v-if="courseItems.length === 0" description="暂无课程需要评分" />
    </div>
    
    <!-- 评分对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      title="学员评分"
      width="600px"
    >
      <div class="student-info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="学员姓名">{{ currentStudent?.personName }}</el-descriptions-item>
          <el-descriptions-item label="课程名称">{{ selectedCourse?.courseName }}</el-descriptions-item>
          <el-descriptions-item label="上课日期">{{ selectedCourse?.classDate }}</el-descriptions-item>
          <el-descriptions-item label="学员自评分">
            <el-tag v-if="currentStudent?.selfScore" type="primary">{{ currentStudent?.selfScore }}</el-tag>
            <span v-else style="color: #909399;">未自评</span>
          </el-descriptions-item>
          <el-descriptions-item label="学员自评内容" :span="2">
            {{ currentStudent?.selfComment || '暂无' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
      
      <el-form :model="gradingForm" label-width="120px" style="margin-top: 20px;">
        <el-form-item label="讲师评语">
          <el-input
            v-model="gradingForm.teacherComment"
            type="textarea"
            :rows="4"
            placeholder="请填写对学员的评价，如学习态度、课堂表现、回答问题情况等"
          />
          <div class="form-tip">
            <el-icon><InfoFilled /></el-icon>
            <span>可以选择填写评语后使用AI生成分数，或直接输入分数</span>
          </div>
        </el-form-item>
        
        <el-form-item label="讲师评分">
          <el-input-number 
            v-model="gradingForm.teacherScore" 
            :min="0" 
            :max="100" 
            :precision="1"
            placeholder="0-100"
          />
          <el-button 
            type="primary" 
            plain 
            style="margin-left: 15px;"
            @click="useAIGrading"
          >
            <el-icon><MagicStick /></el-icon>
            AI生成分数
          </el-button>
        </el-form-item>
        
        <el-form-item label="讲师评分占比">
          <el-slider 
            v-model="gradingForm.scoreRatio" 
            :min="0" 
            :max="1" 
            :step="0.1"
            :format-tooltip="(val) => `${(val * 100).toFixed(0)}%`"
          />
          <div class="ratio-info">
            讲师评分占比：{{ (gradingForm.scoreRatio * 100).toFixed(0) }}%，
            学员自评占比：{{ ((1 - gradingForm.scoreRatio) * 100).toFixed(0) }}%
          </div>
        </el-form-item>
        
        <el-form-item label="预计综合得分" v-if="gradingForm.teacherScore && currentStudent?.selfScore">
          <el-tag size="large" type="success" effect="dark">
            {{ (currentStudent.selfScore * (1 - gradingForm.scoreRatio) + gradingForm.teacherScore * gradingForm.scoreRatio).toFixed(1) }}分
          </el-tag>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitGrading" :loading="submitting">提交评分</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-container {
  padding: 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 8px;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 6px;
}

.stats-info {
  display: flex;
  gap: 15px;
}

.course-title-item {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
}

.course-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.score-value {
  font-weight: 600;
  color: #409EFF;
}

.student-info {
  margin-bottom: 20px;
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.ratio-info {
  color: #606266;
  font-size: 13px;
  margin-top: 8px;
}
</style>
