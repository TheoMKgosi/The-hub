import { ref, computed } from 'vue'

interface PendingOperation {
  id: string
  url: string
  method: string
  headers?: Record<string, string>
  body?: string
  timestamp: number
  store: string
  operation: string
}

interface CachedData {
  key: string
  data: any
  timestamp: number
  expiresAt?: number
}

export const useOffline = () => {
  const isOnline = ref(navigator.onLine)
  const isRegistered = ref(false)
  const dbVersion = 2

  // Check online status
  const updateOnlineStatus = () => {
    isOnline.value = navigator.onLine
  }

  // Listen for online/offline events
  if (process.client) {
    window.addEventListener('online', updateOnlineStatus)
    window.addEventListener('offline', updateOnlineStatus)
  }

  // Open IndexedDB
  const openDB = (): Promise<IDBDatabase> => {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open('TheHubOfflineDB', dbVersion)

      request.onerror = () => reject(request.error)
      request.onsuccess = () => resolve(request.result)

      request.onupgradeneeded = (event) => {
        const db = event.target.result as IDBDatabase

        // Create stores for different data types
        if (!db.objectStoreNames.contains('pendingOperations')) {
          const pendingStore = db.createObjectStore('pendingOperations', { keyPath: 'id' })
          pendingStore.createIndex('store', 'store', { unique: false })
          pendingStore.createIndex('timestamp', 'timestamp', { unique: false })
        }

        if (!db.objectStoreNames.contains('cachedData')) {
          const cacheStore = db.createObjectStore('cachedData', { keyPath: 'key' })
          cacheStore.createIndex('timestamp', 'timestamp', { unique: false })
        }

        // Create stores for each data type
        const stores = [
          'tasks', 'goals', 'categories', 'budgets', 'incomes', 'transactions',
          'learningPaths', 'decks', 'cards', 'resources', 'schedule', 'tags',
          'topics', 'studySessions', 'taskLearning'
        ]

        stores.forEach(storeName => {
          if (!db.objectStoreNames.contains(storeName)) {
            // Use appropriate key path based on data structure
            const keyPathMap: Record<string, string> = {
              tasks: 'task_id',
              goals: 'goal_id',
              categories: 'budget_category_id',
              budgets: 'budget_id',
              incomes: 'income_id',
              transactions: 'transaction_id',
              learningPaths: 'learning_path_id',
              decks: 'deck_id',
              cards: 'card_id',
              resources: 'id',
              schedule: 'id',
              tags: 'tag_id',
              topics: 'topic_id',
              studySessions: 'id',
              taskLearning: 'task_learning_id'
            }
            const keyPath = keyPathMap[storeName] || 'id'
            db.createObjectStore(storeName, { keyPath })
          }
        })
      }
    })
  }

  // Cache data for offline use
  const cacheData = async (store: string, key: string, data: any, ttl?: number): Promise<void> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['cachedData'], 'readwrite')
      const cacheStore = transaction.objectStore('cachedData')

      const cachedData: CachedData = {
        key: `${store}:${key}`,
        data,
        timestamp: Date.now(),
        expiresAt: ttl ? Date.now() + ttl : undefined
      }

      await new Promise((resolve, reject) => {
        const request = cacheStore.put(cachedData)
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error)
      })
    } catch (error) {
      console.error('Failed to cache data:', error)
    }
  }

  // Get cached data
  const getCachedData = async (store: string, key: string): Promise<any | null> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['cachedData'], 'readonly')
      const cacheStore = transaction.objectStore('cachedData')

      const data = await new Promise<CachedData | null>((resolve, reject) => {
        const request = cacheStore.get(`${store}:${key}`)
        request.onsuccess = () => resolve(request.result || null)
        request.onerror = () => reject(request.error)
      })

      if (!data) return null

      // Check if data has expired
      if (data.expiresAt && Date.now() > data.expiresAt) {
        // Remove expired data
        await removeCachedData(store, key)
        return null
      }

      return data.data
    } catch (error) {
      console.error('Failed to get cached data:', error)
      return null
    }
  }

  // Remove cached data
  const removeCachedData = async (store: string, key: string): Promise<void> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['cachedData'], 'readwrite')
      const cacheStore = transaction.objectStore('cachedData')

      await new Promise((resolve, reject) => {
        const request = cacheStore.delete(`${store}:${key}`)
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error)
      })
    } catch (error) {
      console.error('Failed to remove cached data:', error)
    }
  }

  // Add pending operation for background sync
  const addPendingOperation = async (operation: Omit<PendingOperation, 'id' | 'timestamp'>): Promise<void> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['pendingOperations'], 'readwrite')
      const store = transaction.objectStore('pendingOperations')

      const pendingOp: PendingOperation = {
        ...operation,
        id: `${operation.store}-${operation.operation}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: Date.now()
      }

      await new Promise((resolve, reject) => {
        const request = store.add(pendingOp)
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error)
      })

      // Notify service worker
      if ('serviceWorker' in navigator && navigator.serviceWorker.controller) {
        navigator.serviceWorker.controller.postMessage({
          type: 'ADD_PENDING_OPERATION',
          operation: pendingOp
        })
      }
    } catch (error) {
      console.error('Failed to add pending operation:', error)
    }
  }

  // Get pending operations for a store
  const getPendingOperations = async (store?: string): Promise<PendingOperation[]> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['pendingOperations'], 'readonly')
      const opStore = transaction.objectStore('pendingOperations')

      const operations = await new Promise<PendingOperation[]>((resolve, reject) => {
        let request: IDBRequest
        if (store) {
          const index = opStore.index('store')
          request = index.getAll(store)
        } else {
          request = opStore.getAll()
        }

        request.onsuccess = () => resolve(request.result || [])
        request.onerror = () => reject(request.error)
      })

      return operations.sort((a, b) => a.timestamp - b.timestamp)
    } catch (error) {
      console.error('Failed to get pending operations:', error)
      return []
    }
  }

  // Remove pending operation
  const removePendingOperation = async (id: string): Promise<void> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['pendingOperations'], 'readwrite')
      const store = transaction.objectStore('pendingOperations')

      await new Promise((resolve, reject) => {
        const request = store.delete(id)
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error)
      })
    } catch (error) {
      console.error('Failed to remove pending operation:', error)
    }
  }

  // Store data locally for offline access
  const storeDataLocally = async (storeName: string, data: any[]): Promise<void> => {
    try {
      const db = await openDB()
      const transaction = db.transaction([storeName], 'readwrite')
      const store = transaction.objectStore(storeName)

      // Clear existing data
      await new Promise((resolve, reject) => {
        const request = store.clear()
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error)
      })

      // Add new data
      for (const item of data) {
        await new Promise((resolve, reject) => {
          const request = store.add(item)
          request.onsuccess = () => resolve(request.result)
          request.onerror = () => reject(request.error)
        })
      }
    } catch (error) {
      console.error('Failed to store data locally:', error)
    }
  }

  // Get data from local storage
  const getLocalData = async (storeName: string): Promise<any[]> => {
    try {
      const db = await openDB()
      const transaction = db.transaction([storeName], 'readonly')
      const store = transaction.objectStore(storeName)

      return await new Promise<any[]>((resolve, reject) => {
        const request = store.getAll()
        request.onsuccess = () => resolve(request.result || [])
        request.onerror = () => reject(request.error)
      })
    } catch (error) {
      console.error('Failed to get local data:', error)
      return []
    }
  }

  // Sync pending operations when back online
  const syncPendingOperations = async (): Promise<void> => {
    if (!isOnline.value) return

    try {
      const pendingOps = await getPendingOperations()

      for (const op of pendingOps) {
        try {
          const response = await fetch(op.url, {
            method: op.method,
            headers: {
              'Content-Type': 'application/json',
              ...op.headers
            },
            body: op.body
          })

          if (response.ok) {
            await removePendingOperation(op.id)
          } else {
            console.error('Failed to sync operation:', op, response.status)
          }
        } catch (error) {
          console.error('Network error during sync:', op, error)
        }
      }
    } catch (error) {
      console.error('Failed to sync pending operations:', error)
    }
  }

  // Register background sync
  const registerBackgroundSync = async (): Promise<void> => {
    if ('serviceWorker' in navigator && 'sync' in window.ServiceWorkerRegistration.prototype) {
      try {
        const registration = await navigator.serviceWorker.ready
        await registration.sync.register('background-sync')
      } catch (error) {
        console.error('Failed to register background sync:', error)
      }
    }
  }

  // Clean up expired cache entries
  const cleanupExpiredCache = async (): Promise<void> => {
    try {
      const db = await openDB()
      const transaction = db.transaction(['cachedData'], 'readwrite')
      const store = transaction.objectStore('cachedData')
      const index = store.index('timestamp')

      const now = Date.now()
      const range = IDBKeyRange.upperBound(now)

      const expiredEntries = await new Promise<CachedData[]>((resolve, reject) => {
        const request = index.openCursor(range)
        const results: CachedData[] = []

        request.onsuccess = (event) => {
          const cursor = event.target.result
          if (cursor) {
            const data = cursor.value
            if (data.expiresAt && data.expiresAt < now) {
              results.push(data)
            }
            cursor.continue()
          } else {
            resolve(results)
          }
        }
        request.onerror = () => reject(request.error)
      })

      // Remove expired entries
      for (const entry of expiredEntries) {
        await new Promise((resolve, reject) => {
          const request = store.delete(entry.key)
          request.onsuccess = () => resolve(request.result)
          request.onerror = () => reject(request.error)
        })
      }
    } catch (error) {
      console.error('Failed to cleanup expired cache:', error)
    }
  }

  return {
    isOnline: computed(() => isOnline.value),
    isRegistered: computed(() => isRegistered.value),
    cacheData,
    getCachedData,
    removeCachedData,
    addPendingOperation,
    getPendingOperations,
    removePendingOperation,
    storeDataLocally,
    getLocalData,
    syncPendingOperations,
    registerBackgroundSync,
    cleanupExpiredCache
  }
}