<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { Connection, Setting, ArrowRight, Check, Close, Warning } from '@element-plus/icons-vue'
import { uploadComputingResource } from '../api/resources'

const props = defineProps({
  active: Boolean
})

const emit = defineEmits(['cancel', 'success'])

// 步骤
const currentStep = ref(1)
const totalSteps = 3

// 本地存储键名
const COMPUTING_FORM_STORAGE_KEY = 'computingUploadForm'

// 算力资源表单
const computingForm = reactive({
  name: '',
  serverAddress: '',
  port: '',
  cpu: '',
  memory: '',
  gpu: '',
  storage: '',
  availableTimeStart: '',
  availableTimeEnd: '',
  systemType: '',
  description: '',
  tags: []
})

const systemOptions = [
  { value: 'linux', label: 'Linux' },
  { value: 'windows', label: 'Windows' },
  { value: 'macos', label: 'MacOS' }
]

// 测试步骤状态
const testSteps = reactive([
  { name: '服务器连接', status: 'waiting', message: '等待测试' },
  { name: '沙箱环境', status: 'waiting', message: '等待测试' },
  { name: '数据传输', status: 'waiting', message: '等待测试' },
  { name: '插件安装', status: 'waiting', message: '等待测试' },
  { name: '数据处理', status: 'waiting', message: '等待测试' },
  { name: '环境关闭', status: 'waiting', message: '等待测试' }
])

const sandboxVersion = ref('')
const testingInProgress = ref(false)
const testResult = ref({
  success: false,
  message: ''
})

const tagInput = ref('')

// 保存表单状态到本地存储
const saveFormState = () => {
  localStorage.setItem(COMPUTING_FORM_STORAGE_KEY, JSON.stringify(computingForm))
}

// 从本地存储加载表单状态
const loadFormState = () => {
  const savedState = localStorage.getItem(COMPUTING_FORM_STORAGE_KEY)
  if (savedState) {
    try {
      const parsedState = JSON.parse(savedState)
      Object.assign(computingForm, parsedState)
    } catch (e) {
      console.error('解析保存的表单数据失败:', e)
    }
  }
}

// 清除表单数据
const clearFormState = () => {
  Object.assign(computingForm, {
    name: '',
    serverAddress: '',
    port: '',
    cpu: '',
    memory: '',
    gpu: '',
    storage: '',
    availableTimeStart: '',
    availableTimeEnd: '',
    systemType: '',
    description: '',
    tags: []
  })
  localStorage.removeItem(COMPUTING_FORM_STORAGE_KEY)
}

// 验证步骤1
const validateStep1 = () => {
  // 没有实际验证，因为这只是阅读沙箱环境安装教程
  return true
}

// 验证步骤2
const validateStep2 = () => {
  if (!computingForm.name) {
    ElMessage.warning('请输入服务名称')
    return false
  }
  if (!computingForm.serverAddress) {
    ElMessage.warning('请输入服务器地址')
    return false
  }
  if (!computingForm.port) {
    ElMessage.warning('请输入开放端口')
    return false
  }
  // 其他字段不是必需的，可以选填

  return true
}

// 进入下一步
const nextStep = () => {
  if (currentStep.value === 1 && validateStep1()) {
    currentStep.value = 2
    saveFormState()
  } else if (currentStep.value === 2 && validateStep2()) {
    currentStep.value = 3
    saveFormState()
    startConnectionTest()
  }
}

// 返回上一步
const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
    saveFormState()
  }
}

// 取消
const handleCancel = () => {
  ElMessageBox.confirm(
    '确定要取消上传吗？已填写的内容将会保存（如果未提交）。',
    '取消上传',
    { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
  ).then(() => {
    saveFormState() // 保存表单状态
    emit('cancel')
  }).catch(() => {
    // 用户取消对话框，不做任何处理
  })
}

// 执行连接测试（模拟）
const startConnectionTest = async () => {
  testingInProgress.value = true
  testResult.value = { success: false, message: '' }
  
  // 重置所有测试步骤状态
  testSteps.forEach(step => {
    step.status = 'process'
    step.message = '测试中...'
  })

  // 模拟测试过程，每个步骤延时0.5秒
  for (let i = 0; i < testSteps.length; i++) {
    await new Promise(resolve => setTimeout(resolve, 500))
    testSteps[i].status = 'success'
    testSteps[i].message = '测试通过'
  }

  // 测试完成后，随机生成一个沙箱版本号（实际中应该从服务器返回）
  sandboxVersion.value = `v${Math.floor(Math.random() * 2) + 1}.${Math.floor(Math.random() * 10)}.${Math.floor(Math.random() * 10)}`
  
  testingInProgress.value = false
  testResult.value = {
    success: true,
    message: `连接测试成功！检测到沙箱环境版本: ${sandboxVersion.value}`
  }
}

// 添加标签
const addTag = () => {
  const tagValue = tagInput.value.trim()
  if (tagValue && !computingForm.tags.includes(tagValue) && computingForm.tags.length < 5) {
    computingForm.tags.push(tagValue)
  }
  tagInput.value = ''
}

// 删除标签
const removeTag = (tag) => {
  computingForm.tags.splice(computingForm.tags.indexOf(tag), 1)
}

// 提交表单
const handleSubmit = async () => {
  if (!testResult.value.success) {
    ElMessage.warning('请先完成连接测试')
    return
  }

  const loading = ElLoading.service({
    lock: true,
    text: '正在提交算力资源...',
    background: 'rgba(0, 0, 0, 0.7)'
  })

  try {
    // 在表单中添加沙箱版本信息
    const formData = { 
      ...computingForm,
      sandboxVersion: sandboxVersion.value
    }
    
    await uploadComputingResource(formData)
    ElMessage.success('算力资源提交成功！')
    
    // 清除表单数据
    clearFormState()
    
    // 重置步骤和测试状态
    currentStep.value = 1
    testSteps.forEach(step => {
      step.status = 'waiting'
      step.message = '等待测试'
    })
    sandboxVersion.value = ''
    testResult.value = { success: false, message: '' }
    
    // 发送成功事件
    emit('success')
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error(error.message || '提交失败，请重试')
  } finally {
    loading.close()
  }
}

// 自动保存
let autoSaveInterval
onMounted(() => {
  loadFormState()
  autoSaveInterval = setInterval(() => {
    if (props.active) {
      saveFormState()
    }
  }, 30000) // 每30秒保存一次
})

onBeforeUnmount(() => {
  clearInterval(autoSaveInterval)
  // 最后保存一次，确保数据不丢失
  if (props.active) {
    saveFormState()
  }
})
</script>

<template>
  <div class="computing-resource-form">
    <el-steps :active="currentStep" finish-status="success" simple style="margin-bottom: 20px;">
      <el-step title="沙箱安装" :icon="Connection" />
      <el-step title="基本信息" :icon="Setting" />
      <el-step title="连接测试" :icon="Check" />
    </el-steps>

    <!-- 步骤1：沙箱安装教程 -->
    <el-card v-if="currentStep === 1" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <h3>沙箱环境安装配置</h3>
        </div>
      </template>      <div class="tutorial-content">
        <h4>沙箱环境安装与连接检测指南</h4>
        
        <p>本指南帮助您在服务器上安装、配置并验证沙箱环境，以便平台能够安全、稳定地为用户提供交互式计算服务。</p>
        
        <el-alert
          title="注意：您的服务器必须安装沙箱环境才能在平台上提供算力资源。"
          type="warning"
          :closable="false"
          show-icon
          style="margin: 15px 0;"
        />
        
        <h4>步骤 1：准备环境</h4>
        
        <el-collapse accordion>
          <el-collapse-item title="1. 确认系统要求" name="1">
            <ul>
              <li>操作系统：Ubuntu 18.04 及以上 / CentOS 7 及以上</li>
              <li>Python：3.7 及以上</li>
              <li>Node.js：10 及以上（可选，用于单点登录插件）</li>
              <li>建议内存 ≥ 4 GB，磁盘 ≥ 20 GB</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="2. 创建系统用户" name="2">
            <ul>
              <li>建议为 JupyterHub 创建独立用户 <code>jupyter</code>，用于运行服务并隔离权限</li>
              <li>分配合适的用户组与文件夹权限</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="3. 打开必要端口" name="3">
            <ul>
              <li>默认 HTTP：8000</li>
              <li>HTTPS：8443（若启用 TLS/SSL）</li>
              <li>SSH：22（用于远程管理）</li>
            </ul>
          </el-collapse-item>
        </el-collapse>
        
        <h4>步骤 2：安装 JupyterHub 与依赖</h4>
        
        <el-collapse accordion>
          <el-collapse-item title="1. 安装系统包管理器更新" name="4">
            <ul>
              <li>确保 apt/yum 数据库已更新</li>
              <li>安装 Node.js、npm（若需 OAuth 或其他扩展）</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="2. 安装 Python 环境与虚拟环境工具" name="5">
            <ul>
              <li>安装 <code>python3-venv</code> 或 <code>virtualenv</code></li>
              <li>在 <code>jupyter</code> 用户下创建独立虚拟环境</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="3. 安装 JupyterHub 和 Notebook" name="6">
            <ul>
              <li>使用 <code>pip</code> 安装：JupyterHub、jupyterlab、notebook、traitlets 等核心组件</li>
              <li>（可选）安装 <code>dockerspawner</code>、<code>oauthenticator</code> 等扩展</li>
            </ul>
          </el-collapse-item>
        </el-collapse>
        
        <h4>步骤 3：配置 JupyterHub</h4>
        
        <el-collapse accordion>
          <el-collapse-item title="1. 初始化配置文件" name="7">
            <ul>
              <li>运行 <code>jupyterhub --generate-config</code> 生成 <code>jupyterhub_config.py</code></li>
              <li>将配置文件放置于 <code>/etc/jupyterhub/</code></li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="2. 基本配置项调整" name="8">
            <ul>
              <li><code>c.Spawner.default_url</code>：设置默认启动页面（如 <code>/lab</code>）</li>
              <li><code>c.JupyterHub.bind_url</code>：绑定地址与端口（如 <code>http://:8000</code>）</li>
              <li><code>c.Authenticator.admin_users</code>：定义管理员用户列表</li>
              <li><code>c.Spawner.cmd</code>：指定单用户服务器启动命令</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="3. 存储与身份管理" name="9">
            <ul>
              <li>配置用户家目录路径、数据持久化存储</li>
              <li>按需启用本地用户认证或第三方 OAuth/SAML</li>
            </ul>
          </el-collapse-item>
        </el-collapse>
        
        <h4>步骤 4：启动 JupyterHub 服务</h4>
        
        <el-collapse accordion>
          <el-collapse-item title="1. 启动命令" name="10">
            <p>在 <code>jupyter</code> 用户下执行：</p>
            <pre>jupyterhub -f /etc/jupyterhub/jupyterhub_config.py</pre>
          </el-collapse-item>
          
          <el-collapse-item title="2. 守护进程管理（可选）" name="11">
            <ul>
              <li>建议使用 Systemd 或 Supervisor 管理，保证服务重启与日志收集</li>
              <li>将启动命令写入对应服务单元</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="3. 查看启动日志" name="12">
            <p>确认控制台或日志文件中出现以下信息：</p>
            <ul>
              <li>"JupyterHub is now running at…"</li>
              <li>"Spawning [username]" 等关键字</li>
            </ul>
          </el-collapse-item>
        </el-collapse>
        
        <h4>步骤 5：连接测试</h4>
        
        <el-collapse accordion>
          <el-collapse-item title="1. 测试 Web 访问" name="13">
            <p>在浏览器中访问</p>
            <pre>http://服务器IP:8000/hub/login</pre>
            <p>应看到登录页面，输入用户名/密码能正常跳转</p>
          </el-collapse-item>
          
          <el-collapse-item title="2. 单用户 Notebook 测试" name="14">
            <ul>
              <li>登录后，点击 "Start My Server"</li>
              <li>确认能够进入 JupyterLab/Notebook 页面，并能新建 Python 文件</li>
            </ul>
          </el-collapse-item>
          
          <el-collapse-item title="3. API 接口测试" name="15">
            <p>访问内置 REST API：</p>
            <pre>http://服务器IP:8000/hub/api</pre>
            <p>应返回 JSON 格式用户与服务器状态</p>
          </el-collapse-item>
        </el-collapse>
        
        <h4>常见问题与排查</h4>
        <el-collapse>
          <el-collapse-item title="端口无法访问" name="16">
            <ul>
              <li>检查防火墙规则（UFW、iptables）、云平台安全组</li>
              <li>确认 <code>bind_url</code> 中端口设置正确</li>
            </ul>
          </el-collapse-item>
          <el-collapse-item title="服务启动失败" name="17">
            <ul>
              <li>查看 Systemd/日志输出（<code>journalctl -u jupyterhub</code>）</li>
              <li>检查配置文件语法与路径</li>
            </ul>
          </el-collapse-item>
          <el-collapse-item title="用户无权限" name="18">
            <ul>
              <li>确认 <code>jupyter</code> 用户对虚拟环境、数据目录有读写权限</li>
              <li>检查 <code>admin_users</code> 配置是否包含当前用户名</li>
            </ul>
          </el-collapse-item>
          <el-collapse-item title="扩展加载失败" name="19">
            <ul>
              <li>确保对应 Python 包已正确安装且版本兼容</li>
              <li>在配置中添加 <code>configurable_spawner_class</code> 或扩展插件设置</li>
            </ul>
          </el-collapse-item>
        </el-collapse>
        
        <div class="additional-info" style="margin-top: 20px; padding: 15px; background-color: #f8f9fa; border-radius: 4px;">
          <p><strong>提示</strong>：若需启用 HTTPS，需在 <code>jupyterhub_config.py</code> 中配置 <code>c.JupyterHub.ssl_key</code> 与 <code>c.JupyterHub.ssl_cert</code>，并将端口改为 8443 或 443。</p>
          <p><strong>更多参考</strong>：</p>
          <ul>
            <li>JupyterHub 官方文档：<a href="https://jupyterhub.readthedocs.io" target="_blank">https://jupyterhub.readthedocs.io</a></li>
            <li>配置实例与最佳实践：<a href="https://zero-to-jupyterhub.readthedocs.io" target="_blank">https://zero-to-jupyterhub.readthedocs.io</a></li>
          </ul>
        </div>
      </div>

      <div class="form-actions">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="nextStep">
          我已完成沙箱环境配置 <el-icon class="el-icon--right"><ArrowRight /></el-icon>
        </el-button>
      </div>
    </el-card>

    <!-- 步骤2：基本信息 -->
    <el-card v-if="currentStep === 2" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <h3>算力资源基本信息</h3>
        </div>
      </template>

      <el-form :model="computingForm" label-position="top" require-asterisk-position="right" class="computing-form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="服务名称" required>
              <el-input v-model="computingForm.name" placeholder="为您的算力资源取一个明确的名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="系统类型">
              <el-select v-model="computingForm.systemType" placeholder="选择服务器操作系统类型" style="width: 100%">
                <el-option v-for="item in systemOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="服务器地址" required>
              <el-input v-model="computingForm.serverAddress" placeholder="如：192.168.1.100 或 server.example.com" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="开放端口" required>
              <el-input v-model="computingForm.port" placeholder="如：8080" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="CPU配置">
              <el-input v-model="computingForm.cpu" placeholder="如：8核 Intel i7-10700K 3.8GHz" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内存配置">
              <el-input v-model="computingForm.memory" placeholder="如：32GB DDR4 3200MHz" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="GPU配置">
              <el-input v-model="computingForm.gpu" placeholder="如：NVIDIA RTX 3080 10GB" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="存储配置">
              <el-input v-model="computingForm.storage" placeholder="如：1TB NVMe SSD" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="可用时间开始">
              <el-time-picker v-model="computingForm.availableTimeStart" placeholder="选择时间" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="可用时间结束">
              <el-time-picker v-model="computingForm.availableTimeEnd" placeholder="选择时间" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="服务描述">
          <el-input v-model="computingForm.description" type="textarea" :rows="3" placeholder="描述您的算力服务器的特点、用途和适用场景" />
        </el-form-item>

        <el-form-item label="标签 (最多5个)">
          <div class="tags-input-container">
            <el-tag v-for="tag in computingForm.tags" :key="tag" closable @close="removeTag(tag)" class="form-tag">
              {{ tag }}
            </el-tag>
            <el-input v-if="computingForm.tags.length < 5" v-model="tagInput" placeholder="输入标签后按回车" @keyup.enter="addTag" size="small" class="tag-input-field" />
          </div>
          <div class="el-form-item__info">添加有助于其他用户了解您的算力资源的关键字</div>
        </el-form-item>

        <div class="form-actions">
          <el-button @click="prevStep">返回上一步</el-button>
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="nextStep">开始连接测试</el-button>
        </div>
      </el-form>
    </el-card>

    <!-- 步骤3：连接测试 -->
    <el-card v-if="currentStep === 3" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <h3>连接测试</h3>
        </div>
      </template>

      <div class="test-content">
        <p>系统将对您的服务器进行连接测试，确保沙箱环境正常工作。</p>

        <el-timeline class="test-timeline">
          <el-timeline-item
            v-for="(step, index) in testSteps"
            :key="index"
            :type="step.status === 'success' ? 'success' : step.status === 'error' ? 'danger' : 'primary'"
            :icon="step.status === 'success' ? Check : step.status === 'error' ? Close : Warning"
            :hollow="step.status === 'waiting'"
          >
            <h4>{{ step.name }}</h4>
            <p>{{ step.message }}</p>
          </el-timeline-item>
        </el-timeline>

        <div v-if="testResult.success" class="test-result success">
          <el-icon><Check /></el-icon>
          <span>{{ testResult.message }}</span>
        </div>

        <div v-if="testResult.success === false && testResult.message" class="test-result error">
          <el-icon><Close /></el-icon>
          <span>{{ testResult.message }}</span>
        </div>

        <div class="form-actions">
          <el-button @click="prevStep">返回上一步</el-button>
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="startConnectionTest" :loading="testingInProgress" v-if="!testResult.success">
            开始测试
          </el-button>
          <el-button type="success" @click="handleSubmit" v-if="testResult.success">
            提交
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.computing-resource-form {
  margin-bottom: 20px;
}

.form-card {
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 25px;
  margin-bottom: 20px;
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

.tutorial-content {
  margin-bottom: 20px;
}

.tutorial-content h4 {
  margin-top: 20px;
  margin-bottom: 10px;
  color: #333;
}

.tutorial-content pre {
  background-color: #f6f8fa;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
  font-family: monospace;
  margin: 10px 0;
}

.tutorial-content ul {
  padding-left: 20px;
}

.test-timeline {
  margin: 20px 0;
}

.test-result {
  padding: 15px;
  border-radius: 4px;
  margin: 15px 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.test-result.success {
  background-color: #f0f9eb;
  color: #67c23a;
  border: 1px solid #e1f3d8;
}

.test-result.error {
  background-color: #fef0f0;
  color: #f56c6c;
  border: 1px solid #fde2e2;
}

.computing-form .el-form-item__label {
  font-weight: 500;
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