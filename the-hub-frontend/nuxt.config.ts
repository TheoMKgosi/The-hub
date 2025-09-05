// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: false,
  modules: [
    '@pinia/nuxt',
    'pinia-plugin-persistedstate',
    '@sentry/nuxt/module',
    '@vite-pwa/nuxt'
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

  pwa: {
    manifest: {
      name: 'The Hub - Personal Learning Platform',
      short_name: 'The Hub',
      description: 'Your personal learning management platform',
      theme_color: '#3B82F6',
      background_color: '#FFFFFF',
      display: 'standalone',
      orientation: 'portrait',
      scope: '/',
      start_url: '/',
      icons: [
        {
          src: '/icon-192x192.png',
          sizes: '192x192',
          type: 'image/png'
        },
        {
          src: '/icon-512x512.png',
          sizes: '512x512',
          type: 'image/png'
        }
      ]
    },
    workbox: {
      runtimeCaching: [
        {
          urlPattern: '^https://.*',
          handler: 'NetworkFirst',
          options: {
            cacheName: 'api-cache',
            expiration: {
              maxEntries: 100,
              maxAgeSeconds: 60 * 60 * 24 // 24 hours
            }
          }
        }
      ]
    }
  },
})
