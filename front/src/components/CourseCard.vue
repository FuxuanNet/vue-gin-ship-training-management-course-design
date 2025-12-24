<script setup>
const props = defineProps({
  course: {
    type: Object,
    required: true
  }
})
</script>

<template>
  <el-card class="course-card" shadow="hover">
    <div class="course-header">
      <div class="course-time">
        <el-icon><Clock /></el-icon>
        <span>{{ course.class_begin_time }} - {{ course.class_end_time }}</span>
      </div>
      <el-tag :type="course.status === '已完成' ? 'success' : 'primary'" size="small">
        {{ course.status || '待上课' }}
      </el-tag>
    </div>
    <h3 class="course-title">{{ course.course?.course_name }}</h3>
    <div class="course-info">
      <div class="info-item">
        <el-icon><Location /></el-icon>
        <span>{{ course.location }}</span>
      </div>
      <div class="info-item">
        <el-icon><User /></el-icon>
        <span>{{ course.course?.teacher_name || '李老师' }}</span>
      </div>
    </div>
    <div class="course-desc">
      <p>{{ course.course?.course_desc }}</p>
    </div>
    <div class="course-require" v-if="course.course?.course_require">
      <el-icon><InfoFilled /></el-icon>
      <span>{{ course.course?.course_require }}</span>
    </div>
  </el-card>
</template>

<style scoped>
.course-card {
  margin-bottom: 16px;
  transition: all 0.3s;
}

.course-card:hover {
  transform: translateY(-4px);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.course-time {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #409EFF;
  font-weight: 500;
}

.course-title {
  font-size: 18px;
  margin-bottom: 12px;
  color: #303133;
}

.course-info {
  display: flex;
  gap: 20px;
  margin-bottom: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #606266;
  font-size: 14px;
}

.course-desc {
  margin-bottom: 12px;
  color: #909399;
  font-size: 14px;
  line-height: 1.6;
}

.course-require {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
  color: #606266;
  font-size: 13px;
}
</style>
