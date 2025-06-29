import { createFetch } from "@vueuse/core"

export const useMyFetch = createFetch({
  baseUrl: import.meta.env.VITE_BASE_URL,
  options: {
    async beforeFetch({ options }) {
      const token = localStorage.getItem('token')
      if (token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${token}`,
        }
      }
      return { options }
    },
  },
})
