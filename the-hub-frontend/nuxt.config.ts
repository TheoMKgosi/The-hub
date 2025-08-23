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

  // Vitest configuration
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./test/setup.ts']
  },

  css: ['~/assets/css/main.css'],
  vite: {
    plugins: [
      tailwindcss()
    ]
  },

  $development: {
    runtimeConfig: {
      public: {
        apiBase: 'http://localhost:8080'
      }
    },
  },

  $production: {
    runtimeConfig: {
      public: {
        apiBase: process.env.API_BASE
      }
    },
  },

  sentry: {
    sourceMapsUploadOptions: {
      org: 'theo-kgosiemang',
      project: 'javascript-nuxt'
    }
  },

  sourcemap: {
    client: 'hidden'
  },
})
