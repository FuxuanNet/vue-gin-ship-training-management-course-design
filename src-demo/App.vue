<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CopyDocument, ArrowDown, ArrowUp } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const isHelpExpanded = ref(true) // 默认展开状态

// 动态计算当前激活的导航索引（根据路由路径）
const activeIndex = computed(() => {
  const pathMap = {
    '/': '1',
    '/upload': '2',
    '/publish': '3',
    '/market': '4',
    '/profile': '5', // 与个人中心导航项的 index="5" 完全匹配
    '/sandbox': '6',
    '/logs': '7'
  }
  return pathMap[route.path] || '1' // 默认首页
})

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}

// 复制到剪贴板的函数
const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage({
      message: '已复制到剪贴板',
      type: 'success',
      duration: 1500
    })
  }, () => {
    ElMessage.error('复制失败，请手动复制')
  })
}

// 切换帮助面板展开状态
const toggleHelpPanel = () => {
  isHelpExpanded.value = !isHelpExpanded.value
}
</script>

<template>
  <div class="app-container">
    <el-container>
      <el-header>
        <el-menu
          :default-active="activeIndex"
          class="el-menu-demo"
          mode="horizontal"
          :ellipsis="false"
        >
          <!-- 左侧主要导航项 -->
          <el-menu-item index="1" @click="router.push('/')">
            <el-icon><HomeFilled /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="2" @click="router.push('/upload')">
            <el-icon><Upload /></el-icon>
            <span>资源上传</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="3" @click="router.push('/publish')">
            <el-icon><Shop /></el-icon>
            <span>资源发布</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="4" @click="router.push('/market')">
            <el-icon><Goods /></el-icon>
            <span>交易市场</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="6" @click="router.push('/sandbox')">
            <el-icon><Monitor /></el-icon>
            <span>沙箱环境</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="7" @click="router.push('/logs')">
            <el-icon><Document /></el-icon>
            <span>区块链日志</span>
          </el-menu-item>

          <!-- 分隔容器（占据中间剩余空间，将右侧内容推到最右） -->
          <div class="flex-grow" />

          <!-- 右侧登录相关导航项 -->
          <el-menu-item v-if="isLoggedIn" index="5" @click="router.push('/profile')">
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="9" @click="logout">
            <el-icon><Logout /></el-icon>
            <span>退出登录</span>
          </el-menu-item>
          <el-menu-item v-else index="8" @click="router.push('/login')">
            <el-icon><UserFilled /></el-icon>
            <span>登录/注册</span>
          </el-menu-item>
        </el-menu>
      </el-header>
      <el-main>
        <router-view />
        
        <!-- 评审专家信息悬浮球 -->
        <div class="expert-floating-panel" :class="{ 'collapsed': !isHelpExpanded }">
          <div class="expert-panel-header" @click="toggleHelpPanel">
            <span>评审专家账号密码</span>
            <el-icon v-if="isHelpExpanded"><ArrowDown /></el-icon>
            <el-icon v-else><ArrowUp /></el-icon>
          </div>
          
          <div v-show="isHelpExpanded" class="expert-panel-content">
            <h4>登录信息</h4>
            <div class="credential-item">
              <div>用户名：<span class="credential">user</span></div>
              <el-button 
                type="primary" 
                :icon="CopyDocument" 
                circle 
                size="small" 
                @click="copyToClipboard('user')"
                class="copy-btn"
                plain
              ></el-button>
            </div>
            <div class="credential-item">
              <div>密码：<span class="credential">123$%^</span></div>
              <el-button 
                type="primary" 
                :icon="CopyDocument" 
                circle 
                size="small" 
                @click="copyToClipboard('123$%^')"
                class="copy-btn"
                plain
              ></el-button>
            </div>
            
            <h4>沙箱环境信息</h4>
            <div class="credential-item">
              <div>用户名：<span class="credential">user</span></div>
              <el-button 
                type="primary" 
                :icon="CopyDocument" 
                circle 
                size="small" 
                @click="copyToClipboard('user')"
                class="copy-btn"
                plain
              ></el-button>
            </div>
            <div class="credential-item">
              <div>密码：<span class="credential">123$%^</span></div>
              <el-button 
                type="primary" 
                :icon="CopyDocument" 
                circle 
                size="small" 
                @click="copyToClipboard('123$%^')"
                class="copy-btn"
                plain
              ></el-button>
            </div>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.el-container {
  min-height: 100vh;
}

.el-header {
  padding: 0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.el-main {
  padding: 20px;
  flex: 1;
}

.el-footer {
  text-align: center;
  background-color: #f5f7fa;
  color: #888;
  padding: 20px;
}

.flex-grow {
  flex-grow: 1;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

/* 评审专家信息悬浮球样式 */
.expert-floating-panel {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background-color: rgba(250, 250, 250, 0.65);
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  width: 320px;
  transition: all 0.3s cubic-bezier(0.2, 0.8, 0.2, 1);
  z-index: 1000;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.expert-floating-panel.collapsed {
  width: 220px;
  height: 48px;
  border-radius: 24px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.expert-panel-header {
  background: rgba(59, 130, 246, 0.8);
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
  color: white;
  padding: 14px 18px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  font-weight: 500;
  font-size: 15px;
  letter-spacing: 0.3px;
}

.expert-panel-content {
  padding: 20px;
  background-color: rgba(250, 250, 250, 0.5);
}

.expert-panel-content h4 {
  margin-top: 0;
  margin-bottom: 16px;
  color: #1d1d1f;
  font-size: 16px;
  position: relative;
  font-weight: 500;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  padding-bottom: 8px;
  padding-left: 14px;
}

.expert-panel-content h4::before {
  content: "";
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 16px;
  background-color: rgba(0, 122, 255, 0.8);
  border-radius: 2px;
}

.credential-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  background-color: rgba(255, 255, 255, 0.5);
  padding: 10px 14px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.03);
  transition: all 0.2s ease;
}

.credential-item:hover {
  transform: translateY(-2px);
  background-color: rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.credential {
  font-weight: 500;
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Text', sans-serif;
  background-color: rgba(0, 0, 0, 0.03);
  padding: 4px 10px;
  border-radius: 6px;
  color: #007aff;
  border: none;
}

.copy-btn {
  flex-shrink: 0;
  transition: all 0.2s ease;
  background-color: rgba(0, 122, 255, 0.1);
  border: none;
  color: #007aff;
  margin-left: 8px;
}

.copy-btn:hover {
  transform: scale(1.05);
  background-color: rgba(0, 122, 255, 0.2);
}
</style>
