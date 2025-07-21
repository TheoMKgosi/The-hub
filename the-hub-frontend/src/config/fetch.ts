import { createFetch } from "@vueuse/core"
import { useRouter } from "vue-router"

const router = useRouter()

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
    async onFetchError(ctx){
      if(ctx.response?.status === 401 ){
        localStorage.removeItem('token')
        router.push('/login')
      }
      return ctx
    }
  },
})
