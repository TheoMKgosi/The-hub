<script setup lang="ts">
const receiptStore = useReceiptStore()
const categoryStore = useCategoryStore()
const { addToast } = useToast()

const activeReceiptId = ref<string | null>(null)
const imageViewer = reactive({
  show: false,
  receipt: null as any
})
const showDialog = ref(false)
const showReceiptModal = ref(true)
const searchQuery = ref('')
const receiptToDelete = ref<string>('')

const formData = reactive({
  title: '',
  image_data: '',
  amount: 0,
  date: new Date().toISOString().split('T')[0], // Auto-set to today's date
  category_id: ''
})

const filteredReceipts = computed(() => {
  let result = receiptStore.receipts

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(receipt =>
      receipt.title.toLowerCase().includes(query) ||
      receipt.Category?.name.toLowerCase().includes(query)
    )
  }

  return result
})

const isEditing = computed(() => !!activeReceiptId.value)

const groupedReceipts = computed(() => {
  const grouped = new Map()

  filteredReceipts.value.forEach(receipt => {
    // Use created_at date to determine the folder
    const date = new Date(receipt.created_at)
    const yearMonth = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`

    if (!grouped.has(yearMonth)) {
      grouped.set(yearMonth, [])
    }
    grouped.get(yearMonth).push(receipt)
  })

  // Sort folders by date (newest first)
  return new Map([...grouped.entries()].sort(([a], [b]) => b.localeCompare(a)))
})

const expandedFolders = ref(new Set()) // Will be set when receipts are loaded

// Set default expanded folder to current month when receipts are loaded
watch(() => receiptStore.receipts.length, (newLength) => {
  if (newLength > 0 && expandedFolders.value.size === 0) {
    // Expand the most recent folder (first in the sorted map)
    const firstFolder = Array.from(groupedReceipts.value.keys())[0]
    if (firstFolder) {
      expandedFolders.value.add(firstFolder)
    }
  }
})

const toggleFolder = (folderKey) => {
  if (expandedFolders.value.has(folderKey)) {
    expandedFolders.value.delete(folderKey)
  } else {
    expandedFolders.value.add(folderKey)
  }
}

const formatFolderName = (yearMonth) => {
  const [year, month] = yearMonth.split('-')
  const date = new Date(parseInt(year), parseInt(month) - 1, 1)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })
}

const deleteItem = (id: string) => {
  receiptStore.deleteReceipt(id)
}

const submitForm = async () => {
  const dataToSend = {
    ...formData,
    amount: formData.amount || undefined,
    date: formData.date || undefined,
    category_id: formData.category_id || undefined
  }

  if (activeReceiptId.value) {
    // Update existing receipt
    await receiptStore.updateReceipt(activeReceiptId.value, dataToSend)
  } else {
    // Create new receipt
    await receiptStore.submitForm(dataToSend)
  }

  // Reset form
  Object.assign(formData, {
    title: '',
    image_data: '',
    amount: 0,
    date: new Date().toISOString().split('T')[0], // Reset to today's date
    category_id: ''
  })
  activeReceiptId.value = null
  showReceiptModal.value = true
  capturedImage.value = ''
  showCamera.value = false
}

const handleCategorySelect = (category) => {
  formData.category_id = category.budget_category_id
}

const handleCategoryCreate = async (categoryName) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    // The new category should now be available in the store
    // Find it and set it as selected
    const newCategory = categoryStore.categories.find(cat =>
      cat.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) {
      formData.category_id = newCategory.budget_category_id
    }
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

const formatDate = (date: string) => {
  console.log('formatDate called with:', date, typeof date)

  if (!date) {
    console.warn('Date is empty/undefined')
    return 'No Date'
  }

  try {
    const dateObj = new Date(date)

    // Check if date is valid
    if (isNaN(dateObj.getTime())) {
      console.warn('Invalid date string:', date)
      return 'Invalid Date'
    }

    const formatted = dateObj.toLocaleDateString()
    console.log('Successfully formatted date:', date, '->', formatted)
    return formatted
  } catch (error) {
    console.error('Error formatting date:', date, error)
    return 'Error'
  }
}
const formatCurrency = (amount: number) => `$${amount.toFixed(2)}`


const openForm = (id: string) => {
  activeReceiptId.value = id
}

const openImageViewer = (receipt: any) => {
  imageViewer.receipt = receipt
  imageViewer.show = true
}

const closeImageViewer = () => {
  imageViewer.show = false
  imageViewer.receipt = null
}

const openEditForm = (receipt: any) => {
  activeReceiptId.value = receipt.receipt_id
  Object.assign(formData, {
    title: receipt.title,
    image_data: '', // Can't edit image, keep empty
    amount: receipt.amount || 0,
    date: receipt.created_at ? new Date(receipt.created_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
    category_id: receipt.Category?.budget_category_id || ''
  })
  capturedImage.value = getImageUrl(receipt.image_path) // Show existing image
  showReceiptModal.value = false // Show the modal
}



const closeForm = () => {
  activeReceiptId.value = null
}

const cancelForm = () => {
  // Reset form
  Object.assign(formData, {
    title: '',
    image_data: '',
    amount: 0,
    date: new Date().toISOString().split('T')[0], // Reset to today's date
    category_id: ''
  })
  activeReceiptId.value = null
  showReceiptModal.value = true
  capturedImage.value = ''
  showCamera.value = false
}

onMounted(() => {
  if (receiptStore.receipts.length === 0) {
    receiptStore.fetchReceipts()
  }
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})

// Camera functionality
const showCamera = ref(false)
const capturedImage = ref('')
const videoRef = ref<HTMLVideoElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const stream = ref<MediaStream | null>(null)

// Camera functions
const openCamera = async () => {
  showCamera.value = true
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' } // Use back camera on mobile
    })
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
    }
  } catch (error) {
    addToast('Camera access denied or not available', 'error')
    showCamera.value = false
  }
}

const captureImage = () => {
  if (!videoRef.value || !canvasRef.value) return

  const video = videoRef.value
  const canvas = canvasRef.value
  const context = canvas.getContext('2d')

  if (!context) return

  // Set canvas size to video size
  canvas.width = video.videoWidth
  canvas.height = video.videoHeight

  // Draw video frame to canvas
  context.drawImage(video, 0, 0)

  // Convert to base64
  const imageData = canvas.toDataURL('image/jpeg', 0.8)
  capturedImage.value = imageData
  formData.image_data = imageData

  // Stop camera
  stopCamera()
  showCamera.value = false
  addToast('Receipt captured successfully', 'success')
}

const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach(track => track.stop())
    stream.value = null
  }
}

const retakePhoto = () => {
  capturedImage.value = ''
  formData.image_data = ''
  openCamera()
}

// Get the correct image URL for development and production
const getImageUrl = (imagePath) => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  // In development: images are proxied through Nuxt dev server
  // In production: images are served from the API server
  return `${apiBase}/uploads/${imagePath}`
}

// Cleanup on unmount
onUnmounted(() => {
  stopCamera()
})

// Mobile features
const isMobile = ref(false)

onMounted(() => {
  isMobile.value = window.innerWidth < 768
})
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Receipt Management</h2>

    <!-- Filters + Search -->
    <div class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search receipts..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <!-- Floating Action Button -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="showReceiptModal" @click="showReceiptModal = false" class="fixed bottom-4 right-4 cursor-pointer z-40">
          <div class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
            <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- Receipt Modal -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="!showReceiptModal" @click="showReceiptModal = true" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50">
          <div class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark" @click.stop>

            <!-- Modal Header -->
            <div class="flex items-center justify-between p-6 border-b border-surface-light dark:border-surface-dark">
              <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">{{ isEditing ? 'Edit Receipt' : 'Add New Receipt' }}</h2>
              <UiBaseButton @click="showReceiptModal = true" variant="default" size="sm" class="p-2">
                ×
              </UiBaseButton>
            </div>

            <!-- Modal Body -->
            <div class="p-6">
              <form @submit.prevent="submitForm" class="space-y-4">

                <div>
                  <label for="title" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
                    Title
                  </label>
                  <input type="text" id="title" v-model="formData.title" placeholder="e.g., Grocery Store Receipt"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" required />
                </div>

                <!-- Camera Section -->
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
                    Receipt Image
                  </label>
                  <div v-if="!capturedImage" class="space-y-2">
                    <UiBaseButton @click="openCamera" variant="secondary" size="md" class="w-full">
                      📷 Take Photo
                    </UiBaseButton>
                    <p class="text-xs text-text-light dark:text-text-dark/60 text-center">
                      Capture a photo of your receipt
                    </p>
                  </div>
                  <div v-else class="space-y-2">
                    <img :src="capturedImage" alt="Captured receipt" class="w-full h-48 object-cover rounded-md border border-surface-light dark:border-surface-dark" />
                    <div class="flex gap-2">
                      <UiBaseButton @click="retakePhoto" variant="secondary" size="sm" class="flex-1">
                        Retake
                      </UiBaseButton>
                      <UiBaseButton @click="capturedImage = ''; formData.image_data = ''" variant="danger" size="sm" class="flex-1">
                        Remove
                      </UiBaseButton>
                    </div>
                  </div>
                </div>

               <div>
                 <label for="amount" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount (Optional)</label>
                 <input type="number" id="amount" v-model="formData.amount" placeholder="0.00" step="0.01" min="0"
                   class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
               </div>

               <div>
                  <label for="date" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Date <span class="text-xs text-text-light/60 dark:text-text-dark/60">(Auto-set to today)</span></label>
                 <input type="date" id="date" v-model="formData.date"
                   class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
               </div>

               <div>
                 <label for="category" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
                   Category (Optional)
                 </label>
                 <UiBaseComboBox
                   :model-value="formData.category_id"
                   :categories="categoryStore.categories"
                   placeholder="Select or create category..."
                   @select="handleCategorySelect"
                   @create="handleCategoryCreate"
                 />
               </div>

                <!-- Modal Footer -->
                <div class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
                   <UiBaseButton type="button" @click="cancelForm" variant="default" size="md" class="w-full sm:w-auto">
                     Cancel
                   </UiBaseButton>
                   <UiBaseButton type="submit" :disabled="!formData.title || (!capturedImage && !isEditing)" variant="primary" size="md" class="w-full sm:w-auto">
                     {{ isEditing ? 'Update Receipt' : 'Create Receipt' }}
                   </UiBaseButton>
                </div>

              </form>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- Camera Modal -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="showCamera" class="fixed inset-0 bg-black flex items-center justify-center z-50">
          <div class="relative w-full h-full max-w-md max-h-md">
            <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover"></video>
            <canvas ref="canvasRef" class="hidden"></canvas>

            <!-- Camera Controls -->
            <div class="absolute bottom-0 left-0 right-0 p-6 bg-black/50 backdrop-blur-sm">
              <div class="flex justify-center gap-4">
                <button @click="stopCamera(); showCamera = false" class="bg-red-500 hover:bg-red-600 text-white rounded-full p-4 transition-colors">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
                <button @click="captureImage" class="bg-white hover:bg-gray-200 text-black rounded-full p-4 transition-colors">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="12" r="10"></circle>
                    <circle cx="12" cy="12" r="3"></circle>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

     <p class="text-sm text-text-light dark:text-text-dark/60 text-center">
       Click folder headers to expand/collapse • Click images to view • Click edit icon to edit • Double-click receipts to delete
     </p>

    <!-- Folder-based Receipt Organization -->
    <div class="space-y-4">
      <div v-if="receiptStore.receipts.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No receipts added yet</p>
        <p class="text-sm">Create your first receipt above to get started</p>
      </div>

      <!-- Folder Structure -->
      <div v-for="[folderKey, receipts] in groupedReceipts" :key="folderKey"
        class="rounded-lg border border-surface-light dark:border-surface-dark overflow-hidden">

        <!-- Folder Header -->
        <div class="bg-surface-light/50 dark:bg-surface-dark/50 p-4 cursor-pointer hover:bg-surface-light/70 dark:hover:bg-surface-dark/70 transition-colors"
          @click="toggleFolder(folderKey)">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <!-- Folder Icon -->
              <svg :class="['w-6 h-6 transition-transform duration-200', expandedFolders.has(folderKey) ? 'rotate-90' : '']"
                fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 5l7 7-7 7"></path>
              </svg>

              <!-- Folder Icon -->
              <svg class="w-6 h-6 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z"></path>
              </svg>

              <!-- Folder Name -->
              <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">
                {{ formatFolderName(folderKey) }}
              </h3>
            </div>

            <!-- Receipt Count -->
            <div class="flex items-center gap-2">
              <span class="px-2 py-1 bg-primary/10 text-primary rounded-full text-sm font-medium">
                {{ receipts.length }} receipt{{ receipts.length !== 1 ? 's' : '' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Folder Contents -->
        <div v-if="expandedFolders.has(folderKey)"
          class="p-4 bg-surface-light/20 dark:bg-surface-dark/20">
          <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <div v-for="receipt in receipts" :key="receipt.receipt_id"
              class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-sm border border-surface-light dark:border-surface-dark hover:shadow-md transition-shadow duration-200 overflow-hidden"
              @dblclick="(receiptToDelete = receipt.receipt_id, showDialog = true)">

              <!-- Receipt Image -->
              <div class="aspect-video overflow-hidden">
                <img :src="getImageUrl(receipt.image_path)" :alt="receipt.title"
                  class="w-full h-full object-cover cursor-pointer hover:scale-105 transition-transform duration-200"
                  @click="() => { console.log('Click detected on receipt:', receipt?.title); openImageViewer(receipt) }" />
              </div>

               <!-- Receipt Info -->
               <div class="p-4">
                 <div class="flex items-center justify-between mb-1">
                   <h4 class="font-semibold text-text-light dark:text-text-dark truncate flex-1">{{ receipt.title }}</h4>
                   <button @click.stop="openEditForm(receipt)" class="ml-2 p-1 text-text-light dark:text-text-dark/60 hover:text-primary transition-colors">
                     <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                     </svg>
                   </button>
                 </div>
                 <p class="text-sm text-text-light dark:text-text-dark/60 mb-2">
                   {{ formatDate(receipt.created_at) }}
                 </p>

                 <div class="flex items-center gap-2">
                   <span v-if="receipt.amount" class="px-2 py-1 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-400">
                     {{ formatCurrency(receipt.amount) }}
                   </span>
                   <span v-if="receipt.Category" class="px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-400 truncate max-w-30">
                     {{ receipt.Category.name }}
                   </span>
                 </div>
               </div>
            </div>
          </div>
        </div>
      </div>

      <UiConfirmDialog v-model:show="showDialog" :message="`Delete receipt '${receiptToDelete ? receiptStore.receipts.find(r => r.receipt_id === receiptToDelete)?.title : ''}'?`"
        @confirm="deleteItem(receiptToDelete)" />
    </div>
  </div>
  <!-- Image Viewer Modal -->
  <ClientOnly>
    <Teleport to="body">
      <div v-if="imageViewer.show" class="fixed inset-0 bg-black/90 flex items-center justify-center p-4 z-50" @click="closeImageViewer">
      <!-- Debug: Modal should be visible -->
      {{ console.log("Modal should be visible, imageViewer.show:", imageViewer.show) }}
        <div class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-4xl max-h-[90vh] overflow-auto" @click.stop>
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ imageViewer.receipt?.title || 'Receipt' }}</h3>
            <button @click="closeImageViewer" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
          <div class="mb-4">
            <img
              :src="getImageUrl(imageViewer.receipt?.image_path)"
              :alt="imageViewer.receipt?.title || 'Receipt'"
              class="w-full h-auto max-h-[60vh] object-contain rounded-lg" />
          </div>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            <p><strong>Date:</strong> {{ formatDate(imageViewer.receipt?.created_at) }}</p>
            <p v-if="imageViewer.receipt?.amount"><strong>Amount:</strong> ${{ imageViewer.receipt.amount.toFixed(2) }}</p>
            <p v-if="imageViewer.receipt?.Category"><strong>Category:</strong> {{ imageViewer.receipt.Category.name }}</p>
          </div>
        </div>
      </div>
    </Teleport>
  </ClientOnly>

</template>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.25s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.fade-scale-enter-to {
  opacity: 1;
  transform: scale(1);
}

.fade-scale-leave-from {
  opacity: 1;
  transform: scale(1);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
