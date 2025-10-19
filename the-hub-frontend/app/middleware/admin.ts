export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()
  
  // Check if user is authenticated
  if (!auth.isLoggedIn) {
    return navigateTo('/login')
  }
  
  // Check if user has admin role
  if (auth.user?.role !== 'admin') {
    throw createError({
      statusCode: 403,
      statusMessage: 'Access denied. Admin privileges required.'
    })
  }
})