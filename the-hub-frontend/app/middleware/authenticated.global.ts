const publicPages = ['/', '/login', '/register', '/forgot-password', '/privacy-policy', '/terms-of-service', '/admin/login']

export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()

  // If authenticated and trying to go to auth pages, redirect to dashboard
  if (auth.isLoggedIn && ['/login', '/register'].includes(to.path)) {
    return navigateTo('/dashboard')
  }

  // If authenticated and on admin routes
  if (auth.isLoggedIn && to.path.startsWith('/admin')) {
    if (auth.user?.role === 'admin') {
      if (to.path === '/admin/login') {
        return navigateTo('/admin')
      }
      return
    }
    // Non-admin trying to access admin routes
    return navigateTo('/login')
  }

  // If not authenticated and trying to access a private page, redirect to login
  if (!auth.isLoggedIn && !publicPages.includes(to.path)) {
    if (to.path.startsWith('/admin')) {
      return navigateTo('/admin/login')
    }
    return navigateTo('/login')
  }
})
