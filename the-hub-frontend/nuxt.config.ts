// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  devtools: { enabled: true, },
  storybook: {
    enabled: false,
  },
  compatibilityDate: '2025-07-15',
  modules: [
    '@pinia/nuxt',
    'pinia-plugin-persistedstate',
    '@vite-pwa/nuxt',
    '@nuxtjs/storybook',
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

  css: ['./app/assets/css/main.css'],

  vite: {
    plugins: [
      tailwindcss()
    ]
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE || 'http://localhost:8080',
      vapidPublicKey: process.env.NUXT_PUBLIC_VAPID_PUBLIC_KEY,
      posthogPublicKey: 'phc_Q4MgrAFa89TvsLS1z8QIAh1lCr2ScGShzmBJWNahlmx',
      posthogHost: 'https://us.i.posthog.com',
      posthogDefaults: '2025-11-30'
    }
  },
  
  hooks: {
    'imports:sources': (sources) => {
      // Find and remove the Storybook duplicate if it's interfering
      const sbIndex = sources.findIndex(s => s.from?.includes('@storybook-vue/nuxt'));
      if (sbIndex !== -1) {
        // You can filter specific imports here if needed
      }
    }
  },

  sourcemap: {
    client: 'hidden'
  },

  pwa: {
    devOptions: {
      enabled: false
    },
    registerType: 'autoUpdate',
    srcDir: 'public',
    filename: 'sw.js',
    strategies: 'generateSW',
    workbox: {
      // globPatterns: ['**/*.{js,css,html,png,svg,ico}'],
      navigateFallback: '/offline',
      runtimeCaching: [
        {
          urlPattern: ({ url }) => url.origin !== location.origin,
          handler: 'NetworkOnly', // Don't cache external API calls
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
