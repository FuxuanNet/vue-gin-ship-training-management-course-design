import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 使用 sessionId 而不是 token（与后端 Session-ID 对应）
  const sessionId = ref(localStorage.getItem('sessionId') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  
  const isLoggedIn = computed(() => !!sessionId.value)
  const userName = computed(() => userInfo.value.name || '')
  const userRole = computed(() => userInfo.value.role || '')
  
  // 登录：保存会话ID和用户信息
  function login(loginData) {
    sessionId.value = loginData.token // 后端返回的是 sessionId，字段名为 token
    userInfo.value = loginData.user
    localStorage.setItem('sessionId', loginData.token)
    localStorage.setItem('userInfo', JSON.stringify(loginData.user))
  }
  
  // 退出：清除本地数据
  function logout() {
    sessionId.value = ''
    userInfo.value = {}
    localStorage.removeItem('sessionId')
    localStorage.removeItem('userInfo')
  }
  
  return {
    sessionId,
    userInfo,
    isLoggedIn,
    userName,
    userRole,
    login,
    logout
  }
})
