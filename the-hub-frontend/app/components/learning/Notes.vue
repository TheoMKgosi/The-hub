<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useNoteStore, type NoteFormData } from '@/stores/notes'
import { useMarkdown } from '@/composables/useMarkdown'
import SearchIcon from '../ui/svg/SearchIcon.vue'
import PlusIcon from '../ui/svg/PlusIcon.vue'
import DeleteIcon from '../ui/svg/DeleteIcon.vue'
import EditIcon from '../ui/svg/EditIcon.vue'

const noteStore = useNoteStore()
const { renderMarkdown } = useMarkdown()

const isEditing = ref(false)
const isPreview = ref(false)
const newNoteTitle = ref('')
const editingTag = ref('')
const isCreating = ref(false)

const formData = computed(() => ({
  title: noteStore.selectedNote?.title || '',
  content: noteStore.selectedNote?.content || '',
  tags: noteStore.selectedNote?.tags || []
}))

onMounted(() => {
  if (noteStore.notes.length === 0) {
    noteStore.fetchNotes()
  }
})

watch(() => noteStore.selectedNoteId, () => {
  isEditing.value = false
})

const createNote = async () => {
  if (!newNoteTitle.value.trim()) return

  await noteStore.createNote({
    title: newNoteTitle.value,
    content: '',
    tags: []
  })

  newNoteTitle.value = ''
  isCreating.value = false
  isEditing.value = true
}

const saveNote = async () => {
  if (!noteStore.selectedNoteId) return

  await noteStore.updateNote(noteStore.selectedNoteId, {
    title: formData.value.title,
    content: formData.value.content,
    tags: formData.value.tags
  })

  isEditing.value = false
}

const deleteNote = async (id: string) => {
  await noteStore.deleteNote(id)
}

const addTag = () => {
  if (!editingTag.value.trim() || !noteStore.selectedNote) return

  const newTags = [...formData.value.tags, editingTag.value.trim()]
  noteStore.updateNote(noteStore.selectedNoteId!, { tags: newTags })
  editingTag.value = ''
}

const removeTag = (tag: string) => {
  if (!noteStore.selectedNoteId) return

  const newTags = formData.value.tags.filter(t => t !== tag)
  noteStore.updateNote(noteStore.selectedNoteId, { tags: newTags })
}

const renderedContent = computed(() => {
  return renderMarkdown(formData.value.content)
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}
</script>

<template>
  <div class="notes-container flex h-full gap-4">
    <!-- Sidebar -->
    <div class="notes-sidebar w-80 flex flex-col bg-white rounded-lg border shadow-sm">
      <!-- Search -->
      <div class="p-3 border-b">
        <div class="relative">
          <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" />
          <input
            v-model="noteStore.searchQuery"
            @input="noteStore.fetchNotes()"
            type="text"
            placeholder="Search notes..."
            class="w-full pl-9 pr-3 py-2 text-sm border border-black rounded-md outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <!-- Tags filter -->
      <div v-if="noteStore.allTags.length > 0" class="px-3 py-2 border-b">
        <div class="flex flex-wrap gap-1">
          <button
            @click="noteStore.setSelectedTag(null)"
            class="px-2 py-1 text-xs rounded-full transition-colors"
            :class="!noteStore.selectedTag ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'"
          >
            All
          </button>
          <button
            v-for="tag in noteStore.allTags"
            :key="tag"
            @click="noteStore.setSelectedTag(tag)"
            class="px-2 py-1 text-xs rounded-full transition-colors"
            :class="noteStore.selectedTag === tag ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'"
          >
            {{ tag }}
          </button>
        </div>
      </div>

      <!-- Notes list -->
      <div class="flex-1 overflow-y-auto">
        <!-- New note input -->
        <div class="p-3 border-b">
          <div v-if="isCreating" class="space-y-2">
            <input
              v-model="newNoteTitle"
              @keyup.enter="createNote"
              @keyup.escape="isCreating = false"
              type="text"
              placeholder="Note title..."
              class="w-full px-3 py-2 text-sm border rounded-md outline-none focus:ring-2 focus:ring-blue-500"
              autofocus
            />
            <div class="flex gap-2">
              <button
                @click="createNote"
                class="flex-1 px-3 py-1.5 text-sm bg-blue-500 text-white rounded-md hover:bg-blue-600"
              >
                Create
              </button>
              <button
                @click="isCreating = false"
                class="px-3 py-1.5 text-sm text-black hover:text-gray-800"
              >
                Cancel
              </button>
            </div>
          </div>
          <button
            v-else
            @click="isCreating = true"
            class="w-full flex items-center justify-center gap-2 px-3 py-2 text-sm border-2 border-dashed rounded-md text-black hover:border-blue-500 hover:text-blue-500 transition-colors"
          >
            <PlusIcon class="w-4 h-4" />
            New Note
          </button>
        </div>

        <!-- Note items -->
        <div class="divide-y">
          <button
            v-for="note in noteStore.filteredNotes"
            :key="note.id"
            @click="noteStore.selectNote(note.id)"
            class="w-full p-3 text-left text-black hover:bg-gray-50 transition-colors"
            :class="{ 'bg-blue-50 border-l-4 border-blue-500': noteStore.selectedNoteId === note.id }"
          >
            <h3 class="font-medium text-sm truncate">{{ note.title }}</h3>
            <p class="text-xs mt-1 text-black line-clamp-2">
              {{ note.content.slice(0, 100) || 'No content' }}
            </p>
            <div class="flex items-center gap-2 mt-2">
              <span class="text-xs text-gray-400">{{ formatDate(note.updated_at) }}</span>
              <span
                v-for="tag in (note.tags || []).slice(0, 2)"
                :key="tag"
                class="px-1.5 py-0.5 text-xs bg-gray-100 rounded"
              >
                {{ tag }}
              </span>
            </div>
          </button>

          <div v-if="noteStore.filteredNotes.length === 0" class="p-6 text-center text-gray-500 text-sm">
            No notes yet
          </div>
        </div>
      </div>
    </div>

    <!-- Main editor area -->
    <div class="notes-editor flex-1 flex flex-col bg-white text-black rounded-lg border shadow-sm">
      <div v-if="noteStore.selectedNote" class="flex flex-col h-full">
        <!-- Editor header -->
        <div class="flex items-center justify-between p-4 border-b">
          <div class="flex items-center gap-3 flex-1">
            <input
              v-model="formData.title"
              :readonly="!isEditing"
              class="text-lg font-semibold bg-transparent border-none outline-none flex-1"
              :class="{ 'focus:ring-2 focus:ring-blue-500 rounded px-2 -mx-2': isEditing }"
            />
          </div>
          <div class="flex items-center gap-2">
            <button
              v-if="!isEditing"
              @click="isEditing = true"
              class="p-2  hover:text-blue-500 hover:bg-gray-100 rounded"
            >
              <EditIcon class="w-4 h-4" />
            </button>
            <button
              v-if="isEditing"
              @click="saveNote"
              class="px-3 py-1.5 text-sm bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Save
            </button>
            <button
              v-if="isEditing"
              @click="isEditing = false"
              class="px-3 py-1.5 text-sm hover:text-gray-800"
            >
              Cancel
            </button>
            <button
              @click="isPreview = !isPreview"
              class="px-3 py-1.5 text-sm border rounded hover:bg-gray-50"
              :class="isPreview ? 'bg-blue-50 border-blue-500 text-blue-600' : ''"
            >
              {{ isPreview ? 'Edit' : 'Preview' }}
            </button>
            <button
              @click="deleteNote(noteStore.selectedNote!.id)"
              class="p-2 text-gray-500 hover:text-red-500 hover:bg-red-50 rounded"
            >
              <DeleteIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Tags -->
        <div class="px-4 py-2 border-b flex items-center gap-2 flex-wrap">
          <span class="text-sm">Tags:</span>
          <span
            v-for="tag in formData.tags"
            :key="tag"
            class="inline-flex items-center gap-1 px-2 py-1 text-xs bg-gray-100 rounded-full"
          >
            {{ tag }}
            <button
              @click="removeTag(tag)"
              class="text-gray-400 hover:text-red-500"
              :disabled="!isEditing"
            >
              ×
            </button>
          </span>
          <div v-if="isEditing" class="flex items-center gap-1">
            <input
              v-model="editingTag"
              @keyup.enter="addTag"
              type="text"
              placeholder="Add tag..."
              class="px-2 py-1 text-xs border rounded outline-none focus:ring-1 focus:ring-blue-500"
            />
            <button
              @click="addTag"
              class="px-2 py-1 text-xs bg-gray-100 rounded hover:bg-gray-200"
            >
              Add
            </button>
          </div>
        </div>

        <!-- Content -->
        <div class="flex-1 overflow-hidden">
          <textarea
            v-if="!isPreview"
            v-model="formData.content"
            :readonly="!isEditing"
            class="w-full h-full p-4 resize-none outline-none font-mono text-sm"
            :class="{ 'bg-gray-50': !isEditing }"
            placeholder="Start writing in Markdown..."
          ></textarea>
          <div
            v-else
            class="w-full h-full p-4 overflow-y-auto prose prose-sm max-w-none"
            v-html="renderedContent"
          ></div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="flex-1 flex items-center justify-center">
        <div class="text-center">
          <p class="text-lg mb-2">No note selected</p>
          <p class="text-sm">Select a note or create a new one</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notes-container {
  height: calc(100vh - 180px);
  min-height: 500px;
}

.notes-sidebar {
  min-width: 280px;
}

.notes-editor {
  min-width: 0;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.prose h1 { font-size: 1.5em; font-weight: bold; margin: 0.5em 0; }
.prose h2 { font-size: 1.25em; font-weight: bold; margin: 0.5em 0; }
.prose h3 { font-size: 1.1em; font-weight: bold; margin: 0.5em 0; }
.prose p { margin: 0.5em 0; }
.prose ul, .prose ol { margin: 0.5em 0; padding-left: 1.5em; }
.prose li { margin: 0.25em 0; }
.prose code { background: #f3f4f6; padding: 0.1em 0.3em; border-radius: 3px; font-size: 0.9em; }
.prose pre { background: #f3f4f6; padding: 1em; border-radius: 6px; overflow-x: auto; margin: 0.5em 0; }
.prose pre code { background: none; padding: 0; }
.prose blockquote { border-left: 3px solid #d1d5db; padding-left: 1em; color: #6b7280; margin: 0.5em 0; }
.prose a { color: #3b82f6; text-decoration: underline; }
.prose table { border-collapse: collapse; width: 100%; margin: 0.5em 0; }
.prose th, .prose td { border: 1px solid #d1d5db; padding: 0.5em; text-align: left; }
.prose th { background: #f3f4f6; }
</style>
