export default defineNuxtRouteMiddleware((to, from) => {
  const auth = useAuthStore()

  const publicPages = ['/login', '/register']

  // If authenticated and trying to go to login/register, redirect to dashboard
  if (auth.token && publicPages.includes(to.path)) {
    return navigateTo('/dashboard')
  }

  // If not authenticated and trying to access a private page, redirect to login
  if (!auth.token && !publicPages.includes(to.path)) {
    return navigateTo('/login')
  }
})

