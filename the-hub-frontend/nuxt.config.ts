// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  modules: [
    '@pinia/nuxt',
    'pinia-plugin-persistedstate',
    '@sentry/nuxt/module',
    '@vite-pwa/nuxt'
  ],
  ssr: false,

  css: ['~/assets/css/main.css'],
  vite: {
    plugins: [
      tailwindcss()
    ]
  },

  $development: {
    runtimeConfig: {
      public: {
        apiBase: 'http://localhost:8080',
        vapidPublicKey: process.env.NUXT_PUBLIC_VAPID_PUBLIC_KEY,
        sentry: {
          dsn: process.env.NUXT_PUBLIC_SENTRY_DSN || 'https://0f93f48bf9daff70a1730cd729955dc0@o4509804910936064.ingest.de.sentry.io/4509804913557584'
        }
      }
    },
  },

  $production: {
    runtimeConfig: {
      public: {
        apiBase: process.env.API_BASE,
        vapidPublicKey: process.env.NUXT_PUBLIC_VAPID_PUBLIC_KEY,
        sentry: {
          dsn: process.env.NUXT_PUBLIC_SENTRY_DSN
        }
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
     registerType: 'autoUpdate',
     srcDir: 'public',
     filename: 'sw.js',
     strategies: 'injectManifest',
    workbox: {
      globPatterns: ['**/*.{js,css,html,png,svg,ico}'],
      navigateFallback: '/offline',
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
        },
        {
          urlPattern: '^/.*',
          handler: 'StaleWhileRevalidate',
          options: {
            cacheName: 'pages-cache',
            expiration: {
              maxEntries: 50,
              maxAgeSeconds: 60 * 60 * 24 * 7 // 7 days
            }
          }
        }
      ]
    },
    client: {
      installPrompt: true,
      periodicSyncForUpdates: 20
    }
  },
})
