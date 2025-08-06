export default defineNuxtRouteMiddleware((to, from) => {
  const auth = useAuthStore()
  if (auth.user) return navigateTo('/dashboard')
})

