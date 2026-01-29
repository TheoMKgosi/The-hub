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
    '@nuxt/image'
  ],
  components: [
    {
      path: '~/components',
      global: true,
      extensions: ['vue'],
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

  nitro: {
    devProxy: {
      '/uploads': {
        target: 'http://localhost:8080/uploads',
        changeOrigin: true
      }
    }
  },

  // Production Image Serving Setup:
  // The app uses the API_BASE environment variable to construct image URLs
  // Make sure to set API_BASE in production to point to your API server
  // Example: API_BASE=https://your-api-domain.com
  //
  // Alternative production setups:
  // 1. Copy images to public/uploads during build process
  // 2. Use a CDN service (Cloudflare, AWS S3, etc.)
  // 3. Set up nginx/apache to serve images from the same domain



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
