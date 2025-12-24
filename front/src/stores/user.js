import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  
  const isLoggedIn = computed(() => !!token.value)
  const userName = computed(() => userInfo.value.name || '')
  const userRole = computed(() => userInfo.value.role || '')
  
  function login(loginData) {
    token.value = loginData.token
    userInfo.value = loginData.user
    localStorage.setItem('token', loginData.token)
    localStorage.setItem('userInfo', JSON.stringify(loginData.user))
  }
  
  function logout() {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }
  
  return {
    token,
    userInfo,
    isLoggedIn,
    userName,
    userRole,
    login,
    logout
  }
})
