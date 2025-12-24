<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { Edit, User, Clock } from '@element-plus/icons-vue'
import { getCurrentUser, logout, updatePassword, uploadAvatar, getAvatarUrl, depositBalance, withdrawBalance, getTransactions } from '../api/auth'
import { getFavorites, unfavoriteResource, getMarketResourceDetail, getPurchasedResources } from '../api/market'
import { marked } from 'marked'

const router = useRouter()

// 用户信息
const userInfo = ref({
  avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
  name: '张三',
  company: '科技有限公司',
  email: 'zhangsan@example.com',
  phone: '138****1234',
  role: 'data_provider',
  joinTime: '2023-01-15',
  credits: 2350,
  balance: 0
})

// 头像相关
const avatarHover = ref(false)
const avatarUploadRef = ref(null)
const avatarUploading = ref(false)
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'

// 处理头像悬浮
const handleAvatarMouseEnter = () => {
  avatarHover.value = true
}

const handleAvatarMouseLeave = () => {
  avatarHover.value = false
}

// 触发文件选择
const triggerAvatarUpload = () => {
  avatarUploadRef.value.click()
}

// 处理头像上传
const handleAvatarChange = (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  // 验证文件类型
  const acceptedTypes = ['image/jpeg', 'image/png', 'image/gif']
  if (!acceptedTypes.includes(file.type)) {
    ElMessage.error('请上传jpg、png或gif格式的图片')
    return
  }
  
  // 验证文件大小
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过2MB')
    return
  }
  
  // 上传头像
  avatarUploading.value = true
  uploadAvatar(file)
    .then(response => {
      // 更新用户信息中的头像 - 使用完整URL
      const userId = userInfo.value.id
      if (userId) {
        // 使用getAvatarUrl获取完整的头像URL
        userInfo.value.avatar = getAvatarUrl(userId) + '?t=' + new Date().getTime()
      } else {
        // 如果没有用户ID，直接使用响应中的路径
        userInfo.value.avatar = response.data.avatar
      }
      ElMessage.success('头像上传成功')
    })
    .catch(error => {
      console.error('头像上传失败:', error)
      ElMessage.error(error.message || '头像上传失败')
    })
    .finally(() => {
      avatarUploading.value = false
      // 清空文件选择，以便可以选择相同的文件
      event.target.value = ''
    })
}

// 用户角色中文名称
const roleName = computed(() => {
  const map = {
    data_provider: '数据提供方',
    data_user: '数据使用方',
    data_service: '数据服务方'
  }
  return map[userInfo.value.role] || '未知角色'
})

// 活跃标签
const activeTab = ref('favorite')

// 充值/提现对话框
const balanceDialogVisible = ref(false)
const balanceDialogType = ref('deposit') // 'deposit' 或 'withdraw'
const balanceForm = reactive({
  amount: '',
  remark: ''
})
const balanceFormRef = ref(null)
const submittingBalance = ref(false)

// 余额相关规则
const balanceFormRules = {
  amount: [
    { required: true, message: '请输入金额', trigger: 'blur' },
    { validator: (rule, value, callback) => {
      if (value <= 0) {
        callback(new Error('金额必须大于0'))
      } else if (balanceDialogType.value === 'withdraw' && value > userInfo.value.balance) {
        callback(new Error('提现金额不能超过余额'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ]
}

// 打开充值对话框
const openDepositDialog = () => {
  balanceDialogType.value = 'deposit'
  balanceForm.amount = ''
  balanceForm.remark = ''
  balanceDialogVisible.value = true
}

// 打开提现对话框
const openWithdrawDialog = () => {
  balanceDialogType.value = 'withdraw'
  balanceForm.amount = ''
  balanceForm.remark = ''
  balanceDialogVisible.value = true
}

// 处理余额操作
const handleBalanceOperation = async () => {
  if (!balanceFormRef.value) return
  
  try {
    await balanceFormRef.value.validate()
    
    const operation = balanceDialogType.value === 'deposit' ? 
      depositBalance({
        amount: Number(balanceForm.amount),
        remark: balanceForm.remark
      }) : 
      withdrawBalance({
        amount: Number(balanceForm.amount),
        remark: balanceForm.remark
      })
    
    const operationText = balanceDialogType.value === 'deposit' ? '充值' : '提现'
    
    submittingBalance.value = true
    
    try {
      const result = await operation
      if (result.code === 200) {
        ElMessage.success(`${operationText}成功`)
        // 更新用户余额
        userInfo.value.balance = result.data.balance
        // 关闭对话框
        balanceDialogVisible.value = false
        // 重新获取交易记录
        if (activeTab.value === 'transactions') {
          fetchTransactions()
        }
      } else {
        ElMessage.error(result.message || `${operationText}失败`)
      }
    } catch (error) {
      console.error(`${operationText}出错:`, error)
      ElMessage.error(error.message || `${operationText}失败，请稍后重试`)
    } finally {
      submittingBalance.value = false
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 用户上传的资源
const uploadedResources = ref([])
// 用户购买的资源
const purchasedResources = ref([])
// 交易记录
const transactions = ref([])

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

// 获取用户上传的资源
const fetchUploadedResources = () => {
  // 模拟API调用
  setTimeout(() => {
    uploadedResources.value = [
      {
        id: 1,
        name: '城市交通流量数据集',
        type: 'data',
        description: '包含2022年全年主要城市道路交通流量统计数据，每15分钟一次采样',
        createTime: '2023-03-15',
        size: '2.3GB',
        status: 'published',
        tags: ['交通', '城市', '结构化数据']
      },
      {
        id: 2,
        name: '电商用户行为分析模型',
        type: 'algorithm',
        description: '基于深度学习的电商用户行为分析模型，可预测用户购买倾向',
        createTime: '2023-04-20',
        size: '450MB',
        status: 'published',
        tags: ['电商', '用户行为', '深度学习']
      },
      {
        id: 3,
        name: '医疗患者数据',
        type: 'data',
        description: '匿名化的患者诊断和治疗记录数据集，适用于医疗AI研究',
        createTime: '2023-05-18',
        size: '1.5GB',
        status: 'private',
        tags: ['医疗', '患者数据', '结构化数据']
      }
    ]
  }, 1000)
}

// 获取用户购买的资源
const fetchPurchasedResources = async () => {
  try {
    const loading = ElLoading.service({
      lock: true,
      text: '加载购买记录...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    const res = await getPurchasedResources()
    if (res.code === 200) {
      purchasedResources.value = res.data || []
    } else {
      ElMessage.warning(res.message || '获取购买记录为空')
      purchasedResources.value = []
    }
    
    loading.close()
  } catch (error) {
    console.error('获取购买记录失败:', error)
    ElMessage.error(error.message || '获取购买记录失败，请稍后重试')
    purchasedResources.value = []
  }
}

// 获取交易记录
const fetchTransactions = async () => {
  try {
    const loading = ElLoading.service({
      lock: true,
      text: '加载交易记录...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    const res = await getTransactions()
    if (res.code === 200 && res.data) {
      transactions.value = res.data.transactions.map(t => {
        return {
          id: t.id,
          type: t.type,
          amount: t.amount,
          balance: t.balance,
          remark: t.remark || '',
          time: formatDate(t.created_at),
          relatedId: t.related_id || 0
        }
      })
    } else {
      ElMessage.warning(res.message || '获取交易记录为空')
      transactions.value = []
    }
    
    loading.close()
  } catch (error) {
    console.error('获取交易记录失败:', error)
    ElMessage.error(error.message || '获取交易记录失败，请稍后重试')
  }
}

// 加载用户信息
const loadUserInfo = async () => {
  try {
    // 显示加载状态
    const loading = ElLoading.service({
      lock: true,
      text: '正在加载用户信息...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    const response = await getCurrentUser()
    // 更新用户信息
    const userData = response.data
    userInfo.value = {
      ...userData,
      balance: userData.balance || 0, // 确保余额字段存在
    }
    
    // 如果用户有自定义头像，则使用完整URL加载
    if (userData.id && userData.avatar && !userData.avatar.startsWith('http')) {
      // 添加时间戳参数避免缓存
      userInfo.value.avatar = getAvatarUrl(userData.id) + '?t=' + new Date().getTime()
    }
    
    loading.close()
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取角色名称
const getRoleName = (role) => {
  const roleMap = {
    'data_provider': '数据提供方',
    'data_user': '数据使用方',
    'data_service': '数据服务方',
    'admin': '管理员'
  }
  return roleMap[role] || role
}

// 用户收藏的资源
const favoriteResources = ref([])

// 获取用户收藏的资源
const fetchFavoriteResources = async () => {
  try {
    const loading = ElLoading.service({
      lock: true,
      text: '加载收藏资源...',
      background: 'rgba(255, 255, 255, 0.7)'
    })
    
    const res = await getFavorites()
    if (res.code === 200) {
      favoriteResources.value = res.data || []
    } else {
      ElMessage.warning(res.message || '获取收藏资源为空')
      favoriteResources.value = []
    }
    
    loading.close()
  } catch (error) {
    console.error('获取收藏资源失败:', error)
    ElMessage.error(error.message || '获取收藏资源失败，请稍后重试')
    favoriteResources.value = []
  }
}

// 取消收藏
const handleUnfavorite = async (resource) => {
  try {
    const res = await unfavoriteResource(resource.resource_id)
    if (res.code === 200) {
      ElMessage.success('取消收藏成功')
      // 重新获取收藏列表
      fetchFavoriteResources()
    } else {
      ElMessage.error(res.message || '取消收藏失败')
    }
  } catch (error) {
    console.error('取消收藏失败:', error)
    ElMessage.error(error.message || '取消收藏失败，请稍后重试')
  }
}

onMounted(() => {
  loadUserInfo()
  fetchFavoriteResources()
  fetchPurchasedResources()
  fetchTransactions()
})

// 查看交易合约
const contractDialogVisible = ref(false)
const currentContract = ref(null)

const viewContract = (contractId) => {
  // 模拟查询合约详情
  currentContract.value = {
    id: contractId,
    title: '数据资源使用许可协议',
    content: `
    甲方（许可方）：${userInfo.value.name}
    乙方（被许可方）：某公司

    一、许可内容
    甲方授权乙方使用以下资源：${transactions.value.find(t => t.contractId === contractId)?.resourceName || '未知资源'}

    二、许可范围
    乙方可在以下范围内使用该资源：
    1. 仅限于被许可方内部使用
    2. 禁止对外提供数据服务
    3. 禁止二次转售

    三、许可期限
    自协议生效日起一年内有效

    四、使用限制
    乙方不得将数据用于以下用途：
    1. 违反国家法律法规的活动
    2. 侵害第三方合法权益的活动
    3. 其他未经甲方书面许可的用途

    五、知识产权
    甲方保留对所提供资源的所有权利

    六、安全保障
    乙方应采取必要措施保障数据安全，防止未经授权的访问和泄露

    七、终止条款
    如乙方违反本协议任何条款，甲方有权立即终止本协议

    八、争议解决
    因本协议引起的任何争议，双方应协商解决；协商不成的，提交至有管辖权的人民法院诉讼解决

    九、其他
    本协议自双方签字盖章之日起生效

    甲方（签章）：
    日期：

    乙方（签章）：
    日期：
    `,
    date: '2023-04-10',
    status: 'valid'
  }
  contractDialogVisible.value = true
}

// 打开沙箱环境使用购买的资源
const openInSandbox = (resource) => {
  // 将资源ID存储到本地，以便在沙箱环境中读取
  localStorage.setItem('sandboxResource', JSON.stringify(resource))
  // 跳转到沙箱环境
  router.push('/sandbox')
}

// 用户资料编辑对话框
const editProfileDialog = ref(false)
const userForm = reactive({
  name: '',
  email: '',
  phone: '',
  company: '',
  avatar: ''
})

const openEditProfile = () => {
  userForm.name = userInfo.value.name
  userForm.email = userInfo.value.email
  userForm.phone = userInfo.value.phone
  userForm.company = userInfo.value.company
  userForm.avatar = userInfo.value.avatar
  editProfileDialog.value = true
}

const submitUserForm = () => {
  // 模拟保存用户信息
  userInfo.value = {
    ...userInfo.value,
    name: userForm.name,
    email: userForm.email,
    phone: userForm.phone,
    company: userForm.company,
    avatar: userForm.avatar
  }
  editProfileDialog.value = false
  ElMessage.success('个人信息更新成功')
}

// 修改密码相关
const passwordFormRef = ref(null)
const changePasswordVisible = ref(false)
const changingPassword = ref(false)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmNewPassword: ''
})

// 验证新密码与确认密码是否一致
const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入新密码'))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
  ],
  confirmNewPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// 显示修改密码对话框
const showChangePasswordDialog = () => {
  changePasswordVisible.value = true
  // 重置表单
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmNewPassword = ''
  if (passwordFormRef.value) {
    passwordFormRef.value.resetFields()
  }
}

// 处理修改密码
const handleChangePassword = () => {
  passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        changingPassword.value = true
        await updatePassword({
          oldPassword: passwordForm.oldPassword,
          newPassword: passwordForm.newPassword
        })
        
        changePasswordVisible.value = false
        ElMessage.success('密码修改成功')
      } catch (error) {
        console.error('修改密码失败:', error)
        // 显示后端返回的具体错误消息
        ElMessage.error(error.message || '修改密码失败，请稍后重试')
      } finally {
        changingPassword.value = false
      }
    }
  })
}

// 处理退出登录
const handleLogout = () => {
  ElMessageBox.confirm(
    '确定要退出登录吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await logout()
      // 清除本地存储的用户信息和令牌
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      ElMessage.success('已退出登录')
      router.push('/login')
    } catch (error) {
      console.error('退出登录失败:', error)
      ElMessage.error('退出登录失败')
    }
  }).catch(() => {})
}

// 资源详情对话框
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
      target: '.resource-tabs',
      text: '加载资源详情...'
    })
    
    const res = await getMarketResourceDetail(resource.resource_id)
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
</script>

<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="profile-header">
          <h2>个人信息</h2>
        </div>
      </template>

      <div class="profile-content">
        <div class="profile-avatar">
          <div class="avatar-container" 
            @mouseenter="handleAvatarMouseEnter" 
            @mouseleave="handleAvatarMouseLeave" 
            @click="triggerAvatarUpload"
            v-loading="avatarUploading"
            element-loading-text="上传中..."
            element-loading-background="rgba(255, 255, 255, 0.8)">
            <el-avatar :size="100" :src="userInfo.avatar || defaultAvatar"></el-avatar>
            <div class="avatar-overlay" v-show="avatarHover">
              <el-icon :size="24"><Edit /></el-icon>
              <span>更换头像</span>
            </div>
          </div>
          <input type="file" ref="avatarUploadRef" @change="handleAvatarChange" accept="image/jpeg,image/png,image/gif" style="display: none;">
        </div>

        <div class="profile-info">
          <div class="info-item">
            <span class="label">用户名</span>
            <span class="value">{{ userInfo.username }}</span>
          </div>

          <div class="info-item">
            <span class="label">邮箱</span>
            <span class="value">{{ userInfo.email }}</span>
          </div>

          <div class="info-item">
            <span class="label">所属企业</span>
            <span class="value">{{ userInfo.company }}</span>
          </div>

          <div class="info-item">
            <span class="label">用户角色</span>
            <span class="value">{{ getRoleName(userInfo.role) }}</span>
          </div>

          <div class="info-item">
            <span class="label">账户余额</span>
            <span class="value">
              ¥{{ userInfo.balance.toFixed(2) }}
              <el-button type="primary" size="small" @click="openDepositDialog" style="margin-left: 10px;">充值</el-button>
              <el-button type="warning" size="small" @click="openWithdrawDialog" :disabled="userInfo.balance <= 0">提现</el-button>
            </span>
          </div>

          <div class="info-item">
            <span class="label">注册时间</span>
            <span class="value">{{ formatDate(userInfo.created_at) }}</span>
          </div>
          
          <div class="info-item" v-if="userInfo.last_login">
            <span class="label">最后登录</span>
            <span class="value">{{ formatDate(userInfo.last_login) }}</span>
          </div>
        </div>
        </div>
        
      <div class="profile-actions">
        <el-button type="primary" @click="showChangePasswordDialog">修改密码</el-button>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </el-card>
    
    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="changePasswordVisible"
      title="修改密码"
      width="400px"
      center
    >
      <el-form
        :model="passwordForm"
        :rules="passwordRules"
        ref="passwordFormRef"
        label-width="100px"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            placeholder="请输入旧密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmNewPassword">
          <el-input
            v-model="passwordForm.confirmNewPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="changePasswordVisible = false">取消</el-button>
          <el-button type="primary" @click="handleChangePassword" :loading="changingPassword">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 充值/提现对话框 -->
    <el-dialog
      v-model="balanceDialogVisible"
      :title="balanceDialogType === 'deposit' ? '充值余额' : '提现余额'"
      width="400px"
      center
    >
      <el-form
        :model="balanceForm"
        :rules="balanceFormRules"
        ref="balanceFormRef"
        label-width="80px"
      >
        <el-form-item label="金额" prop="amount">
          <el-input
            v-model.number="balanceForm.amount"
            type="number"
            min="0.01"
            :placeholder="balanceDialogType === 'deposit' ? '请输入充值金额' : '请输入提现金额'"
          >
            <template #prefix>¥</template>
          </el-input>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="balanceForm.remark"
            type="text"
            :placeholder="balanceDialogType === 'deposit' ? '充值用途（选填）' : '提现备注（选填）'"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="balanceDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleBalanceOperation" :loading="submittingBalance">
            {{ balanceDialogType === 'deposit' ? '确认充值' : '确认提现' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 添加交易记录Tab -->
    <el-tabs v-model="activeTab" class="resource-tabs" style="margin-top: 20px;">
      <el-tab-pane label="我的收藏" name="favorite">
        <el-empty v-if="!favoriteResources?.length" description="暂无收藏资源" />
        <el-table v-else :data="favoriteResources" style="width: 100%">
          <el-table-column prop="name" label="资源名称" min-width="200">
            <template #default="scope">
              <div class="resource-name">
                <el-tag :type="typeColorMap[scope.row.resource_type]" size="small" style="margin-right: 8px">
                  {{ typeNameMap[scope.row.resource_type] }}
                </el-tag>
                {{ scope.row.name }}
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="publisher" label="发布者" width="120" />
          <el-table-column prop="price" label="价格" width="120">
            <template #default="scope">
              ¥{{ scope.row.price.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="tags" label="标签" min-width="200">
            <template #default="scope">
              <el-tag
                v-for="tag in scope.row.tags.slice(0, 3)"
                :key="tag"
                size="small"
                effect="plain"
                style="margin-right: 4px"
              >
                {{ tag }}
              </el-tag>
              <el-tag
                v-if="scope.row.tags.length > 3"
                size="small"
                effect="plain"
              >
                +{{ scope.row.tags.length - 3 }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="favorite_time" label="收藏时间" width="160" />
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="openResourceDetail(scope.row)">
                详情
              </el-button>
              <el-button type="danger" size="small" @click="handleUnfavorite(scope.row)">
                取消收藏
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="我的购买" name="purchase">
        <el-empty v-if="!purchasedResources?.length" description="暂无购买记录" />
        <el-table v-else :data="purchasedResources" style="width: 100%">
          <el-table-column prop="name" label="资源名称" min-width="200">
            <template #default="scope">
              <div class="resource-name">
                <el-tag :type="typeColorMap[scope.row.resource_type]" size="small" style="margin-right: 8px">
                  {{ typeNameMap[scope.row.resource_type] }}
                </el-tag>
                {{ scope.row.name }}
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="publisher" label="发布者" width="120" />
          <el-table-column prop="price" label="购买价格" width="120">
            <template #default="scope">
              ¥{{ scope.row.price.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="purchase_time" label="购买时间" width="160" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
                {{ scope.row.status === 'active' ? '可用' : '已过期' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="openResourceDetail(scope.row)">
                详情
              </el-button>
              <el-button 
                type="success" 
                size="small" 
                @click="$router.push('/sandbox')"
                :disabled="scope.row.status !== 'active'"
              >
                使用
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="交易记录" name="transactions">
        <el-table :data="transactions" style="width: 100%">
          <el-table-column prop="id" label="交易号" width="80" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="scope">
              <el-tag
                :type="scope.row.type === 'deposit' ? 'success' : 
                      scope.row.type === 'withdraw' ? 'warning' :
                      scope.row.type === 'purchase' ? 'danger' : 'info'"
              >
                {{ 
                  scope.row.type === 'deposit' ? '充值' : 
                  scope.row.type === 'withdraw' ? '提现' :
                  scope.row.type === 'purchase' ? '购买' :
                  scope.row.type === 'sale' ? '销售' : '其他'
                }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="金额" width="120">
            <template #default="scope">
              <span :class="scope.row.amount >= 0 ? 'positive-amount' : 'negative-amount'">
                {{ scope.row.amount >= 0 ? '+' : '' }}{{ scope.row.amount.toFixed(2) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="balance" label="余额" width="120">
            <template #default="scope">
              {{ scope.row.balance.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" />
          <el-table-column prop="time" label="交易时间" width="180" />
        </el-table>
      </el-tab-pane>
    </el-tabs>

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
            <span>{{ detailResource.published_at }}</span>
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
            <el-button type="primary" @click="viewSampleData">查看示例数据</el-button>
            <el-button type="success" @click="viewDocument">查看说明文档</el-button>
          </template>
          
          <!-- 算力资源按钮 -->
          <template v-else-if="detailResource.resource_type === 'computing'">
            <el-button type="primary">按量计费</el-button>
            <el-button type="primary">按月计费</el-button>
            <el-button type="primary">按年计费</el-button>
            <el-button type="success" @click="viewDocument">查看说明文档</el-button>
          </template>
          
          <!-- 算法资源按钮 -->
          <template v-else-if="detailResource.resource_type === 'algorithm'">
            <el-button type="primary">立即购买</el-button>
            <el-button type="success" @click="viewDocument">查看说明文档</el-button>
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
.profile-container {
  max-width: 800px;
  margin: 40px auto;
  padding: 0 20px;
}

.profile-card {
  margin-bottom: 20px;
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.profile-avatar {
  margin-bottom: 30px;
}

.avatar-container {
  position: relative;
  display: inline-block;
  border-radius: 50%;
  cursor: pointer;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #fff;
}

.avatar-overlay i {
  font-size: 24px;
  margin-bottom: 5px;
}

.avatar-overlay span {
  font-size: 14px;
}

.profile-info {
  width: 100%;
  max-width: 500px;
}

.info-item {
  display: flex;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #EBEEF5;
}

.label {
  flex: 0 0 100px;
  color: #606266;
  font-size: 14px;
}

.value {
  flex: 1;
  color: #303133;
  font-size: 14px;
}

.profile-actions {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  gap: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.positive-amount {
  color: #67c23a;
  font-weight: bold;
}

.negative-amount {
  color: #f56c6c;
  font-weight: bold;
}

.resource-tabs {
  margin-top: 20px;
}

.balance-info {
  display: flex;
  align-items: center;
}

.balance-actions {
  margin-left: 20px;
}

.resource-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.resource-card {
  width: calc(33.33% - 20px);
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

.resource-type-tag {
  margin-bottom: 10px;
}

.resource-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.resource-tags {
  margin-bottom: 10px;
}

.resource-tag {
  margin-right: 5px;
}

.resource-description {
  margin-bottom: 10px;
}

.resource-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.resource-provider {
  display: flex;
  align-items: center;
}

.resource-provider i {
  margin-right: 5px;
}

.resource-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.resource-price {
  font-weight: bold;
}

.resource-name {
  display: flex;
  align-items: center;
}

.resource-detail {
  padding: 20px;
}

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