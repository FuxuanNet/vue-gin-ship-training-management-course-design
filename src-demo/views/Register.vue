<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '../api/auth'

const router = useRouter()

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  company: '',
  role: ''
})

const loading = ref(false)
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  company: [
    { required: true, message: '请输入公司/组织名称', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

const registerFormRef = ref(null)

const handleRegister = () => {
  registerFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        await register(registerForm)
        
        ElMessage.success('注册成功，请登录')
        router.push('/login')
      } catch (error) {
        console.error('注册失败:', error)
        // 显示后端返回的具体错误消息
        ElMessage.error(error.message || '注册失败，请稍后再试')
      } finally {
        loading.value = false
      }
    } else {
      console.log('表单验证失败')
      return false
    }
  })
}

const goToLogin = () => {
  router.push('/login')
}

const roleOptions = [
  { label: '数据提供方', value: 'data_provider', description: '上传数据、发布数据到交易市场' },
  { label: '数据使用方', value: 'data_user', description: '购买数据、使用数据、联系数据服务方定制等' },
  { label: '数据服务方', value: 'data_service', description: '提供算力服务、托管服务、定制服务、算法插件包开发、中介服务' }
]
</script>

<template>
  <div class="register-container">
    <div class="register-card">
      <h2 class="title">注册 TrustedSpace</h2>
      <el-form :model="registerForm" :rules="rules" ref="registerFormRef" label-width="0" class="register-form">
        <el-form-item prop="username">
          <el-input v-model="registerForm.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="确认密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="registerForm.email" placeholder="邮箱" prefix-icon="Message" />
        </el-form-item>
        <el-form-item prop="company">
          <el-input v-model="registerForm.company" placeholder="公司/组织" prefix-icon="Office-Building" />
        </el-form-item>
        <el-form-item prop="role">
          <el-select v-model="registerForm.role" placeholder="选择角色" class="role-select">
            <el-option label="数据提供方" value="data_provider"></el-option>
            <el-option label="数据使用方" value="data_user"></el-option>
            <el-option label="数据服务方" value="data_service"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="register-button" @click="handleRegister">注册</el-button>
        </el-form-item>
        <div class="form-footer">
          <span>已有账号？</span>
          <router-link to="/login">去登录</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 40px 0;
}

.register-card {
  width: 400px;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
  background-color: #fff;
}

.title {
  text-align: center;
  margin-bottom: 30px;
  color: #303133;
  font-weight: 500;
}

.register-form {
  margin-top: 20px;
}

.role-select {
  width: 100%;
}

.register-button {
  width: 100%;
  margin-top: 10px;
}

.form-footer {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
  color: #606266;
}

.form-footer a {
  color: #409EFF;
  text-decoration: none;
  margin-left: 5px;
}
</style> 