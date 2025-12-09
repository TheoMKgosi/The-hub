export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()
  
  // If user is already logged in, redirect based on role
  if (auth.isLoggedIn) {
    if (auth.user?.role === 'admin' && to.path.startsWith('/admin')) {
      // Allow admin to access admin routes
      return
    } else if (auth.user?.role === 'admin' && to.path === '/admin/login') {
      // Admin trying to access admin login, redirect to admin dashboard
      return navigateTo('/admin')
    } else if (auth.user?.role !== 'admin' && to.path.startsWith('/admin')) {
      // Non-admin trying to access admin routes, redirect to regular login
      return navigateTo('/login')
    } else if (to.path === '/login' || to.path === '/register') {
      // Regular user trying to access auth pages, redirect to dashboard
      return navigateTo('/dashboard')
    }
  }
  
  // Allow access to login, register, admin login, home page, and legal pages for non-authenticated users
  const publicRoutes = ['/', '/login', '/register', '/admin/login', '/privacy-policy', '/terms-of-service']
  if (!publicRoutes.includes(to.path) && !auth.isLoggedIn) {
    // For admin routes, redirect to admin login
    if (to.path.startsWith('/admin')) {
      return navigateTo('/admin/login')
    }
    // For other routes, redirect to regular login
    return navigateTo('/login')
  }
})