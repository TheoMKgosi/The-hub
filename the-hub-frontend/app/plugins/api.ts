export default defineNuxtPlugin((nuxtApp) => {
  const auth = useAuthStore()
  const runtimeConfig = useRuntimeConfig()

  const api = $fetch.create({
    baseURL: runtimeConfig.public.apiBase,
    async onRequest({ options }) {
      // Get current access token (will refresh if needed)
      const token = await auth.getAccessToken()
      if (token) {
        options.headers.set('Authorization', `Bearer ${token}`)
      }
    },

    async onRequestError({ error, request, options }) {
      // Handle network errors (Failed to fetch, etc.)
      if (error.message?.includes('Failed to fetch') || error.message?.includes('NetworkError') || error.message?.includes('ERR_NETWORK')) {
        const isOnline = navigator.onLine
        if (!isOnline) {
          throw new Error('No internet connection. Please check your network and try again.')
        } else {
          throw new Error('Unable to connect to the server. Please check your connection and try again.')
        }
      }

      // Handle timeout errors
      if (error.message?.includes('timeout') || error.message?.includes('TimeoutError')) {
        throw new Error('Request timed out. Please check your connection and try again.')
      }

      // Handle CORS errors
      if (error.message?.includes('CORS') || error.message?.includes('Access-Control')) {
        throw new Error('Connection blocked. Please try again or contact support if the problem persists.')
      }

      // Re-throw other errors
      throw error
    },

    async onResponseError({ response, request, options }) {
      // Handle authentication errors
      if (response?.status === 401) {
        // Try to refresh token if we have a refresh token
        if (auth.refreshToken && !auth.isRefreshing) {
          try {
            await auth.refreshAccessToken()
            // Retry the original request with new token
            const newToken = await auth.getAccessToken()
            if (newToken) {
              options.headers.set('Authorization', `Bearer ${newToken}`)
              // Retry the request
              return $fetch(request, options)
            }
          } catch (refreshError) {
            console.warn('Token refresh failed:', refreshError)
            // Fall through to logout
          }
        }

        // If refresh failed or no refresh token, logout
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

