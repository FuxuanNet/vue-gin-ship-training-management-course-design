<script setup>
import { computed, ref } from 'vue'
import { getScores, getCourseTypeScores } from '../../api/employee'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { onMounted, nextTick } from 'vue'

const loading = ref(true)
const activeTab = ref('detail')
const radarChartRef = ref(null)
const myScores = ref([])
const courseTypeScores = ref([])

// 获取员工成绩
const fetchScores = async () => {
  loading.value = true
  try {
    const response = await getScores()
    if (response.code === 200) {
      myScores.value = response.data.scores || []
    }
  } catch (error) {
    ElMessage.error('获取成绩失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 获取课程类型成绩
const fetchCourseTypeScores = async () => {
  try {
    const response = await getCourseTypeScores()
    if (response.code === 200) {
      courseTypeScores.value = response.data.courseTypeScores || []
      // 重新初始化雷达图
      if (activeTab.value === 'analysis') {
        initRadarChart()
      }
    }
  } catch (error) {
    ElMessage.error('获取类型成绩失败：' + (error.message || '未知错误'))
  }
}

// 平均分
const averageScore = computed(() => {
  if (myScores.value.length === 0) return 0
  const sum = myScores.value.reduce((acc, item) => acc + parseFloat(item.weightedScore), 0)
  return (sum / myScores.value.length).toFixed(2)
})

// 初始化雷达图
const initRadarChart = () => {
  nextTick(() => {
    if (!radarChartRef.value) return
    if (courseTypeScores.value.length === 0) return
    
    const chart = echarts.init(radarChartRef.value)
    
    const indicator = courseTypeScores.value.map(item => ({
      name: item.courseClass,
      max: 100
    }))
    
    const data = courseTypeScores.value.map(item => parseFloat(item.avgScore))
    
    const option = {
      title: {
        text: '课程类型掌握度分析',
        left: 'center',
        top: 10,
        textStyle: {
          fontSize: 16,
          fontWeight: 600
        }
      },
      tooltip: {
        trigger: 'item'
      },
      radar: {
        indicator: indicator,
        radius: '65%',
        splitNumber: 5,
        name: {
          textStyle: {
            color: '#606266',
            fontSize: 13
          }
        },
        splitArea: {
          areaStyle: {
            color: ['rgba(64, 158, 255, 0.1)', 'rgba(64, 158, 255, 0.05)']
          }
        },
        axisLine: {
          lineStyle: {
            color: 'rgba(64, 158, 255, 0.3)'
          }
        },
        splitLine: {
          lineStyle: {
            color: 'rgba(64, 158, 255, 0.3)'
          }
        }
      },
      series: [{
        type: 'radar',
        data: [{
          value: data,
          name: '我的掌握度',
          areaStyle: {
            color: 'rgba(64, 158, 255, 0.3)'
          },
          lineStyle: {
            color: '#409EFF',
            width: 2
          },
          itemStyle: {
            color: '#409EFF'
          }
        }]
      }]
    }
    
    chart.setOption(option)
    
    // 响应式
    window.addEventListener('resize', () => {
      chart.resize()
    })
  })
}

// 监听标签切换
const handleTabChange = (tab) => {
  if (tab === 'analysis') {
    if (courseTypeScores.value.length === 0) {
      fetchCourseTypeScores()
    } else {
      initRadarChart()
    }
  }
}

onMounted(() => {
  fetchScores()
  if (activeTab.value === 'analysis') {
    fetchCourseTypeScores()
  }
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>我的成绩</h1>
      <div class="score-summary">
        <el-statistic title="平均成绩" :value="averageScore" suffix="分">
          <template #prefix>
            <el-icon><TrendCharts /></el-icon>
          </template>
        </el-statistic>
        <el-statistic title="已完成课程" :value="myScores.length" suffix="门">
          <template #prefix>
            <el-icon><Reading /></el-icon>
          </template>
        </el-statistic>
      </div>
    </div>
    
    <div class="content-wrapper" v-loading="loading">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="成绩明细" name="detail">
          <el-empty v-if="myScores.length === 0" description="暂无成绩记录">
            <el-button type="primary" @click="$router.push('/employee/evaluation')">去评价</el-button>
          </el-empty>
          <el-table v-else :data="myScores" border>
            <el-table-column prop="classDate" label="上课日期" width="120" />
            <el-table-column label="课程名称" min-width="180">
              <template #default="{ row }">
                {{ row.courseName }}
              </template>
            </el-table-column>
            <el-table-column label="课程类型" width="120">
              <template #default="{ row }">
                <el-tag size="small">{{ row.courseClass }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="selfScore" label="自评分" width="100" align="center">
              <template #default="{ row }">
                <span class="score-value">{{ row.selfScore || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="teacherScore" label="讲师评分" width="100" align="center">
              <template #default="{ row }">
                <span class="score-value">{{ row.teacherScore || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="weightedScore" label="综合得分" width="100" align="center">
              <template #default="{ row }">
                <el-tag 
                  v-if="row.weightedScore"
                  :type="row.weightedScore >= 90 ? 'success' : row.weightedScore >= 80 ? '' : row.weightedScore >= 60 ? 'warning' : 'danger'"
                  effect="dark"
                >
                  {{ row.weightedScore.toFixed(1) }}
                </el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="selfComment" label="自评内容" min-width="200" show-overflow-tooltip />
            <el-table-column prop="teacherComment" label="讲师评语" min-width="200" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="能力分析" name="analysis">
          <div class="analysis-container">
            <div class="chart-container">
              <div ref="radarChartRef" class="radar-chart"></div>
            </div>
            
            <div class="type-scores">
              <h3>各类型平均分</h3>
              <el-empty v-if="courseTypeScores.length === 0" description="暂无数据" />
              <el-row v-else :gutter="20">
                <el-col :span="12" v-for="item in courseTypeScores" :key="item.courseClass">
                  <el-card class="type-score-card" shadow="hover">
                    <div class="type-name">{{ item.courseClass }}</div>
                    <div class="type-score">{{ item.avgScore }}分</div>
                    <el-progress 
                      :percentage="parseFloat(item.avgScore)" 
                      :color="parseFloat(item.avgScore) >= 85 ? '#67C23A' : '#409EFF'"
                    />
                  </el-card>
                </el-col>
              </el-row>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: 28px;
  color: #303133;
}

.score-summary {
  display: flex;
  gap: 40px;
}

.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.score-value {
  font-weight: 600;
  color: #409EFF;
}

.analysis-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.chart-container {
  display: flex;
  justify-content: center;
}

.radar-chart {
  width: 600px;
  height: 450px;
}

.type-scores h3 {
  margin-bottom: 20px;
  color: #303133;
  font-size: 18px;
}

.type-score-card {
  margin-bottom: 15px;
  text-align: center;
}

.type-name {
  font-size: 16px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 10px;
}

.type-score {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 15px;
}
</style>
