<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { computed, ref, onMounted } from 'vue'
import { getStatistics } from '../api/home'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.isLoggedIn)
const userRole = computed(() => userStore.userRole)

// 平台统计数据（从后端获取）
const platformStats = ref([
  { icon: 'Reading', label: '培训课程', value: '0', color: '#409EFF' },
  { icon: 'User', label: '讲师团队', value: '0', color: '#67C23A' },
  { icon: 'Notebook', label: '培训计划', value: '0', color: '#E6A23C' },
  { icon: 'TrendCharts', label: '平均满意度', value: '0%', color: '#F56C6C' }
])

// 加载统计数据
const loadStatistics = async () => {
  try {
    const response = await getStatistics()
    const data = response.data
    
    // 更新平台统计数据
    platformStats.value = [
      { icon: 'Reading', label: '培训课程', value: `${data.courseCount}+`, color: '#409EFF' },
      { icon: 'User', label: '讲师团队', value: `${data.teacherCount}+`, color: '#67C23A' },
      { icon: 'Notebook', label: '培训计划', value: `${data.planCount}+`, color: '#E6A23C' },
      { icon: 'TrendCharts', label: '平均满意度', value: `${data.averageSatisfaction}%`, color: '#F56C6C' }
    ]
  } catch (error) {
    console.error('加载统计数据失败:', error)
    // 错误已在拦截器中处理，这里静默处理
  }
}

// 页面加载时获取统计数据
onMounted(() => {
  loadStatistics()
})

const features = [
  {
    icon: 'Calendar',
    title: '智能课程表',
    description: '清晰展示每日、每周课程安排，支持多端同步查看',
    path: userRole.value === 'employee' ? '/employee/schedule' : userRole.value === 'teacher' ? '/teacher/schedule' : '/planner/plans'
  },
  {
    icon: 'Document',
    title: '在线评价系统',
    description: '员工自评与讲师评价相结合，全面评估培训效果',
    path: userRole.value === 'employee' ? '/employee/evaluation' : userRole.value === 'teacher' ? '/teacher/grading' : '/planner/analytics'
  },
  {
    icon: 'DataAnalysis',
    title: '数据分析看板',
    description: '可视化展示培训数据，助力决策优化',
    path: userRole.value === 'planner' ? '/planner/analytics' : '/employee/scores'
  }
]

const advantages = [
  {
    icon: 'Management',
    title: '培训管理规范',
    description: '统一制定和管理各类培训计划，确保培训有序进行'
  },
  {
    icon: 'Tickets',
    title: '课程信息集中',
    description: '集中展示课程信息、时间、讲师等详细内容'
  },
  {
    icon: 'Medal',
    title: '评估机制完善',
    description: '通过自评和讲师评价，建立合理的评估机制'
  }
]

const navigateTo = (path) => {
  if (isLoggedIn.value) {
    router.push(path)
  } else {
    router.push('/login')
  }
}
</script>

<template>
  <div class="home-container">
    <!-- 头部大屏背景 -->
    <div class="hero-section">
      <div class="hero-content">
        <div class="title-line"></div>
        <h1>船舶企业培训管理系统</h1>
        <div class="subtitle-wrapper">
          <div class="subtitle-line"></div>
          <p>提升员工技能水平，助力企业持续发展</p>
          <div class="subtitle-line"></div>
        </div>
        <div class="hero-actions" v-if="!isLoggedIn">
          <el-button type="primary" size="large" @click="router.push('/login')">立即登录</el-button>
        </div>
      </div>
    </div>

    <div class="content-container">
      <!-- 平台统计数据 -->
      <div class="stats-section">
        <el-row :gutter="20" justify="center">
          <el-col :xs="12" :sm="12" :md="6" v-for="stat in platformStats" :key="stat.label">
            <div class="stat-card">
              <el-icon class="stat-icon" :style="{ color: stat.color }">
                <component :is="stat.icon" />
              </el-icon>
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 平台主要功能模块 -->
      <div class="features-section" v-if="isLoggedIn">
        <h2 class="section-title">快速导航</h2>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" v-for="(feature, index) in features" :key="index">
            <el-card class="feature-card" @click="navigateTo(feature.path)" shadow="hover">
              <div class="feature-icon">
                <el-icon><component :is="feature.icon" /></el-icon>
              </div>
              <h3>{{ feature.title }}</h3>
              <p>{{ feature.description }}</p>
              <el-button type="primary" plain size="small">立即使用</el-button>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 平台优势 -->
      <div class="advantages-section">
        <h2 class="section-title">系统特色</h2>
        <el-row :gutter="30">
          <el-col :xs="24" :sm="12" :md="8" v-for="advantage in advantages" :key="advantage.title">
            <div class="advantage-item">
              <el-icon class="advantage-icon"><component :is="advantage.icon" /></el-icon>
              <h3>{{ advantage.title }}</h3>
              <p>{{ advantage.description }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  width: 100%;
  margin: 0;
  padding: 0;
  min-height: 100vh;
}

/* 大屏样式 */
.hero-section {
  width: 100%;
  height: 520px;
  background-image: url('/image.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  position: relative;
  margin-bottom: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(2px);
  pointer-events: none;
}

.hero-content {
  text-align: center;
  color: #2c3e50;
  padding: 40px;
  position: relative;
  z-index: 2;
}

.title-line {
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, transparent, #2c3e50, transparent);
  margin: 0 auto 20px;
}

.hero-content h1 {
  font-size: 48px;
  margin: 0 0 20px;
  font-weight: 600;
  letter-spacing: 3px;
  text-shadow: 0 2px 8px rgba(255, 255, 255, 0.5);
}

.subtitle-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  max-width: 800px;
  margin: 0 auto 30px;
}

.subtitle-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(44, 62, 80, 0.5), transparent);
}

.hero-content p {
  font-size: 20px;
  margin: 0;
  font-weight: 400;
  color: #34495e;
  white-space: nowrap;
}

.hero-actions {
  margin-top: 30px;
}

/* 内容容器 */
.content-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px 60px;
}

/* 统计数据样式 */
.stats-section {
  margin-bottom: 60px;
}

.stat-card {
  background-color: #fff;
  padding: 30px 20px;
  border-radius: 8px;
  text-align: center;
  transition: all 0.3s;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  font-size: 36px;
  margin-bottom: 12px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

/* 功能模块样式 */
.features-section {
  margin-bottom: 60px;
}

.section-title {
  text-align: center;
  font-size: 28px;
  margin-bottom: 40px;
  color: #303133;
  font-weight: 600;
  position: relative;
}

.section-title::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 3px;
  background: #409EFF;
}

.feature-card {
  cursor: pointer;
  text-align: center;
  padding: 20px;
  min-height: 240px;
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.feature-icon {
  font-size: 48px;
  color: #409EFF;
  margin-bottom: 15px;
}

.feature-card h3 {
  font-size: 18px;
  margin-bottom: 12px;
  color: #303133;
}

.feature-card p {
  color: #606266;
  font-size: 14px;
  margin-bottom: 20px;
  line-height: 1.6;
}

/* 平台优势样式 */
.advantages-section {
  margin-bottom: 40px;
}

.advantage-item {
  text-align: center;
  padding: 30px 20px;
  background-color: #fff;
  border-radius: 8px;
  height: 100%;
  transition: all 0.3s;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.advantage-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.advantage-icon {
  font-size: 42px;
  color: #409EFF;
  margin-bottom: 15px;
}

.advantage-item h3 {
  font-size: 18px;
  margin-bottom: 12px;
  color: #303133;
  font-weight: 600;
}

.advantage-item p {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .hero-content h1 {
    font-size: 32px;
  }
  
  .hero-content p {
    font-size: 16px;
    white-space: normal;
  }
  
  .subtitle-wrapper {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
