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

    async onResponseError({ response }) {
      // Handle authentication errors
      if (response?.status === 401) {
        await nuxtApp.runWithContext(() => auth.logout())
        throw new Error('Your session has expired. Please log in again.')
      }

      // Extract error message from response
      let errorMessage = 'An unexpected error occurred'

      try {
        const errorData = await response.json()
        if (errorData?.error) {
          errorMessage = errorData.error
        } else if (errorData?.message) {
          errorMessage = errorData.message
        }
      } catch {
        // If we can't parse the response, use status-based messages
        switch (response?.status) {
          case 400:
            errorMessage = 'Invalid request. Please check your input.'
            break
          case 403:
            errorMessage = 'You do not have permission to perform this action.'
            break
          case 404:
            errorMessage = 'The requested resource was not found.'
            break
          case 409:
            errorMessage = 'This action conflicts with existing data.'
            break
          case 422:
            errorMessage = 'Validation failed. Please check your input.'
            break
          case 500:
            errorMessage = 'Server error. Please try again later.'
            break
          default:
            errorMessage = `Request failed with status ${response?.status}`
        }
      }

      // Create and throw a new error with the extracted message
      const error = new Error(errorMessage)
      error.statusCode = response?.status
      throw error
    },
  })

  return {
    provide: {
      api
    }
  }
})

