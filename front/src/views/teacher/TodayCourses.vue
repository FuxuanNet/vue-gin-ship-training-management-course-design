<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTodayCourses } from '../../api/teacher'
import { ElMessage } from 'element-plus'
import CourseCard from '../../components/CourseCard.vue'

const loading = ref(true)
const coursesData = ref({
  date: '',
  courseCount: 0,
  courses: []
})

const currentDate = computed(() => {
  if (!coursesData.value.date) {
    const date = new Date()
    const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
    return date.toLocaleDateString('zh-CN', options)
  }
  
  const date = new Date(coursesData.value.date)
  const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
  return date.toLocaleDateString('zh-CN', options)
})

const fetchTodayCourses = async () => {
  loading.value = true
  try {
    const response = await getTodayCourses()
    if (response.code === 200) {
      coursesData.value = response.data
    }
  } catch (error) {
    ElMessage.error('获取今日授课失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTodayCourses()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>今日授课</h1>
      <p class="date-info">{{ currentDate }}</p>
    </div>
    
    <div class="content-wrapper" v-loading="loading">
      <el-empty v-if="!loading && coursesData.courseCount === 0" description="今天没有授课安排">
        <el-button type="primary" @click="$router.push('/teacher/schedule')">查看授课表</el-button>
      </el-empty>
      
      <el-row :gutter="20" v-else-if="!loading">
        <el-col :xs="24" :sm="12" :lg="8" v-for="course in coursesData.courses" :key="course.itemId">
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
