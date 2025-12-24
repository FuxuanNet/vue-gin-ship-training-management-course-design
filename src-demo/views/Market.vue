<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { 
  getMarketResources, 
  getMarketResourceDetail, 
  getResourceCategories, 
  getPopularTags, 
  favoriteResource, 
  unfavoriteResource,
  getResourceSample,
  purchaseResource,
  getPurchasedResources
} from '../api/market'
import { ElMessage, ElLoading } from 'element-plus'
import { marked } from 'marked'

// 配置marked选项
marked.setOptions({
  breaks: true, // 将回车转换为<br>
  gfm: true,    // 启用GitHub风格Markdown
  headerIds: false, // 关闭自动添加header ids
  mangle: false // 关闭转义
})

// 资源类型
const resourceTypes = [
  { value: 'all', label: '全部资源' },
  { value: 'data', label: '数据资源', icon: 'DataLine' },
  { value: 'computing', label: '算力资源', icon: 'Connection' },
  { value: 'algorithm', label: '算法资源', icon: 'SetUp' }
]

// 当前选中的资源类型
const activeType = ref('all')

// 资源列表
const resourceList = ref([])

// 搜索条件
const searchForm = reactive({
  keyword: '',
  priceRange: [0, 100000],
  tags: [],
  sortBy: 'newest',
  ratingAbove: 0
})

// 高级搜索抽屉是否可见
const advancedSearchVisible = ref(false)

// 标签选项
const tagOptions = ref([])

// 排序选项
const sortOptions = [
  { value: 'newest', label: '最新发布' },
  { value: 'popular', label: '最多销量' },
  { value: 'price-asc', label: '价格从低到高' },
  { value: 'price-desc', label: '价格从高到低' }
]

// 加载状态
const loading = ref(false)

// 根据当前选中类型过滤资源
const filteredResources = computed(() => {
  // 已经在API请求时进行了筛选，前端不需要再次筛选
  return resourceList.value
})

// 资源类型标签颜色映射
const typeColorMap = {
  data: 'success',
  computing: 'warning',
  algorithm: 'danger'
}

// 资源类型中文映射
const typeNameMap = {
  data: '数据资源',
  computing: '算力资源',
  algorithm: '算法资源'
}

// 添加购买相关状态
const purchasing = ref(false)

// 获取市场资源
const fetchMarketResources = async () => {
  loading.value = true
  try {
    const params = {
      type: activeType.value === 'all' ? '' : activeType.value,
      keyword: searchForm.keyword,
      price_min: searchForm.priceRange[0],
      price_max: searchForm.priceRange[1],
      tags: searchForm.tags.join(','),
      sort_by: searchForm.sortBy
    }
    
    const res = await getMarketResources(params)
    if (res.code === 200) {
      // 对资源列表进行去重处理
      const uniqueResources = [];
      const resourceMap = new Map();
      
      res.data.forEach(resource => {
        const key = `${resource.resource_id || resource.id}-${resource.resource_type}`;
        if (!resourceMap.has(key)) {
          resourceMap.set(key, true);
          uniqueResources.push(resource);
        }
      });
      
      resourceList.value = uniqueResources;
      console.log('去重后的资源列表:', uniqueResources.length);
    } else {
      ElMessage.error(res.message || '获取市场资源失败')
    }
  } catch (error) {
    console.error('获取市场资源错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 获取标签选项
const fetchTags = async () => {
  try {
    const res = await getPopularTags(activeType.value, 50)
    if (res.code === 200) {
      tagOptions.value = res.data.map(tag => ({
        value: tag,
        label: tag
      }))
    }
  } catch (error) {
    console.error('获取热门标签失败:', error)
  }
}

// 查看资源详情
const resourceDetailDialog = ref(false)
const detailResource = ref(null)

// 样例数据对话框状态
const sampleDataDialog = ref(false)
const sampleData = ref(null)
const sampleLoading = ref(false)

// 文档对话框状态
const documentDialog = ref(false)
const renderedDocument = ref({
  documentation: '',
  usageGuide: '',
  apiDoc: ''
})

// 打开资源详情对话框
const openResourceDetail = async (resource) => {
  try {
    const loadingInstance = ElLoading.service({
      target: '.resources-section',
      text: '加载资源详情...'
    })
    
    const res = await getMarketResourceDetail(resource.resource_id || resource.id)
    loadingInstance.close()
    
    if (res.code === 200) {
      detailResource.value = res.data
      resourceDetailDialog.value = true
    } else {
      ElMessage.error(res.message || '获取资源详情失败')
    }
  } catch (error) {
    console.error('获取资源详情失败:', error)
    ElMessage.error('网络错误，请稍后重试')
  }
}

// 重置搜索条件
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.priceRange = [0, 100000]
  searchForm.tags = []
  searchForm.sortBy = 'newest'
  searchForm.ratingAbove = 0
  advancedSearchVisible.value = false
}

// 搜索按钮点击
const handleSearch = () => {
  advancedSearchVisible.value = false
  fetchMarketResources()
}

// 查看样例数据
const viewSampleData = async () => {
  if (!detailResource.value) return
  
  const resourceId = detailResource.value.id
  sampleLoading.value = true
  
  try {
    const response = await getResourceSample(resourceId)
    
    // 检查响应类型
    const contentType = response.headers?.['content-type']
    console.log('响应ContentType:', contentType)
    
    if (contentType && contentType.includes('application/json')) {
      // JSON响应表示出错
      sampleDataDialog.value = true
      sampleData.value = "该资源没有提供示例数据"
    } else {
      // 如果是文件，触发下载
      const contentDisposition = response.headers?.['content-disposition']
      let fileName = `${detailResource.value.name}_sample`
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename=([^;]+)/)
        if (filenameMatch && filenameMatch[1]) {
          fileName = filenameMatch[1].replace(/"/g, '')
        }
      }
      
      const blob = new Blob([response.data], { type: contentType || 'application/octet-stream' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', fileName)
      document.body.appendChild(link)
      link.click()
      
      setTimeout(() => {
        window.URL.revokeObjectURL(url)
        document.body.removeChild(link)
      }, 100)
      
      ElMessage.success('示例数据下载已开始')
    }
  } catch (error) {
    console.error('获取示例数据失败:', error)
    
    // 打开对话框显示错误
    sampleDataDialog.value = true
    sampleData.value = "该资源没有提供示例数据"
  } finally {
    sampleLoading.value = false
  }
}

// 渲染Markdown内容
const renderMarkdown = (content) => {
  if (!content) return '没有提供文档内容';
  return marked.parse(content);
}

// 查看说明文档
const viewDocument = () => {
  if (!detailResource.value) return;
  
  documentDialog.value = true;
  const resource = detailResource.value;
  
  // 根据资源类型渲染不同的文档
  if (resource.resource_type === 'data') {
    renderedDocument.value.documentation = renderMarkdown(resource.documentation);
  } else if (resource.resource_type === 'algorithm') {
    renderedDocument.value.usageGuide = renderMarkdown(resource.usage_guide);
    renderedDocument.value.apiDoc = renderMarkdown(resource.api_documentation);
  }
}

// 收藏资源
const toggleFavorite = async (resource, event) => {
  event.stopPropagation(); // 阻止冒泡，避免触发卡片的点击事件
  
  try {
    const resourceId = resource.resource_id || resource.id;
    let res;
    
    if (resource.is_favorite) {
      res = await unfavoriteResource(resourceId);
      if (res.code === 200) {
        ElMessage.success('已取消收藏');
        resource.is_favorite = false;
      }
    } else {
      res = await favoriteResource(resourceId);
      if (res.code === 200) {
        ElMessage.success('收藏成功');
        resource.is_favorite = true;
      }
    }
  } catch (error) {
    console.error('操作收藏失败:', error);
    ElMessage.error('操作失败，请稍后重试');
  }
}

// 类型切换
const handleTypeChange = (type) => {
  activeType.value = type;
  fetchMarketResources();
  fetchTags();
}

onMounted(() => {
  fetchMarketResources();
  fetchTags();
})

// 检查是否为自己发布的资源
const isSelfPublished = (resource) => {
  // 这里实际应该根据当前登录用户信息判断
  return false; // 默认为false，表示不是自己发布的
}

// 处理资源购买
const handlePurchase = async (resource) => {
  try {
    purchasing.value = true
    const res = await purchaseResource(resource.resource_id || resource.id)
    if (res.code === 200) {
      ElMessage.success('购买成功')
      // 更新资源状态
      resource.is_purchased = true
      // 关闭详情对话框
      resourceDetailDialog.value = false
      // 刷新资源列表
      fetchMarketResources()
    } else {
      ElMessage.error(res.message || '购买失败')
    }
  } catch (error) {
    console.error('购买失败:', error)
    ElMessage.error(error.message || '购买失败，请稍后重试')
  } finally {
    purchasing.value = false
  }
}
</script>

<template>
  <div class="market-container">
    <div class="market-header">
      <h1>交易市场</h1>
      <p>探索和购买高质量的数据、算力和算法资源</p>
    </div>
    
    <!-- 搜索区域 -->
    <div class="search-section">
      <div class="search-bar">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索资源名称、描述或标签"
          clearable
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
          <template #append>
            <el-button @click="handleSearch">搜索</el-button>
          </template>
        </el-input>
        <el-button type="primary" @click="advancedSearchVisible = true">
          高级筛选
          <el-icon><Filter /></el-icon>
        </el-button>
      </div>
      <div class="market-tabs">
        <div 
          v-for="type in resourceTypes" 
          :key="type.value" 
          class="market-tab-item"
          :class="{ active: activeType === type.value }"
          @click="handleTypeChange(type.value)"
        >
          <el-icon v-if="type.icon"><component :is="type.icon" /></el-icon>
          <span>{{ type.label }}</span>
        </div>
      </div>
    </div>
    
    <!-- 资源列表 -->
    <div v-loading="loading" class="resources-section">
      <div class="resources-header">
        <div class="resources-count">
          找到 <span class="highlight">{{ filteredResources.length }}</span> 个资源
        </div>
        <div class="resources-sort">
          <span>排序：</span>
          <el-select v-model="searchForm.sortBy" size="small">
            <el-option
              v-for="item in sortOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </div>
      </div>
      
      <el-empty v-if="filteredResources.length === 0" description="没有找到匹配的资源" />
      
      <div v-else class="resource-grid">
        <el-card 
          v-for="resource in filteredResources" 
          :key="resource.resource_id || resource.id" 
          class="resource-card"
          shadow="hover"
          @click="openResourceDetail(resource)"
        >
          <div class="resource-type-tag">
            <el-tag :type="typeColorMap[resource.resource_type]">{{ typeNameMap[resource.resource_type] }}</el-tag>
          </div>
          <h3 class="resource-title">{{ resource.name }}</h3>
          <div class="resource-tags" v-if="resource.tags && resource.tags.length > 0">
            <el-tag 
              v-for="tag in resource.tags.slice(0, 3)" 
              :key="tag" 
              size="small"
              effect="plain"
              class="resource-tag"
            >
              {{ tag }}
            </el-tag>
            <el-tag 
              v-if="resource.tags && resource.tags.length > 3" 
              size="small"
              effect="plain"
              class="resource-tag"
            >
              +{{ resource.tags.length - 3 }}
            </el-tag>
          </div>
          <p class="resource-description">{{ resource.description }}</p>
          <div class="resource-meta">
            <div class="resource-provider">
              <el-icon><User /></el-icon>
              <span>{{ resource.publisher }}</span>
            </div>
            <div class="resource-sales">
              <el-icon><Shop /></el-icon>
              <span>{{ resource.sales_count || 0 }}次</span>
            </div>
          </div>
          <div class="resource-rating">
            <el-rate v-model="resource.rating" disabled text-color="#ff9900" />
            <span class="rating-value">{{ resource.rating }}</span>
          </div>
          <div class="resource-footer">
            <div class="resource-price">¥{{ resource.price }}</div>
            <el-button type="primary" size="small">详情</el-button>
          </div>
        </el-card>
      </div>
    </div>
    
    <!-- 高级搜索抽屉 -->
    <el-drawer
      v-model="advancedSearchVisible"
      title="高级筛选"
      direction="rtl"
      size="400px"
    >
      <div class="advanced-search-form">
        <h3>价格范围</h3>
        <el-slider
          v-model="searchForm.priceRange"
          range
          :min="0"
          :max="50000"
          :step="100"
        />
        <div class="price-range-display">
          <span>¥{{ searchForm.priceRange[0] }}</span>
          <span>¥{{ searchForm.priceRange[1] }}</span>
        </div>
        
        <h3>资源标签</h3>
        <el-select
          v-model="searchForm.tags"
          multiple
          collapse-tags
          collapse-tags-tooltip
          placeholder="请选择标签"
          style="width: 100%"
        >
          <el-option
            v-for="item in tagOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
        
        <h3>最低评分</h3>
        <el-rate v-model="searchForm.ratingAbove" />
        
        <div class="drawer-footer">
          <el-button @click="resetSearch">重置</el-button>
          <el-button type="primary" @click="handleSearch">应用筛选</el-button>
        </div>
      </div>
    </el-drawer>
    
    <!-- 资源详情对话框 -->
    <el-dialog
      v-model="resourceDetailDialog"
      title="资源详情"
      width="700px"
    >
      <div v-if="detailResource" class="resource-detail">
        <div class="resource-detail-header">
          <h2>{{ detailResource.name }}</h2>
          <el-tag :type="typeColorMap[detailResource.resource_type]">{{ typeNameMap[detailResource.resource_type] }}</el-tag>
        </div>
        
        <div class="resource-meta">
          <div class="meta-item">
            <span class="meta-label">发布者：</span>
            <span>{{ detailResource.publisher }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">发布时间：</span>
            <span>{{ detailResource.publishTime }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">销量：</span>
            <span>{{ detailResource.sales_count || 0 }} 次</span>
          </div>
          <div class="meta-item price-item">
            <span class="meta-label">价格：</span>
            <span class="detail-price">¥{{ detailResource.price }}</span>
          </div>
        </div>
        
        <div class="detail-rating">
          <span class="meta-label">评分：</span>
          <el-rate v-model="detailResource.rating" disabled text-color="#ff9900" />
          <span class="rating-value">({{ detailResource.rating }})</span>
        </div>
        
        <div class="resource-tags">
          <span class="meta-label">标签：</span>
          <el-tag 
            v-for="tag in detailResource.tags" 
            :key="tag" 
            size="small"
            effect="plain"
            class="resource-tag"
          >
            {{ tag }}
          </el-tag>
        </div>
        
        <div class="resource-description">
          <div class="meta-label">资源描述：</div>
          <p>{{ detailResource.description }}</p>
        </div>
        
        <div class="resource-detail-actions">
          <!-- 数据资源按钮 -->
          <template v-if="detailResource.resource_type === 'data'">
            <el-button 
              v-if="!detailResource.is_purchased" 
              type="primary" 
              :loading="purchasing"
              @click="handlePurchase(detailResource)"
            >
              立即购买
            </el-button>
            <el-button 
              v-else 
              type="success" 
              @click="$router.push('/sandbox')"
            >
              使用资源
            </el-button>
            <el-button type="primary" @click="viewSampleData">查看示例数据</el-button>
            <el-button type="success" @click="viewDocument">查看说明文档</el-button>
            <el-button type="info" @click="toggleFavorite(detailResource, $event)">
              {{ detailResource.is_favorite ? '取消收藏' : '加入收藏' }}
            </el-button>
          </template>
          
          <!-- 算力资源按钮 -->
          <template v-else-if="detailResource.resource_type === 'computing'">
            <el-button 
              v-if="!detailResource.is_purchased" 
              type="primary" 
              :loading="purchasing"
              @click="handlePurchase(detailResource)"
            >
              按量计费
            </el-button>
            <el-button 
              v-else 
              type="success" 
              @click="$router.push('/sandbox')"
            >
              使用资源
            </el-button>
            <el-button type="success" @click="viewDocument">查看说明文档</el-button>
            <el-button type="warning" @click="toggleFavorite(detailResource, $event)">
              {{ detailResource.is_favorite ? '取消收藏' : '加入收藏' }}
            </el-button>
          </template>
          
          <!-- 算法资源按钮 -->
          <template v-else-if="detailResource.resource_type === 'algorithm'">
            <el-button 
              v-if="!detailResource.is_purchased" 
              type="primary" 
              :loading="purchasing"
              @click="handlePurchase(detailResource)"
            >
              立即购买
            </el-button>
            <el-button 
              v-else 
              type="success" 
              @click="$router.push('/sandbox')"
            >
              使用资源
            </el-button>
            <el-button type="success" @click="viewDocument">查看说明文档</el-button>
            <el-button type="info" @click="toggleFavorite(detailResource, $event)">
              {{ detailResource.is_favorite ? '取消收藏' : '加入收藏' }}
            </el-button>
          </template>
        </div>
      </div>
    </el-dialog>
    
    <!-- 样例数据对话框 -->
    <el-dialog
      v-model="sampleDataDialog"
      title="样例数据预览"
      width="800px"
    >
      <pre v-if="sampleData" class="sample-data-preview">{{ sampleData }}</pre>
      <div v-else class="empty-data">
        <el-empty description="暂无样例数据" />
      </div>
    </el-dialog>
    
    <!-- 文档对话框 -->
    <el-dialog
      v-model="documentDialog"
      title="使用说明"
      width="800px"
    >
      <div v-if="renderedDocument.documentation" class="documentation-preview">
        <h3>文档内容</h3>
        <div v-html="renderedDocument.documentation"></div>
      </div>
      <div v-else-if="renderedDocument.usageGuide" class="usage-guide-preview">
        <h3>使用指南</h3>
        <div v-html="renderedDocument.usageGuide"></div>
      </div>
      <div v-else-if="renderedDocument.apiDoc" class="api-doc-preview">
        <h3>API文档</h3>
        <div v-html="renderedDocument.apiDoc"></div>
      </div>
      <div v-else class="empty-data">
        <el-empty description="暂无文档内容" />
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.market-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.market-header {
  text-align: center;
  margin-bottom: 30px;
}

.market-header h1 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 10px;
}

.market-header p {
  font-size: 16px;
  color: #606266;
}

/* 搜索区域样式 */
.search-section {
  margin-bottom: 20px;
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.search-input {
  flex-grow: 1;
}

.market-tabs {
  display: flex;
  border-bottom: 1px solid #e4e7ed;
  padding-bottom: 10px;
}

.market-tab-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.market-tab-item.active {
  color: #409EFF;
  background-color: #ecf5ff;
  font-weight: bold;
}

.market-tab-item .el-icon {
  margin-right: 6px;
}

/* 资源列表样式 */
.resources-section {
  margin-bottom: 40px;
}

.resources-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.resources-count {
  font-size: 14px;
  color: #606266;
}

.resources-count .highlight {
  color: #409EFF;
  font-weight: bold;
}

.resources-sort {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.resource-card {
  height: 100%;
  cursor: pointer;
  transition: transform 0.3s;
  position: relative;
}

.resource-card:hover {
  transform: translateY(-5px);
}

.resource-type-tag {
  position: absolute;
  top: 10px;
  right: 10px;
}

.resource-title {
  font-size: 16px;
  margin: 10px 0;
  height: 44px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.resource-tags {
  margin-bottom: 10px;
  min-height: 32px;
}

.resource-tag {
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
  margin-bottom: 10px;
  font-size: 12px;
  color: #909399;
}

.resource-provider,
.resource-sales {
  display: flex;
  align-items: center;
}

.resource-provider .el-icon,
.resource-sales .el-icon {
  margin-right: 5px;
}

.resource-rating {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.rating-value {
  margin-left: 8px;
  font-size: 14px;
  color: #ff9900;
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

/* 高级搜索抽屉样式 */
.advanced-search-form {
  padding: 0 20px;
}

.advanced-search-form h3 {
  margin-top: 20px;
  margin-bottom: 10px;
  font-size: 16px;
  color: #303133;
}

.price-range-display {
  display: flex;
  justify-content: space-between;
  margin-top: 5px;
  color: #606266;
}

.drawer-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20px;
  background-color: #fff;
  border-top: 1px solid #e4e7ed;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 资源详情对话框样式 */
.resource-detail-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.resource-detail-header h2 {
  margin-right: 10px;
  margin-bottom: 0;
}

.resource-meta {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 20px;
}

.meta-item {
  display: flex;
  align-items: center;
}

.meta-label {
  font-weight: bold;
  margin-right: 5px;
  color: #606266;
}

.detail-price {
  font-weight: bold;
  color: #f56c6c;
  font-size: 18px;
}

.detail-rating {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.resource-description {
  margin-top: 20px;
  margin-bottom: 30px;
}

.resource-detail-actions {
  margin-top: 30px;
  display: flex;
  gap: 10px;
}

/* 样例数据和协议预览样式 */
.sample-data-preview,
.usage-guide-preview,
.api-doc-preview {
  white-space: pre-wrap;
  word-break: break-all;
  background-color: #f8f8f8;
  padding: 15px;
  border-radius: 5px;
  font-family: monospace;
  max-height: 500px;
  overflow-y: auto;
}

.usage-guide-preview {
  white-space: pre-line;
  font-family: system-ui, -apple-system, sans-serif;
}

.empty-data {
  padding: 20px;
  text-align: center;
}

.documentation-preview {
  padding: 20px;
}

.documentation-preview h3 {
  margin-bottom: 10px;
}
</style> 