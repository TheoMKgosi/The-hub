// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: false,
  modules: [
    '@pinia/nuxt',
    'pinia-plugin-persistedstate',
    '@sentry/nuxt/module'
  ],

  css: ['~/assets/css/main.css'],
  vite: {
    plugins: [
      tailwindcss()
    ]
  },

  runtimeConfig: {
    public: {
      apiBase: 'http://localhost:8080'
    }
  },

  sentry: {
    sourceMapsUploadOptions: {
      org: 'theo-kgosiemang',
      project: 'javascript-nuxt'
    }
  },

  sourcemap: {
    client: 'hidden'
  }
})
