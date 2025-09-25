// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
   modules: [
     '@pinia/nuxt',
     'pinia-plugin-persistedstate',
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

      }
    },
  },

  $production: {
    runtimeConfig: {
      public: {
        apiBase: process.env.API_BASE,
        vapidPublicKey: process.env.NUXT_PUBLIC_VAPID_PUBLIC_KEY,

      }
    },
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
