<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from './stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.isLoggedIn)
const userRole = computed(() => userStore.userRole)
const userName = computed(() => userStore.userName)

// 动态计算当前激活的导航索引
const activeIndex = computed(() => {
  return route.path
})

const logout = () => {
  userStore.logout()
  router.push('/login')
}

// 根据用户角色获取所有菜单项
const allMenuItems = computed(() => {
  if (!isLoggedIn.value) return []
  
  const roleMenus = {
    employee: [
      { index: '/employee/today', label: '今日课程', path: '/employee/today', icon: 'Calendar' },
      { index: '/employee/schedule', label: '课程表', path: '/employee/schedule', icon: 'Grid' },
      { index: '/employee/evaluation', label: '课程自评', path: '/employee/evaluation', icon: 'EditPen' },
      { index: '/employee/scores', label: '我的成绩', path: '/employee/scores', icon: 'TrophyBase' }
    ],
    teacher: [
      { index: '/teacher/today', label: '今日授课', path: '/teacher/today', icon: 'Calendar' },
      { index: '/teacher/schedule', label: '授课表', path: '/teacher/schedule', icon: 'Grid' },
      { index: '/teacher/grading', label: '学员评分', path: '/teacher/grading', icon: 'Edit' },
      { index: '/teacher/scores', label: '成绩查看', path: '/teacher/scores', icon: 'DataLine' }
    ],
    planner: [
      { index: '/planner/plans', label: '培训计划', path: '/planner/plans', icon: 'Document' },
      { index: '/planner/courses', label: '课程管理', path: '/planner/courses', icon: 'Reading' },
      { index: '/planner/analytics', label: '数据分析', path: '/planner/analytics', icon: 'DataAnalysis' }
    ]
  }
  
  return roleMenus[userRole.value] || []
})
</script>

<template>
  <div class="app-container">
    <el-container>
      <el-header v-if="isLoggedIn">
        <el-menu
          :default-active="activeIndex"
          class="el-menu-demo"
          mode="horizontal"
          :ellipsis="false"
        >
          <!-- 左侧主要导航项 -->
          <el-menu-item index="/" @click="router.push('/')">
            <el-icon><HomeFilled /></el-icon>
            <span>首页</span>
          </el-menu-item>
          
          <!-- 平铺显示所有导航项 -->
          <el-menu-item 
            v-for="item in allMenuItems" 
            :key="item.index" 
            :index="item.index"
            @click="router.push(item.path)"
            class="nav-menu-item"
          >
            <span>{{ item.label }}</span>
          </el-menu-item>

          <!-- 更多菜单（移动端显示） -->
          <el-sub-menu index="more" class="more-menu">
            <template #title>
              <el-icon><More /></el-icon>
              <span>更多</span>
            </template>
            <el-menu-item 
              v-for="item in allMenuItems" 
              :key="'more-' + item.index" 
              :index="item.index"
              @click="router.push(item.path)"
            >
              {{ item.label }}
            </el-menu-item>
          </el-sub-menu>

          <!-- 分隔容器 -->
          <div class="flex-grow" />

          <!-- 右侧用户信息 -->
          <el-menu-item index="/user" disabled class="user-info-item">
            <el-icon><Avatar /></el-icon>
            <span class="user-name">{{ userName }}</span>
          </el-menu-item>
          <el-menu-item index="/logout" @click="logout">
            <el-icon><SwitchButton /></el-icon>
            <span>退出</span>
          </el-menu-item>
        </el-menu>
      </el-header>
      
      <el-main>
        <router-view />
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
  background-color: #fff;
}

.el-main {
  padding: 0;
  flex: 1;
  background-color: #f5f7fa;
}

.flex-grow {
  flex-grow: 1;
}

.el-menu-demo {
  border-bottom: none;
}

/* 默认隐藏更多菜单，只在移动端显示 */
.more-menu {
  display: none;
}

/* 响应式设计 - 在小屏幕时隐藏部分导航项，显示更多菜单 */
@media (max-width: 768px) {
  /* 隐藏平铺的导航项 */
  .nav-menu-item {
    display: none;
  }
  
  /* 显示更多菜单 */
  .more-menu {
    display: flex;
  }
  
  .user-name {
    display: none;
  }
}

@media (max-width: 576px) {
  .el-menu-item span:not(.user-name) {
    font-size: 14px;
  }
}
</style>
