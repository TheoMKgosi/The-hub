<script setup lang="ts">
import { useNoteStore, type NoteFormData } from '@/stores/notes'
import { useMarkdown } from '@/composables/useMarkdown'
import type { EditorToolbarItem } from '@nuxt/ui'

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

const items: EditorToolbarItem[] = [
  { kind: 'mark', mark: 'bold', icon: 'i-lucide-bold' },
  { kind: 'mark', mark: 'italic', icon: 'i-lucide-italic' },
  { kind: 'heading', level: 1, icon: 'i-lucide-heading-1' },
  { kind: 'heading', level: 2, icon: 'i-lucide-heading-2' },
  { kind: 'textAlign', align: 'left', icon: 'i-lucide-align-left' },
  { kind: 'textAlign', align: 'center', icon: 'i-lucide-align-center' },
  { kind: 'bulletList', icon: 'i-lucide-list' },
  { kind: 'orderedList', icon: 'i-lucide-list-ordered' },
  { kind: 'blockquote', icon: 'i-lucide-quote' },
  { kind: 'link', icon: 'i-lucide-link' }
]

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
  <div class="notes-container flex flex-col md:flex-row h-full gap-4">
    <!-- Sidebar -->
    <div class="notes-sidebar w-full md:w-80 max-h-[40vh] md:max-h-none flex flex-col bg-slate-200 dark:bg-slate-800 rounded-lg border shadow-sm">
      <!-- Search -->
      <div class="p-3 border-b">
        <UInput v-model="noteStore.searchQuery" @input="noteStore.fetchNotes()" placeholder="Search notes..." size="lg"
          class="w-full" />
      </div>

      <!-- Tags filter -->
      <div v-if="noteStore.allTags.length > 0" class="px-3 py-2 border-b">
        <div class="flex flex-wrap gap-1">
          <button @click="noteStore.setSelectedTag(null)" class="px-2 py-1 text-xs rounded-full transition-colors"
            :class="!noteStore.selectedTag ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'">
            All
          </button>
          <button v-for="tag in noteStore.allTags" :key="tag" @click="noteStore.setSelectedTag(tag)"
            class="px-2 py-1 text-xs rounded-full transition-colors"
            :class="noteStore.selectedTag === tag ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'">
            {{ tag }}
          </button>
        </div>
      </div>

      <!-- Notes list -->
      <div class="flex-1 overflow-y-auto ">
        <!-- New note input -->
        <div class="p-3 border-b">
          <div v-if="isCreating" class="space-y-2">
            <UInput v-model="newNoteTitle" @keyup.enter="createNote" @keyup.escape="isCreating = false"
              placeholder="Note title..." class="w-full" size="lg" />
            <div class="flex gap-2">
              <UButton label="Create" @click="createNote" />
              <UButton label="Cancel" @click="isCreating = false" color="error" />
            </div>
          </div>
          <UButton v-else label="New Note" icon="i-lucide-plus" @click="isCreating = true" block />
        </div>

        <!-- Note items -->
        <div class="divide-y">
          <button v-for="note in noteStore.filteredNotes" :key="note.id" @click="noteStore.selectNote(note.id)"
            class="w-full p-3 md:p-3 text-left hover:bg-gray-50 transition-colors"
            :class="{ 'bg-blue-50 border-l-4 border-blue-500': noteStore.selectedNoteId === note.id }">
            <h3 class="font-medium text-sm truncate">{{ note.title }}</h3>
            <p class="text-xs mt-1 line-clamp-2">
              {{ note.content.slice(0, 100) || 'No content' }}
            </p>
            <div class="flex items-center gap-2 mt-2">
              <span class="text-xs text-gray-400">{{ formatDate(note.updated_at) }}</span>
              <span v-for="tag in (note.tags || []).slice(0, 2)" :key="tag"
                class="px-1.5 py-0.5 text-xs bg-gray-100 rounded">
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
    <div class="notes-editor flex-1 min-h-[60vh] md:min-h-0 flex flex-col rounded-lg border shadow-sm">
      <div v-if="noteStore.selectedNote" class="flex flex-col h-full">
        <!-- Editor header -->
        <div class="flex items-center justify-between p-2 md:p-4 border-b">
          <div class="flex items-center gap-3 flex-1 min-w-0">
            <UInput v-model="formData.title" variant="none" size="lg" class="min-w-0" />
          </div>
          <div class="flex items-center gap-1 md:gap-2 shrink-0">
            <UButton icon="i-lucide-save" @click="saveNote" variant="outline" />
            <UButton label="Export Note" variant="subtle" color="neutral"
              @click="noteStore.exportNote(noteStore.selectedNoteId!)" class="hidden md:inline-flex" />
            <UButton icon="i-lucide-trash-2" @click="deleteNote(noteStore.selectedNoteId!)" color="error"
              variant="subtle" />
          </div>
        </div>

        <!-- Tags -->
        <div class="px-2 md:px-4 py-2 border-b flex items-center gap-2 flex-wrap">
          <span class="text-sm">Tags:</span>
          <span v-for="tag in formData.tags" :key="tag">
            <UBadge :label="tag" @click="removeTag(tag)" variant="outline"/>
          </span>
          <div class="flex items-center gap-1">
            <input v-model="editingTag" @keyup.enter="addTag" type="text" placeholder="Add tag..."
              class="px-2 py-1 text-xs border rounded outline-none focus:ring-1 focus:ring-blue-500" />
            <UButton label="Add" @click="addTag" size="sm" />
          </div>
        </div>

        <!-- Content -->
        <div class="flex-1 overflow-y-auto p-2 md:p-4">
          <UEditor v-slot="{ editor }" v-model="formData.content" placeholder="Start writing..."
            content-type="markdown">
            <UEditorToolbar :editor="editor" :items="items" layout="floating" />
          </UEditor>
          <!--   <textarea v-if="!isPreview" v-model="formData.content" :readonly="!isEditing" -->
          <!--     class="w-full h-full p-4 resize-none outline-none font-mono text-sm" -->
          <!--     :class="{ 'dark:bg-black bg-gray-50': !isEditing }" placeholder="Start writing in Markdown..."></textarea> -->
          <!--   <div v-else class="w-full h-full p-4 overflow-y-auto prose prose-sm max-w-none" v-html="renderedContent"> -->
          <!--   </div> -->
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
  min-width: 0;
}

@media (min-width: 768px) {
  .notes-sidebar {
    min-width: 280px;
  }
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
</style>
