import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/pages/dashboard.vue'
import LoginView from '@/pages/login.vue'

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
      component: () => import('../pages/register.vue'),
    },
    {
      path: '/plan',
      name: 'plan',
      component: () => import('../pages/plan.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/time',
      name: 'time',
      component: () => import('../pages/time.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/learning',
      name: 'learning',
      component: () => import('../pages/learning.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/learning/cards/:deck_id',
      name: 'cards',
      component: () => import('../pages/cards.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/finance',
      name: 'finance',
      component: () => import('../pages/finace.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../pages/setting.vue'),
      meta: { requiresAuth: true }
    },
  ],
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    // not logged in, redirect to login
    next({ name: 'login' })
  } else if ((to.name === 'login' || to.name === 'register') && token) {
    // already logged in, redirect to dashboard
    next({ name: 'dashboard' })
  } else {
    next()
  }
})


export default router
