export const useMyFetch = (url: string, options = {}) => {
  const auth = useAuthStore()

  return useFetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: auth.token ? `Bearer ${auth.token}` : ''
    }
  })
}

