<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { useMockDataStore } from '../stores/mockData'

const router = useRouter()
const userStore = useUserStore()
const mockDataStore = useMockDataStore()
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
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 处理登录
const handleLogin = () => {
  loginFormRef.value.validate((valid) => {
    if (valid) {
      loading.value = true
      
      // 模拟登录验证
      const account = mockDataStore.accounts.find(
        acc => acc.login_name === loginForm.username && acc.password_hash === loginForm.password
      )
      
      if (account) {
        const person = mockDataStore.persons.find(p => p.person_id === account.person_id)
        
        // 角色映射
        const roleMap = {
          '员工': 'employee',
          '讲师': 'teacher',
          '课程大纲制定者': 'planner'
        }
        
        userStore.login({
          token: 'mock-token-' + Date.now(),
          user: {
            id: person.person_id,
            name: person.name,
            role: roleMap[person.role]
          }
        })
        
        ElMessage.success('登录成功')
        loading.value = false
        router.push('/')
      } else {
        ElMessage.error('用户名或密码错误')
        loading.value = false
      }
    }
  })
}

// 快速填充测试账号
const fillTestAccount = (role) => {
  const testAccounts = {
    employee: { username: 'employee', password: '123456' },
    teacher: { username: 'teacher', password: '123456' },
    planner: { username: 'planner', password: '123456' }
  }
  loginForm.username = testAccounts[role].username
  loginForm.password = testAccounts[role].password
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="title">船舶企业培训管理系统</h2>
      <p class="subtitle">Training Management System</p>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="0" class="login-form">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" placeholder="用户名" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password size="large" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="login-button" @click="handleLogin" size="large">登录</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 测试账号提示 -->
      <div class="test-accounts">
        <el-divider>测试账号</el-divider>
        <div class="account-buttons">
          <el-button size="small" @click="fillTestAccount('employee')">员工账号</el-button>
          <el-button size="small" @click="fillTestAccount('teacher')">讲师账号</el-button>
          <el-button size="small" @click="fillTestAccount('planner')">管理员账号</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-image: url('/image.jpg');
  background-size: cover;
  background-position: center;
  position: relative;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(5px);
}

.login-card {
  width: 420px;
  padding: 45px 40px;
  border-radius: 12px;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
  background-color: #fff;
  position: relative;
  z-index: 1;
}

.title {
  text-align: center;
  margin-bottom: 8px;
  color: #303133;
  font-weight: 600;
  font-size: 24px;
}

.subtitle {
  text-align: center;
  margin-bottom: 30px;
  color: #909399;
  font-size: 13px;
  letter-spacing: 1px;
}

.login-form {
  margin-top: 30px;
}

.login-button {
  width: 100%;
  margin-top: 10px;
  font-weight: 500;
}

.test-accounts {
  margin-top: 30px;
}

.account-buttons {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.account-buttons .el-button {
  flex: 1;
}
</style>
