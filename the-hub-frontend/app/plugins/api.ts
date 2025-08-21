export default defineNuxtPlugin((nuxtApp) => {
  const auth = useAuthStore()
  const runtimeConfig = useRuntimeConfig()

  const api = $fetch.create({
    baseURL: runtimeConfig.public.apiBase,
    onRequest({ options }) {
      if (auth.token) {
        options.headers.set('Authorization', `Bearer ${auth.token}`)
      }
    },

    async onRequestError({ response }) {
      if (response?.status === 401) {
        await nuxtApp.runWithContext(() => navigateTo('/login'))
      }
    },
    async onResponseError({ response }) {
      if (response?.status === 401) {
        await nuxtApp.runWithContext(() => navigateTo('/login'))
      }
    },
  })

  return {
    provide: {
      api
    }
  }
})

