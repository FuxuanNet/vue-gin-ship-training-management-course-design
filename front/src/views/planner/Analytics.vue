<template>
  <div class="analytics-container">
    <el-card class="page-header">
      <h2>数据分析</h2>
      <el-button type="primary" :loading="loading" @click="fetchAnalytics">
        <el-icon><Refresh /></el-icon>
        刷新数据
      </el-button>
    </el-card>

    <el-row :gutter="20" v-loading="loading">
      <!-- 课程排名 -->
      <el-col :span="12">
        <el-card class="chart-card">
          <div ref="courseRankingChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <!-- 培训计划排名 -->
      <el-col :span="12">
        <el-card class="chart-card">
          <div ref="planRankingChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <!-- 课程类型分布 -->
      <el-col :span="12">
        <el-card class="chart-card">
          <div ref="courseClassChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <!-- 计划状态统计 -->
      <el-col :span="12">
        <el-card class="chart-card">
          <div ref="planStatusChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <!-- 员工排名 -->
      <el-col :span="12">
        <el-card class="chart-card">
          <div ref="employeeRankingChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <!-- 讲师统计 -->
      <el-col :span="12">
        <el-card class="chart-card">
          <div ref="teacherStatsChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { getAnalytics } from '../../api/planner'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

const courseRankingChartRef = ref(null)
const planRankingChartRef = ref(null)
const courseClassChartRef = ref(null)
const planStatusChartRef = ref(null)
const employeeRankingChartRef = ref(null)
const teacherStatsChartRef = ref(null)

const analyticsData = ref({
  courseRankings: [],
  planRankings: [],
  courseClassDistribution: [],
  planStatusStatistics: [],
  employeeRankings: [],
  teacherStatistics: []
})

const loading = ref(false)

// 获取分析数据
const fetchAnalytics = async () => {
  loading.value = true
  try {
    const res = await getAnalytics(10)
    if (res.code === 200) {
      analyticsData.value = res.data
      nextTick(() => {
        initAllCharts()
      })
    } else {
      ElMessage.error(res.message || '获取数据分析失败')
    }
  } catch (error) {
    console.error('获取数据分析失败:', error)
    ElMessage.error('获取数据分析失败')
  } finally {
    loading.value = false
  }
}

// 初始化所有图表
const initAllCharts = () => {
  initCourseRankingChart()
  initPlanRankingChart()
  initCourseClassChart()
  initPlanStatusChart()
  initEmployeeRankingChart()
  initTeacherStatsChart()
}

// 课程排名图表
const initCourseRankingChart = () => {
  if (!courseRankingChartRef.value) return
  
  const chart = echarts.init(courseRankingChartRef.value)
  const data = analyticsData.value.courseRankings || []
  
  const option = {
    title: {
      text: '课程评分排名 Top 10',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    grid: {
      left: '10%',
      right: '10%',
      bottom: '10%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.courseName),
      axisLabel: {
        rotate: 30,
        interval: 0,
        fontSize: 10
      }
    },
    yAxis: {
      type: 'value',
      name: '平均分',
      min: 0,
      max: 100
    },
    series: [
      {
        type: 'bar',
        data: data.map(item => item.courseAvgScore.toFixed(2)),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 1, color: '#188df0' }
          ])
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}分'
        }
      }
    ]
  }
  
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

// 培训计划排名图表
const initPlanRankingChart = () => {
  if (!planRankingChartRef.value) return
  
  const chart = echarts.init(planRankingChartRef.value)
  const data = analyticsData.value.planRankings || []
  
  const option = {
    title: {
      text: '培训计划评分排名 Top 10',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    grid: {
      left: '10%',
      right: '10%',
      bottom: '10%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.planName),
      axisLabel: {
        rotate: 30,
        interval: 0,
        fontSize: 10
      }
    },
    yAxis: {
      type: 'value',
      name: '平均分',
      min: 0,
      max: 100
    },
    series: [
      {
        type: 'bar',
        data: data.map(item => item.planAvgScore.toFixed(2)),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#fbc252' },
            { offset: 1, color: '#f69846' }
          ])
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}分'
        }
      }
    ]
  }
  
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

// 课程类型分布图表
const initCourseClassChart = () => {
  if (!courseClassChartRef.value) return
  
  const chart = echarts.init(courseClassChartRef.value)
  const data = analyticsData.value.courseClassDistribution || []
  
  const option = {
    title: {
      text: '课程类型分布',
      left: 'center'
    },
    tooltip: {
      trigger: 'item',
      formatter: (params) => {
        return `${params.name}<br/>课程数: ${params.value}<br/>占比: ${params.percent}%`
      }
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      top: 'middle'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['60%', '50%'],
        data: data.map(item => ({
          name: `${item.courseClass} (${item.avgScore.toFixed(2)}分)`,
          value: item.courseCount
        })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
  
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

// 计划状态统计图表
const initPlanStatusChart = () => {
  if (!planStatusChartRef.value) return
  
  const chart = echarts.init(planStatusChartRef.value)
  const data = analyticsData.value.planStatusStatistics || []
  
  const option = {
    title: {
      text: '培训计划状态统计',
      left: 'center'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      bottom: 10,
      left: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: '60%',
        data: data.map(item => ({
          name: item.planStatus,
          value: item.count
        })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
  
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

// 员工排名图表
const initEmployeeRankingChart = () => {
  if (!employeeRankingChartRef.value) return
  
  const chart = echarts.init(employeeRankingChartRef.value)
  const data = analyticsData.value.employeeRankings || []
  
  const option = {
    title: {
      text: '员工成绩排名 Top 10',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter: (params) => {
        const item = data[params[0].dataIndex]
        return `${item.personName}<br/>平均分: ${item.avgScore.toFixed(2)}分<br/>完成课程: ${item.courseCount}节`
      }
    },
    grid: {
      left: '10%',
      right: '10%',
      bottom: '10%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.personName),
      axisLabel: {
        rotate: 30,
        interval: 0,
        fontSize: 10
      }
    },
    yAxis: {
      type: 'value',
      name: '平均分',
      min: 0,
      max: 100
    },
    series: [
      {
        type: 'bar',
        data: data.map(item => item.avgScore.toFixed(2)),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#91cc75' },
            { offset: 1, color: '#5ab33c' }
          ])
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}分'
        }
      }
    ]
  }
  
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

// 讲师统计图表
const initTeacherStatsChart = () => {
  if (!teacherStatsChartRef.value) return
  
  const chart = echarts.init(teacherStatsChartRef.value)
  const data = analyticsData.value.teacherStatistics || []
  
  const option = {
    title: {
      text: '讲师授课统计',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' }
    },
    legend: {
      data: ['授课数量', '平均分', '学员数量'],
      top: 30
    },
    grid: {
      left: '10%',
      right: '10%',
      bottom: '10%',
      top: '20%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.teacherName),
      axisLabel: {
        rotate: 30,
        interval: 0,
        fontSize: 10
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '数量',
        position: 'left'
      },
      {
        type: 'value',
        name: '平均分',
        position: 'right',
        min: 0,
        max: 100
      }
    ],
    series: [
      {
        name: '授课数量',
        type: 'bar',
        data: data.map(item => item.courseCount),
        itemStyle: { color: '#5470c6' }
      },
      {
        name: '学员数量',
        type: 'bar',
        data: data.map(item => item.studentCount),
        itemStyle: { color: '#91cc75' }
      },
      {
        name: '平均分',
        type: 'line',
        yAxisIndex: 1,
        data: data.map(item => item.avgScore.toFixed(2)),
        itemStyle: { color: '#fac858' }
      }
    ]
  }
  
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

onMounted(() => {
  fetchAnalytics()
})
</script>

<style scoped>
.analytics-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h2 {
  margin: 0;
}

.chart-card {
  margin-bottom: 20px;
}
</style>
