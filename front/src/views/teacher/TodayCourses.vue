<script setup>
import { computed } from 'vue'
import { useMockDataStore } from '../../stores/mockData'
import { useUserStore } from '../../stores/user'
import CourseCard from '../../components/CourseCard.vue'

const mockDataStore = useMockDataStore()
const userStore = useUserStore()

const todayCourses = computed(() => {
  return mockDataStore.getTodayCoursesForTeacher(userStore.userInfo.id)
})

const currentDate = computed(() => {
  const date = new Date()
  const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
  return date.toLocaleDateString('zh-CN', options)
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>今日授课</h1>
      <p class="date-info">{{ currentDate }}</p>
    </div>
    
    <div class="content-wrapper">
      <el-empty v-if="todayCourses.length === 0" description="今天没有授课安排">
        <el-button type="primary" @click="$router.push('/teacher/schedule')">查看授课表</el-button>
      </el-empty>
      
      <el-row :gutter="20" v-else>
        <el-col :xs="24" :sm="12" :lg="8" v-for="course in todayCourses" :key="course.item_id">
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
