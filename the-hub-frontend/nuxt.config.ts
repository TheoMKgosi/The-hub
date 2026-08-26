// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  devtools: { enabled: true, },
  modules: [
    '@pinia/nuxt',
    'pinia-plugin-persistedstate',
    // '@vite-pwa/nuxt',
    '@nuxt/image',
    '@nuxt/ui',
    'dayjs-nuxt',
    // '@nuxt/hints',
  ],

  components: [
    {
      path: '~/components',
      extensions: ['vue'],
      pathPrefix: false,
    }
  ],

  ssr: false,

  experimental: {
    appManifest: false
  },

  css: ['~/assets/css/main.css'],

  vite: {
    plugins: [
      tailwindcss()
    ],
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE || 'http://localhost:8080',
      vapidPublicKey: process.env.NUXT_PUBLIC_VAPID_PUBLIC_KEY,
    }
  },

  nitro: {
    devProxy: {
      '/uploads': {
        target: 'http://localhost:8080/uploads',
        changeOrigin: true
      }
    }
  },

  sourcemap: {
    client: 'hidden'
  },

  dayjs: {
    plugins: ['relativeTime']
  },

  // pwa: {
  //   devOptions: {
  //     enabled: false
  //   },
  //   registerType: 'autoUpdate',
  //   srcDir: 'public',
  //   filename: 'sw.js',
  //   strategies: 'generateSW',
  //   workbox: {
  //     // globPatterns: ['**/*.{js,css,html,png,svg,ico}'],
  //     navigateFallback: '/offline',
  //     runtimeCaching: [
  //       {
  //         urlPattern: ({ url }) => url.origin !== location.origin,
  //         handler: 'NetworkOnly', // Don't cache external API calls
  //       },
  //       {
  //         urlPattern: '^/.*',
  //         handler: 'StaleWhileRevalidate',
  //         options: {
  //           cacheName: 'pages-cache',
  //           expiration: {
  //             maxEntries: 50,
  //             maxAgeSeconds: 60 * 60 * 24 * 7 // 7 days
  //           }
  //         }
  //       }
  //     ]
  //   },
  //   client: {
  //     installPrompt: true,
  //     periodicSyncForUpdates: 20
  //   }
  // },
})
