<script>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '../api/auth'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const loginFormRef = ref(null)
    const loading = ref(false)

    // 登录表单数据
    const loginForm = reactive({
      username: '',
      password: ''
    })

    // 表单验证规则
    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
      ]
    }

    // 处理登录
    const handleLogin = () => {
      loginFormRef.value.validate(async (valid) => {
        if (valid) {
          try {
            loading.value = true
            const response = await login(loginForm)
            
            // 登录成功，存储token和用户信息
            const { token, user } = response.data
            localStorage.setItem('token', token)
            localStorage.setItem('user', JSON.stringify(user))
            
            ElMessage.success('登录成功')
            router.push('/')
          } catch (error) {
            console.error('登录失败:', error)
            // 显示具体错误信息
            ElMessage.error(error.message || '登录失败，请检查用户名和密码')
          } finally {
            loading.value = false
          }
        } else {
          console.log('表单验证失败')
          return false
        }
      })
    }

    return {
      loginFormRef,
      loginForm,
      rules,
      loading,
      handleLogin
    }
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="title">登录 TrustedSpace</h2>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="0" class="login-form">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="login-button" @click="handleLogin">登录</el-button>
        </el-form-item>
        <div class="form-footer">
          <span>还没有账号？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
}

.login-card {
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

.login-form {
  margin-top: 20px;
}

.login-button {
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