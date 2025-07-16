import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
    },
    {
      path: '/plan',
      name: 'plan',
      // route level code-splitting
      // this generates a separate chunk (Goal.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/PlanView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/time',
      name: 'time',
      component: () => import('../views/TimeView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/learning',
      name: 'learning',
      component: () => import('../views/LearningView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/learning/cards/:deck_id',
      name: 'cards',
      component: () => import('../views/CardsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/finance',
      name: 'finance',
      component: () => import('../views/FinanceView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue'),
      meta: { requiresAuth: true }
    },
  ],
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
