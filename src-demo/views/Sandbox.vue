<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { 
  Document, 
  Setting, 
  Link, 
  Plus, 
  DataLine, 
  Monitor, 
  RefreshRight, 
  Download,
  Notebook,
  Connection
} from '@element-plus/icons-vue'
// 移除原有的 API 导入
// import { getPurchasedResources, getUserResources } from '../api/market'
// import { previewResource as fetchResourcePreview } from '../api/resources'
import { marked } from 'marked'
import Papa from 'papaparse'  // 添加 CSV 解析库
import axios from 'axios'  // 直接导入 axios

// 激活的面板
const activePanel = ref('ui')

// 用户资源列表
const dataResources = ref([])
const algorithmResources = ref([])
const computingResources = ref([])

// 可用算法插件
const availablePlugins = ref([])

// 当前选择的资源
const selectedResource = ref(null)

// 当前选择的算力资源
const selectedComputing = ref(null)

// 添加目录导航相关的状态
const currentPath = ref('/')
const pathStack = ref(['/'])

// 修改命令历史记录的初始内容
const commandHistory = ref([
  { type: 'system', content: '欢迎使用智信链枢沙箱环境终端' },
  { type: 'system', content: '版本: v1.0.0' },
  { type: 'system', content: '输入 help 获取可用命令列表' }
])

// 当前命令
const currentCommand = ref('')

// 数据可视化结果 (模拟)
const visualizationResult = ref({
  chart: 'bar',
  title: '样本数据分析结果',
  data: {
    labels: ['类别A', '类别B', '类别C', '类别D', '类别E'],
    datasets: [
      {
        label: '样本数',
        data: [120, 65, 80, 45, 100],
        backgroundColor: [
          'rgba(255, 99, 132, 0.5)',
          'rgba(54, 162, 235, 0.5)',
          'rgba(255, 206, 86, 0.5)',
          'rgba(75, 192, 192, 0.5)',
          'rgba(153, 102, 255, 0.5)'
        ]
      }
    ]
  },
  options: {}
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

// 添加预览数据状态
const previewData = ref({
  headers: [],
  rows: [],
  loading: false,
  error: null
})

// 添加 Jupyter URL 配置
const jupyterUrl = ref('http://172.191.73.102:80')

// 打开Jupyter窗口的方法
const openJupyterWindow = () => {
  const windowFeatures = 'width=1200,height=800,left=100,top=100,resizable=yes,scrollbars=yes,toolbar=yes,menubar=no,location=yes,directories=no,status=yes'
  window.open(jupyterUrl.value, 'JupyterWindow', windowFeatures)
}

// 监听activePanel的变化
watch(activePanel, (newValue) => {
  if (newValue === 'jupyter') {
    // 当切换到Jupyter面板时，自动打开新窗口
    openJupyterWindow()
  }
})

// 获取用户的数据资源
const fetchDataResources = async () => {
  try {
    const loading = ElLoading.service({
      target: '.resource-section',
      text: '加载数据资源...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    // 直接使用 axios 进行请求，不依赖封装的 API
    const apiBaseUrl = 'http://localhost:8080'
    
    // 获取用户已购买的资源
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    
    const purchasedRes = await axios.get(`${apiBaseUrl}/api/v1/market/purchased`, { headers })
    const userRes = await axios.get(`${apiBaseUrl}/api/v1/user/resources`, { headers })
    
    if (purchasedRes.data.code === 200 && userRes.data.code === 200) {
      // 合并并过滤出数据资源
      const allResources = [...(purchasedRes.data.data || []), ...(userRes.data.data || [])]
      dataResources.value = allResources.filter(r => r.resource_type === 'data')
    }
    
    loading.close()
  } catch (error) {
    console.error('获取数据资源失败:', error)
    ElMessage.error('获取数据资源失败，请稍后重试')
    loading?.close()
  }
}

// 获取用户的算法资源
const fetchAlgorithmResources = async () => {
  try {
    const loading = ElLoading.service({
      target: '.resource-section',
      text: '加载算法资源...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    // 直接使用 axios 进行请求，不依赖封装的 API
    const apiBaseUrl = 'http://localhost:8080'
    
    // 获取用户已购买的资源
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    
    const purchasedRes = await axios.get(`${apiBaseUrl}/api/v1/market/purchased`, { headers })
    const userRes = await axios.get(`${apiBaseUrl}/api/v1/user/resources`, { headers })
    
    if (purchasedRes.data.code === 200 && userRes.data.code === 200) {
      // 合并并过滤出算法资源
      const allResources = [...(purchasedRes.data.data || []), ...(userRes.data.data || [])]
      algorithmResources.value = allResources.filter(r => r.resource_type === 'algorithm')
    }
    
    loading.close()
  } catch (error) {
    console.error('获取算法资源失败:', error)
    ElMessage.error('获取算法资源失败，请稍后重试')
    loading?.close()
  }
}

// 获取用户的算力资源
const fetchComputingResources = async () => {
  try {
    const loading = ElLoading.service({
      target: '.resource-section',
      text: '加载算力资源...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    // 直接使用 axios 进行请求，不依赖封装的 API
    const apiBaseUrl = 'http://localhost:8080'
    
    // 获取用户已购买的资源
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    
    const purchasedRes = await axios.get(`${apiBaseUrl}/api/v1/market/purchased`, { headers })
    const userRes = await axios.get(`${apiBaseUrl}/api/v1/user/resources`, { headers })
    
    if (purchasedRes.data.code === 200 && userRes.data.code === 200) {
      // 合并并过滤出算力资源
      const allResources = [...(purchasedRes.data.data || []), ...(userRes.data.data || [])]
      computingResources.value = allResources.filter(r => r.resource_type === 'computing')
    }
    
    loading.close()
  } catch (error) {
    console.error('获取算力资源失败:', error)
    ElMessage.error('获取算力资源失败，请稍后重试')
    loading?.close()
  }
}

// 安装插件
const installPlugin = (plugin) => {
  ElMessage.success(`插件 ${plugin.name} 安装成功！`)
  // 在实际应用中，这里应该调用API完成插件安装
}

// 连接第三方算力
const connectComputing = (computing) => {
  if (computing.availability === 'maintenance') {
    ElMessage.warning(`该算力资源正在维护中，预计维护结束时间: ${computing.maintenanceEndTime}`)
    return
  }
  
  selectedComputing.value = computing
  ElMessage.success(`成功连接到算力资源: ${computing.name}`)
  
  // 添加命令行记录
  commandHistory.value.push({
    type: 'system',
    content: `已连接到远程算力服务：${computing.name}`
  })
  commandHistory.value.push({
    type: 'system',
    content: `CPU: ${computing.performance.cpu} | GPU: ${computing.performance.gpu || 'N/A'} | 内存: ${computing.performance.memory}`
  })
}

// 获取当前目录下的内容
const getCurrentDirContent = () => {
  if (currentPath.value === '/') {
    return ['data/', 'plugins/']
  } else if (currentPath.value === '/data') {
    return dataResources.value.map(r => r.name)
  } else if (currentPath.value === '/plugins') {
    return algorithmResources.value.map(r => r.name)
  }
  return []
}

// 执行命令
const executeCommand = () => {
  if (!currentCommand.value.trim()) return
  
  // 添加用户命令到历史记录
  commandHistory.value.push({
    type: 'user',
    content: currentCommand.value
  })
  
  // 处理命令
  const command = currentCommand.value.trim().toLowerCase()
  const args = command.split(' ')
  
  if (command === 'help') {
    commandHistory.value.push({
      type: 'system',
      content: `
可用命令:
- help: 显示帮助信息
- clear: 清除终端
- ls: 列出当前目录文件
- cd <dir>: 切换目录
- python <file.py>: 运行Python脚本（自动跳转到Jupyter）
- jupyter: 启动Jupyter Notebook
- cat <file>: 查看文件内容（自动跳转到Jupyter）
- resources: 显示可用资源
- plugins: 显示已安装插件
      `
    })
  } else if (command === 'clear') {
    commandHistory.value = []
  } else if (command === 'ls') {
    const content = getCurrentDirContent()
    commandHistory.value.push({
      type: 'system',
      content: content.join('\n')
    })
  } else if (args[0] === 'cd') {
    const target = args[1] || ''
    
    if (target === '..') {
      if (currentPath.value === '/') {
        commandHistory.value.push({
          type: 'system',
          content: '权限不足：无法访问上级目录'
        })
      } else {
        pathStack.value.pop()
        currentPath.value = pathStack.value[pathStack.value.length - 1]
        commandHistory.value.push({
          type: 'system',
          content: `当前目录: ${currentPath.value}`
        })
      }
    } else if (target === 'data' && currentPath.value === '/') {
      currentPath.value = '/data'
      pathStack.value.push('/data')
      commandHistory.value.push({
        type: 'system',
        content: '已进入数据目录'
      })
    } else if (target === 'plugins' && currentPath.value === '/') {
      currentPath.value = '/plugins'
      pathStack.value.push('/plugins')
      commandHistory.value.push({
        type: 'system',
        content: '已进入插件目录'
      })
    } else {
      commandHistory.value.push({
        type: 'system',
        content: '无效的目录'
      })
    }
  } else if (command.startsWith('python ') || command === 'jupyter' || command.startsWith('cat ')) {
    activePanel.value = 'jupyter'
    commandHistory.value.push({
      type: 'system',
      content: '已切换到Jupyter界面'
    })
    // 自动打开Jupyter窗口
    openJupyterWindow()
  } else if (command === 'resources') {
    const resourceList = dataResources.value.map(r => `- ${r.name}`).join('\n')
    commandHistory.value.push({
      type: 'system',
      content: `可用资源:\n${resourceList}`
    })
  } else if (command === 'plugins') {
    const pluginList = algorithmResources.value.map(r => `- ${r.name} v${r.version}`).join('\n')
    commandHistory.value.push({
      type: 'system',
      content: `已安装插件:\n${pluginList}`
    })
  } else {
    commandHistory.value.push({
      type: 'system',
      content: `未知命令: ${command}. 使用 help 查看可用命令。`
    })
  }
  
  // 清空当前命令
  currentCommand.value = ''
  
  // 滚动到底部
  setTimeout(() => {
    const terminalElement = document.querySelector('.terminal-content')
    if (terminalElement) {
      terminalElement.scrollTop = terminalElement.scrollHeight
    }
  }, 100)
}

// 用图表类型
const chartTypes = [
  { label: '柱状图', value: 'bar' },
  { label: '折线图', value: 'line' },
  { label: '饼图', value: 'pie' },
  { label: '散点图', value: 'scatter' }
]

// 当前选择的图表类型
const selectedChartType = ref('bar')

// 改变图表类型
const changeChartType = (type) => {
  selectedChartType.value = type
  // 在实际应用中，这里应该更新可视化内容
}

// 修改打开资源的函数，添加自动切换到可视化界面的功能
const openResource = async (resource) => {
  selectedResource.value = resource
  ElMessage.success(`已打开资源: ${resource.name}`)
  
  // 如果当前在终端或Jupyter界面，自动切换到可视化界面
  if (activePanel.value === 'terminal' || activePanel.value === 'jupyter') {
    activePanel.value = 'ui'
  }
  
  // 如果是数据资源，尝试预览
  if (resource.resource_type === 'data') {
    try {
      previewData.value.loading = true
      previewData.value.error = null
      
      // 直接使用 axios 调用预览API
      const apiBaseUrl = 'http://localhost:8080'
      const token = localStorage.getItem('token')
      const headers = token ? { Authorization: `Bearer ${token}` } : {}
      
      // 调用预览API
      const response = await axios.get(`${apiBaseUrl}/api/v1/resources/preview/${resource.id}`, { 
        headers, 
        responseType: 'blob'
      })
      
      // 检查响应类型
      const contentType = response.headers?.['content-type']
      
      if (contentType && contentType.includes('text/csv')) {
        // 如果是CSV文件，使用FileReader读取内容
        const reader = new FileReader()
        reader.onload = (e) => {
          try {
            // 使用Papa Parse解析CSV数据
            const result = Papa.parse(e.target.result, {
              header: true,
              dynamicTyping: true
            })
            
            if (result.data && result.data.length > 0) {
              previewData.value.headers = Object.keys(result.data[0])
              previewData.value.rows = result.data
            } else {
              previewData.value.error = '无法解析CSV数据'
            }
          } catch (error) {
            console.error('解析CSV失败:', error)
            previewData.value.error = '解析CSV数据失败'
          }
        }
        reader.onerror = () => {
          console.error('读取CSV失败')
          previewData.value.error = '读取CSV数据失败'
        }
        reader.readAsText(response.data)
      } else if (contentType && contentType.includes('application/json')) {
        // 如果是JSON响应（可能是MD5值或错误信息）
        const reader = new FileReader()
        reader.onload = (e) => {
          try {
            const result = JSON.parse(e.target.result)
            if (result.code !== 200) {
              previewData.value.error = result.message || '无法预览该资源'
            }
          } catch (error) {
            console.error('解析JSON失败:', error)
            previewData.value.error = '解析预览数据失败'
          }
        }
        reader.readAsText(response.data)
      } else {
        previewData.value.error = '不支持的文件类型'
      }
    } catch (error) {
      console.error('预览资源失败:', error)
      previewData.value.error = '预览失败: ' + (error.message || '未知错误')
    } finally {
      previewData.value.loading = false
    }
  }
}

onMounted(() => {
  fetchDataResources()
  fetchAlgorithmResources()
  fetchComputingResources()
})
</script>

<template>
  <div class="sandbox-container">
    <div class="sandbox-header">
      <h1>沙箱环境</h1>
      <p>在安全隔离的环境中使用数据和算法资源</p>
    </div>
    
    <div class="sandbox-main">
      <!-- 左侧资源面板 -->
      <div class="resource-panel">
        <h3 class="panel-title">资源管理</h3>
        
        <!-- 我的资源 -->
        <div class="resource-section">
          <div class="section-header">
            <h4>可用资源</h4>
            <el-button size="small" type="primary" @click="$router.push('/upload')">
              <el-icon><Plus /></el-icon>添加资源
            </el-button>
          </div>
          
          <el-empty v-if="!dataResources.length" description="没有可用资源" />
          
          <div v-else class="resource-list">
            <div 
              v-for="resource in dataResources" 
              :key="resource.id" 
              class="resource-item"
              :class="{ active: selectedResource && selectedResource.id === resource.id }"
              @click="openResource(resource)"
            >
              <div class="resource-icon">
                <el-icon v-if="resource.type === 'data'"><Document /></el-icon>
                <el-icon v-else-if="resource.type === 'algorithm'"><Setting /></el-icon>
                <el-icon v-else><Link /></el-icon>
              </div>
              <div class="resource-info">
                <div class="resource-name">
                  {{ resource.name }}
                  <el-tag size="small" :type="typeColorMap[resource.type]">{{ typeNameMap[resource.type] }}</el-tag>
                </div>
                <div class="resource-description">{{ resource.description }}</div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 算法插件库 -->
        <div class="resource-section">
          <div class="section-header">
            <h4>算法插件库</h4>
          </div>
          
          <el-empty v-if="!algorithmResources.length" description="没有可用插件" />
          
          <div v-else class="resource-list">
            <div 
              v-for="plugin in algorithmResources" 
              :key="plugin.id" 
              class="resource-item"
            >
              <div class="resource-icon">
                <el-icon><Setting /></el-icon>
              </div>
              <div class="resource-info">
                <div class="resource-name">
                  {{ plugin.name }}
                  <span class="plugin-version">v{{ plugin.version }}</span>
                </div>
                <div class="resource-description">{{ plugin.description }}</div>
                <div class="plugin-tags">
                  <el-tag 
                    v-for="tag in plugin.tags" 
                    :key="tag" 
                    size="small"
                    effect="plain"
                    class="resource-tag"
                  >
                    {{ tag }}
                  </el-tag>
                </div>
              </div>
              <div class="resource-actions">
                <el-button size="small" type="primary" @click="installPlugin(plugin)">安装</el-button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 算力资源 -->
        <div class="resource-section">
          <div class="section-header">
            <h4>第三方算力资源</h4>
          </div>
          
          <el-empty v-if="!computingResources.length" description="没有可用算力资源" />
          
          <div v-else class="resource-list">
            <div 
              v-for="computing in computingResources" 
              :key="computing.id" 
              class="resource-item"
              :class="{ active: selectedComputing && selectedComputing.id === computing.id }"
            >
              <div class="resource-icon">
                <el-icon><Link /></el-icon>
              </div>
              <div class="resource-info">
                <div class="resource-name">
                  {{ computing.name }}
                  <el-tag 
                    size="small" 
                    type="success"
                  >
                    免费
                  </el-tag>
                </div>
                <div class="resource-description">{{ computing.description }}</div>
                <div class="computing-details">
                  <div>提供商: {{ computing.provider }}</div>
                  <div>价格: 免费</div>
                </div>
              </div>
              <div class="resource-actions">
                <el-button 
                  size="small" 
                  type="success" 
                  :icon="Connection"
                  disabled
                >
                  已连接
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 右侧工作区 -->
      <div class="workspace-panel">
        <div class="workspace-tabs">
          <div class="tab-header">
            <div 
              class="tab-item" 
              :class="{ active: activePanel === 'ui' }"
              @click="activePanel = 'ui'"
            >
              <el-icon><DataLine /></el-icon>
              <span>可视化界面</span>
            </div>
            <div 
              class="tab-item" 
              :class="{ active: activePanel === 'terminal' }"
              @click="activePanel = 'terminal'"
            >
              <el-icon><Monitor /></el-icon>
              <span>命令行终端</span>
            </div>
            <div 
              class="tab-item" 
              :class="{ active: activePanel === 'jupyter' }"
              @click="activePanel = 'jupyter'"
            >
              <el-icon><Notebook /></el-icon>
              <span>Jupyter</span>
            </div>
          </div>
          
          <div class="tab-content">
            <!-- 可视化界面 -->
            <div v-if="activePanel === 'ui'" class="visual-panel">
              <div class="visual-toolbar">
                <div class="toolbar-section">
                  <span class="toolbar-label">图表类型:</span>
                  <el-radio-group v-model="selectedChartType" size="small">
                    <el-radio-button v-for="type in chartTypes" :key="type.value" :label="type.value">
                      {{ type.label }}
                    </el-radio-button>
                  </el-radio-group>
                </div>
                
                <div class="toolbar-section">
                  <el-button type="primary" size="small">
                    <el-icon><RefreshRight /></el-icon>
                    刷新
                  </el-button>
                  <el-button type="success" size="small">
                    <el-icon><Download /></el-icon>
                    导出
                  </el-button>
                </div>
              </div>
              
              <div class="visual-content">
                <div v-if="selectedResource" class="visualization-result">
                  <h3>{{ selectedResource.name }}</h3>
                  
                  <!-- 数据预览部分 -->
                  <div v-if="selectedResource.resource_type === 'data'" class="data-preview">
                    <div v-loading="previewData.loading">
                      <template v-if="previewData.error">
                        <el-alert
                          :title="previewData.error"
                          type="error"
                          :closable="false"
                          show-icon
                        />
                      </template>
                      
                      <template v-else-if="previewData.rows.length > 0">
                        <el-table
                          :data="previewData.rows"
                          style="width: 100%"
                          height="400"
                          border
                          stripe
                        >
                          <el-table-column
                            v-for="header in previewData.headers"
                            :key="header"
                            :prop="header"
                            :label="header"
                            sortable
                            show-overflow-tooltip
                          />
                        </el-table>
                        
                        <div class="table-info">
                          <p>共 {{ previewData.rows.length }} 条数据</p>
                        </div>
                      </template>
                      
                      <el-empty v-else description="暂无预览数据" />
                    </div>
                  </div>
                  
                  <!-- 原有的图表显示部分 -->
                  <div v-else class="chart-placeholder">
                    <div class="chart-bars">
                      <div 
                        v-for="(value, index) in visualizationResult.data.datasets[0].data" 
                        :key="index"
                        class="chart-bar"
                        :style="{ 
                          height: `${value / 2}px`, 
                          backgroundColor: visualizationResult.data.datasets[0].backgroundColor[index]
                        }"
                      ></div>
                    </div>
                    <div class="chart-labels">
                      <div 
                        v-for="(label, index) in visualizationResult.data.labels" 
                        :key="index"
                        class="chart-label"
                      >
                        {{ label }}
                      </div>
                    </div>
                  </div>
                  
                  <div class="visualization-data">
                    <h4>数据表格视图</h4>
                    <el-table :data="visualizationResult.data.labels.map((label, index) => ({
                      category: label,
                      value: visualizationResult.data.datasets[0].data[index]
                    }))" style="width: 100%">
                      <el-table-column prop="category" label="类别" />
                      <el-table-column prop="value" label="数值" />
                    </el-table>
                  </div>
                </div>
                
                <el-empty 
                  v-else 
                  description="请从左侧选择一个资源" 
                />
              </div>
            </div>
            
            <!-- 命令行终端 -->
            <div v-else-if="activePanel === 'terminal'" class="terminal-panel">
              <div class="terminal-content" ref="terminalRef">
                <div 
                  v-for="(item, index) in commandHistory" 
                  :key="index"
                  class="terminal-line"
                  :class="{ 'user-command': item.type === 'user' }"
                >
                  <span v-if="item.type === 'user'" class="command-prompt">$ </span>
                  <span class="command-text">{{ item.content }}</span>
                </div>
                
                <div class="terminal-input">
                  <span class="command-prompt">$ </span>
                  <input 
                    ref="commandInput"
                    v-model="currentCommand"
                    type="text"
                    class="command-input"
                    placeholder="输入命令..."
                    @keyup.enter="executeCommand"
                  />
                </div>
              </div>
            </div>
            
            <!-- Jupyter 界面 -->
            <div v-else-if="activePanel === 'jupyter'" class="jupyter-panel">
              <div class="jupyter-toolbar">
                <el-button type="primary" size="small" @click="openJupyterWindow">
                  <el-icon><Link /></el-icon>
                  在新窗口打开
                </el-button>
              </div>
              <div class="jupyter-container">
                <iframe
                  :src="jupyterUrl"
                  frameborder="0"
                  class="jupyter-iframe"
                  sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
                  allow="clipboard-read; clipboard-write"
                ></iframe>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sandbox-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.sandbox-header {
  text-align: center;
  margin-bottom: 30px;
}

.sandbox-header h1 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 10px;
}

.sandbox-header p {
  font-size: 16px;
  color: #606266;
}

.sandbox-main {
  display: flex;
  gap: 20px;
  min-height: 700px;
}

/* 资源面板样式 */
.resource-panel {
  width: 320px;
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 15px;
  overflow-y: auto;
}

.panel-title {
  font-size: 18px;
  margin-top: 0;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e4e7ed;
}

.resource-section {
  margin-bottom: 25px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.section-header h4 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.resource-list {
  max-height: 300px;
  overflow-y: auto;
}

.resource-item {
  display: flex;
  padding: 10px;
  background-color: white;
  border-radius: 6px;
  margin-bottom: 10px;
  cursor: pointer;
  border: 1px solid #e4e7ed;
  transition: all 0.3s;
}

.resource-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.resource-item.active {
  border-color: #409EFF;
  background-color: #ecf5ff;
}

.resource-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background-color: #f0f2f5;
  border-radius: 6px;
  margin-right: 10px;
}

.resource-icon .el-icon {
  font-size: 20px;
  color: #409EFF;
}

.resource-info {
  flex: 1;
  min-width: 0;
}

.resource-name {
  display: flex;
  align-items: center;
  font-weight: bold;
  margin-bottom: 5px;
}

.resource-name .el-tag {
  margin-left: 8px;
}

.plugin-version {
  margin-left: 8px;
  font-size: 12px;
  color: #909399;
  font-weight: normal;
}

.resource-description {
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 5px;
}

.plugin-tags {
  display: flex;
  flex-wrap: wrap;
}

.resource-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.computing-details {
  font-size: 12px;
  color: #606266;
  margin-top: 5px;
}

.resource-actions {
  display: flex;
  align-items: center;
}

/* 工作区样式 */
.workspace-panel {
  flex: 1;
  background-color: white;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  overflow: hidden;
}

.workspace-tabs {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tab-header {
  display: flex;
  background-color: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.tab-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  cursor: pointer;
  transition: all 0.3s;
  border-bottom: 2px solid transparent;
}

.tab-item .el-icon {
  margin-right: 5px;
}

.tab-item.active {
  color: #409EFF;
  border-bottom-color: #409EFF;
  background-color: white;
}

.tab-content {
  flex: 1;
  padding: 15px;
  overflow: auto;
}

/* 可视化面板样式 */
.visual-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.visual-toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e4e7ed;
}

.toolbar-section {
  display: flex;
  align-items: center;
}

.toolbar-label {
  margin-right: 10px;
  color: #606266;
}

.visual-content {
  flex: 1;
  overflow: auto;
}

.visualization-result {
  padding: 20px;
}

.visualization-result h3 {
  margin-top: 0;
  margin-bottom: 20px;
  text-align: center;
}

.chart-placeholder {
  height: 300px;
  background-color: #f8f9fa;
  border-radius: 6px;
  padding: 20px;
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
}

.chart-bars {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  flex: 1;
}

.chart-bar {
  width: 50px;
  border-radius: 4px 4px 0 0;
  transition: height 0.5s;
}

.chart-labels {
  display: flex;
  justify-content: space-around;
  margin-top: 10px;
}

.chart-label {
  width: 50px;
  text-align: center;
  font-size: 12px;
  color: #606266;
}

.visualization-data {
  margin-top: 30px;
}

.visualization-data h4 {
  margin-bottom: 15px;
}

/* 终端样式 */
.terminal-panel {
  height: 100%;
}

.terminal-content {
  height: 600px;
  background-color: #1e1e1e;
  color: #f0f0f0;
  padding: 15px;
  font-family: 'Courier New', Courier, monospace;
  overflow-y: auto;
  border-radius: 6px;
}

.terminal-line {
  margin-bottom: 5px;
  white-space: pre-wrap;
  line-height: 1.5;
}

.terminal-line.user-command {
  color: #4fc08d;
}

.command-prompt {
  color: #4fc08d;
}

.command-text {
  word-break: break-all;
}

.terminal-input {
  display: flex;
  margin-top: 10px;
}

.command-input {
  flex: 1;
  background: transparent;
  border: none;
  color: #f0f0f0;
  font-family: 'Courier New', Courier, monospace;
  outline: none;
  font-size: inherit;
}

.data-preview {
  margin: 20px 0;
}

.table-info {
  margin-top: 10px;
  text-align: right;
  color: #606266;
  font-size: 14px;
}

/* Jupyter 面板样式 */
.jupyter-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.jupyter-toolbar {
  padding: 10px;
  border-bottom: 1px solid #e4e7ed;
  background-color: #f5f7fa;
}

.jupyter-container {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.jupyter-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: none;
  background-color: white;
}

.jupyter-fallback {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  background-color: #f5f7fa;
}

@media (max-width: 1200px) {
  .sandbox-main {
    flex-direction: column;
  }
  
  .resource-panel {
    width: 100%;
    height: 300px;
  }
}
</style>