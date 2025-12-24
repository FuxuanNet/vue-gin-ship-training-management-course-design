<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, computed } from 'vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { SetUp, Edit, Plus, Delete, UploadFilled, Document, View, ArrowRight } from '@element-plus/icons-vue'
import { uploadAlgorithmResource } from '../api/resources'
import { marked } from 'marked'

const props = defineProps({
  active: Boolean
})

const emit = defineEmits(['cancel', 'success'])

// 当前步骤
const currentStep = ref(1)

// 本地存储键名
const ALGORITHM_FORM_STORAGE_KEY = 'algorithmUploadForm'
const ALGORITHM_USAGE_GUIDE_KEY = 'algorithmUsageGuide'
const ALGORITHM_API_DOC_KEY = 'algorithmApiDoc'

// 算法资源表单
const algorithmForm = reactive({
  name: '',
  version: '',
  category: '',
  dependencies: '',
  description: '',
  tags: []
})

// Markdown 编辑器相关
const usageGuideContent = ref('')
const apiDocContent = ref('')
const usageGuidePreviewMode = ref(false)
const apiDocPreviewMode = ref(false)

const compiledUsageGuide = computed(() => marked(usageGuideContent.value || ''))
const compiledApiDoc = computed(() => marked(apiDocContent.value || ''))

const algorithmFile = ref(null)
const docFile = ref(null)
const algorithmFileEl = ref(null)
const docFileEl = ref(null)

const categoryOptions = [
  { value: 'ml', label: '机器学习' },
  { value: 'dl', label: '深度学习' },
  { value: 'data_processing', label: '数据处理' },
  { value: 'image_processing', label: '图像处理' },
  { value: 'nlp', label: '自然语言处理' },
  { value: 'statistics', label: '统计分析' },
  { value: 'optimization', label: '优化算法' },
  { value: 'visualization', label: '数据可视化' },
  { value: 'other', label: '其他' }
]

const tagInput = ref('')

// 保存表单状态到本地存储
const saveFormState = () => {
  localStorage.setItem(ALGORITHM_FORM_STORAGE_KEY, JSON.stringify(algorithmForm))
}

// 保存使用指南内容
const saveUsageGuide = () => {
  localStorage.setItem(ALGORITHM_USAGE_GUIDE_KEY, usageGuideContent.value)
}

// 保存API文档内容
const saveApiDoc = () => {
  localStorage.setItem(ALGORITHM_API_DOC_KEY, apiDocContent.value)
}

// 从本地存储加载表单状态
const loadFormState = () => {
  const savedState = localStorage.getItem(ALGORITHM_FORM_STORAGE_KEY)
  if (savedState) {
    try {
      const parsedState = JSON.parse(savedState)
      Object.assign(algorithmForm, parsedState)
    } catch (e) {
      console.error('解析保存的表单数据失败:', e)
    }
  }
  
  // 加载使用指南
  const savedUsageGuide = localStorage.getItem(ALGORITHM_USAGE_GUIDE_KEY)
  if (savedUsageGuide) {
    usageGuideContent.value = savedUsageGuide
  }
  
  // 加载API文档
  const savedApiDoc = localStorage.getItem(ALGORITHM_API_DOC_KEY)
  if (savedApiDoc) {
    apiDocContent.value = savedApiDoc
  }
}

// 清除表单数据
const clearFormState = () => {
  Object.assign(algorithmForm, {
    name: '',
    version: '',
    category: '',
    dependencies: '',
    description: '',
    tags: []
  })
  usageGuideContent.value = ''
  apiDocContent.value = ''
  clearAlgorithmFile()
  clearDocFile()
  localStorage.removeItem(ALGORITHM_FORM_STORAGE_KEY)
  localStorage.removeItem(ALGORITHM_USAGE_GUIDE_KEY)
  localStorage.removeItem(ALGORITHM_API_DOC_KEY)
}

const handleAlgorithmFileChange = (fileInfo) => {
  if (!fileInfo || !fileInfo.raw) {
    algorithmFile.value = null
    if (algorithmFileEl.value && typeof algorithmFileEl.value.clearFiles === 'function') {
      algorithmFileEl.value.clearFiles()
    }
    return
  }
  const file = fileInfo.raw

  const maxSize = 200 * 1024 * 1024 // 200MB
  if (file.size > maxSize) {
    ElMessage.warning('算法文件大小超过200MB')
    if (algorithmFileEl.value && typeof algorithmFileEl.value.clearFiles === 'function') {
      algorithmFileEl.value.clearFiles()
    }
    algorithmFile.value = null
    return
  }
  algorithmFile.value = file
}

const handleDocFileChange = (fileInfo) => {
  if (!fileInfo || !fileInfo.raw) {
    docFile.value = null
    if (docFileEl.value && typeof docFileEl.value.clearFiles === 'function') {
      docFileEl.value.clearFiles()
    }
    return
  }
  const file = fileInfo.raw

  const maxSize = 50 * 1024 * 1024 // 50MB
  if (file.size > maxSize) {
    ElMessage.warning('文档文件大小超过50MB')
    if (docFileEl.value && typeof docFileEl.value.clearFiles === 'function') {
      docFileEl.value.clearFiles()
    }
    docFile.value = null
    return
  }
  docFile.value = file
}

const clearAlgorithmFile = () => {
  algorithmFile.value = null
  if (algorithmFileEl.value && typeof algorithmFileEl.value.clearFiles === 'function') {
    algorithmFileEl.value.clearFiles()
  }
}

const clearDocFile = () => {
  docFile.value = null
  if (docFileEl.value && typeof docFileEl.value.clearFiles === 'function') {
    docFileEl.value.clearFiles()
  }
}

// 验证步骤1
const validateStep1 = () => {
  if (!algorithmForm.name) {
    ElMessage.warning('请输入算法名称')
    return false
  }
  if (!algorithmForm.version) {
    ElMessage.warning('请输入算法版本')
    return false
  }
  if (!algorithmForm.category) {
    ElMessage.warning('请选择算法分类')
    return false
  }
  if (!algorithmFile.value) {
    ElMessage.warning('请上传算法文件')
    return false
  }
  return true
}

// 验证步骤2 - 使用指南
const validateStep2 = () => {
  if (!usageGuideContent.value.trim()) {
    ElMessage.warning('请填写使用指南内容')
    return false
  }
  return true
}

// 验证步骤3 - API文档
const validateStep3 = () => {
  // API文档可以选填
  return true
}

// 进入下一步
const nextStep = () => {
  if (currentStep.value === 1) {
    if (validateStep1()) {
      saveFormState()
      currentStep.value = 2
      // 初始化使用指南模板
      if (!usageGuideContent.value) {
        usageGuideContent.value = `# ${algorithmForm.name || '算法'} 使用指南\n\n## 简介\n请在此处描述算法的基本功能和用途...\n\n## 环境要求\n请描述算法运行所需的环境和依赖...\n\n## 参数说明\n请描述算法的输入参数...\n\n## 使用方法\n请提供算法的使用示例...\n\n## 注意事项\n请描述使用算法时需要注意的问题...\n`
      }
    }
  } else if (currentStep.value === 2) {
    if (validateStep2()) {
      saveUsageGuide()
      currentStep.value = 3
      // 初始化API文档模板
      if (!apiDocContent.value) {
        apiDocContent.value = `# ${algorithmForm.name || '算法'} API文档\n\n## API概述\n请在此处描述API的基本功能和用途...\n\n## 接口列表\n请列出算法提供的所有接口...\n\n## 接口详情\n### 接口1\n- **函数名**: example_function\n- **参数**: \n  - param1: 参数1说明\n  - param2: 参数2说明\n- **返回值**: 返回值说明\n- **示例**:\n\`\`\`python\nresult = example_function(param1, param2)\n\`\`\`\n\n## 错误码\n请列出可能的错误码和对应的说明...\n`
      }
    }
  }
}

// 返回上一步
const prevStep = () => {
  if (currentStep.value > 1) {
    if (currentStep.value === 2) {
      saveUsageGuide()
    } else if (currentStep.value === 3) {
      saveApiDoc()
    }
    currentStep.value--
  }
}

// 取消
const handleCancel = () => {
  ElMessageBox.confirm(
    '确定要取消上传吗？已填写的内容将会保存（如果未提交）。',
    '取消上传',
    { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
  ).then(() => {
    // 保存当前步骤的内容
    if (currentStep.value === 1) {
      saveFormState()
    } else if (currentStep.value === 2) {
      saveUsageGuide()
    } else if (currentStep.value === 3) {
      saveApiDoc()
    }
    emit('cancel')
  }).catch(() => {
    // 用户取消对话框，不做任何处理
  })
}

// 提交
const handleSubmit = async () => {
  // 最终验证
  if (!validateStep3()) return
  
  // 组装完整表单数据
  const completeForm = {
    ...algorithmForm,
    usageGuide: usageGuideContent.value,
    apiDocumentation: apiDocContent.value
  }

  const loading = ElLoading.service({
    lock: true,
    text: '正在上传算法资源...',
    background: 'rgba(0, 0, 0, 0.7)'
  })

  try {
    await uploadAlgorithmResource(completeForm, algorithmFile.value, docFile.value)
    ElMessage.success('算法资源上传成功！')
    
    // 清除表单数据
    clearFormState()
    
    // 发送成功事件
    emit('success')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error(error.message || '上传失败，请重试')
  } finally {
    loading.close()
  }
}

// 添加标签
const addTag = () => {
  const tagValue = tagInput.value.trim()
  if (tagValue && !algorithmForm.tags.includes(tagValue) && algorithmForm.tags.length < 5) {
    algorithmForm.tags.push(tagValue)
  }
  tagInput.value = ''
}

// 删除标签
const removeTag = (tag) => {
  algorithmForm.tags.splice(algorithmForm.tags.indexOf(tag), 1)
}

// 自动保存
let autoSaveInterval
onMounted(() => {
  loadFormState()
  autoSaveInterval = setInterval(() => {
    if (props.active) {
      if (currentStep.value === 1) {
        saveFormState()
      } else if (currentStep.value === 2) {
        saveUsageGuide()
      } else if (currentStep.value === 3) {
        saveApiDoc()
      }
    }
  }, 30000) // 每30秒保存一次
})

onBeforeUnmount(() => {
  clearInterval(autoSaveInterval)
  // 最后保存一次，确保数据不丢失
  if (props.active) {
    if (currentStep.value === 1) {
      saveFormState()
    } else if (currentStep.value === 2) {
      saveUsageGuide()
    } else if (currentStep.value === 3) {
      saveApiDoc()
    }
  }
})
</script>

<template>
  <div class="algorithm-resource-form">
    <el-steps :active="currentStep" finish-status="success" simple style="margin-bottom: 20px;">
      <el-step title="基本信息" :icon="Document" />
      <el-step title="使用指南" :icon="SetUp" />
      <el-step title="API文档" :icon="Edit" />
    </el-steps>

    <!-- 步骤1：基本信息 -->
    <el-card v-if="currentStep === 1" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <h3>算法基本信息</h3>
        </div>
      </template>

      <el-form :model="algorithmForm" label-position="top" require-asterisk-position="right" class="algorithm-form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="算法名称" required>
              <el-input v-model="algorithmForm.name" placeholder="为您的算法资源取一个明确的名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="算法版本" required>
              <el-input v-model="algorithmForm.version" placeholder="如：1.0.0" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="算法分类" required>
          <el-select v-model="algorithmForm.category" placeholder="选择一个最合适的分类" style="width: 100%">
            <el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="标签 (最多5个)">
          <div class="tags-input-container">
            <el-tag v-for="tag in algorithmForm.tags" :key="tag" closable @close="removeTag(tag)" class="form-tag">
              {{ tag }}
            </el-tag>
            <el-input v-if="algorithmForm.tags.length < 5" v-model="tagInput" placeholder="输入标签后按回车" @keyup.enter="addTag" size="small" class="tag-input-field" />
          </div>
          <div class="el-form-item__info">添加助于其他用户了解您的算法资源的关键字</div>
        </el-form-item>

        <el-form-item label="依赖项">
          <el-input v-model="algorithmForm.dependencies" type="textarea" :rows="3" placeholder="列出算法运行所需的依赖项，如：numpy==1.19.5, pandas>=1.0.0, tensorflow>=2.4.0" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="算法文件" required>
              <el-upload ref="algorithmFileEl" action="#" :auto-upload="false" :show-file-list="false" :on-change="handleAlgorithmFileChange" :limit="1" class="full-width-upload">
                <el-button type="primary" :icon="UploadFilled">选择算法文件</el-button>
                <div class="el-upload__tip">支持 .py, .jar, .zip, .tar.gz 等文件类型，最大 200MB</div>
              </el-upload>
              <div v-if="algorithmFile" class="file-preview">
                <el-icon><SetUp /></el-icon>
                <span class="file-name">{{ algorithmFile.name }}</span>
                <span class="file-size">({{ (algorithmFile.size / 1024 / 1024).toFixed(2) }} MB)</span>
                <el-button type="danger" :icon="Delete" text circle @click="clearAlgorithmFile" />
              </div>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="文档文件（可选）">
              <el-upload ref="docFileEl" action="#" :auto-upload="false" :show-file-list="false" :on-change="handleDocFileChange" :limit="1" class="full-width-upload">
                <el-button :icon="UploadFilled">选择文档文件</el-button>
                <div class="el-upload__tip">支持 .md, .pdf, .docx 等文档类型，最大 50MB</div>
              </el-upload>
              <div v-if="docFile" class="file-preview">
                <el-icon><Edit /></el-icon>
                <span class="file-name">{{ docFile.name }}</span>
                <span class="file-size">({{ (docFile.size / 1024 / 1024).toFixed(2) }} MB)</span>
                <el-button type="danger" :icon="Delete" text circle @click="clearDocFile" />
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="算法描述">
          <el-input v-model="algorithmForm.description" type="textarea" :rows="4" placeholder="描述您的算法的功能、特点、适用场景和性能等" />
        </el-form-item>

        <div class="form-actions">
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="nextStep">
            下一步：填写使用指南 <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </div>
      </el-form>
    </el-card>

    <!-- 步骤2：使用指南 -->
    <el-card v-if="currentStep === 2" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <h3>使用指南 (Markdown)</h3>
          <el-button-group>
            <el-button :type="usageGuidePreviewMode ? '' : 'primary'" @click="usageGuidePreviewMode = false" :icon="Edit">编辑</el-button>
            <el-button :type="usageGuidePreviewMode ? 'primary' : ''" @click="usageGuidePreviewMode = true" :icon="View">预览</el-button>
          </el-button-group>
        </div>
      </template>

      <div class="markdown-editor">
        <textarea v-if="!usageGuidePreviewMode" v-model="usageGuideContent" placeholder="请使用Markdown格式编写算法使用指南，包括环境需求、输入输出参数、使用方法等..."></textarea>
        <div v-else class="markdown-preview styled-markdown" v-html="compiledUsageGuide"></div>
      </div>

      <div class="form-actions">
        <el-button @click="prevStep">返回上一步</el-button>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="nextStep">
          下一步：填写API文档 <el-icon class="el-icon--right"><ArrowRight /></el-icon>
        </el-button>
      </div>
    </el-card>

    <!-- 步骤3：API文档 -->
    <el-card v-if="currentStep === 3" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <h3>API文档 (Markdown)</h3>
          <el-button-group>
            <el-button :type="apiDocPreviewMode ? '' : 'primary'" @click="apiDocPreviewMode = false" :icon="Edit">编辑</el-button>
            <el-button :type="apiDocPreviewMode ? 'primary' : ''" @click="apiDocPreviewMode = true" :icon="View">预览</el-button>
          </el-button-group>
        </div>
      </template>

      <div class="markdown-editor">
        <textarea v-if="!apiDocPreviewMode" v-model="apiDocContent" placeholder="请使用Markdown格式编写API文档，包括函数接口、参数说明、返回值、使用示例等..."></textarea>
        <div v-else class="markdown-preview styled-markdown" v-html="compiledApiDoc"></div>
      </div>

      <div class="form-actions">
        <el-button @click="prevStep">返回上一步</el-button>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="success" @click="handleSubmit">完成并提交</el-button>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.algorithm-resource-form {
  margin-bottom: 20px;
}

.form-card {
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 25px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.form-actions {
  margin-top: 25px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.algorithm-form .el-form-item__label {
  font-weight: 500;
}

.full-width-upload .el-upload {
  width: 100%;
}

.full-width-upload .el-upload-dragger {
  width: 100%;
}

.file-preview {
  margin-top: 8px;
  padding: 8px 10px;
  background-color: #f9f9f9;
  border: 1px solid #eee;
  border-radius: 4px;
  display: flex;
  align-items: center;
  font-size: 13px;
  color: #555;
}

.file-preview .el-icon {
  margin-right: 6px;
  color: #409EFF;
}

.file-preview .file-name {
  flex-grow: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 8px;
}

.file-preview .file-size {
  color: #888;
  margin-right: 8px;
}

/* Markdown编辑器样式 */
.markdown-editor {
  height: 450px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
}

.markdown-editor textarea {
  width: 100%;
  height: 100%;
  padding: 15px;
  border: none;
  resize: none;
  outline: none;
  font-family: 'Menlo', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.6;
  background-color: #fdfdfd;
}

.markdown-preview {
  width: 100%;
  height: 100%;
  padding: 15px;
  overflow-y: auto;
  background-color: #fff;
}

/* GitHub Markdown 基础样式 */
.styled-markdown h1, .styled-markdown h2, .styled-markdown h3 {
  border-bottom: 1px solid #eaecef;
  padding-bottom: .3em;
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
}
.styled-markdown h1 { font-size: 2em; }
.styled-markdown h2 { font-size: 1.5em; }
.styled-markdown h3 { font-size: 1.25em; }
.styled-markdown p { margin-top: 0; margin-bottom: 16px; }
.styled-markdown ul, .styled-markdown ol {
  margin-top: 0;
  margin-bottom: 16px;
  padding-left: 2em;
}
.styled-markdown blockquote {
  margin: 0 0 16px 0;
  padding: 0 1em;
  color: #6a737d;
  border-left: .25em solid #dfe2e5;
}
.styled-markdown pre {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
  margin-bottom: 16px;
}
.styled-markdown code:not(pre > code) {
  padding: .2em .4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27,31,35,.05);
  border-radius: 3px;
}
.styled-markdown table {
  border-collapse: collapse;
  margin-bottom: 16px;
  width: 100%;
}
.styled-markdown th, .styled-markdown td {
  border: 1px solid #dfe2e5;
  padding: 6px 13px;
}

.tags-input-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.form-tag {
  font-size: 13px;
}

.tag-input-field {
  width: 180px;
}
</style> 