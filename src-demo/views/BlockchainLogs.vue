<script setup>
import { ref, reactive, onMounted, computed } from 'vue'

// 加载状态
const loading = ref(false)

// 区块链日志数据
const blockchainLogs = ref([])

// 当前页
const currentPage = ref(1)
// 每页显示数
const pageSize = ref(10)
// 总记录数
const total = ref(0)

// 选中的日志类型筛选
const selectedLogTypes = ref([])
// 选中的操作类型筛选
const selectedOperationTypes = ref([])
// 时间范围筛选
const dateRange = ref([])

// 高级搜索表单
const advancedSearchForm = reactive({
  userId: '',
  resourceId: '',
  transactionHash: '',
  blockHeight: ''
})

// 显示高级搜索
const showAdvancedSearch = ref(false)

// 获取区块链日志数据
const fetchBlockchainLogs = () => {
  loading.value = true
  
  // 模拟API调用
  setTimeout(() => {
    // 模拟数据
    blockchainLogs.value = [
      {
        id: 1,
        timestamp: '2023-05-30 14:35:22',
        blockHeight: 1056824,
        transactionHash: '0x7d2a4f3d5c9b8a2e6f1d0c7b4a5e3f2d1c0b9a8e7d6c5b4a3f2e1d0c9b8a7f6e5',
        logType: 'transaction',
        operationType: 'purchase',
        userId: 'user123',
        userName: '张三',
        resourceId: 'res456',
        resourceName: '医疗影像数据集',
        details: {
          price: 12800,
          seller: '医疗大数据公司',
          contractId: 'TC202302200001'
        }
      },
      {
        id: 2,
        timestamp: '2023-05-30 15:22:48',
        blockHeight: 1056830,
        transactionHash: '0x8e3b5f4e2d1c0b9a8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5e4',
        logType: 'resource',
        operationType: 'upload',
        userId: 'user456',
        userName: '李四',
        resourceId: 'res789',
        resourceName: '金融数据分析模型',
        details: {
          size: '235MB',
          format: 'ZIP',
          type: 'algorithm'
        }
      },
      {
        id: 3,
        timestamp: '2023-05-30 16:05:33',
        blockHeight: 1056835,
        transactionHash: '0x9f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5e4d3',
        logType: 'sandbox',
        operationType: 'use',
        userId: 'user789',
        userName: '王五',
        resourceId: 'res123',
        resourceName: '城市交通流量数据集',
        details: {
          usage: '数据分析',
          duration: '45分钟',
          result: '生成报告'
        }
      },
      {
        id: 4,
        timestamp: '2023-05-30 16:38:12',
        blockHeight: 1056840,
        transactionHash: '0xa0f1e2d3c4b5a6f7e8d9c0b1a2f3e4d5c6b7a8f9e0d1c2b3a4f5e6d7c8b9a0f1',
        logType: 'transaction',
        operationType: 'sell',
        userId: 'user123',
        userName: '张三',
        resourceId: 'res321',
        resourceName: '电商用户行为分析模型',
        details: {
          price: 8200,
          buyer: '电商科技有限公司',
          contractId: 'TC202304250005'
        }
      },
      {
        id: 5,
        timestamp: '2023-05-30 17:12:55',
        blockHeight: 1056845,
        transactionHash: '0xb1f2e3d4c5b6a7f8e9d0c1b2a3f4e5d6c7b8a9f0e1d2c3b4a5f6e7d8c9b0a1f2',
        logType: 'sandbox',
        operationType: 'compute',
        userId: 'user456',
        userName: '李四',
        resourceId: 'res654',
        resourceName: '高性能GPU计算服务',
        details: {
          computeType: 'AI训练',
          duration: '2小时15分钟',
          cost: 450
        }
      },
      {
        id: 6,
        timestamp: '2023-05-30 18:02:29',
        blockHeight: 1056850,
        transactionHash: '0xc2f3e4d5c6b7a8f9e0d1c2b3a4f5e6d7c8b9a0f1e2d3c4b5a6f7e8d9c0b1a2f3',
        logType: 'resource',
        operationType: 'publish',
        userId: 'user789',
        userName: '王五',
        resourceId: 'res987',
        resourceName: '智能医疗诊断算法',
        details: {
          price: 15600,
          type: 'algorithm',
          category: '医疗'
        }
      }
    ]
    
    total.value = 120 // 模拟总记录数
    loading.value = false
  }, 1000)
}

// 日志类型选项
const logTypeOptions = [
  { label: '交易记录', value: 'transaction' },
  { label: '资源操作', value: 'resource' },
  { label: '沙箱使用', value: 'sandbox' }
]

// 操作类型选项
const operationTypeOptions = [
  { label: '购买', value: 'purchase' },
  { label: '销售', value: 'sell' },
  { label: '上传', value: 'upload' },
  { label: '发布', value: 'publish' },
  { label: '使用', value: 'use' },
  { label: '计算', value: 'compute' }
]

// 切换高级搜索显示
const toggleAdvancedSearch = () => {
  showAdvancedSearch.value = !showAdvancedSearch.value
}

// 重置筛选条件
const resetFilters = () => {
  selectedLogTypes.value = []
  selectedOperationTypes.value = []
  dateRange.value = []
  advancedSearchForm.userId = ''
  advancedSearchForm.resourceId = ''
  advancedSearchForm.transactionHash = ''
  advancedSearchForm.blockHeight = ''
  
  fetchBlockchainLogs()
}

// 应用筛选
const applyFilters = () => {
  currentPage.value = 1
  fetchBlockchainLogs()
}

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  fetchBlockchainLogs()
}

// 格式化日志详情
const formatDetails = (details, logType, operationType) => {
  if (!details) return ''
  
  let result = ''
  if (logType === 'transaction') {
    if (operationType === 'purchase') {
      result = `价格: ¥${details.price} | 卖家: ${details.seller} | 合约ID: ${details.contractId}`
    } else if (operationType === 'sell') {
      result = `价格: ¥${details.price} | 买家: ${details.buyer} | 合约ID: ${details.contractId}`
    }
  } else if (logType === 'resource') {
    if (operationType === 'upload') {
      result = `大小: ${details.size} | 格式: ${details.format} | 类型: ${details.type}`
    } else if (operationType === 'publish') {
      result = `价格: ¥${details.price} | 类型: ${details.type} | 分类: ${details.category}`
    }
  } else if (logType === 'sandbox') {
    if (operationType === 'use') {
      result = `用途: ${details.usage} | 时长: ${details.duration} | 结果: ${details.result}`
    } else if (operationType === 'compute') {
      result = `计算类型: ${details.computeType} | 时长: ${details.duration} | 费用: ¥${details.cost}`
    }
  }
  
  return result
}

// 查看交易详情
const viewTransactionDetails = (transactionHash) => {
  // 在实际应用中，这里应该跳转到区块链浏览器或显示详情对话框
  ElMessage.info(`查看交易 ${transactionHash} 的详情`)
}

// 操作类型标签样式
const getOperationTypeTagType = (type) => {
  const map = {
    'purchase': 'danger',
    'sell': 'success',
    'upload': 'info',
    'publish': 'warning',
    'use': 'primary',
    'compute': ''
  }
  return map[type] || 'info'
}

// 日志类型标签样式
const getLogTypeTagType = (type) => {
  const map = {
    'transaction': 'danger',
    'resource': 'success',
    'sandbox': 'warning'
  }
  return map[type] || 'info'
}

// 日志类型中文名称
const getLogTypeName = (type) => {
  const map = {
    'transaction': '交易记录',
    'resource': '资源操作',
    'sandbox': '沙箱使用'
  }
  return map[type] || type
}

// 操作类型中文名称
const getOperationTypeName = (type) => {
  const map = {
    'purchase': '购买',
    'sell': '销售',
    'upload': '上传',
    'publish': '发布',
    'use': '使用',
    'compute': '计算'
  }
  return map[type] || type
}

onMounted(() => {
  fetchBlockchainLogs()
})

// 显示的日志数据，应用过滤
const filteredLogs = computed(() => {
  return blockchainLogs.value
})

// 区块链统计数据
const blockchainStats = ref({
  totalBlocks: 1056850,
  totalTransactions: 5692,
  dailyTransactions: 128,
  averageBlockTime: '15秒'
})

// 总计数据
const totalData = ref({
  transactions: 2845,
  resources: 1520,
  users: 687,
  daily: 98
})
</script>

<template>
  <div class="logs-container">
    <div class="logs-header">
      <h1>区块链日志查询</h1>
      <p>查询平台所有操作的区块链记录，实现全程可追溯和监管</p>
    </div>
    
    <!-- 区块链统计信息 -->
    <el-card class="stats-card">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-icon">
            <el-icon><Link /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ blockchainStats.totalBlocks }}</div>
            <div class="stat-label">区块总数</div>
          </div>
        </div>
        
        <div class="stat-item">
          <div class="stat-icon">
            <el-icon><Histogram /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ blockchainStats.totalTransactions }}</div>
            <div class="stat-label">交易总数</div>
          </div>
        </div>
        
        <div class="stat-item">
          <div class="stat-icon">
            <el-icon><Calendar /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ blockchainStats.dailyTransactions }}</div>
            <div class="stat-label">日交易量</div>
          </div>
        </div>
        
        <div class="stat-item">
          <div class="stat-icon">
            <el-icon><Timer /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ blockchainStats.averageBlockTime }}</div>
            <div class="stat-label">平均出块时间</div>
          </div>
        </div>
      </div>
    </el-card>
    
    <!-- 区块链可视化 -->
    <el-card class="visualization-card">
      <template #header>
        <div class="card-header">
          <div>区块链可视化</div>
        </div>
      </template>
      
      <div class="blockchain-visualization">
        <!-- 实际应用中，这里应该是使用图表库的可视化组件 -->
        <div class="visualization-placeholder">
          <div class="chain-blocks">
            <div 
              v-for="i in 10" 
              :key="i" 
              class="chain-block"
              :class="{ 'latest-block': i === 1 }"
            >
              <div class="block-header">
                <div class="block-number">#{{ blockchainStats.totalBlocks - (i - 1) }}</div>
                <div class="block-time">{{ new Date().toLocaleTimeString() }}</div>
              </div>
              <div class="block-transactions">
                <div 
                  v-for="j in (Math.floor(Math.random() * 5) + 1)" 
                  :key="j"
                  class="block-tx"
                >
                  TX
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
    
    <!-- 日志查询过滤器 -->
    <el-card class="filter-card">
      <div class="filter-header">
        <div class="filter-title">日志查询</div>
        <el-button type="primary" plain size="small" @click="toggleAdvancedSearch">
          {{ showAdvancedSearch ? '收起高级搜索' : '显示高级搜索' }}
        </el-button>
      </div>
      
      <div class="filter-content">
        <el-row :gutter="20">
          <el-col :lg="6" :md="12">
            <div class="filter-item">
              <div class="filter-label">日志类型</div>
              <el-select 
                v-model="selectedLogTypes" 
                multiple 
                collapse-tags 
                placeholder="选择日志类型"
                style="width: 100%"
              >
                <el-option 
                  v-for="item in logTypeOptions" 
                  :key="item.value" 
                  :label="item.label" 
                  :value="item.value" 
                />
              </el-select>
            </div>
          </el-col>
          
          <el-col :lg="6" :md="12">
            <div class="filter-item">
              <div class="filter-label">操作类型</div>
              <el-select 
                v-model="selectedOperationTypes" 
                multiple 
                collapse-tags 
                placeholder="选择操作类型"
                style="width: 100%"
              >
                <el-option 
                  v-for="item in operationTypeOptions" 
                  :key="item.value" 
                  :label="item.label" 
                  :value="item.value" 
                />
              </el-select>
            </div>
          </el-col>
          
          <el-col :lg="6" :md="12">
            <div class="filter-item">
              <div class="filter-label">时间范围</div>
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="width: 100%"
              />
            </div>
          </el-col>
          
          <el-col :lg="6" :md="12">
            <div class="filter-item">
              <div class="filter-label">&nbsp;</div>
              <div class="filter-buttons">
                <el-button type="primary" @click="applyFilters">查询</el-button>
                <el-button @click="resetFilters">重置</el-button>
              </div>
            </div>
          </el-col>
        </el-row>
        
        <!-- 高级搜索 -->
        <el-collapse-transition>
          <div v-if="showAdvancedSearch" class="advanced-search">
            <el-divider>高级搜索</el-divider>
            <el-row :gutter="20">
              <el-col :lg="6" :md="12">
                <div class="filter-item">
                  <div class="filter-label">用户ID</div>
                  <el-input 
                    v-model="advancedSearchForm.userId" 
                    placeholder="输入用户ID"
                  />
                </div>
              </el-col>
              
              <el-col :lg="6" :md="12">
                <div class="filter-item">
                  <div class="filter-label">资源ID</div>
                  <el-input 
                    v-model="advancedSearchForm.resourceId" 
                    placeholder="输入资源ID"
                  />
                </div>
              </el-col>
              
              <el-col :lg="6" :md="12">
                <div class="filter-item">
                  <div class="filter-label">交易哈希</div>
                  <el-input 
                    v-model="advancedSearchForm.transactionHash" 
                    placeholder="输入交易哈希"
                  />
                </div>
              </el-col>
              
              <el-col :lg="6" :md="12">
                <div class="filter-item">
                  <div class="filter-label">区块高度</div>
                  <el-input 
                    v-model="advancedSearchForm.blockHeight" 
                    placeholder="输入区块高度"
                  />
                </div>
              </el-col>
            </el-row>
          </div>
        </el-collapse-transition>
      </div>
    </el-card>
    
    <!-- 日志表格 -->
    <el-card class="logs-table-card">
      <div v-loading="loading">
        <el-table
          :data="filteredLogs"
          style="width: 100%"
          border
        >
          <el-table-column prop="id" label="序号" width="70" />
          
          <el-table-column prop="timestamp" label="时间" width="160" />
          
          <el-table-column prop="blockHeight" label="区块高度" width="120" align="center" />
          
          <el-table-column label="交易哈希" width="120">
            <template #default="{ row }">
              <el-tooltip 
                :content="row.transactionHash" 
                placement="top" 
                :show-after="500"
              >
                <el-button 
                  link 
                  type="primary" 
                  @click="viewTransactionDetails(row.transactionHash)"
                >
                  查看交易
                </el-button>
              </el-tooltip>
            </template>
          </el-table-column>
          
          <el-table-column label="日志类型" width="120">
            <template #default="{ row }">
              <el-tag 
                :type="getLogTypeTagType(row.logType)"
                effect="dark"
              >
                {{ getLogTypeName(row.logType) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="操作类型" width="100">
            <template #default="{ row }">
              <el-tag 
                :type="getOperationTypeTagType(row.operationType)"
                effect="plain"
              >
                {{ getOperationTypeName(row.operationType) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="用户" width="150">
            <template #default="{ row }">
              <div>{{ row.userName }}</div>
              <div class="small-text">ID: {{ row.userId }}</div>
            </template>
          </el-table-column>
          
          <el-table-column label="资源" width="180">
            <template #default="{ row }">
              <div>{{ row.resourceName }}</div>
              <div class="small-text">ID: {{ row.resourceId }}</div>
            </template>
          </el-table-column>
          
          <el-table-column label="详情" min-width="250">
            <template #default="{ row }">
              {{ formatDetails(row.details, row.logType, row.operationType) }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchBlockchainLogs"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.logs-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.logs-header {
  text-align: center;
  margin-bottom: 30px;
}

.logs-header h1 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 10px;
}

.logs-header p {
  font-size: 16px;
  color: #606266;
}

/* 统计卡片样式 */
.stats-card {
  margin-bottom: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  background-color: #ecf5ff;
  border-radius: 8px;
  margin-right: 15px;
}

.stat-icon .el-icon {
  font-size: 24px;
  color: #409EFF;
}

.stat-content {
  flex: 1;
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

/* 可视化卡片样式 */
.visualization-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.blockchain-visualization {
  padding: 20px 0;
}

.visualization-placeholder {
  overflow-x: auto;
}

.chain-blocks {
  display: flex;
  padding: 10px;
  min-width: 100%;
}

.chain-block {
  width: 120px;
  height: 120px;
  background-color: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  margin-right: 10px;
  padding: 10px;
  position: relative;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}

.chain-block::after {
  content: '';
  position: absolute;
  top: 50%;
  right: -10px;
  width: 10px;
  height: 2px;
  background-color: #e4e7ed;
  z-index: 1;
}

.chain-block:last-child::after {
  display: none;
}

.chain-block.latest-block {
  background-color: #ecf5ff;
  border-color: #409EFF;
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.block-number {
  font-weight: bold;
  color: #303133;
}

.block-time {
  font-size: 12px;
  color: #909399;
}

.block-transactions {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  align-content: flex-start;
}

.block-tx {
  font-size: 10px;
  padding: 2px 5px;
  background-color: #f0f2f5;
  border-radius: 3px;
  color: #606266;
}

.latest-block .block-tx {
  background-color: #d9ecff;
  color: #409EFF;
}

/* 过滤卡片样式 */
.filter-card {
  margin-bottom: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-title {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.filter-content {
  padding: 10px 0;
}

.filter-item {
  margin-bottom: 20px;
}

.filter-label {
  font-size: 14px;
  margin-bottom: 8px;
  color: #606266;
}

.filter-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.advanced-search {
  margin-top: 10px;
}

/* 日志表格样式 */
.logs-table-card {
  margin-bottom: 30px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.small-text {
  font-size: 12px;
  color: #909399;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style> 