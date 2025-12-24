<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { Document, Upload, View, Edit, Delete, Plus, Connection, SetUp } from '@element-plus/icons-vue' // 引入所需图标
import { marked } from 'marked'
import SparkMD5 from 'spark-md5' // <-- 添加 SparkMD5 导入
import { uploadDataResource, getMockCategories } from '../api/resources'
import ComputingResourceForm from '../components/ComputingResourceForm.vue'
import AlgorithmResourceForm from '../components/AlgorithmResourceForm.vue'

const router = useRouter()

// 资源类型定义 (复用您原有的)
const resourceTypes = [
  { 
    value: 'data', 
    label: '数据资源', 
    icon: Document, // 直接使用导入的图标组件
    description: '上传数据集、文件等数据资源，支持多种文件格式'
  },
  { 
    value: 'computing', 
    label: '算力资源', 
    icon: Connection,
    description: '注册您的算力服务器，提供远程计算服务'
  },
  { 
    value: 'algorithm', 
    label: '算法资源', 
    icon: SetUp,
    description: '上传算法、模型或插件包，为数据处理提供算法支持'
  }
]

const activeTab = ref('data') // 默认激活数据资源上传

// 表单步骤
const currentStep = ref(1)

// 本地存储键名
const DATA_FORM_STORAGE_KEY = 'dataUploadForm'
const DATA_DOCS_STORAGE_KEY = 'dataUploadDocs'

// 数据资源表单
const dataForm = reactive({
  name: '',
  category: '',
  tags: [],
  fileExtensions: [],
  dataStructure: '',
  isConfidential: false,
  needHosting: false,
  description: '', // 简短描述
  documentation: '' // 详细说明文档（Markdown格式）
})

const dataFile = ref(null)
const sampleFile = ref(null)
const dataFileEl = ref(null) // ElUpload 组件的引用
const sampleFileEl = ref(null) // ElUpload 组件的引用
const dataFileHash = ref('') // <-- 新增: 存储数据文件的MD5哈希值

const categoryOptions = ref([])
const tagInput = ref('')
const fileExtInput = ref('')

const loadCategories = async () => {
  try {
    const res = await getMockCategories() //  先使用mock数据，后续可替换为真实API
    categoryOptions.value = res.data
  } catch (error) {
    console.error('加载分类失败:', error)
    ElMessage.error('加载分类数据失败')
  }
}

const addTag = () => {
  const tagValue = tagInput.value.trim()
  if (tagValue && !dataForm.tags.includes(tagValue) && dataForm.tags.length < 5) {
    dataForm.tags.push(tagValue)
  }
  tagInput.value = ''
}

const removeTag = (tag) => {
  dataForm.tags.splice(dataForm.tags.indexOf(tag), 1)
}

const addFileExtension = () => {
  const ext = fileExtInput.value.trim().toLowerCase()
  if (ext && !dataForm.fileExtensions.includes(ext)) {
    dataForm.fileExtensions.push(ext)
  }
  fileExtInput.value = ''
}

const removeFileExtension = (ext) => {
  dataForm.fileExtensions.splice(dataForm.fileExtensions.indexOf(ext), 1)
}

// <-- 新增: MD5 计算辅助函数 -->
const calculateFileMD5 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    const spark = new SparkMD5.ArrayBuffer();
    const chunkSize = 2097152; // 2MB chunks
    let currentChunk = 0;
    const chunks = Math.ceil(file.size / chunkSize);

    reader.onload = (e) => {
      spark.append(e.target.result);
      currentChunk++;
      if (currentChunk < chunks) {
        loadNextChunk();
      } else {
        resolve(spark.end());
      }
    };
    
    reader.onerror = (e) => {
      reject(e);
    };
    
    function loadNextChunk() {
      const start = currentChunk * chunkSize;
      const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;
      reader.readAsArrayBuffer(file.slice(start, end));
    }

    loadNextChunk();
  });
};

const handleDataFileChange = async (fileInfo) => {
  if (!fileInfo || !fileInfo.raw) {
    dataFile.value = null;
    dataFileHash.value = '';
    // 确保 dataFileEl.value 存在且 clearFiles 是一个函数
    if (dataFileEl.value && typeof dataFileEl.value.clearFiles === 'function') {
        dataFileEl.value.clearFiles();
    }
    return;
  }
  const file = fileInfo.raw;

  const maxSize = 500 * 1024 * 1024 // 500MB
  if (file.size > maxSize) {
    ElMessage.warning('数据文件大小超过500MB')
    if (dataFileEl.value && typeof dataFileEl.value.clearFiles === 'function') {
        dataFileEl.value.clearFiles();
    }
    dataFile.value = null
    dataFileHash.value = '' // 清除哈希
    return
  }
  dataFile.value = file

  // 计算 MD5
  const loadingMD5 = ElLoading.service({ text: '正在计算文件校验码...', background: 'rgba(0, 0, 0, 0.7)', lock: true });
  try {
    dataFileHash.value = await calculateFileMD5(file);
    // ElMessage.success(`文件校验码 (MD5): ${dataFileHash.value}`); // 可选：显示哈希给用户
  } catch (error) {
    console.error('MD5 calculation error:', error);
    ElMessage.error('计算文件校验码失败');
    dataFileHash.value = ''; // 发生错误时清除哈希
  } finally {
    loadingMD5.close();
  }
}

const handleSampleFileChange = (fileInfo) => { // <-- 修改参数为 fileInfo
  if (!fileInfo || !fileInfo.raw) { // <-- 检查 fileInfo 和 fileInfo.raw
    sampleFile.value = null;
    if (sampleFileEl.value && typeof sampleFileEl.value.clearFiles === 'function') {
        sampleFileEl.value.clearFiles();
    }
    return;
  }
  const file = fileInfo.raw; // <-- 获取实际 File 对象

  const maxSize = 50 * 1024 * 1024 // 50MB
  if (file.size > maxSize) {
    ElMessage.warning('样例文件大小超过50MB')
    if (sampleFileEl.value && typeof sampleFileEl.value.clearFiles === 'function') {
        sampleFileEl.value.clearFiles()
    }
    sampleFile.value = null
    return
  }
  sampleFile.value = file
}

const clearDataFile = () => {
  dataFile.value = null
  dataFileHash.value = '' // <-- 同时清除哈希值
  if (dataFileEl.value && typeof dataFileEl.value.clearFiles === 'function') {
    dataFileEl.value.clearFiles()
  }
}

const clearSampleFile = () => {
  sampleFile.value = null
  if (sampleFileEl.value && typeof sampleFileEl.value.clearFiles === 'function') {
    sampleFileEl.value.clearFiles()
  }
}

const validateStep1 = () => {
  if (!dataForm.name) {
    ElMessage.warning('请输入数据名称'); return false
  }
  if (!dataForm.category) {
    ElMessage.warning('请选择数据分类'); return false
  }
  if (!dataFile.value) {
    ElMessage.warning('请上传数据文件'); return false
  }
  if (!dataFileHash.value && dataFile.value) { // <-- 如果有文件但没有哈希（例如计算失败或未完成）
    ElMessage.warning('文件校验码正在计算或计算失败，请稍候或重新选择文件'); return false;
  }
  return true
}

const goToStep2 = () => {
  if (validateStep1()) {
    saveDataFormState() // 保存第一步的状态
    currentStep.value = 2
    // 加载或初始化Markdown内容
    const savedDocs = localStorage.getItem(DATA_DOCS_STORAGE_KEY)
    if (savedDocs) {
      markdownContent.value = savedDocs
    } else {
      markdownContent.value = `# ${dataForm.name || '数据资源'} 说明文档\n\n## 数据概述\n请在此处描述数据的基本情况...\n\n## 数据来源\n请描述数据的来源...\n\n## 数据结构\n请描述数据的结构和字段...\n\n## 数据样例\n请提供数据样例（可选）...\n\n## 使用方法\n请描述数据的使用方法...\n\n## 注意事项\n请描述使用数据时需要注意的问题...\n`
    }
  }
}

const goToStep1 = () => {
  saveDataDocsState() // 从第二步返回第一步时，保存Markdown内容
  currentStep.value = 1
}

// Markdown 编辑器相关
const markdownContent = ref('')
const previewMode = ref(false)
const compiledMarkdown = computed(() => marked(markdownContent.value || ''))

const saveDataFormState = () => {
  const stateToSave = {
    name: dataForm.name,
    category: dataForm.category,
    tags: dataForm.tags,
    fileExtensions: dataForm.fileExtensions,
    dataStructure: dataForm.dataStructure,
    isConfidential: dataForm.isConfidential,
    needHosting: dataForm.needHosting,
  }
  localStorage.setItem(DATA_FORM_STORAGE_KEY, JSON.stringify(stateToSave))
}

const loadDataFormState = () => {
  const savedState = localStorage.getItem(DATA_FORM_STORAGE_KEY)
  if (savedState) {
    const parsedState = JSON.parse(savedState)
    Object.assign(dataForm, parsedState)
  }
}

const saveDataDocsState = () => {
  localStorage.setItem(DATA_DOCS_STORAGE_KEY, markdownContent.value)
}

let autoSaveInterval
onMounted(() => {
  loadCategories()
  loadDataFormState() // 加载表单的第一部分数据
  // 如果之前停留在第二步，直接加载文档内容
  const savedDocs = localStorage.getItem(DATA_DOCS_STORAGE_KEY)
  if (currentStep.value === 2 && savedDocs) { // 假设currentStep也持久化了，或者根据其他逻辑判断
    markdownContent.value = savedDocs
  }
  autoSaveInterval = setInterval(() => {
    if (currentStep.value === 1) {
      saveDataFormState()
    } else if (currentStep.value === 2) {
      saveDataDocsState()
    }
  }, 30000) // 每30秒保存一次
})

onBeforeUnmount(() => {
  clearInterval(autoSaveInterval)
  // 最后保存一次，确保数据不丢失
  if (currentStep.value === 1) {
    saveDataFormState()
  } else if (currentStep.value === 2) {
    saveDataDocsState()
  }
})

const handleSubmitAll = async () => {
  if (!markdownContent.value.trim()) {
    ElMessage.warning('请编写说明文档内容')
    return
  }
  dataForm.documentation = markdownContent.value // 将Markdown内容赋给documentation字段

  if (!dataFileHash.value && dataFile.value) {
      ElMessage.warning('数据文件校验码未生成，无法提交。请重新选择文件或等待计算完成。');
      return;
  }

  const loading = ElLoading.service({
    lock: true,
    text: '正在上传数据资源...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  try {
    await uploadDataResource(dataForm, dataFile.value, sampleFile.value, dataFileHash.value)
    ElMessage.success('数据资源上传成功！')
    localStorage.removeItem(DATA_FORM_STORAGE_KEY)
    localStorage.removeItem(DATA_DOCS_STORAGE_KEY)
    
    // 清理操作：先清理文件和相关状态，再清理表单，最后切换步骤
    clearDataFile() // 这会清除 dataFile 和 dataFileHash
    clearSampleFile()
    
    Object.assign(dataForm, { name: '', category: '', tags: [], fileExtensions: [], dataStructure: '', isConfidential: false, needHosting: false, description: '', documentation: '' })
    markdownContent.value = ''

    currentStep.value = 1 // 最后改变步骤
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error(error.message || '上传失败，请重试')
  } finally {
    loading.close()
  }
}

const handleCancel = () => {
  ElMessageBox.confirm(
    '确定要取消上传吗？已填写的内容将会保存（如果未提交）。',
    '取消上传',
    { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
  ).then(() => {
    // 在 onBeforeUnmount 中已经处理了保存逻辑
    router.push('/') // 或其他目标页面
  }).catch(() => { /* 用户取消对话框 */ })
}

// 算力资源和算法资源表单 (使用组件化的方式)
const handleComputingSuccess = () => {
  ElMessage.success('算力资源上传成功！')
}

const handleAlgorithmSuccess = () => {
  ElMessage.success('算法资源上传成功！')
}

</script>

<template>
  <div class="upload-container">
    <div class="upload-header">
      <h1>资源上传</h1>
      <p>上传您的数据资源、算力资源或算法资源，让它们在平台创造更大价值</p>
    </div>

    <el-tabs v-model="activeTab" type="card" class="resource-type-selector">
        <el-tab-pane v-for="type in resourceTypes" :key="type.value" :name="type.value">
          <template #label>
            <div class="tab-label">
              <el-icon><component :is="type.icon" /></el-icon>
              <span>{{ type.label }}</span>
            </div>
          </template>
          
        <div class="tab-content-wrapper">
          <el-alert :title="type.description" type="info" :closable="false" show-icon class="tab-description" />

          <!-- 数据资源上传表单 -->
          <div v-if="activeTab === 'data'">
            <el-steps :active="currentStep" finish-status="success" simple style="margin-bottom: 20px;">
              <el-step title="基本信息" :icon="Document" />
              <el-step title="详细说明" :icon="Edit" />
            </el-steps>

            <!-- 步骤一：基本信息 -->
            <el-card v-if="currentStep === 1" shadow="never" class="form-card">
              <el-form :model="dataForm" label-position="top" require-asterisk-position="right">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="数据名称" required>
                      <el-input v-model="dataForm.name" placeholder="为您的数据资源取一个明确的名称" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="数据分类" required>
                      <el-select v-model="dataForm.category" placeholder="选择一个最合适的分类" style="width: 100%;">
                        <el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" />
                      </el-select>
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="数据结构简述">
                  <el-input v-model="dataForm.dataStructure" type="textarea" :rows="2" placeholder="例如：CSV文件，包含id, name, age字段；或 JSON对象数组" />
                </el-form-item>

                <el-form-item label="标签 (最多5个)">
                  <div class="tags-input-container">
                    <el-tag v-for="tag in dataForm.tags" :key="tag" closable @close="removeTag(tag)" class="form-tag">
                      {{ tag }}
                    </el-tag>
                    <el-input v-if="dataForm.tags.length < 5" v-model="tagInput" placeholder="输入标签后按回车" @keyup.enter="addTag" size="small" class="tag-input-field" />
                  </div>
                </el-form-item>

                <el-form-item label="文件类型标签">
                  <div class="tags-input-container">
                    <el-tag v-for="ext in dataForm.fileExtensions" :key="ext" closable @close="removeFileExtension(ext)" type="success" class="form-tag">
                      {{ ext }}
                    </el-tag>
                    <el-input v-model="fileExtInput" placeholder="如 csv, xlsx, zip, png" @keyup.enter="addFileExtension" size="small" class="tag-input-field" />
                  </div>
                  <div class="el-form-item__info">告知用户文件或压缩包内文件的后缀名，例如：csv、xlsx、zip、png等。</div>
                </el-form-item>

                <el-row :gutter="20">
                  <el-col :xs="24" :sm="12">
                    <el-form-item label="数据主文件" required>
                      <el-upload ref="dataFileEl" action="#" :auto-upload="false" :show-file-list="false" :on-change="handleDataFileChange" :limit="1" class="full-width-upload">
                        <el-button type="primary" :icon="Upload">选择数据文件</el-button>
                      </el-upload>
                      <div v-if="dataFile" class="file-preview">
                        <el-icon><Document /></el-icon>
                        <span class="file-name">{{ dataFile.name }}</span>
                        <span class="file-size">({{ (dataFile.size / 1024 / 1024).toFixed(2) }} MB)</span>
                        <el-button type="danger" :icon="Delete" text circle @click="clearDataFile" />
                      </div>
                    </el-form-item>
                  </el-col>
                  <el-col :xs="24" :sm="12">
                    <el-form-item label="样例数据文件 (可选)">
                      <el-upload ref="sampleFileEl" action="#" :auto-upload="false" :show-file-list="false" :on-change="handleSampleFileChange" :limit="1" class="full-width-upload">
                        <el-button :icon="Upload">选择样例文件</el-button>
                      </el-upload>
                      <div v-if="sampleFile" class="file-preview">
                        <el-icon><Document /></el-icon>
                        <span class="file-name">{{ sampleFile.name }}</span>
                        <span class="file-size">({{ (sampleFile.size / 1024 / 1024).toFixed(2) }} MB)</span>
                        <el-button type="danger" :icon="Delete" text circle @click="clearSampleFile" />
                      </div>
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="额外选项" style="margin-bottom: 0;">
                  <el-checkbox v-model="dataForm.isConfidential" label="设为保密数据" />
                  <el-checkbox v-model="dataForm.needHosting" label="需要平台托管存储" />
                </el-form-item>

                <el-form-item label="数据描述">
                  <el-input v-model="dataForm.description" type="textarea" :rows="3" placeholder="请输入简短的数据描述（不使用Markdown格式）" />
                </el-form-item>

                <el-form-item class="form-actions-step1">
                  <el-button @click="handleCancel">取消</el-button>
                  <el-button type="primary" @click="goToStep2">添加说明文档</el-button>
                </el-form-item>
              </el-form>
            </el-card>

            <!-- 步骤二：详细说明 (Markdown) -->
            <el-card v-if="currentStep === 2" shadow="never" class="form-card">
              <div class="markdown-editor-header">
                 <h3>详细说明文档 (Markdown)</h3>
                <el-button-group>
                  <el-button :type="previewMode ? '' : 'primary'" @click="previewMode = false" :icon="Edit">编辑</el-button>
                  <el-button :type="previewMode ? 'primary' : ''" @click="previewMode = true" :icon="View">预览</el-button>
                </el-button-group>
                </div>
              <div class="markdown-editor">
                <textarea v-if="!previewMode" v-model="markdownContent" placeholder="请使用Markdown格式编写详细的数据说明、来源、结构、使用方法等..."></textarea>
                <div v-else class="markdown-preview styled-markdown" v-html="compiledMarkdown"></div>
                </div>
              <el-form-item class="form-actions-step2">
                <el-button @click="goToStep1">返回上一步</el-button>
                <el-button @click="handleCancel">取消</el-button>
                <el-button type="success" @click="handleSubmitAll">完成并提交</el-button>
              </el-form-item>
            </el-card>
          </div>

          <!-- 算力资源表单 -->
          <div v-if="activeTab === 'computing'">
            <ComputingResourceForm 
              :active="activeTab === 'computing'"
              @cancel="activeTab = 'data'"
              @success="handleComputingSuccess"
            />
          </div>
          
          <!-- 算法资源表单 -->
          <div v-if="activeTab === 'algorithm'">
            <AlgorithmResourceForm
              :active="activeTab === 'algorithm'"
              @cancel="activeTab = 'data'"
              @success="handleAlgorithmSuccess"
            />
                    </div>

          </div>
        </el-tab-pane>
      </el-tabs>
  </div>
</template>

<style scoped>
.upload-container {
  max-width: 900px; /* 稍微调整宽度 */
  margin: 20px auto;
  padding: 20px;
  background-color: #fff; /* 背景设为白色 */
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.upload-header {
  text-align: center;
  margin-bottom: 25px;
}

.upload-header h1 {
  font-size: 26px; /* 调整标题大小 */
  color: #2c3e50; /* 深蓝灰色 */
  font-weight: 600;
}

.upload-header p {
  font-size: 15px;
  color: #5f6368; /* 中灰色 */
  margin-top: 8px;
}

.resource-type-selector {
  margin-bottom: 25px;
}

.tab-label {
  display: flex;
  align-items: center;
  font-size: 15px;
  padding: 0 10px; /* 增加标签页内边距 */
}

.tab-label .el-icon {
  margin-right: 8px;
  font-size: 18px; /* 图标稍大 */
}

.tab-content-wrapper {
    padding: 5px 15px 15px 15px; /* 给tab内容区域增加内边距 */
}

.tab-description {
  margin-bottom: 20px;
  border-radius: 6px;
}

.form-card {
  border: 1px solid #e0e0e0; /* 卡片边框颜色 */
  border-radius: 6px;
  padding: 25px;
}

.el-form-item__label {
  font-weight: 500; /* 标签字体加粗 */
  color: #333;
}

.el-form-item__info {
    font-size: 12px;
    color: #888;
    line-height: 1.4;
    margin-top: 4px;
}

.tags-input-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.form-tag {
  font-size: 13px;
}

.tag-input-field {
  width: 180px;
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

.form-actions-step1,
.form-actions-step2 {
  margin-top: 25px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.markdown-editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.markdown-editor-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.markdown-editor {
  height: 450px; /* 编辑器高度 */
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden; /* 确保子元素在圆角内 */
  display: flex; /* 使用flex布局 */
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

/* GitHub Markdown 基础样式 (styled-markdown) */
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

</style>