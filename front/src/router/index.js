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
    // 员工端路由
    {
      path: '/employee/schedule',
      name: 'employee-schedule',
      component: () => import('../views/employee/Schedule.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    {
      path: '/employee/evaluation',
      name: 'employee-evaluation',
      component: () => import('../views/employee/Evaluation.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    {
      path: '/employee/scores',
      name: 'employee-scores',
      component: () => import('../views/employee/Scores.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    // 讲师端路由
    {
      path: '/teacher/schedule',
      name: 'teacher-schedule',
      component: () => import('../views/teacher/Schedule.vue'),
      meta: { requiresAuth: true, role: 'teacher' }
    },
    {
      path: '/teacher/grading',
      name: 'teacher-grading',
      component: () => import('../views/teacher/Grading.vue'),
      meta: { requiresAuth: true, role: 'teacher' }
    },
    {
      path: '/teacher/scores',
      name: 'teacher-scores',
      component: () => import('../views/teacher/Scores.vue'),
      meta: { requiresAuth: true, role: 'teacher' }
    },
    // 培训大纲制定者路由
    {
      path: '/planner/plans',
      name: 'planner-plans',
      component: () => import('../views/planner/Plans.vue'),
      meta: { requiresAuth: true, role: 'planner' }
    },
    {
      path: '/planner/courses',
      name: 'planner-courses',
      component: () => import('../views/planner/Courses.vue'),
      meta: { requiresAuth: true, role: 'planner' }
    },
    {
      path: '/planner/analytics',
      name: 'planner-analytics',
      component: () => import('../views/planner/Analytics.vue'),
      meta: { requiresAuth: true, role: 'planner' }
    }
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
