<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMockDataStore } from '../../stores/mockData'
import { useUserStore } from '../../stores/user'

const mockDataStore = useMockDataStore()
const userStore = useUserStore()

// 获取讲师授课的课程安排
const myCourseItems = computed(() => {
  return mockDataStore.courseItems
    .filter(item => {
      const course = mockDataStore.courses.find(c => c.course_id === item.course_id)
      return course && course.teacher_id === userStore.userInfo.id
    })
    .map(item => ({
      ...item,
      course: mockDataStore.courses.find(c => c.course_id === item.course_id),
      students: mockDataStore.attendanceEvaluations
        .filter(ae => ae.item_id === item.item_id)
        .map(ae => ({
          ...ae,
          person: mockDataStore.persons.find(p => p.person_id === ae.person_id)
        }))
    }))
})

const selectedCourse = ref(null)
const dialogVisible = ref(false)
const currentStudent = ref(null)
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
    teacherScore: student.teacher_score || null,
    teacherComment: student.teacher_comment || '',
    scoreRatio: student.score_ratio || 0.7
  }
  
  dialogVisible.value = true
}

// 提交评分
const submitGrading = () => {
  if (gradingForm.value.teacherScore === null) {
    ElMessage.warning('请输入评分')
    return
  }
  
  if (gradingForm.value.teacherScore < 0 || gradingForm.value.teacherScore > 100) {
    ElMessage.warning('评分应在0-100之间')
    return
  }
  
  ElMessage.success('评分提交成功')
  dialogVisible.value = false
  
  // 实际应该调用API保存评分
  console.log('提交评分:', {
    item_id: selectedCourse.value.item_id,
    person_id: currentStudent.value.person_id,
    ...gradingForm.value
  })
}

// 使用AI评分
const useAIGrading = () => {
  if (!gradingForm.value.teacherComment.trim()) {
    ElMessage.warning('请先填写评语，AI将根据评语生成分数')
    return
  }
  
  // 模拟AI评分
  const aiScore = Math.floor(Math.random() * 20) + 75
  gradingForm.value.teacherScore = aiScore
  ElMessage.success(`AI已根据评语生成评分：${aiScore}分`)
}
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>学员评分</h1>
      <p class="subtitle">对学员的学习表现进行评价和打分</p>
    </div>
    
    <div class="content-wrapper">
      <el-collapse accordion>
        <el-collapse-item 
          v-for="courseItem in myCourseItems" 
          :key="courseItem.item_id"
          :name="courseItem.item_id"
        >
          <template #title>
            <div class="course-title-item">
              <el-tag type="primary">{{ courseItem.class_date }}</el-tag>
              <span class="course-name">{{ courseItem.course?.course_name }}</span>
              <el-tag size="small">{{ courseItem.students.length }}人</el-tag>
            </div>
          </template>
          
          <el-table :data="courseItem.students" border>
            <el-table-column prop="person.name" label="学员姓名" width="120" />
            <el-table-column prop="self_score" label="自评分" width="100" align="center">
              <template #default="{ row }">
                <span class="score-value">{{ row.self_score || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="self_comment" label="自评内容" min-width="200" show-overflow-tooltip />
            <el-table-column prop="teacher_score" label="我的评分" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.teacher_score" type="success">{{ row.teacher_score }}</el-tag>
                <span v-else style="color: #909399;">未评分</span>
              </template>
            </el-table-column>
            <el-table-column label="综合得分" width="100" align="center">
              <template #default="{ row }">
                <el-tag 
                  v-if="row.teacher_score && row.self_score"
                  effect="dark"
                >
                  {{ (row.self_score * (1 - row.score_ratio) + row.teacher_score * row.score_ratio).toFixed(1) }}
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
                  {{ row.teacher_score ? '修改' : '评分' }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-collapse-item>
      </el-collapse>
      
      <el-empty v-if="myCourseItems.length === 0" description="暂无课程需要评分" />
    </div>
    
    <!-- 评分对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      title="学员评分"
      width="600px"
    >
      <div class="student-info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="学员姓名">{{ currentStudent?.person?.name }}</el-descriptions-item>
          <el-descriptions-item label="课程名称">{{ selectedCourse?.course?.course_name }}</el-descriptions-item>
          <el-descriptions-item label="上课日期">{{ selectedCourse?.class_date }}</el-descriptions-item>
          <el-descriptions-item label="学员自评分">
            <el-tag v-if="currentStudent?.self_score" type="primary">{{ currentStudent?.self_score }}</el-tag>
            <span v-else style="color: #909399;">未自评</span>
          </el-descriptions-item>
          <el-descriptions-item label="学员自评内容" :span="2">
            {{ currentStudent?.self_comment || '暂无' }}
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
        
        <el-form-item label="预计综合得分" v-if="gradingForm.teacherScore && currentStudent?.self_score">
          <el-tag size="large" type="success" effect="dark">
            {{ (currentStudent.self_score * (1 - gradingForm.scoreRatio) + gradingForm.teacherScore * gradingForm.scoreRatio).toFixed(1) }}分
          </el-tag>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGrading">提交评分</el-button>
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
