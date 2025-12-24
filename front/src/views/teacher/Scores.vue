<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useMockDataStore } from '../../stores/mockData'
import { useUserStore } from '../../stores/user'
import * as echarts from 'echarts'

const mockDataStore = useMockDataStore()
const userStore = useUserStore()
const chartRef = ref(null)
const selectedCourseId = ref(null)

// 获取讲师教授的课程
const myCourses = computed(() => {
  return mockDataStore.courses.filter(c => c.teacher_id === userStore.userInfo.id)
})

// 获取选中课程的所有学员成绩
const courseScores = computed(() => {
  if (!selectedCourseId.value) return []
  
  // 找出该课程的所有课程安排
  const courseItems = mockDataStore.courseItems.filter(
    item => item.course_id === selectedCourseId.value
  )
  
  // 获取所有学员的成绩
  const scores = []
  courseItems.forEach(item => {
    const evaluations = mockDataStore.attendanceEvaluations.filter(
      ae => ae.item_id === item.item_id
    )
    evaluations.forEach(ae => {
      const person = mockDataStore.persons.find(p => p.person_id === ae.person_id)
      const weightedScore = ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
      scores.push({
        ...ae,
        item,
        person,
        weighted_score: weightedScore.toFixed(2)
      })
    })
  })
  
  return scores
})

// 课程平均分
const courseAverageScore = computed(() => {
  if (courseScores.value.length === 0) return 0
  const sum = courseScores.value.reduce((acc, item) => acc + parseFloat(item.weighted_score), 0)
  return (sum / courseScores.value.length).toFixed(2)
})

// 课程统计信息
const courseStats = computed(() => {
  const scores = courseScores.value.map(s => parseFloat(s.weighted_score))
  if (scores.length === 0) return { max: 0, min: 0, avg: 0, count: 0 }
  
  return {
    max: Math.max(...scores).toFixed(2),
    min: Math.min(...scores).toFixed(2),
    avg: courseAverageScore.value,
    count: scores.length
  }
})

// 初始化成绩分布图表
const initChart = () => {
  nextTick(() => {
    if (!chartRef.value || courseScores.value.length === 0) return
    
    const chart = echarts.init(chartRef.value)
    
    // 按学员分组数据
    const studentData = {}
    courseScores.value.forEach(score => {
      const name = score.person.name
      if (!studentData[name]) {
        studentData[name] = []
      }
      studentData[name].push({
        date: score.item.class_date,
        score: parseFloat(score.weighted_score)
      })
    })
    
    // 准备图表数据
    const students = Object.keys(studentData)
    const series = students.map(student => ({
      name: student,
      type: 'line',
      data: studentData[student].map(d => d.score),
      smooth: true,
      symbol: 'circle',
      symbolSize: 8
    }))
    
    // 获取所有日期
    const dates = [...new Set(courseScores.value.map(s => s.item.class_date))].sort()
    
    const option = {
      title: {
        text: '学员成绩趋势',
        left: 'center',
        top: 10
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        }
      },
      legend: {
        data: students,
        top: 40,
        type: 'scroll'
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
        data: dates,
        boundaryGap: false
      },
      yAxis: {
        type: 'value',
        min: 0,
        max: 100,
        name: '分数'
      },
      series: series
    }
    
    chart.setOption(option)
    
    window.addEventListener('resize', () => {
      chart.resize()
    })
  })
}

// 监听课程选择变化
const handleCourseChange = () => {
  initChart()
}

onMounted(() => {
  if (myCourses.value.length > 0) {
    selectedCourseId.value = myCourses.value[0].course_id
    initChart()
  }
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
      >
        <el-option 
          v-for="course in myCourses" 
          :key="course.course_id"
          :label="course.course_name"
          :value="course.course_id"
        />
      </el-select>
    </div>
    
    <div class="stats-section" v-if="selectedCourseId">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card>
            <el-statistic title="参与人数" :value="courseStats.count" suffix="人">
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic title="平均分" :value="courseStats.avg" suffix="分">
              <template #prefix>
                <el-icon><TrendCharts /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic title="最高分" :value="courseStats.max" suffix="分">
              <template #prefix>
                <el-icon><Top /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic title="最低分" :value="courseStats.min" suffix="分">
              <template #prefix>
                <el-icon><Bottom /></el-icon>
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
      <div class="table-section">
        <h3>成绩明细</h3>
        <el-table :data="courseScores" border>
          <el-table-column prop="person.name" label="学员姓名" width="120" />
          <el-table-column prop="item.class_date" label="上课日期" width="120" />
          <el-table-column prop="self_score" label="自评分" width="100" align="center">
            <template #default="{ row }">
              <span class="score-value">{{ row.self_score }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="teacher_score" label="讲师评分" width="100" align="center">
            <template #default="{ row }">
              <span class="score-value">{{ row.teacher_score }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="weighted_score" label="综合得分" width="100" align="center">
            <template #default="{ row }">
              <el-tag 
                :type="row.weighted_score >= 90 ? 'success' : row.weighted_score >= 80 ? '' : row.weighted_score >= 60 ? 'warning' : 'danger'"
                effect="dark"
              >
                {{ row.weighted_score }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="self_comment" label="自评内容" min-width="180" show-overflow-tooltip />
          <el-table-column prop="teacher_comment" label="讲师评语" min-width="180" show-overflow-tooltip />
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
