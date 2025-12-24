<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMockDataStore } from '../../stores/mockData'
import { useUserStore } from '../../stores/user'

const mockDataStore = useMockDataStore()
const userStore = useUserStore()

// 获取已完成的课程（模拟）
const completedCourses = computed(() => {
  return mockDataStore.courseItems
    .slice(0, 3)
    .map(item => ({
      ...item,
      course: mockDataStore.courses.find(c => c.course_id === item.course_id),
      evaluated: mockDataStore.attendanceEvaluations.some(
        ae => ae.item_id === item.item_id && ae.person_id === userStore.userInfo.id
      )
    }))
})

const dialogVisible = ref(false)
const currentCourse = ref(null)
const evaluationForm = ref({
  selfComment: '',
  understanding: 5,
  difficulty: 5,
  satisfaction: 5
})

// 打开评价对话框
const openEvaluationDialog = (course) => {
  currentCourse.value = course
  
  // 如果已评价，加载已有评价
  const existingEval = mockDataStore.attendanceEvaluations.find(
    ae => ae.item_id === course.item_id && ae.person_id === userStore.userInfo.id
  )
  
  if (existingEval) {
    evaluationForm.value.selfComment = existingEval.self_comment
  } else {
    evaluationForm.value = {
      selfComment: '',
      understanding: 5,
      difficulty: 5,
      satisfaction: 5
    }
  }
  
  dialogVisible.value = true
}

// 提交评价
const submitEvaluation = () => {
  if (!evaluationForm.value.selfComment.trim()) {
    ElMessage.warning('请填写学习心得')
    return
  }
  
  // 模拟AI评分
  const aiScore = Math.floor(Math.random() * 20) + 75
  
  ElMessage.success(`评价提交成功！AI评估分数：${aiScore}分`)
  dialogVisible.value = false
  
  // 实际应该调用API保存评价
  console.log('提交评价:', {
    item_id: currentCourse.value.item_id,
    person_id: userStore.userInfo.id,
    ...evaluationForm.value,
    ai_score: aiScore
  })
}
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>课程自评</h1>
      <p class="subtitle">填写学习心得，系统将通过AI生成您的掌握度评分</p>
    </div>
    
    <div class="content-wrapper">
      <el-table :data="completedCourses" border>
        <el-table-column prop="class_date" label="上课日期" width="120" />
        <el-table-column label="课程名称" min-width="180">
          <template #default="{ row }">
            {{ row.course?.course_name }}
          </template>
        </el-table-column>
        <el-table-column label="课程类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.course?.course_class }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="上课地点" width="150" />
        <el-table-column label="上课时间" width="150">
          <template #default="{ row }">
            {{ row.class_begin_time.slice(0, 5) }} - {{ row.class_end_time.slice(0, 5) }}
          </template>
        </el-table-column>
        <el-table-column label="评价状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.evaluated ? 'success' : 'info'" size="small">
              {{ row.evaluated ? '已评价' : '未评价' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small"
              @click="openEvaluationDialog(row)"
            >
              {{ row.evaluated ? '修改评价' : '开始评价' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 评价对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="currentCourse?.course?.course_name + ' - 学习评价'"
      width="600px"
    >
      <el-form :model="evaluationForm" label-width="100px">
        <el-form-item label="学习心得">
          <el-input
            v-model="evaluationForm.selfComment"
            type="textarea"
            :rows="6"
            placeholder="请详细描述您的学习收获、遇到的困难、掌握情况等，内容越详细，AI评分越准确"
          />
          <div class="form-tip">
            <el-icon><InfoFilled /></el-icon>
            <span>系统将根据您的描述通过AI分析生成掌握度评分（0-100分）</span>
          </div>
        </el-form-item>
        
        <el-form-item label="内容理解度">
          <el-rate v-model="evaluationForm.understanding" show-score />
        </el-form-item>
        
        <el-form-item label="课程难度">
          <el-rate v-model="evaluationForm.difficulty" show-score />
        </el-form-item>
        
        <el-form-item label="课程满意度">
          <el-rate v-model="evaluationForm.satisfaction" show-score />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEvaluation">提交评价</el-button>
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
  font-size: 12px;
}
</style>
