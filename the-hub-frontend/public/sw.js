// Custom service worker for The Hub PWA
// This file will be enhanced by Workbox injectManifest

// Import Workbox
importScripts('https://storage.googleapis.com/workbox-cdn/releases/7.2.0/workbox-sw.js');

// Skip waiting and claim clients
self.skipWaiting();
self.clients.claim();

// Basic caching for static assets
workbox.routing.registerRoute(
  /\.(?:png|jpg|jpeg|svg|gif|ico|css|js)$/,
  new workbox.strategies.StaleWhileRevalidate({
    cacheName: 'assets-cache',
    plugins: [
      new workbox.expiration.ExpirationPlugin({
        maxEntries: 100,
        maxAgeSeconds: 30 * 24 * 60 * 60, // 30 days
      }),
    ],
  })
);

// Network-first for API calls (but only for same-origin)
workbox.routing.registerRoute(
  ({ url }) => url.origin === self.location.origin && url.pathname.startsWith('/api/'),
  new workbox.strategies.NetworkFirst({
    cacheName: 'api-cache',
    plugins: [
      new workbox.expiration.ExpirationPlugin({
        maxEntries: 50,
        maxAgeSeconds: 5 * 60, // 5 minutes
      }),
    ],
  })
);

// Stale-while-revalidate for pages
workbox.routing.registerRoute(
  ({ request }) => request.destination === 'document',
  new workbox.strategies.StaleWhileRevalidate({
    cacheName: 'pages-cache',
    plugins: [
      new workbox.expiration.ExpirationPlugin({
        maxEntries: 20,
        maxAgeSeconds: 24 * 60 * 60, // 24 hours
      }),
    ],
  })
);