export default defineNuxtRouteMiddleware((to, from) => {
  const auth = useAuthStore()

  const publicPages = ['/login', '/register', '/forgot-password']

  // If authenticated and trying to go to login/register, redirect to dashboard
  if (auth.isLoggedIn && publicPages.includes(to.path)) {
    return navigateTo('/dashboard')
  }

  // If not authenticated and trying to access a private page, redirect to login
  if (!auth.isLoggedIn && !publicPages.includes(to.path)) {
    return navigateTo('/login')
  }
})

