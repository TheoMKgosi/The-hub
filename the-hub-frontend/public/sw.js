// Custom service worker for The Hub PWA
// This extends the auto-generated workbox service worker

// Import workbox from CDN (will be available when workbox is loaded)
importScripts('https://storage.googleapis.com/workbox-cdn/releases/6.5.4/workbox-sw.js');

// Only use workbox if it's available
if (workbox) {
  // Set workbox to debug mode in development
  workbox.setConfig({
    debug: false
  });

  // Cache the app shell
  workbox.routing.registerRoute(
    ({ request }) => request.destination === 'document',
    new workbox.strategies.NetworkFirst({
      cacheName: 'pages-cache',
      plugins: [
        new workbox.expiration.ExpirationPlugin({
          maxEntries: 50,
          maxAgeSeconds: 7 * 24 * 60 * 60, // 7 days
        }),
      ],
    })
  );

  // Cache static assets
  workbox.routing.registerRoute(
    ({ request }) => request.destination === 'script' ||
                     request.destination === 'style' ||
                     request.destination === 'image' ||
                     request.destination === 'font',
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

  // Cache API responses with network-first strategy
  workbox.routing.registerRoute(
    ({ url }) => url.pathname.startsWith('/api/') || url.hostname.includes('localhost'),
    new workbox.strategies.NetworkFirst({
      cacheName: 'api-cache',
      plugins: [
        new workbox.expiration.ExpirationPlugin({
          maxEntries: 100,
          maxAgeSeconds: 24 * 60 * 60, // 24 hours
        }),
        new workbox.cacheableResponse.CacheableResponsePlugin({
          statuses: [0, 200],
        }),
      ],
    })
  );

  // Handle background sync for offline operations
  if ('serviceWorker' in navigator && 'sync' in window.ServiceWorkerRegistration.prototype) {
    self.addEventListener('sync', function(event) {
      if (event.tag === 'background-sync') {
        event.waitUntil(doBackgroundSync());
      }
    });
  }

  // Background sync function
  async function doBackgroundSync() {
    try {
      // Get pending operations from IndexedDB
      const pendingOps = await getPendingOperations();

      for (const op of pendingOps) {
        try {
          await syncOperation(op);
          await removePendingOperation(op.id);
        } catch (error) {
          console.error('Failed to sync operation:', op, error);
        }
      }
    } catch (error) {
      console.error('Background sync failed:', error);
    }
  }

  // IndexedDB helper functions for offline operations
  function openDB() {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open('TheHubOfflineDB', 1);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result);

      request.onupgradeneeded = (event) => {
        const db = event.target.result;

        // Create stores for different data types
        if (!db.objectStoreNames.contains('pendingOperations')) {
          db.createObjectStore('pendingOperations', { keyPath: 'id' });
        }

        if (!db.objectStoreNames.contains('cachedData')) {
          db.createObjectStore('cachedData', { keyPath: 'key' });
        }
      };
    });
  }

  async function getPendingOperations() {
    const db = await openDB();
    const transaction = db.transaction(['pendingOperations'], 'readonly');
    const store = transaction.objectStore('pendingOperations');

    return new Promise((resolve, reject) => {
      const request = store.getAll();
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }

  async function removePendingOperation(id) {
    const db = await openDB();
    const transaction = db.transaction(['pendingOperations'], 'readwrite');
    const store = transaction.objectStore('pendingOperations');

    return new Promise((resolve, reject) => {
      const request = store.delete(id);
      request.onsuccess = () => resolve(true);
      request.onerror = () => reject(request.error);
    });
  }

  async function syncOperation(operation) {
    // This would be called when online to sync pending operations
    const response = await fetch(operation.url, {
      method: operation.method,
      headers: operation.headers,
      body: operation.body
    });

    if (!response.ok) {
      throw new Error(`Sync failed: ${response.status}`);
    }

    return response.json();
  }

  // Listen for messages from the main thread
  self.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'SKIP_WAITING') {
      self.skipWaiting();
    }

    if (event.data && event.data.type === 'CACHE_DATA') {
      // Cache data for offline use
      cacheDataForOffline(event.data.key, event.data.data);
    }

    if (event.data && event.data.type === 'ADD_PENDING_OPERATION') {
      // Add operation to pending queue for background sync
      addPendingOperation(event.data.operation);
    }
  });

  async function cacheDataForOffline(key, data) {
    try {
      const db = await openDB();
      const transaction = db.transaction(['cachedData'], 'readwrite');
      const store = transaction.objectStore('cachedData');

      await new Promise((resolve, reject) => {
        const request = store.put({ key, data, timestamp: Date.now() });
        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject(request.error);
      });
    } catch (error) {
      console.error('Failed to cache data:', error);
    }
  }

  async function addPendingOperation(operation) {
    try {
      const db = await openDB();
      const transaction = db.transaction(['pendingOperations'], 'readwrite');
      const store = transaction.objectStore('pendingOperations');

      await new Promise((resolve, reject) => {
        const request = store.add(operation);
        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject(request.error);
      });
    } catch (error) {
      console.error('Failed to add pending operation:', error);
    }
  }

  // Handle push notifications
  self.addEventListener('push', function(event) {
    if (!event.data) return;

    const data = event.data.json();

    const options = {
      body: data.body,
      icon: '/icon-192x192.png',
      badge: '/icon-192x192.png',
      vibrate: [100, 50, 100],
      data: {
        dateOfArrival: Date.now(),
        primaryKey: data.id
      },
      actions: [
        {
          action: 'view',
          title: 'View',
          icon: '/icon-192x192.png'
        },
        {
          action: 'dismiss',
          title: 'Dismiss'
        }
      ]
    };

    event.waitUntil(
      self.registration.showNotification(data.title, options)
    );
  });

  // Handle notification clicks
  self.addEventListener('notificationclick', function(event) {
    event.notification.close();

    if (event.action === 'view') {
      // Open the app
      event.waitUntil(
        clients.openWindow('/')
      );
    }
  });

  // Clean up old caches on activation
  self.addEventListener('activate', (event) => {
    event.waitUntil(
      caches.keys().then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (cacheName !== 'pages-cache' && cacheName !== 'assets-cache' && cacheName !== 'api-cache') {
              return caches.delete(cacheName);
            }
          })
        );
      })
    );
  });

} else {
  console.log('Workbox could not be loaded. No offline functionality available.');
}