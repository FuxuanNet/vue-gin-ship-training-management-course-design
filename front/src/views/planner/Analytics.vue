<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useMockDataStore } from '../../stores/mockData'
import * as echarts from 'echarts'

const mockDataStore = useMockDataStore()
const planChartRef = ref(null)
const courseChartRef = ref(null)
const typeChartRef = ref(null)

// 计算培训计划平均分
const planScores = computed(() => {
  return mockDataStore.trainingPlans.map(plan => {
    const planItems = mockDataStore.courseItems.filter(item => item.plan_id === plan.plan_id)
    const allScores = []
    
    planItems.forEach(item => {
      const evaluations = mockDataStore.attendanceEvaluations.filter(ae => ae.item_id === item.item_id)
      evaluations.forEach(ae => {
        const weightedScore = ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
        allScores.push(weightedScore)
      })
    })
    
    const avgScore = allScores.length > 0 
      ? (allScores.reduce((a, b) => a + b, 0) / allScores.length).toFixed(2)
      : 0
    
    return {
      plan_name: plan.plan_name,
      avg_score: parseFloat(avgScore),
      student_count: new Set(mockDataStore.attendanceEvaluations.filter(ae => 
        planItems.some(item => item.item_id === ae.item_id)
      ).map(ae => ae.person_id)).size
    }
  })
})

// 计算课程平均分
const courseScores = computed(() => {
  return mockDataStore.courses.map(course => {
    const courseItems = mockDataStore.courseItems.filter(item => item.course_id === course.course_id)
    const allScores = []
    
    courseItems.forEach(item => {
      const evaluations = mockDataStore.attendanceEvaluations.filter(ae => ae.item_id === item.item_id)
      evaluations.forEach(ae => {
        const weightedScore = ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
        allScores.push(weightedScore)
      })
    })
    
    const avgScore = allScores.length > 0 
      ? (allScores.reduce((a, b) => a + b, 0) / allScores.length).toFixed(2)
      : 0
    
    return {
      course_name: course.course_name,
      course_class: course.course_class,
      avg_score: parseFloat(avgScore),
      student_count: allScores.length
    }
  })
})

// 按课程类型统计
const typeScores = computed(() => {
  const typeMap = {}
  
  courseScores.value.forEach(course => {
    if (!typeMap[course.course_class]) {
      typeMap[course.course_class] = []
    }
    if (course.avg_score > 0) {
      typeMap[course.course_class].push(course.avg_score)
    }
  })
  
  return Object.entries(typeMap).map(([type, scores]) => ({
    course_class: type,
    avg_score: scores.length > 0 
      ? (scores.reduce((a, b) => a + b, 0) / scores.length).toFixed(2)
      : 0,
    course_count: scores.length
  }))
})

// 初始化培训计划对比图
const initPlanChart = () => {
  nextTick(() => {
    if (!planChartRef.value) return
    
    const chart = echarts.init(planChartRef.value)
    
    const option = {
      title: {
        text: '培训计划评分对比',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: planScores.value.map(p => p.plan_name),
        axisLabel: {
          rotate: 30,
          interval: 0
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
          name: '平均分',
          type: 'bar',
          data: planScores.value.map(p => p.avg_score),
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#83bff6' },
              { offset: 0.5, color: '#188df0' },
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
    
    window.addEventListener('resize', () => {
      chart.resize()
    })
  })
}

// 初始化课程评分图
const initCourseChart = () => {
  nextTick(() => {
    if (!courseChartRef.value) return
    
    const chart = echarts.init(courseChartRef.value)
    
    const option = {
      title: {
        text: '课程评分分析',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        }
      },
      legend: {
        data: ['平均分', '参与人数'],
        top: 30
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: 70,
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: courseScores.value.map(c => c.course_name),
        axisLabel: {
          rotate: 30,
          interval: 0
        }
      },
      yAxis: [
        {
          type: 'value',
          name: '平均分',
          min: 0,
          max: 100,
          position: 'left'
        },
        {
          type: 'value',
          name: '参与人数',
          position: 'right'
        }
      ],
      series: [
        {
          name: '平均分',
          type: 'bar',
          data: courseScores.value.map(c => c.avg_score),
          itemStyle: {
            color: '#5470c6'
          }
        },
        {
          name: '参与人数',
          type: 'line',
          yAxisIndex: 1,
          data: courseScores.value.map(c => c.student_count),
          itemStyle: {
            color: '#91cc75'
          }
        }
      ]
    }
    
    chart.setOption(option)
    
    window.addEventListener('resize', () => {
      chart.resize()
    })
  })
}

// 初始化课程类型饼图
const initTypeChart = () => {
  nextTick(() => {
    if (!typeChartRef.value) return
    
    const chart = echarts.init(typeChartRef.value)
    
    const option = {
      title: {
        text: '课程类型平均分分布',
        left: 'center'
      },
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c}分 ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        top: 'middle'
      },
      series: [
        {
          name: '平均分',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}\n{c}分'
          },
          labelLine: {
            show: true
          },
          data: typeScores.value.map(t => ({
            name: t.course_class,
            value: parseFloat(t.avg_score)
          }))
        }
      ]
    }
    
    chart.setOption(option)
    
    window.addEventListener('resize', () => {
      chart.resize()
    })
  })
}

onMounted(() => {
  initPlanChart()
  initCourseChart()
  initTypeChart()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>数据分析</h1>
      <p class="subtitle">全面分析培训效果，优化培训策略</p>
    </div>
    
    <!-- 统计卡片 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :xs="12" :sm="12" :md="6" :lg="6">
          <el-card>
            <el-statistic title="培训计划数" :value="mockDataStore.trainingPlans.length" suffix="个">
              <template #prefix>
                <el-icon><Notebook /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :xs="12" :sm="12" :md="6" :lg="6">
          <el-card>
            <el-statistic title="培训课程数" :value="mockDataStore.courses.length" suffix="门">
              <template #prefix>
                <el-icon><Reading /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :xs="12" :sm="12" :md="6" :lg="6">
          <el-card>
            <el-statistic title="参训员工数" :value="mockDataStore.persons.filter(p => p.role === '员工').length" suffix="人">
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :xs="12" :sm="12" :md="6" :lg="6">
          <el-card>
            <el-statistic title="讲师团队" :value="mockDataStore.persons.filter(p => p.role === '讲师').length" suffix="人">
              <template #prefix>
                <el-icon><UserFilled /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <!-- 图表区域 -->
    <div class="charts-section">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24" :md="24" :lg="24">
          <el-card class="chart-card">
            <div ref="planChartRef" class="chart"></div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :xs="24" :sm="24" :md="16" :lg="16">
          <el-card class="chart-card">
            <div ref="courseChartRef" class="chart"></div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="24" :md="8" :lg="8">
          <el-card class="chart-card" style="margin-top: 0;">
            <div ref="typeChartRef" class="chart chart-type"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <!-- 数据表格 -->
    <div class="tables-section">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24" :md="12" :lg="12">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>培训计划详情</span>
              </div>
            </template>
            <el-table :data="planScores" border max-height="400">
              <el-table-column prop="plan_name" label="计划名称" min-width="180" />
              <el-table-column prop="avg_score" label="平均分" width="100" align="center">
                <template #default="{ row }">
                  <el-tag 
                    :type="row.avg_score >= 85 ? 'success' : row.avg_score >= 75 ? '' : 'warning'"
                  >
                    {{ row.avg_score }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="student_count" label="参与人数" width="100" align="center" />
            </el-table>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="24" :md="12" :lg="12">
          <el-card style="margin-top: 0;">
            <template #header>
              <div class="card-header">
                <span>课程类型统计</span>
              </div>
            </template>
            <el-table :data="typeScores" border max-height="400">
              <el-table-column prop="course_class" label="课程类型" min-width="120" />
              <el-table-column prop="avg_score" label="平均分" width="100" align="center">
                <template #default="{ row }">
                  <el-tag 
                    :type="row.avg_score >= 85 ? 'success' : row.avg_score >= 75 ? '' : 'warning'"
                  >
                    {{ row.avg_score }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="course_count" label="课程数" width="100" align="center" />
            </el-table>
          </el-card>
        </el-col>
      </el-row>
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

.stats-section {
  margin-bottom: 30px;
}

.charts-section {
  margin-bottom: 30px;
}

.chart-card {
  height: 100%;
}

.chart {
  width: 100%;
  height: 400px;
}

.tables-section {
  margin-bottom: 30px;
}

.card-header {
  font-weight: 600;
  font-size: 16px;
}

.chart-type {
  height: 400px;
}

/* 响应式设计 - 移动端 */
@media (max-width: 768px) {
  .page-container {
    padding: 15px;
  }
  
  .page-header h1 {
    font-size: 22px;
  }
  
  /* 统计卡片在手机端两列布局 */
  .stats-section :deep(.el-col) {
    margin-bottom: 15px;
  }
  
  /* 图表在手机端一行一个 */
  .charts-section :deep(.el-col) {
    margin-bottom: 15px;
  }
  
  .chart {
    height: 300px;
  }
  
  .chart-type {
    height: 300px;
  }
  
  /* 表格在手机端一行一个 */
  .tables-section :deep(.el-col) {
    margin-bottom: 15px;
  }
}

@media (max-width: 576px) {
  .page-header h1 {
    font-size: 20px;
  }
  
  .chart {
    height: 250px;
  }
  
  .chart-type {
    height: 250px;
  }
}
</style>
