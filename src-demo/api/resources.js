import request from './request'

// 上传数据资源
export async function uploadDataResource(dataForm, dataFile, sampleFile, dataFileHash) {
  const formData = new FormData()
  
  // 添加表单数据
  formData.append('name', dataForm.name)
  formData.append('category', dataForm.category)
  formData.append('dataStructure', dataForm.dataStructure || '')
  formData.append('isConfidential', String(dataForm.isConfidential))
  formData.append('needHosting', String(dataForm.needHosting))
  formData.append('description', dataForm.description)
  
  // 添加文档内容
  if (dataForm.documentation) {
    formData.append('documentation', dataForm.documentation)
  }
  
  // 添加标签
  if (dataForm.tags && dataForm.tags.length > 0) {
    dataForm.tags.forEach((tag) => {
      formData.append('tags', tag) // 后端Gin可以处理同名多值的form字段为slice
    })
  }
  
  // 添加文件类型标签
  if (dataForm.fileExtensions && dataForm.fileExtensions.length > 0) {
    dataForm.fileExtensions.forEach((ext) => {
      formData.append('fileExtensions', ext) // 同上
    })
  }
  
  // 添加文件
  if (dataFile) {
    formData.append('dataFile', dataFile)
  }
  
  // 添加样例文件（如果有）
  if (sampleFile) {
    formData.append('sampleFile', sampleFile)
  }
  
  // 添加文件哈希
  if (dataFileHash) {
    formData.append('fileHash', dataFileHash)
  }
  
  return request({
    url: '/api/v1/resources/data',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取用户上传的所有资源
export function getUserResources() {
  return request({
    url: '/api/v1/user/resources',
    method: 'get'
  })
}

// 获取资源详情
export function getResourceDetail(id) {
  return request({
    url: `/api/v1/resources/${id}`,
    method: 'get'
  })
}

// 预览资源
export function previewResource(id) {
  return request({
    url: `/api/v1/resources/${id}/preview`,
    method: 'get',
    responseType: 'blob', // 使用blob方式接收响应
    timeout: 30000 // 增加超时时间
  })
}

// 发布资源
export function publishResource(id, data) {
  return request({
    url: `/api/v1/resources/${id}/publish`,
    method: 'post',
    data
  })
}

// 删除资源
export function deleteResource(id) {
  return request({
    url: `/api/v1/resources/${id}`,
    method: 'delete'
  })
}

// 获取用户已发布的资源
export function getUserPublishedResources() {
  return request({
    url: '/api/v1/user/published-resources',
    method: 'get'
  })
}

// 下架资源
export function unpublishResource(id) {
  return request({
    url: `/api/v1/resources/${id}/unpublish`,
    method: 'post'
  })
}

// 获取资源分类列表 (真实API，如果后端已实现)
export function getResourceCategories() {
  return request({
    url: '/api/v1/market/categories', // 假设后端有此接口
    method: 'get'
  })
}

// 模拟获取资源分类列表（前端开发时使用）
export function getMockCategories() {
  return Promise.resolve({
    code: 200,
    message: '获取成功',
    data: [
      { value: 'structured', label: '结构化数据' },
      { value: 'unstructured', label: '非结构化数据' },
      { value: 'semi_structured', label: '半结构化数据' },
      { value: 'time_series', label: '时间序列数据' },
      { value: 'spatial', label: '空间数据' },
      { value: 'graph', label: '图数据' },
      { value: 'text', label: '文本数据' },
      { value: 'image', label: '图像数据' },
      { value: 'audio', label: '音频数据' },
      { value: 'video', label: '视频数据' },
      { value: 'sensor', label: '传感器数据' },
      { value: 'mixed', label: '混合数据' }
    ]
  })
}

// 上传算力资源
export async function uploadComputingResource(formData) {
  const data = new FormData()
  
  // 添加基本信息
  data.append('name', formData.name)
  data.append('serverAddress', formData.serverAddress)
  data.append('port', formData.port)
  data.append('cpu', formData.cpu || '')
  data.append('memory', formData.memory || '')
  data.append('gpu', formData.gpu || '')
  data.append('storage', formData.storage || '')
  
  // 添加可用时间
  if (formData.availableTimeStart) {
    data.append('availableTimeStart', formData.availableTimeStart)
  }
  if (formData.availableTimeEnd) {
    data.append('availableTimeEnd', formData.availableTimeEnd)
  }
  
  // 添加系统类型
  if (formData.systemType) {
    data.append('systemType', formData.systemType)
  }
  
  // 添加沙箱版本
  if (formData.sandboxVersion) {
    data.append('sandboxVersion', formData.sandboxVersion)
  }
  
  // 添加描述
  if (formData.description) {
    data.append('description', formData.description)
  }
  
  // 添加标签
  if (formData.tags && formData.tags.length > 0) {
    formData.tags.forEach((tag) => {
      data.append('tags', tag)
    })
  }
  
  return request({
    url: '/api/v1/resources/computing',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传算法资源
export async function uploadAlgorithmResource(formData, algorithmFile, docFile) {
  const data = new FormData()
  
  // 添加基本信息
  data.append('name', formData.name)
  data.append('version', formData.version)
  data.append('category', formData.category)
  
  // 添加可选字段
  if (formData.dependencies) {
    data.append('dependencies', formData.dependencies)
  }
  if (formData.usageGuide) {
    data.append('usageGuide', formData.usageGuide)
  }
  if (formData.apiDocumentation) {
    data.append('apiDocumentation', formData.apiDocumentation)
  }
  if (formData.description) {
    data.append('description', formData.description)
  }
  
  // 添加标签
  if (formData.tags && formData.tags.length > 0) {
    formData.tags.forEach((tag) => {
      data.append('tags', tag)
    })
  }
  
  // 添加文件
  if (algorithmFile) {
    data.append('algorithmFile', algorithmFile)
  }
  if (docFile) {
    data.append('docFile', docFile)
  }
  
  return request({
    url: '/api/v1/resources/algorithm',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
} 