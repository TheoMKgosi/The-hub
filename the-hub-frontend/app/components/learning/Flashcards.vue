<script setup lang="ts">
import PlusIcon from '../ui/svg/PlusIcon.vue'
import BoltIcon from '../ui/svg/BoltIcon.vue'
import CrossIcon from '../ui/svg/CrossIcon.vue'
import CheckMarkIcon from '../ui/svg/CheckMarkIcon.vue'
const deckStore = useDeckStore()
const cardStore = useCardStore()

onMounted(() => {
  if (deckStore.decks.length === 0) {
    deckStore.fetchDecks()
  }
})

const formData = reactive({
  name: ''
})

const addDeck = () => {
  if (!formData.name.trim()) return // no empty names pls
  deckStore.submitForm({ ...formData })
  formData.name = '' // reset input
}

// Edit deck name state
const editingDeckId = ref<string | null>(null)
const editFormData = reactive({
  name: ''
})

const editInput = ref<HTMLInputElement | null>(null)

const startEditingDeck = (deck: any) => {
  editingDeckId.value = deck.deck_id
  editFormData.name = deck.name
  nextTick(() => {
    editInput.value?.focus()
    editInput.value?.select()
  })
}

const cancelEditingDeck = () => {
  editingDeckId.value = null
  editFormData.name = ''
}

const saveDeckName = async () => {
  if (!editingDeckId.value) return

  const trimmedName = editFormData.name.trim()
  if (!trimmedName) {
    // Show error for empty name
    return
  }

  // Check if name actually changed
  const currentDeck = deckStore.decks.find(d => d.deck_id === editingDeckId.value)
  if (currentDeck && currentDeck.name === trimmedName) {
    // No change, just cancel editing
    cancelEditingDeck()
    return
  }

  try {
    const deckToUpdate = {
      deck_id: editingDeckId.value,
      name: trimmedName
    }

    await deckStore.editDeck(deckToUpdate)
    editingDeckId.value = null
    editFormData.name = ''
  } catch (error) {
    // Error handling is done in the store with toast notifications
    console.error('Failed to update deck name:', error)
  }
}

const removeDeck = (id: string) => {
  deckStore.deleteDeck(id)
}

const editDeck = (deck_id: string) => {
  navigateTo(`/learning/cards/${deck_id}`)
}

const reviewDeck = (deck_id: string) => {
  navigateTo(`/learning/review/${deck_id}`)
}

const browseDeck = (deck_id: string) => {
  navigateTo(`/learning/browse/${deck_id}`)
}

const exportDeck = async (deckId: string, format: 'json' | 'csv') => {
  try {
    await cardStore.exportCards(deckId, format)
  } catch (error) {
    console.error('Export failed:', error)
  }
}

// Import modal state
const showImport = ref(false)
const selectedDeckId = ref('')
const importFile = ref<File | null>(null)
const importFormat = ref<'json' | 'csv'>('json')
const importLoading = ref(false)
const importResult = ref<any>(null)

const showImportModal = (deckId: string) => {
  selectedDeckId.value = deckId
  showImport.value = true
  importFile.value = null
  importFormat.value = 'json'
  importResult.value = null
}

const closeImportModal = () => {
  showImport.value = false
  importFile.value = null
  importResult.value = null
}

const fileInput = ref<HTMLInputElement | null>(null)

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    importFile.value = target.files[0]
  }
}

const performImport = async () => {
  if (!importFile.value || !selectedDeckId.value) return

  importLoading.value = true
  try {
    const result = await cardStore.importCards(selectedDeckId.value, importFile.value, importFormat.value)
    importResult.value = result
  } catch (error) {
    console.error('Import failed:', error)
  } finally {
    importLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <div class="max-w-6xl mx-auto px-4 py-8">
      <!-- Header Section -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-text-light dark:text-text-dark mb-2">Flashcard Decks</h1>
        <p class="text-text-light/70 dark:text-text-dark/70">Create and manage your learning decks</p>
      </div>

      <!-- Add New Deck Form -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-6 mb-8">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Create New Deck
        </h2>
        <div class="flex gap-4">
          <div class="grow">
            <input type="text" v-model="formData.name"
              placeholder="Enter deck name (e.g., 'Spanish Vocabulary', 'Chemistry Terms')"
              class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200" />
          </div>
          <BaseButton @click="addDeck" text="Create Deck" :icon="PlusIcon" variant="primary" size="md" :disabled="!formData.name.trim()" class="px-6" />
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="deckStore.decks.length === 0" class="text-center py-16">
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-12 max-w-md mx-auto">
          <svg class="w-16 h-16 text-text-light/30 dark:text-text-dark/30 mx-auto mb-4" fill="none"
            stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1"
              d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No decks yet</h3>
          <p class="text-text-light/70 dark:text-text-dark/70 mb-6">Start building your knowledge by creating your first
            flashcard deck.</p>
          <BaseButton variant="primary" text="Create Your First Deck" size="lg" :icon="PlusIcon"
            @click="$refs.deckNameInput.focus()" />
        </div>
      </div>

      <!-- Decks Grid -->
      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="deck in deckStore.decks" :key="deck.deck_id"
          class="group bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark hover:shadow-xl hover:border-primary/20 dark:hover:border-primary/20 transition-all duration-300 overflow-hidden">
          <!-- Deck Header -->
          <div class="p-6 border-b border-surface-light dark:border-surface-dark">
            <div class="flex items-start justify-between mb-3">
              <div class="flex items-center gap-3 grow min-w-0">
                <div
                  class="w-10 h-10 bg-primary/10 dark:bg-primary/20 rounded-lg flex items-center justify-center shrink-0">
                  <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  </svg>
                </div>
                <div class="min-w-0 grow">
                  <!-- Edit Mode -->
                  <div v-if="editingDeckId === deck.deck_id" class="flex items-center gap-2">
                    <input v-model="editFormData.name" type="text"
                      class="grow border border-primary bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded px-2 py-1 text-lg font-semibold focus:outline-none focus:ring-2 focus:ring-primary"
                      @keyup.enter="saveDeckName" @keyup.escape="cancelEditingDeck" ref="editInput" />
                    <BaseButton @click="saveDeckName" variant="primary" size="sm" :iconOnly="true" :icon="CheckMarkIcon"
                      :disabled="!editFormData.name.trim() || editFormData.name.trim() === deckStore.decks.find(d => d.deck_id === editingDeckId)?.name" />
                    <BaseButton @click="cancelEditingDeck" variant="default" size="sm" :iconOnly="true"
                      :icon="CrossIcon" />
                  </div>
                  <!-- Display Mode -->
                  <div v-else class="flex items-center gap-2 min-w-0">
                    <h3
                      class="text-lg font-semibold text-text-light dark:text-text-dark truncate group-hover:text-primary transition-colors cursor-pointer"
                      @click="editDeck(deck.deck_id)">
                      {{ deck.name }}
                    </h3>
                    <button @click="startEditingDeck(deck)"
                      class="opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-text-light/60 dark:text-text-dark/60 hover:text-primary shrink-0"
                      title="Edit deck name">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                  </div>
                  <p class="text-sm text-text-light/60 dark:text-text-dark/60">Flashcard deck</p>
                </div>
              </div>
              <BaseButton @click="removeDeck(deck.deck_id)" variant="danger" size="sm" :iconOnly="false"
                icon="DeleteIcon"
                class="opacity-100 md:opacity-0 md:group-hover:opacity-100 transition-opacity duration-200 shrink-0" />
            </div>
          </div>

          <!-- Deck Actions -->
          <div class="p-6">
            <div class="grid grid-cols-2 gap-3 mb-3">
              <UiButton @click="editDeck(deck.deck_id)" variant="default" size="sm" class="justify-center">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Edit Cards
              </UiButton>
              <UiButton @click="browseDeck(deck.deck_id)" variant="secondary" size="sm" class="justify-center">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                Browse Cards
              </UiButton>
            </div>
            <div class="grid grid-cols-3 gap-3">
              <UiButton @click="exportDeck(deck.deck_id, 'json')" variant="outline" size="sm" class="justify-center"
                title="Export as JSON">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                JSON
              </UiButton>
              <UiButton @click="exportDeck(deck.deck_id, 'csv')" variant="outline" size="sm" class="justify-center"
                title="Export as CSV">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                CSV
              </UiButton>
              <UiButton @click="showImportModal(deck.deck_id)" variant="outline" size="sm" class="justify-center"
                title="Import cards">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                Import
              </UiButton>
            </div>
            <div class="mt-3">
              <BaseButton @click="reviewDeck(deck.deck_id)" text="Review" :icon="BoltIcon" variant="primary" size="sm"
                class="w-full justify-center" />
            </div>
          </div>
        </div>
      </div>

      <!-- Import Modal -->
      <div v-if="showImport" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex justify-center items-center z-50">
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-xl max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
          <div class="p-6">
            <div class="flex items-center justify-between mb-6">
              <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">Import Cards</h3>
              <button @click="closeImportModal"
                class="text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="space-y-6">
              <!-- Format Selection -->
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
                  Import Format
                </label>
                <div class="flex gap-4">
                  <label class="flex items-center">
                    <input v-model="importFormat" type="radio" value="json" class="text-primary focus:ring-primary" />
                    <span class="ml-2 text-sm text-text-light dark:text-text-dark">JSON</span>
                  </label>
                  <label class="flex items-center">
                    <input v-model="importFormat" type="radio" value="csv" class="text-primary focus:ring-primary" />
                    <span class="ml-2 text-sm text-text-light dark:text-text-dark">CSV</span>
                  </label>
                </div>
              </div>

              <!-- File Upload -->
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
                  Select File
                </label>
                <div
                  class="border-2 border-dashed border-surface-light dark:border-surface-dark rounded-lg p-6 text-center hover:border-primary transition-colors">
                  <input type="file" :accept="importFormat === 'json' ? '.json' : '.csv'" @change="handleFileSelect"
                    class="hidden" ref="fileInput" />
                  <div v-if="!importFile" @click="fileInput?.click()" class="cursor-pointer">
                    <svg class="w-12 h-12 text-text-light/50 dark:text-text-dark/50 mx-auto mb-4" fill="none"
                      stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                    <p class="text-text-light dark:text-text-dark mb-2">Click to select a {{ importFormat.toUpperCase()
                      }} file</p>
                    <p class="text-sm text-text-light/70 dark:text-text-dark/70">
                      Maximum file size: 10MB
                    </p>
                  </div>
                  <div v-else class="flex items-center justify-center gap-3">
                    <svg class="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <div class="text-left">
                      <p class="text-text-light dark:text-text-dark font-medium">{{ importFile.name }}</p>
                      <p class="text-sm text-text-light/70 dark:text-text-dark/70">
                        {{ (importFile.size / 1024 / 1024).toFixed(2) }} MB
                      </p>
                    </div>
                    <BaseButton @click="fileInput.value = ''; importFile = null" text="Change" variant="outline"
                      size="sm" />
                  </div>
                </div>
              </div>

              <!-- Import Result -->
              <div v-if="importResult" class="bg-surface-light dark:bg-surface-dark rounded-lg p-4">
                <div class="flex items-center gap-2 mb-2">
                  <svg v-if="importResult.success_count > 0" class="w-5 h-5 text-green-500" fill="none"
                    stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg v-else class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <h4 class="font-medium text-text-light dark:text-text-dark">Import Results</h4>
                </div>
                <p class="text-sm text-text-light dark:text-text-dark">
                  Successfully imported {{ importResult.success_count }} cards
                  <span v-if="importResult.error_count > 0" class="text-red-500">
                    ({{ importResult.error_count }} errors)
                  </span>
                </p>
                <div v-if="importResult.errors && importResult.errors.length > 0" class="mt-2">
                  <details class="text-sm">
                    <summary class="cursor-pointer text-red-500 hover:text-red-600">View errors</summary>
                    <ul class="mt-2 space-y-1">
                      <li v-for="error in importResult.errors" :key="error.row || error.error" class="text-red-500">
                        <span v-if="error.row">Row {{ error.row }}: </span>{{ error.error }}
                      </li>
                    </ul>
                  </details>
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="flex justify-end gap-3">
                <BaseButton @click="closeImportModal" text="Cancel" variant="default" />
                <BaseButton @click="performImport" text="Import Cards" :disabled="!importFile || importLoading"
                  variant="primary" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
