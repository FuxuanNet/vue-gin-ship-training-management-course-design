<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { getSchedule } from '../../api/teacher'
import { ElMessage } from 'element-plus'

const loading = ref(true)
const scheduleData = ref([])
const currentWeekOffset = ref(0)

// 时间段定义
const timeSlots = [
  { label: '08:00-10:00', start: '08:00', end: '10:00' },
  { label: '10:00-12:00', start: '10:00', end: '12:00' },
  { label: '14:00-16:00', start: '14:00', end: '16:00' },
  { label: '16:00-18:00', start: '16:00', end: '18:00' }
]

// 获取当前周的日期范围
const weekDates = computed(() => {
  const today = new Date()
  const currentDay = today.getDay()
  const monday = new Date(today)
  monday.setDate(today.getDate() - currentDay + 1 + currentWeekOffset.value * 7)
  
  const dates = []
  for (let i = 0; i < 7; i++) {
    const date = new Date(monday)
    date.setDate(monday.getDate() + i)
    dates.push({
      date: date,
      dateStr: date.toISOString().split('T')[0],
      dayLabel: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'][i],
      monthDay: `${date.getMonth() + 1}/${date.getDate()}`
    })
  }
  return dates
})

// 当前周的描述
const weekDescription = computed(() => {
  const firstDay = weekDates.value[0]
  const lastDay = weekDates.value[6]
  return `${firstDay.monthDay} - ${lastDay.monthDay}`
})

// 获取授课表
const fetchSchedule = async () => {
  loading.value = true
  const firstDay = weekDates.value[0]
  const lastDay = weekDates.value[6]
  
  try {
    const response = await getSchedule({
      startDate: firstDay.dateStr,
      endDate: lastDay.dateStr
    })
    if (response.code === 200) {
      scheduleData.value = response.data.schedule || []
    }
  } catch (error) {
    ElMessage.error('获取授课表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 课程表数据（按时间段组织）
const organizedScheduleData = computed(() => {
  const schedule = {}
  weekDates.value.forEach(day => {
    schedule[day.dateStr] = {}
    timeSlots.forEach(slot => {
      schedule[day.dateStr][slot.label] = []
    })
  })
  
  if (scheduleData.value && Array.isArray(scheduleData.value)) {
    scheduleData.value.forEach(daySchedule => {
      const dateStr = daySchedule.date
      if (schedule[dateStr] && daySchedule.courses) {
        daySchedule.courses.forEach(course => {
          const timeKey = findTimeSlot(course.classBeginTime, course.classEndTime)
          if (timeKey && schedule[dateStr][timeKey]) {
            schedule[dateStr][timeKey].push(course)
          }
        })
      }
    })
  }
  
  return schedule
})

// 查找对应的时间段
const findTimeSlot = (beginTime, endTime) => {
  const begin = beginTime.slice(0, 5)
  for (const slot of timeSlots) {
    if (begin >= slot.start && begin < slot.end) {
      return slot.label
    }
  }
  return null
}

// 切换周
const changeWeek = (offset) => {
  currentWeekOffset.value += offset
}

// 回到本周
const goToCurrentWeek = () => {
  currentWeekOffset.value = 0
}

// 监听周变化
watch(currentWeekOffset, () => {
  fetchSchedule()
})

// 页面加载时获取数据
onMounted(() => {
  fetchSchedule()
})

// 判断是否是今天
const isToday = (dateStr) => {
  const today = new Date().toISOString().split('T')[0]
  return dateStr === today
}
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-title">
        <h1>我的授课表</h1>
      </div>
      <div class="week-selector">
        <el-button @click="changeWeek(-1)" icon="ArrowLeft" class="week-btn">上一周</el-button>
        <span class="week-info">{{ weekDescription }}</span>
        <el-button @click="changeWeek(1)" icon="ArrowRight" class="week-btn">下一周</el-button>
      </div>
      <div class="current-week-btn">
        <el-button @click="goToCurrentWeek" type="primary" plain>本周</el-button>
      </div>
    </div>
    
    <div class="schedule-container" v-loading="loading">
      <el-table :data="timeSlots" border class="schedule-table">
        <el-table-column label="时间" width="120" align="center">
          <template #default="{ row }">
            <div class="time-label">{{ row.label }}</div>
          </template>
        </el-table-column>
        
        <el-table-column 
          v-for="day in weekDates" 
          :key="day.dateStr" 
          :label="day.dayLabel"
          :class-name="isToday(day.dateStr) ? 'today-column' : ''"
          align="center"
        >
          <template #header>
            <div class="day-header" :class="{ 'today': isToday(day.dateStr) }">
              <div class="day-name">{{ day.dayLabel }}</div>
              <div class="day-date">{{ day.monthDay }}</div>
            </div>
          </template>
          <template #default="{ row }">
            <div class="course-cell">
              <div 
                v-for="course in organizedScheduleData[day.dateStr]?.[row.label] || []" 
                :key="course.itemId"
                class="course-item teacher-course"
              >
                <div class="course-name">{{ course.courseName }}</div>
                <div class="course-location">
                  <el-icon><Location /></el-icon>
                  {{ course.location }}
                </div>
                <div class="course-time">
                  {{ course.classBeginTime.slice(0, 5) }} - {{ course.classEndTime.slice(0, 5) }}
                </div>
              </div>
              <div v-if="!organizedScheduleData[day.dateStr]?.[row.label] || organizedScheduleData[day.dateStr][row.label].length === 0" class="no-course">
                -
              </div>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  padding: 30px;
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.header-title h1 {
  font-size: 28px;
  color: #303133;
  margin: 0;
}

.week-selector {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
  justify-content: center;
}

.week-info {
  font-size: 16px;
  color: #606266;
  font-weight: 500;
  min-width: 120px;
  text-align: center;
}

.current-week-btn {
  display: flex;
}

/* 响应式设计 - 移动端 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    padding: 15px;
  }
  
  .header-title {
    text-align: center;
  }
  
  .header-title h1 {
    font-size: 20px;
  }
  
  .week-selector {
    justify-content: space-between;
    flex: none;
    gap: 8px;
  }
  
  .week-btn {
    flex: 1;
    padding: 8px 12px;
  }
  
  .week-info {
    font-size: 14px;
    min-width: 80px;
  }
  
  .current-week-btn {
    justify-content: center;
  }
  
  .current-week-btn .el-button {
    width: 100%;
  }
  
  .schedule-container {
    overflow-x: auto;
  }
  
  .schedule-table {
    min-width: 800px;
  }
}

.schedule-container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.schedule-table {
  width: 100%;
}

.time-label {
  font-weight: 500;
  color: #606266;
}

.day-header {
  padding: 8px;
}

.day-header.today {
  background-color: #ecf5ff;
  border-radius: 4px;
}

.day-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}

.day-date {
  font-size: 12px;
  color: #909399;
}

.course-cell {
  min-height: 80px;
  padding: 8px;
}

.course-item {
  color: white;
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.teacher-course {
  background: #67C23A;
}

.course-item:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.4);
}

.course-name {
  font-weight: 600;
  margin-bottom: 6px;
  font-size: 14px;
}

.course-location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  margin-bottom: 4px;
  opacity: 0.9;
}

.course-time {
  font-size: 12px;
  opacity: 0.9;
}

.no-course {
  color: #dcdfe6;
  font-size: 20px;
}

:deep(.today-column) {
  background-color: #fafafa;
}

:deep(.el-table) {
  font-size: 13px;
}

:deep(.el-table th) {
  background-color: #f5f7fa;
}
</style>
