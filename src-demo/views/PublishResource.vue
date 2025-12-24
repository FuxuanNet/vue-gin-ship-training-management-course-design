<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserResources, getResourceDetail, previewResource as fetchResourcePreview, publishResource, deleteResource as removeResource, getUserPublishedResources, unpublishResource } from '../api/resources'
import { marked } from 'marked'
import { useRouter } from 'vue-router'

// 配置marked选项
marked.setOptions({
  breaks: true, // 将回车转换为<br>
  gfm: true,    // 启用GitHub风格Markdown
  headerIds: false, // 关闭自动添加header ids
  mangle: false // 关闭转义
})

// 用户已上传但尚未发布的资源列表
const privateResources = ref([])
// 市场上已发布的资源列表
const marketResources = ref([])

// 加载状态
const loading = ref(false)
const detailLoading = ref(false)

// 搜索查询
const searchQuery = ref('')

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

// 资源类型筛选
const selectedType = ref('all') // 默认显示所有类型
const typeOptions = [
  { value: 'all', label: '全部类型' },
  { value: 'data', label: '数据资源' },
  { value: 'computing', label: '算力资源' },
  { value: 'algorithm', label: '算法资源' }
]

// 筛选资源列表
const filteredPrivateResources = computed(() => {
  if (!privateResources.value) return []
  if (selectedType.value === 'all') {
    return privateResources.value
  } else {
    return privateResources.value.filter(item => item.resource_type === selectedType.value)
  }
})

// 在script setup内部正确使用onMounted
onMounted(() => {
  console.log('组件已挂载，开始获取资源数据')
  // 使用Promise.all同时获取私人和市场资源，降低警告风险
  Promise.all([
    fetchPrivateResources(),
    fetchMarketResources()
  ]).catch(error => {
    console.error('获取资源数据失败:', error)
    ElMessage.error('数据加载失败，请刷新页面重试')
  })
})

// 获取用户已上传但未发布的资源
const fetchPrivateResources = async () => {
  loading.value = true
  try {
    const res = await getUserResources()
    if (res.code === 200) {
      // 只显示未发布的资源
      privateResources.value = res.data.filter(item => !item.is_published && item.status !== 'published')
      
      // 调试输出
      console.log('未发布资源列表:', JSON.parse(JSON.stringify(privateResources.value)))
      console.log('资源类型统计:', privateResources.value.reduce((acc, item) => {
        acc[item.resource_type] = (acc[item.resource_type] || 0) + 1
        return acc
      }, {}))
    } else {
      ElMessage.error(res.message || '获取资源列表失败')
    }
  } catch (error) {
    console.error('获取资源列表错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 获取市场上已发布的资源
const fetchMarketResources = async () => {
  loading.value = true
  try {
    const res = await getUserPublishedResources()
    if (res.code === 200) {
      marketResources.value = res.data
    } else {
      ElMessage.error(res.message || '获取已发布资源列表失败')
    }
  } catch (error) {
    console.error('获取已发布资源列表错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 发布资源对话框
const publishDialog = ref(false)
const currentResource = ref(null)

const publishForm = reactive({
  price: '',
  licenseType: '',
  validPeriod: '',
  usageRestrictions: '',
  agreementTerms: false
})

// 许可证类型选项
const licenseOptions = [
  { label: '永久授权', value: 'permanent' },
  { label: '订阅授权', value: 'subscription' },
  { label: '单次使用', value: 'oneTime' },
  { label: '自定义授权', value: 'custom' }
]

// 发布表单规则
const publishFormRules = {
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', message: '价格必须是数字', trigger: 'blur' },
    { validator: (rule, value, callback) => {
      if (value <= 0) {
        callback(new Error('价格必须大于0'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ],
  licenseType: [
    { required: true, message: '请选择授权类型', trigger: 'change' }
  ],
  agreementTerms: [
    { validator: (rule, value, callback) => {
      if (value !== true) {
        callback(new Error('必须同意协议条款'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ]
}

// 发布资源
const publishFormRef = ref(null)

// 打开发布对话框
const openPublishDialog = (resource) => {
  currentResource.value = resource
  publishDialog.value = true
}

// 发布资源
const handlePublish = async () => {
  if (!publishFormRef.value) return
  
  try {
    // 表单验证
    await publishFormRef.value.validate()
    
    ElMessage.info('正在提交发布请求...')
    
    // 调试输出价格信息
    console.log('发布价格:', publishForm.price, typeof publishForm.price)
    
    // 调用发布API
    const res = await publishResource(currentResource.value.id, {
      price: parseFloat(publishForm.price),
      licenseType: publishForm.licenseType,
      validPeriod: publishForm.validPeriod,
      usageRestrictions: publishForm.usageRestrictions
    })
    
    if (res.code === 200) {
      ElMessage.success('资源发布成功！')
      publishDialog.value = false
      
      // 刷新资源列表
      await Promise.all([
        fetchPrivateResources(),
        fetchMarketResources()
      ])
      
      // 重置表单
      publishForm.price = ''
      publishForm.licenseType = ''
      publishForm.validPeriod = ''
      publishForm.usageRestrictions = ''
      publishForm.agreementTerms = false
    } else {
      ElMessage.error(res.message || '发布失败')
    }
  } catch (error) {
    console.error('发布资源错误:', error)
    ElMessage.error('发布失败: ' + (error.message || '请检查表单信息'))
  }
}

// 资源详情对话框
const resourceDetailDialog = ref(false)
const detailResource = ref(null)
const renderedMarkdown = ref({
  usageGuide: '',
  apiDoc: '',
  description: '',
  documentation: ''
})

// 打开资源详情对话框
const openResourceDetail = async (resource) => {
  detailLoading.value = true
  try {
    // 获取完整的资源详情
    const res = await getResourceDetail(resource.resource_id || resource.id)
    if (res.code === 200) {
      detailResource.value = res.data
      // 如果是已发布资源，添加发布相关信息
      if (resource.publication_id) {
        detailResource.value.price = resource.price
        detailResource.value.license_type = resource.license_type
        detailResource.value.valid_period = resource.valid_period
        detailResource.value.published_at = resource.published_at
        detailResource.value.sales_count = resource.sales_count
      }
      // 预处理Markdown内容
      renderMarkdownContent()
    } else {
      ElMessage.error(res.message || '获取资源详情失败')
    }
  } catch (error) {
    console.error('获取资源详情错误:', error)
    ElMessage.error('网络错误，请稍后重试')
    // 如果API调用失败，使用列表中的简要信息
    detailResource.value = resource
  } finally {
    detailLoading.value = false
    resourceDetailDialog.value = true
  }
}

// 渲染Markdown内容
const renderMarkdownContent = () => {
  if (!detailResource.value) return
  
  const resource = detailResource.value
  
  // 渲染算法资源的Markdown内容
  if (resource.resource_type === 'algorithm') {
    if (resource.usage_guide) {
      renderedMarkdown.value.usageGuide = marked.parse(resource.usage_guide)
    } else {
      renderedMarkdown.value.usageGuide = ''
    }
    
    if (resource.api_documentation) {
      renderedMarkdown.value.apiDoc = marked.parse(resource.api_documentation)
    } else {
      renderedMarkdown.value.apiDoc = ''
    }
  }
  
  // 渲染数据资源的Markdown内容
  if (resource.resource_type === 'data') {
    if (resource.documentation) {
      renderedMarkdown.value.documentation = marked.parse(resource.documentation)
    } else if (resource.description && typeof resource.description === 'string' &&
              (resource.description.includes('*') || 
               resource.description.includes('#') || 
               resource.description.includes('```') ||
               resource.description.includes('>'))) {
      renderedMarkdown.value.description = marked.parse(resource.description)
    } else {
      renderedMarkdown.value.documentation = ''
      renderedMarkdown.value.description = ''
    }
  }
}

// 关闭资源详情对话框
const closeResourceDetail = () => {
  resourceDetailDialog.value = false
  detailResource.value = null
  renderedMarkdown.value = {
    usageGuide: '',
    apiDoc: '',
    description: '',
    documentation: ''
  }
}

// 预览资源
const previewResource = async (resource) => {
  try {
    // 显示加载状态
    ElMessage.info('正在获取预览内容，请稍候...')
    
    if (!resource || (!resource.resource_id && !resource.id)) {
      ElMessage.error('资源ID无效，无法预览')
      return
    }
    
    console.log('准备预览资源:', resource.name, '资源ID:', resource.resource_id || resource.id)
    
    const response = await fetchResourcePreview(resource.resource_id || resource.id)
    
    // 检查响应头中是否有MD5值
    const md5Header = response.headers?.['x-md5']
    
    // 获取文件名（如果有的话）
    const contentDisposition = response.headers?.['content-disposition']
    let fileName = `${resource.name}_preview`
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename=([^;]+)/)
      if (filenameMatch && filenameMatch[1]) {
        fileName = filenameMatch[1].replace(/"/g, '')
      }
    }
    
    // 检查响应类型
    const contentType = response.headers?.['content-type']
    console.log('响应ContentType:', contentType)
    
    // 优化响应类型检测逻辑
    if (contentType && contentType.includes('application/json')) {
      // 如果是JSON响应，说明只返回了MD5值或错误信息
      const reader = new FileReader()
      reader.onload = (e) => {
        try {
          const result = JSON.parse(e.target.result)
          console.log('解析JSON响应:', result)
          if (result.code === 200 && result.data?.md5) {
            showMd5Dialog(result.data.md5, resource.name)
          } else if (result.message) {
            ElMessage.warning(result.message || '资源没有可用的预览内容')
          } else {
            ElMessage.warning('资源没有可用的预览内容')
          }
        } catch (error) {
          console.error('解析JSON失败:', error)
          ElMessage.error('解析预览内容失败')
        }
      }
      reader.onerror = () => {
        console.error('读取响应失败')
        ElMessage.error('读取响应内容失败')
      }
      reader.readAsText(response.data)
    } else {
      // 如果是文件，触发下载
      try {
        console.log('下载文件类型:', contentType)
        
        // 检查响应数据是否有效
        if (!response.data || response.data.size === 0) {
          ElMessage.error('获取到的文件内容为空，请稍后重试')
          return
        }
        
        const blob = new Blob([response.data], { type: contentType || 'application/octet-stream' })
        
        // 使用URL.createObjectURL创建下载链接
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', fileName)
        document.body.appendChild(link)
        
        // 触发点击下载
        link.click()
        
        // 清理
        setTimeout(() => {
          window.URL.revokeObjectURL(url)
          document.body.removeChild(link)
        }, 100)
        
        ElMessage.success('文件下载已开始')
        
        // 如果同时有MD5值，显示MD5对话框
        if (md5Header) {
          showMd5Dialog(md5Header, resource.name)
        }
      } catch (error) {
        console.error('文件下载错误:', error)
        ElMessage.error('文件下载失败, 错误详情: ' + error.message)
      }
    }
  } catch (error) {
    console.error('预览资源错误:', error)
    const errorMessage = error.message || (error.code ? `错误码: ${error.code}` : '未知错误')
    ElMessage.error('预览资源失败: ' + errorMessage)
  }
}

// MD5对话框
const md5DialogVisible = ref(false)
const currentMd5 = ref('')
const currentResourceName = ref('')

// 显示MD5对话框
const showMd5Dialog = (md5, resourceName) => {
  currentMd5.value = md5
  currentResourceName.value = resourceName
  md5DialogVisible.value = true
}

// 复制MD5值
const copyMd5 = () => {
  navigator.clipboard.writeText(currentMd5.value)
    .then(() => {
      ElMessage.success('MD5已复制到剪贴板')
    })
    .catch(() => {
      ElMessage.error('复制失败，请手动复制')
    })
}

// 删除资源
const deleteResource = (resource) => {
  ElMessageBox.confirm(`确定要删除资源 "${resource.name}" 吗？此操作不可恢复！`, '删除确认', {
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
    type: 'warning'
  })
  .then(async () => {
    try {
      // 调用删除API
      const res = await removeResource(resource.id)
      if (res.code === 200) {
        ElMessage.success(`资源 ${resource.name} 已成功删除`)
        // 更新列表
        privateResources.value = privateResources.value.filter(r => r.id !== resource.id)
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    } catch (error) {
      console.error('删除资源错误:', error)
      ElMessage.error('删除失败，请稍后重试')
    }
  })
  .catch(() => {
    // 用户取消删除
  })
}

// 分页设置
const privateCurrentPage = ref(1)
const marketCurrentPage = ref(1)
const pageSize = ref(5)

// 下架资源
const removeFromMarket = (resource) => {
  ElMessageBox.confirm(`确定要下架资源 "${resource.name}" 吗？下架后可重新发布。`, '下架确认', {
    confirmButtonText: '确认下架',
    cancelButtonText: '取消',
    type: 'warning'
  })
  .then(async () => {
    try {
      // 调用下架API
      const res = await unpublishResource(resource.resource_id)
      if (res.code === 200) {
        ElMessage.success(`资源 ${resource.name} 已成功下架`)
        // 更新列表
        marketResources.value = marketResources.value.filter(r => r.resource_id !== resource.resource_id)
        // 重新获取未发布资源列表
        fetchPrivateResources()
      } else {
        ElMessage.error(res.message || '下架失败')
      }
    } catch (error) {
      console.error('下架资源错误:', error)
      ElMessage.error('下架失败，请稍后重试')
    }
  })
  .catch(() => {
    // 用户取消下架
  })
}

// 更改合约
const changeContract = (resource) => {
  ElMessage.info(`准备更改 ${resource.name} 的合约条款`)
  // 实际开发中这里应该打开合约更改对话框
}

// 更改价格
const changePrice = (resource) => {
  ElMessage.info(`准备更改 ${resource.name} 的价格`)
  // 实际开发中这里应该打开价格更改对话框
}

// 重新上架
const relistResource = (resource) => {
  ElMessage.success(`${resource.name} 已重新上架`)
  // 实际开发中这里应该调用API重新上架资源
}

// 导入Router
const router = useRouter()

// 根据资源类型处理查看详细的行为
const handleViewDetail = (resourceInDialog) => {
  if (!resourceInDialog) {
    console.error('handleViewDetail: resourceInDialog is null or undefined');
    return;
  }

  const localResourceType = resourceInDialog.resource_type;
  console.log('handleViewDetail - Resource Type:', localResourceType, 'Resource ID:', resourceInDialog.id);
  console.log('handleViewDetail - Full resourceInDialog:', JSON.stringify(resourceInDialog, null, 2));


  if (localResourceType === 'computing') {
    console.log('Navigating to sandbox for computing resource:', resourceInDialog.id);
    router.push({
      path: '/sandbox',
      query: { 
        resource_id: resourceInDialog.id,
        name: resourceInDialog.name
      }
    });
    resourceDetailDialog.value = false; // Close the current detail dialog
    return;
  }

  let doc = null;
  let usage = null;
  let api = null;
  let tab = 'documentation';

  if (localResourceType === 'data') {
    doc = resourceInDialog.documentation;
    console.log('Data resource documentation:', doc);
    console.log('Data resource documentation length:', doc ? doc.length : 'N/A');
    console.log('Data resource documentation type:', doc ? typeof doc : 'N/A');
    
    if (doc && doc.length > 0) {
      // 检查前20字符和后20字符，确认内容正确传递
      console.log('Documentation 开始部分:', doc.substring(0, 20));
      console.log('Documentation 结束部分:', doc.substring(doc.length - 20));
    } else {
      console.warn('数据资源的documentation字段为空，请检查后端返回数据');
    }
  } else if (localResourceType === 'algorithm') {
    usage = resourceInDialog.usage_guide;
    api = resourceInDialog.api_documentation;
    tab = 'usageGuide'; // Default to usage guide for algorithm
    console.log('Algorithm usage guide length:', usage ? usage.length : 'N/A');
    console.log('Algorithm API doc length:', api ? api.length : 'N/A');
  }

  currentDocument.value = {
    resourceType: localResourceType,
    resourceName: resourceInDialog.name,
    documentation: doc,
    usageGuide: usage,
    apiDoc: api,
    activeDocTab: tab
  };
  
  console.log('CurrentDocument for dialog:', JSON.parse(JSON.stringify(currentDocument.value)));

  documentDialog.value = true;
  resourceDetailDialog.value = false; // Close the first dialog before opening the second
  
  // Ensure DOM is updated before rendering markdown, especially if dialogs are re-used
  nextTick(() => {
    renderDocumentContent();
  });
}

// 渲染Markdown内容 (for documentDialog)
const renderDocumentContent = () => {
  if (!currentDocument.value) {
    console.error('renderDocumentContent: currentDocument is null');
    renderedDocument.value = { documentation: '', usageGuide: '', apiDoc: '' };
    return;
  }
  console.log('Rendering document content for type:', currentDocument.value.resourceType);

  // 检查documentation字段是否存在且非空
  if (currentDocument.value.resourceType === 'data') {
    if (currentDocument.value.documentation) {
      renderedDocument.value.documentation = marked.parse(currentDocument.value.documentation);
      console.log('已渲染数据资源文档，原始内容长度:', currentDocument.value.documentation.length);
      console.log('文档内容前50字符:', currentDocument.value.documentation.substring(0,50));
    } else {
      renderedDocument.value.documentation = '该资源没有详细说明文档。';
      console.warn('数据资源文档内容为空');
    }
  } else {
    renderedDocument.value.documentation = currentDocument.value.documentation ? marked.parse(currentDocument.value.documentation) : '该资源没有详细说明文档。';
  }
  
  renderedDocument.value.usageGuide = currentDocument.value.usageGuide ? marked.parse(currentDocument.value.usageGuide) : '该资源没有使用指南。';
  renderedDocument.value.apiDoc = currentDocument.value.apiDoc ? marked.parse(currentDocument.value.apiDoc) : '该资源没有API文档。';
  
  console.log('Rendered content prepared:', {
      doc: renderedDocument.value.documentation.substring(0,50) + '...',
      guide: renderedDocument.value.usageGuide.substring(0,50) + '...',
      api: renderedDocument.value.apiDoc.substring(0,50) + '...'
  });
}

// 文档对话框
const documentDialog = ref(false)
const currentDocument = ref(null)
const renderedDocument = ref({
  documentation: '',
  usageGuide: '',
  apiDoc: ''
})
</script>

<template>
  <div class="publish-container">
    <div class="publish-header">
      <h1>资源发布</h1>
      <p>将您的资源发布到交易市场，与其他用户共享并获取收益</p>
    </div>
    
    <!-- 我的资源 -->
    <div class="my-resources-section">
      <h2 class="section-title">我的资源</h2>
      <el-card class="resource-list-card">
        <template #header>
          <div class="card-header">
            <div>可发布资源列表</div>
            <div class="header-actions">
              <el-select v-model="selectedType" placeholder="资源类型" size="small" style="width: 120px; margin-right: 10px;">
                <el-option
                  v-for="item in typeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
              <el-button type="primary" size="small" @click="$router.push('/upload')">
                <el-icon><Plus /></el-icon>上传新资源
              </el-button>
            </div>
          </div>
        </template>
        
        <div v-loading="loading">
          <el-empty v-if="!privateResources || privateResources.length === 0" description="暂无可发布的资源" />
          
          <el-table v-else-if="privateResources && privateResources.length > 0" :data="filteredPrivateResources.slice((privateCurrentPage - 1) * pageSize, privateCurrentPage * pageSize)" style="width: 100%">
            <el-table-column prop="name" label="资源名称" min-width="180">
              <template #default="{ row }">
                <div class="resource-name">
                  <el-tag size="small" :type="typeColorMap[row.resource_type]">{{ typeNameMap[row.resource_type] }}</el-tag>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
            
            <el-table-column label="标签" width="200">
              <template #default="{ row }">
                <el-tag 
                  v-for="tag in (row.tags || [])" 
                  :key="tag" 
                  size="small"
                  effect="plain"
                  class="resource-tag"
                >
                  {{ tag }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column label="创建时间" width="160">
              <template #default="{ row }">
                {{ row.created_at }}
              </template>
            </el-table-column>
            
            <el-table-column label="大小" width="100">
              <template #default="{ row }">
                {{ row.file_size ? (row.file_size / 1024 / 1024).toFixed(2) + ' MB' : '-' }}
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="300" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="info" @click="openResourceDetail(row)">详情</el-button>
                <el-button size="small" type="success" @click="previewResource(row)">预览</el-button>
                <el-button size="small" type="primary" @click="openPublishDialog(row)">发布</el-button>
                <el-button size="small" type="danger" @click="deleteResource(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination-container" v-if="privateResources.length > 0">
            <el-pagination
              v-model:current-page="privateCurrentPage"
              :page-size="pageSize"
              layout="prev, pager, next, jumper"
              :total="privateResources.length"
              background
            />
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 市场资源 -->
    <div class="published-resources-section">
      <h2 class="section-title">已发布资源</h2>
      <el-card class="resource-list-card">
        <template #header>
          <div class="card-header">
            <div>已发布资源</div>
            <el-input
              v-model="searchQuery"
              placeholder="搜索资源"
              prefix-icon="Search"
              style="width: 200px;"
            />
          </div>
        </template>
        
        <div>
          <el-empty v-if="!marketResources || marketResources.length === 0" description="暂无已发布的资源" />
          
          <el-table v-else-if="marketResources && marketResources.length > 0" :data="marketResources.slice((marketCurrentPage - 1) * pageSize, marketCurrentPage * pageSize)" style="width: 100%">
            <el-table-column prop="name" label="资源名称" min-width="150" />
            <el-table-column prop="resource_type" label="资源类型" width="100">
              <template #default="scope">
                <el-tag :type="typeColorMap[scope.row.resource_type]">
                  {{ typeNameMap[scope.row.resource_type] }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="tags" label="标签" min-width="150">
              <template #default="scope">
                <div class="table-tags">
                  <el-tag 
                    v-for="tag in (scope.row.tags || [])" 
                    :key="tag" 
                    size="small" 
                    effect="plain"
                    class="table-tag"
                  >
                    {{ tag }}
                  </el-tag>
                  <span v-if="!scope.row.tags || scope.row.tags.length === 0">-</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="publisher" label="发布者" width="150" />
            <el-table-column prop="price" label="价格" width="120">
              <template #default="{ row }">
                <span class="price">¥{{ row.price }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="sales_count" label="销量" width="100" />
            <el-table-column prop="published_at" label="发布时间" width="160" />
            <el-table-column label="操作" width="280">
              <template #default="{ row }">
                <el-button size="small" type="info" @click="openResourceDetail(row)">详情</el-button>
                <el-button size="small" type="success" @click="previewResource(row)">预览</el-button>
                <el-button size="small" type="danger" @click="removeFromMarket(row)">下架</el-button>
                <el-button size="small" type="primary" @click="changeContract(row)">更改合约</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination-container" v-if="marketResources.length > 0">
            <el-pagination
              v-model:current-page="marketCurrentPage"
              :page-size="pageSize"
              layout="prev, pager, next, jumper"
              :total="marketResources.length"
              background
            />
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 发布资源对话框 -->
    <el-dialog
      v-model="publishDialog"
      title="发布资源到交易市场"
      width="600px"
    >
      <div v-if="currentResource">
        <div class="publish-resource-info">
          <h3>{{ currentResource.name }}</h3>
          <p>{{ currentResource.description }}</p>
          <div class="resource-tags">
            <el-tag 
              v-for="tag in (currentResource.tags || [])" 
              :key="tag" 
              size="small"
              effect="plain"
              class="resource-tag"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
        
        <el-form :model="publishForm" label-width="100px" class="publish-form" ref="publishFormRef" :rules="publishFormRules">
          <el-form-item label="定价" required>
            <el-input-number v-model="publishForm.price" :min="0" :step="100" :precision="2" />
            <span class="price-unit">元</span>
          </el-form-item>
          
          <el-form-item label="授权类型" required>
            <el-select v-model="publishForm.licenseType" placeholder="请选择授权类型">
              <el-option 
                v-for="option in licenseOptions" 
                :key="option.value" 
                :label="option.label" 
                :value="option.value" 
              />
            </el-select>
          </el-form-item>
          
          <el-form-item label="有效期限">
            <el-input v-model="publishForm.validPeriod" placeholder="例如：1年、永久等" />
          </el-form-item>
          
          <el-form-item label="使用限制">
            <el-input 
              v-model="publishForm.usageRestrictions" 
              type="textarea" 
              :rows="3" 
              placeholder="请输入使用限制条款"
            />
          </el-form-item>
          
          <el-form-item label="协议条款" required>
            <el-checkbox v-model="publishForm.agreementTerms">
              我已阅读并同意《资源交易服务协议》和《平台规范》
            </el-checkbox>
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="publishDialog = false">取消</el-button>
          <el-button type="primary" @click="handlePublish">发布</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 资源详情对话框 -->
    <el-dialog
      v-model="resourceDetailDialog"
      title="资源详情"
      width="700px"
      destroy-on-close
    >
      <div v-loading="detailLoading">
        <div v-if="detailResource" class="resource-detail">
          <div class="resource-detail-header">
            <h2>{{ detailResource.name }}</h2>
            <el-tag :type="typeColorMap[detailResource.resource_type || detailResource.type]">
              {{ typeNameMap[detailResource.resource_type || detailResource.type] }}
            </el-tag>
          </div>
          
          <el-tabs>
            <el-tab-pane label="基本信息">
              <div class="resource-meta">
                <div class="meta-item">
                  <span class="meta-label">创建时间：</span>
                  <span>{{ detailResource.created_at || detailResource.publishTime }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.file_size">
                  <span class="meta-label">文件大小：</span>
                  <span>{{ (detailResource.file_size / 1024 / 1024).toFixed(2) }} MB</span>
                </div>
                <div class="meta-item" v-if="detailResource.is_confidential !== undefined">
                  <span class="meta-label">保密数据：</span>
                  <span>{{ detailResource.is_confidential ? '是' : '否' }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.file_hash">
                  <span class="meta-label">MD5：</span>
                  <span>{{ detailResource.file_hash }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.publisher">
                  <span class="meta-label">发布者：</span>
                  <span>{{ detailResource.publisher }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.sales !== undefined">
                  <span class="meta-label">销量：</span>
                  <span>{{ detailResource.sales }} 次</span>
                </div>
                <div class="meta-item price-item" v-if="detailResource.price !== undefined">
                  <span class="meta-label">价格：</span>
                  <span class="detail-price">¥{{ detailResource.price }}</span>
                </div>
              </div>
              
              <div class="resource-tags" v-if="detailResource.tags && detailResource.tags.length > 0">
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
            </el-tab-pane>
            
            <!-- 数据资源特有的信息 -->
            <el-tab-pane v-if="detailResource.resource_type === 'data'" label="数据详情">
              <div class="resource-meta">
                <div class="meta-item" v-if="detailResource.category">
                  <span class="meta-label">数据分类：</span>
                  <span>{{ detailResource.category }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.data_structure">
                  <span class="meta-label">数据结构：</span>
                  <span>{{ detailResource.data_structure }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.need_hosting !== undefined">
                  <span class="meta-label">需要托管：</span>
                  <span>{{ detailResource.need_hosting ? '是' : '否' }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.file_extensions && detailResource.file_extensions.length > 0">
                  <span class="meta-label">文件类型：</span>
                  <span>{{ Array.isArray(detailResource.file_extensions) ? detailResource.file_extensions.join(', ') : detailResource.file_extensions }}</span>
                </div>
              </div>
              <div v-if="detailResource.documentation" class="resource-documentation-preview">
                <div class="meta-label">文档预览：</div>
                <div class="documentation-preview-content">{{ detailResource.documentation.substring(0, 200) }}{{ detailResource.documentation.length > 200 ? '...' : '' }}</div>
              </div>
            </el-tab-pane>
            
            <!-- 算力资源特有的信息 -->
            <el-tab-pane v-if="detailResource.resource_type === 'computing'" label="算力详情">
              <div class="resource-meta">
                <div class="meta-item"><span class="meta-label">服务名称：</span><span>{{ detailResource.name }}</span></div>
                <div class="meta-item"><span class="meta-label">服务描述：</span><span>{{ detailResource.description }}</span></div>
                <div class="meta-item" v-if="detailResource.system_type"><span class="meta-label">系统类型：</span><span>{{ detailResource.system_type }}</span></div>
                <div class="meta-item" v-if="detailResource.server_address"><span class="meta-label">服务器地址：</span><span>{{ detailResource.server_address }}</span></div>
                <div class="meta-item" v-if="detailResource.port"><span class="meta-label">开放端口：</span><span>{{ detailResource.port }}</span></div>
                <div class="meta-item" v-if="detailResource.cpu"><span class="meta-label">CPU配置：</span><span>{{ detailResource.cpu }}</span></div>
                <div class="meta-item" v-if="detailResource.memory"><span class="meta-label">内存配置：</span><span>{{ detailResource.memory }}</span></div>
                <div class="meta-item" v-if="detailResource.gpu"><span class="meta-label">GPU配置：</span><span>{{ detailResource.gpu }}</span></div>
                <div class="meta-item" v-if="detailResource.storage"><span class="meta-label">存储配置：</span><span>{{ detailResource.storage }}</span></div>
                <div class="meta-item" v-if="detailResource.available_time_start || detailResource.available_time_end">
                  <span class="meta-label">可用时间：</span>
                  <span>{{ detailResource.available_time_start || '未指定' }} - {{ detailResource.available_time_end || '未指定' }}</span>
                </div>
                <div class="meta-item" v-if="detailResource.sandbox_version"><span class="meta-label">沙箱版本：</span><span>{{ detailResource.sandbox_version }}</span></div>
              </div>
            </el-tab-pane>
            
            <!-- 算法资源特有的信息 -->
            <el-tab-pane v-if="detailResource.resource_type === 'algorithm'" label="算法详情">
              <div class="resource-meta">
                 <div class="meta-item"><span class="meta-label">算法名称：</span><span>{{ detailResource.name }}</span></div>
                <div class="meta-item"><span class="meta-label">算法描述：</span><span>{{ detailResource.description }}</span></div>
                <div class="meta-item" v-if="detailResource.version"><span class="meta-label">算法版本：</span><span>{{ detailResource.version }}</span></div>
                <div class="meta-item" v-if="detailResource.category"><span class="meta-label">算法分类：</span><span>{{ detailResource.category }}</span></div>
                <div class="meta-item" v-if="detailResource.dependencies"><span class="meta-label">依赖项：</span><pre>{{ detailResource.dependencies }}</pre></div>
              </div>
              <!-- 使用指南预览 -->
              <div v-if="detailResource.usage_guide" class="resource-documentation-preview">
                <div class="meta-label">使用指南（预览）：</div>
                <div class="documentation-preview-content">{{ detailResource.usage_guide.substring(0, 150) }}{{ detailResource.usage_guide.length > 150 ? '...' : '' }}</div>
              </div>
              
              <!-- API文档预览 -->
              <div v-if="detailResource.api_documentation" class="resource-documentation-preview">
                <div class="meta-label">API文档（预览）：</div>
                <div class="documentation-preview-content">{{ detailResource.api_documentation.substring(0, 150) }}{{ detailResource.api_documentation.length > 150 ? '...' : '' }}</div>
              </div>
            </el-tab-pane>
          </el-tabs>
          
          <div class="resource-detail-actions">
            <el-button type="primary" size="small" @click="handleViewDetail(detailResource)">查看详细</el-button>
          </div>
        </div>
        <el-empty v-else description="资源详情加载失败" />
      </div>
    </el-dialog>
    
    <!-- 文档查看对话框 -->
    <el-dialog
      v-model="documentDialog"
      :title="currentDocument?.resourceName || '详细文档'"
      width="800px"
      destroy-on-close
      custom-class="document-dialog"
    >
      <div class="document-dialog-content" v-if="currentDocument">
        <!-- 算法资源 - 显示两个Tab：使用指南和API文档 -->
        <el-tabs v-if="currentDocument.resourceType === 'algorithm'" v-model="currentDocument.activeDocTab">
          <el-tab-pane label="使用指南" name="usageGuide">
            <div class="markdown-content styled-markdown" v-html="renderedDocument.usageGuide || '暂无使用指南'"></div>
          </el-tab-pane>
          <el-tab-pane label="API文档" name="apiDoc">
            <div class="markdown-content styled-markdown" v-html="renderedDocument.apiDoc || '暂无API文档'"></div>
          </el-tab-pane>
        </el-tabs>
        
        <!-- 数据资源 - 显示Markdown文档 -->
        <div v-else-if="currentDocument.resourceType === 'data'" class="markdown-content styled-markdown" v-html="renderedDocument.documentation || '暂无详细文档'"></div>

        <!-- 其他类型资源 -->
        <el-empty v-else description="该类型资源无详细文档" />
        
        <!-- 文档操作按钮 -->
        <div class="document-actions">
          <el-button @click="documentDialog = false">关闭</el-button>
        </div>
      </div>
      <el-empty v-else description="暂无详细文档" />
    </el-dialog>
    
    <!-- MD5对话框 -->
    <el-dialog
      v-model="md5DialogVisible"
      :title="`${currentResourceName} - MD5值`"
      width="500px"
      destroy-on-close
    >
      <div class="md5-dialog-content">
        <p>此资源为保密资源，仅提供MD5校验值供参考：</p>
        <el-input
          v-model="currentMd5"
          readonly
          class="md5-input"
        >
          <template #append>
            <el-button @click="copyMd5">复制</el-button>
          </template>
        </el-input>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.publish-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.publish-header {
  text-align: center;
  margin-bottom: 30px;
}

.publish-header h1 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 10px;
}

.publish-header p {
  font-size: 16px;
  color: #606266;
}

.section-title {
  font-size: 20px;
  margin-bottom: 20px;
  color: #303133;
  position: relative;
  padding-left: 12px;
}

.section-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 20px;
  background-color: #409EFF;
  border-radius: 2px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

.resource-list-card {
  margin-bottom: 30px;
}

.resource-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.resource-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.publish-resource-info {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.publish-resource-info h3 {
  margin-bottom: 10px;
}

.price-unit {
  margin-left: 5px;
}

.resource-detail {
  padding: 10px;
}

.resource-detail-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 1px dashed #ebeef5;
  padding-bottom: 15px;
}

.resource-detail-header h2 {
  margin: 0;
}

.resource-meta {
  margin-bottom: 15px;
}

.meta-item {
  margin-bottom: 8px;
}

.meta-label {
  font-weight: bold;
  color: #606266;
  margin-right: 5px;
}

.detail-price {
  color: #f56c6c;
  font-size: 18px;
  font-weight: bold;
}

.resource-description {
  margin-top: 15px;
}

.resource-detail-actions {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.markdown-content) {
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  background-color: #f9f9f9;
  min-height: 100px;
  max-height: 400px;
  overflow-y: auto;
}

:deep(.markdown-content pre) {
  background-color: #f1f1f1;
  padding: 10px;
  border-radius: 5px;
  overflow-x: auto;
}

:deep(.markdown-content code) {
  background-color: #f1f1f1;
  padding: 2px 5px;
  border-radius: 3px;
}

:deep(.markdown-content table) {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 15px;
}

:deep(.markdown-content th),
:deep(.markdown-content td) {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

:deep(.markdown-content th) {
  background-color: #f2f2f2;
}

:deep(.markdown-content img) {
  max-width: 100%;
}

.md5-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.md5-input {
  font-family: monospace;
}

/* 文档对话框样式 */
.document-dialog {
  display: flex;
  flex-direction: column;
}

.document-dialog .el-dialog__body {
  height: 60vh;
  overflow: auto;
  padding: 0;
}

.document-dialog-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px;
}

.markdown-content.styled-markdown {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #fafafa;
  border: 1px solid #eaeaea;
  border-radius: 4px;
  line-height: 1.6;
}

.document-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.resource-documentation-preview {
  margin-top: 15px;
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  background-color: #f9fafb;
}

.documentation-preview-content {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.6;
  color: #606266;
  white-space: pre-wrap;
}

/* Markdown样式 */
.styled-markdown h1 {
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 2em;
}

.styled-markdown h2 {
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 1.5em;
}

.styled-markdown h3 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 1.25em;
}

.styled-markdown code {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
}

.styled-markdown pre {
  margin-bottom: 16px;
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
}

.styled-markdown pre code {
  background-color: transparent;
  padding: 0;
  margin: 0;
  font-size: inherit;
  word-break: normal;
  white-space: pre;
  border: 0;
}

.styled-markdown blockquote {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin: 16px 0;
}

.styled-markdown table {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
}

.styled-markdown table th,
.styled-markdown table td {
  padding: 6px 13px;
  border: 1px solid #dfe2e5;
}

.styled-markdown table tr {
  background-color: #fff;
  border-top: 1px solid #c6cbd1;
}

.styled-markdown table tr:nth-child(2n) {
  background-color: #f6f8fa;
}

.styled-markdown img {
  max-width: 100%;
  box-sizing: initial;
}

.styled-markdown ul,
.styled-markdown ol {
  padding-left: 2em;
  margin-top: 0;
  margin-bottom: 16px;
}

.styled-markdown li {
  word-wrap: break-all;
}

.styled-markdown hr {
  height: 0.25em;
  padding: 0;
  margin: 24px 0;
  background-color: #e1e4e8;
  border: 0;
}

.table-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.table-tag {
  margin-right: 5px;
}
</style> 