<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getCourseStatistics, getTeachingStatistics } from '../../api/teacher'
import * as echarts from 'echarts'

const chartRef = ref(null)
const selectedCourseId = ref(null)
const loading = ref(true)
const statisticsData = ref(null)
const teachingStats = ref(null)

// 获取讲师整体授课统计
const fetchTeachingStatistics = async () => {
  try {
    const response = await getTeachingStatistics({})
    if (response.code === 200) {
      teachingStats.value = response.data
      // 如果有课程，默认选择第一个
      if (response.data.recentClasses && response.data.recentClasses.length > 0) {
        // 从最近授课中提取课程ID（需要从后端返回courseId）
        // 暂时使用第一门课程
      }
    }
  } catch (error) {
    ElMessage.error('获取授课统计失败')
  }
}

// 获取课程成绩统计
const fetchCourseStatistics = async () => {
  if (!selectedCourseId.value) return
  
  loading.value = true
  try {
    const response = await getCourseStatistics({
      courseId: selectedCourseId.value
    })
    if (response.code === 200) {
      statisticsData.value = response.data
      nextTick(() => {
        initChart()
      })
    }
  } catch (error) {
    ElMessage.error('获取课程统计失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 获取讲师教授的课程列表（从授课统计中提取）
const myCourses = computed(() => {
  if (!teachingStats.value) return []
  
  // 从最近授课中提取唯一课程
  const coursesMap = new Map()
  if (teachingStats.value.recentClasses) {
    teachingStats.value.recentClasses.forEach(item => {
      if (!coursesMap.has(item.courseId)) {
        coursesMap.set(item.courseId, {
          courseId: item.courseId,
          courseName: item.courseName
        })
      }
    })
  }
  
  return Array.from(coursesMap.values())
})

// 课程平均分
const courseAverageScore = computed(() => {
  return statisticsData.value?.statistics?.averageScore?.toFixed(2) || '0.00'
})

// 课程统计信息
const courseStats = computed(() => {
  if (!statisticsData.value?.statistics) {
    return { max: 0, min: 0, avg: 0, count: 0 }
  }
  
  const stats = statisticsData.value.statistics
  return {
    max: stats.maxScore?.toFixed(2) || '0.00',
    min: stats.minScore?.toFixed(2) || '0.00',
    avg: stats.averageScore?.toFixed(2) || '0.00',
    count: stats.totalStudents || 0,
    passRate: stats.passRate?.toFixed(1) || '0.0',
    excellentRate: stats.excellentRate?.toFixed(1) || '0.0'
  }
})

// 初始化成绩分布图表
const initChart = () => {
  if (!chartRef.value || !statisticsData.value) return
  
  const chart = echarts.init(chartRef.value)
  
  // 按学员准备数据
  const studentScores = statisticsData.value.studentScores || []
  const studentNames = studentScores.map(s => s.personName)
  const averageScores = studentScores.map(s => s.averageScore)
  const latestScores = studentScores.map(s => s.latestScore)
  
  const option = {
    title: {
      text: '学员成绩分析',
      left: 'center',
      top: 10
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['平均分', '最新成绩'],
      top: 40
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: 80,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: studentNames,
      axisLabel: {
        interval: 0,
        rotate: 45
      }
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: 100,
      name: '分数'
    },
    series: [
      {
        name: '平均分',
        type: 'bar',
        data: averageScores,
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '最新成绩',
        type: 'line',
        data: latestScores,
        smooth: true,
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 监听课程选择变化
const handleCourseChange = () => {
  fetchCourseStatistics()
}

// 监听课程列表变化，自动选择第一个课程
watch(myCourses, (newCourses) => {
  if (newCourses.length > 0 && !selectedCourseId.value) {
    selectedCourseId.value = newCourses[0].courseId
    fetchCourseStatistics()
  }
}, { immediate: true })

onMounted(async () => {
  await fetchTeachingStatistics()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>成绩查看</h1>
      <el-select 
        v-model="selectedCourseId" 
        placeholder="选择课程"
        @change="handleCourseChange"
        style="width: 300px;"
        :loading="loading"
      >
        <el-option 
          v-for="course in myCourses" 
          :key="course.courseId"
          :label="course.courseName"
          :value="course.courseId"
        />
      </el-select>
    </div>
    
    <div class="stats-section" v-if="selectedCourseId && statisticsData">
      <el-row :gutter="20">
        <el-col :span="4">
          <el-card>
            <el-statistic title="参与人数" :value="courseStats.count" suffix="人">
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card>
            <el-statistic title="平均分" :value="courseStats.avg" suffix="分">
              <template #prefix>
                <el-icon><TrendCharts /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card>
            <el-statistic title="最高分" :value="courseStats.max" suffix="分">
              <template #prefix>
                <el-icon><Top /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card>
            <el-statistic title="最低分" :value="courseStats.min" suffix="分">
              <template #prefix>
                <el-icon><Bottom /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card>
            <el-statistic title="及格率" :value="courseStats.passRate" suffix="%">
              <template #prefix>
                <el-icon><Check /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card>
            <el-statistic title="优秀率" :value="courseStats.excellentRate" suffix="%">
              <template #prefix>
                <el-icon><Star /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <div class="content-wrapper" v-if="selectedCourseId">
      <!-- 成绩趋势图 -->
      <div class="chart-section">
        <div ref="chartRef" class="chart"></div>
      </div>
      
      <!-- 成绩明细表 -->
      <div class="table-section" v-loading="loading">
        <h3>学员成绩详情</h3>
        <el-table :data="statisticsData?.studentScores || []" border v-if="statisticsData">
          <el-table-column prop="personName" label="学员姓名" width="120" />
          <el-table-column prop="classCount" label="上课次数" width="100" align="center" />
          <el-table-column prop="averageScore" label="平均分" width="100" align="center">
            <template #default="{ row }">
              <span class="score-value">{{ row.averageScore?.toFixed(2) || '0.00' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="latestScore" label="最新得分" width="100" align="center">
            <template #default="{ row }">
              <el-tag 
                :type="row.latestScore >= 90 ? 'success' : row.latestScore >= 80 ? '' : row.latestScore >= 60 ? 'warning' : 'danger'"
                effect="dark"
              >
                {{ row.latestScore?.toFixed(2) || '0.00' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="trend" label="趋势" width="100" align="center">
            <template #default="{ row }">
              <el-icon v-if="row.trend === 'up'" color="#67C23A"><CaretTop /></el-icon>
              <el-icon v-else-if="row.trend === 'down'" color="#F56C6C"><CaretBottom /></el-icon>
              <el-icon v-else color="#909399"><Minus /></el-icon>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 课次统计表 -->
      <div class="table-section" v-loading="loading" style="margin-top: 20px;">
        <h3>课次统计</h3>
        <el-table :data="statisticsData?.classStat || []" border v-if="statisticsData">
          <el-table-column prop="classDate" label="上课日期" width="120" />
          <el-table-column prop="location" label="上课地点" width="150" />
          <el-table-column prop="studentCount" label="学员人数" width="100" align="center" />
          <el-table-column prop="evaluatedCount" label="已评分" width="100" align="center" />
          <el-table-column prop="averageScore" label="平均分" width="100" align="center">
            <template #default="{ row }">
              <span class="score-value">{{ row.averageScore?.toFixed(2) || '0.00' }}</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
    
    <el-empty v-else description="请先选择一门课程" />
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: 28px;
  color: #303133;
}

.stats-section {
  margin-bottom: 30px;
}

.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.chart-section {
  margin-bottom: 30px;
}

.chart {
  width: 100%;
  height: 400px;
}

.table-section h3 {
  margin-bottom: 15px;
  color: #303133;
  font-size: 18px;
}

.score-value {
  font-weight: 600;
  color: #409EFF;
}
</style>
