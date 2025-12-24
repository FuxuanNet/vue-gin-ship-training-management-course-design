import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/Register.vue')
    },
    {
      path: '/upload',
      name: 'upload',
      component: () => import('../views/ResourceUpload.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/publish',
      name: 'publish',
      component: () => import('../views/PublishResource.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/market',
      name: 'market',
      component: () => import('../views/Market.vue')
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/UserProfile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/sandbox',
      name: 'sandbox',
      component: () => import('../views/Sandbox.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/logs',
      name: 'logs',
      component: () => import('../views/BlockchainLogs.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// 路由导航守卫
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router 