<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTodayCourses } from '../../api/employee'
import { ElMessage, ElLoading } from 'element-plus'
import CourseCard from '../../components/CourseCard.vue'

// 今日课程数据
const todayCourses = ref([])
const loading = ref(false)

// 当前日期
const currentDate = computed(() => {
  const date = new Date()
  const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
  return date.toLocaleDateString('zh-CN', options)
})

// 加载今日课程
const loadTodayCourses = async () => {
  loading.value = true
  try {
    const response = await getTodayCourses()
    todayCourses.value = response.data.courses || []
  } catch (error) {
    console.error('加载今日课程失败:', error)
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  loadTodayCourses()
})
</script>

<template>
  <div class="page-container" v-loading="loading">
    <div class="page-header">
      <h1>今日课程</h1>
      <p class="date-info">{{ currentDate }}</p>
    </div>
    
    <div class="content-wrapper">
      <el-empty v-if="todayCourses.length === 0" description="今天没有安排课程">
        <el-button type="primary" @click="$router.push('/employee/schedule')">查看课程表</el-button>
      </el-empty>
      
      <el-row :gutter="20" v-else>
        <el-col :xs="24" :sm="12" :lg="8" v-for="course in todayCourses" :key="course.itemId">
          <CourseCard :course="course" />
        </el-col>
      </el-row>
    </div>
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

.date-info {
  color: #909399;
  font-size: 14px;
}

.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}
</style>
