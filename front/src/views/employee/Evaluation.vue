<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPendingEvaluations, submitEvaluation } from '../../api/employee'

const loading = ref(true)
const pendingCourses = ref([])

// 获取待自评课程
const fetchPendingEvaluations = async () => {
  loading.value = true
  try {
    const response = await getPendingEvaluations()
    if (response.code === 200) {
      pendingCourses.value = response.data.courses || []
    }
  } catch (error) {
    ElMessage.error('获取待评价课程失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const dialogVisible = ref(false)
const currentCourse = ref(null)
const submitting = ref(false)
const evaluationForm = ref({
  selfComment: '',
  understanding: 5,
  difficulty: 5,
  satisfaction: 5
})

// 打开评价对话框
const openEvaluationDialog = (course) => {
  currentCourse.value = course
  
  // 重置表单
  evaluationForm.value = {
    selfComment: '',
    understanding: 5,
    difficulty: 5,
    satisfaction: 5
  }
  
  dialogVisible.value = true
}

// 提交评价
const handleSubmitEvaluation = async () => {
  if (!evaluationForm.value.selfComment.trim()) {
    ElMessage.warning('请填写学习心得')
    return
  }

  if (evaluationForm.value.selfComment.trim().length < 50) {
    ElMessage.warning('学习心得至少需要50个字符')
    return
  }
  
  submitting.value = true
  try {
    const response = await submitEvaluation({
      itemId: currentCourse.value.itemId,
      selfComment: evaluationForm.value.selfComment,
      understanding: evaluationForm.value.understanding,
      difficulty: evaluationForm.value.difficulty,
      satisfaction: evaluationForm.value.satisfaction
    })
    
    if (response.code === 200) {
      ElMessage.success(`评价提交成功！AI评估分数：${response.data.aiScore.toFixed(1)}分`)
      dialogVisible.value = false
      // 重新获取待评价课程列表
      await fetchPendingEvaluations()
    }
  } catch (error) {
    ElMessage.error('提交评价失败：' + (error.message || '未知错误'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchPendingEvaluations()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>课程自评</h1>
      <p class="subtitle">填写学习心得，系统将通过AI生成您的掌握度评分</p>
    </div>
    
    <div class="content-wrapper" v-loading="loading">
      <el-empty v-if="!loading && pendingCourses.length === 0" description="暂无待评价的课程">
        <el-button type="primary" @click="$router.push('/employee/today-courses')">查看今日课程</el-button>
      </el-empty>
      
      <el-table v-else :data="pendingCourses" border>
        <el-table-column prop="classDate" label="上课日期" width="120" />
        <el-table-column prop="courseName" label="课程名称" min-width="180" />
        <el-table-column label="课程类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.courseClass }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="上课地点" width="150" />
        <el-table-column label="上课时间" width="150">
          <template #default="{ row }">
            {{ row.classBeginTime.slice(0, 5) }} - {{ row.classEndTime.slice(0, 5) }}
          </template>
        </el-table-column>
        <el-table-column prop="teacherName" label="讲师" width="100" />
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small"
              @click="openEvaluationDialog(row)"
            >
              开始评价
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 评价对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="currentCourse?.courseName + ' - 学习评价'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="evaluationForm" label-width="100px">
        <el-form-item label="学习心得" required>
          <el-input
            v-model="evaluationForm.selfComment"
            type="textarea"
            :rows="6"
            placeholder="请详细描述您的学习收获、遇到的困难、掌握情况等，内容越详细，AI评分越准确（至少50字）"
            maxlength="1000"
            show-word-limit
          />
          <div class="form-tip">
            <el-icon><InfoFilled /></el-icon>
            <span>系统将根据您的描述通过AI分析生成掌握度评分（0-100分）</span>
          </div>
        </el-form-item>
        
        <el-form-item label="内容理解度">
          <el-rate v-model="evaluationForm.understanding" show-score />
          <span class="rate-desc">(1分:完全不理解 - 5分:完全理解)</span>
        </el-form-item>
        
        <el-form-item label="课程难度">
          <el-rate v-model="evaluationForm.difficulty" show-score />
          <span class="rate-desc">(1分:很简单 - 5分:很难)</span>
        </el-form-item>
        
        <el-form-item label="课程满意度">
          <el-rate v-model="evaluationForm.satisfaction" show-score />
          <span class="rate-desc">(1分:不满意 - 5分:非常满意)</span>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false" :disabled="submitting">取消</el-button>
        <el-button type="primary" @click="handleSubmitEvaluation" :loading="submitting">提交评价</el-button>
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

.form-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}

.rate-desc {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}
</style>
