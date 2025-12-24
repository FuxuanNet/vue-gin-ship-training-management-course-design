<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const carouselItems = [
  {
    id: 1,
    title: '数据资源交易平台',
    description: '安全可信的数据资源交易，让数据流动创造更大价值',
    image: 'https://via.placeholder.com/1200x400?text=数据资源交易平台'
  },
  {
    id: 2,
    title: '智能沙箱环境',
    description: '安全隔离的数据处理环境，保障数据使用安全合规',
    image: 'https://via.placeholder.com/1200x400?text=智能沙箱环境'
  },
  {
    id: 3,
    title: '区块链存证技术',
    description: '全程数据操作记录上链，确保数据流通全程可监管可追溯',
    image: 'https://via.placeholder.com/1200x400?text=区块链存证技术'
  }
]

const featuredResources = ref([
  {
    id: 1,
    title: '金融行业客户画像数据集',
    description: '包含10万+金融客户的完整画像数据，适用于信用评估模型训练',
    provider: '金融科技有限公司',
    price: '5000',
    downloads: 325,
    rating: 4.8,
    tags: ['金融', '客户画像', '结构化数据']
  },
  {
    id: 2,
    title: '智能医疗影像分析算法',
    description: '基于深度学习的医疗CT影像分析算法，准确率达98.6%',
    provider: '智能医疗科技公司',
    price: '12000',
    downloads: 156,
    rating: 4.9,
    tags: ['医疗', '算法', 'AI']
  },
  {
    id: 3,
    title: '高性能分布式计算资源',
    description: '提供高性能GPU计算集群，支持大规模深度学习模型训练',
    provider: '云算力服务提供商',
    price: '200/小时',
    downloads: 89,
    rating: 4.7,
    tags: ['算力', 'GPU', '分布式计算']
  },
  {
    id: 4,
    title: '智能交通数据分析包',
    description: '包含城市道路实时流量、红绿灯配时等数据，适用于智慧城市建设',
    provider: '城市大数据中心',
    price: '8800',
    downloads: 213,
    rating: 4.6,
    tags: ['交通', '智慧城市', '实时数据']
  }
])

const platformStats = [
  { icon: 'DataLine', label: '数据资源', value: '10,580+' },
  { icon: 'Connection', label: '算力资源', value: '256+' },
  { icon: 'SetUp', label: '算法插件', value: '1,342+' },
  { icon: 'User', label: '平台用户', value: '8,670+' },
  { icon: 'Promotion', label: '成交额', value: '¥3.52亿' }
]

const navigateTo = (path) => {
  router.push(path)
}
</script>

<template>
  <div class="home-container">
    <!-- 头部大屏背景 -->
    <div class="hero-section">
      <div class="hero-content">
        <div class="title-line"></div>
        <h1>智信链枢</h1>
        <div class="subtitle-wrapper">
          <div class="subtitle-line"></div>
          <p>基于聚合交易与沙箱协作的可信智能监管资源平台</p>
          <div class="subtitle-line"></div>
        </div>
      </div>
    </div>

    <div class="content-container">
      <!-- 平台统计数据 -->
      <div class="stats-section">
        <el-row :gutter="20" justify="center">
          <el-col :span="24 / platformStats.length" v-for="stat in platformStats" :key="stat.label">
            <div class="stat-card">
              <el-icon class="stat-icon"><component :is="stat.icon" /></el-icon>
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 平台主要功能模块 -->
      <div class="features-section">
        <h2 class="section-title">平台主要功能</h2>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="feature-card" @click="navigateTo('/upload')">
              <template #header>
                <div class="feature-header">
                  <el-icon><Upload /></el-icon>
                  <span>资源上传</span>
                </div>
              </template>
              <div class="feature-content">
                <p>上传数据资源、算力资源和算法资源，支持多种数据类型和格式。</p>
                <el-button type="primary" plain>立即上传</el-button>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="feature-card" @click="navigateTo('/publish')">
              <template #header>
                <div class="feature-header">
                  <el-icon><Shop /></el-icon>
                  <span>交易市场</span>
                </div>
              </template>
              <div class="feature-content">
                <p>将您的资源发布到交易市场，或浏览购买其他用户提供的高质量资源。</p>
                <el-button type="primary" plain @click.stop="navigateTo('/market')">进入市场</el-button>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="feature-card" @click="navigateTo('/sandbox')">
              <template #header>
                <div class="feature-header">
                  <el-icon><Monitor /></el-icon>
                  <span>沙箱环境</span>
                </div>
              </template>
              <div class="feature-content">
                <p>在安全隔离的环境中使用您的数据和算法，保障数据使用安全合规。</p>
                <el-button type="primary" plain>进入沙箱</el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 精选资源 -->
      <div class="resources-section">
        <h2 class="section-title">精选资源</h2>
        <el-row :gutter="20">
          <el-col :lg="6" :md="12" v-for="resource in featuredResources" :key="resource.id">
            <el-card class="resource-card" shadow="hover">
              <h3 class="resource-title">{{ resource.title }}</h3>
              <div class="resource-tags">
                <el-tag v-for="tag in resource.tags" :key="tag" size="small" effect="plain">{{ tag }}</el-tag>
              </div>
              <p class="resource-description">{{ resource.description }}</p>
              <div class="resource-meta">
                <div class="resource-provider">
                  <el-icon><User /></el-icon>
                  <span>{{ resource.provider }}</span>
                </div>
                <div class="resource-rating">
                  <el-rate v-model="resource.rating" disabled text-color="#ff9900" />
                </div>
              </div>
              <div class="resource-footer">
                <div class="resource-price">￥{{ resource.price }}</div>
                <el-button type="primary" size="small" @click="navigateTo('/market')">查看详情</el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 平台优势 -->
      <div class="advantages-section">
        <h2 class="section-title">平台优势</h2>
        <el-row :gutter="30">
          <el-col :span="8">
            <div class="advantage-item">
              <el-icon class="advantage-icon"><Lock /></el-icon>
              <h3>数据安全可控</h3>
              <p>基于区块链和沙箱技术，确保数据使用全程可控可溯源</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="advantage-item">
              <el-icon class="advantage-icon"><Connection /></el-icon>
              <h3>资源汇聚共享</h3>
              <p>打破数据孤岛，实现多方数据和算力资源的高效连接</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="advantage-item">
              <el-icon class="advantage-icon"><Money /></el-icon>
              <h3>智能合约交易</h3>
              <p>基于区块链智能合约的交易机制，确保交易公平透明</p>
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
}

/* 大屏样式保持全宽 */
.hero-section {
  width: 100%;
  height: 600px;
  background-image: url('/主页大屏.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  position: relative;
  margin-bottom: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

/* 新增内容容器样式 */
.content-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%);
  pointer-events: none;
}

.hero-content {
  text-align: center;
  color: #2b3a6d;
  padding: 40px;
  position: relative;
  z-index: 2;
}

.title-line {
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, transparent, #2b3a6d, transparent);
  margin: 0 auto 20px;
}

.hero-content h1 {
  font-size: 56px;
  margin: 0 0 20px;
  font-weight: 600;
  letter-spacing: 4px;
  text-shadow: 0 2px 10px rgba(255,255,255,0.3);
  background: linear-gradient(180deg, #2b3a6d, #1e4d92);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  position: relative;
}

.subtitle-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.subtitle-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(43,58,109,0.5), transparent);
}

.hero-content p {
  font-size: 24px;
  margin: 0;
  font-weight: 500;
  color: #2b3a6d;
  text-shadow: 0 1px 4px rgba(255,255,255,0.3);
  white-space: nowrap;
}

@media (max-width: 768px) {
  .hero-content h1 {
    font-size: 40px;
  }
  
  .hero-content p {
    font-size: 18px;
    white-space: normal;
  }
  
  .subtitle-wrapper {
    flex-direction: column;
    gap: 10px;
  }
  
  .subtitle-line {
    width: 100px;
  }
}

/* 平台统计数据样式 */
.stats-section {
  margin-bottom: 50px;
  display: flex;
  justify-content: center;
}

.stat-card {
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 6px;
  text-align: center;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 32px;
  color: #409EFF;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

/* 功能模块样式 */
.features-section {
  margin-bottom: 50px;
}

.feature-card {
  cursor: pointer;
  height: 100%;
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.feature-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
}

.feature-header .el-icon {
  margin-right: 10px;
  font-size: 20px;
  color: #409EFF;
}

.feature-content {
  padding: 10px 0;
  height: 120px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

/* 精选资源样式 */
.resources-section {
  margin-bottom: 50px;
}

.resource-card {
  height: 100%;
  margin-bottom: 20px;
  position: relative;
}

.resource-title {
  font-size: 16px;
  margin-bottom: 10px;
  height: 44px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.resource-tags {
  margin-bottom: 10px;
}

.resource-tags .el-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.resource-description {
  color: #606266;
  font-size: 14px;
  margin-bottom: 15px;
  height: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.resource-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.resource-provider {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.resource-provider .el-icon {
  margin-right: 5px;
}

.resource-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.resource-price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

/* 平台优势样式 */
.advantages-section {
  margin-bottom: 40px;
}

.advantage-item {
  text-align: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 6px;
  height: 100%;
  transition: all 0.3s;
}

.advantage-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.advantage-icon {
  font-size: 40px;
  color: #409EFF;
  margin-bottom: 15px;
}

.advantage-item h3 {
  font-size: 18px;
  margin-bottom: 10px;
  color: #303133;
}

.advantage-item p {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}
</style>