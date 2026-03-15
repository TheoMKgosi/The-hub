import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useToast } from '@/composables/useToast'


export interface Note {
  id: string
  user_id: string
  title: string
  content: string
  tags: string[]
  created_at: string
  updated_at: string
}

interface NotesResponse {
  notes: Note[]
}

export interface NoteFormData {
  title: string
  content: string
  tags: string[]
}

export const useNoteStore = defineStore('note', () => {
  const notes = ref<Note[]>([])
  const selectedNoteId = ref<string | null>(null)
  const searchQuery = ref('')
  const selectedTag = ref<string | null>(null)
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  const selectedNote = computed(() => {
    if (!selectedNoteId.value) return null
    return notes.value.find(n => n.id === selectedNoteId.value) || null
  })

  const allTags = computed(() => {
    const tagSet = new Set<string>()
    console.log(notes.value)
    notes.value.forEach(note => {
      (note.tags || []).forEach(tag => tagSet.add(tag))
    })
    return Array.from(tagSet).sort()
  })

  const filteredNotes = computed(() => {
    let filtered = [...notes.value]

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter(
        note =>
          note.title.toLowerCase().includes(query) ||
          note.content.toLowerCase().includes(query)
      )
    }

    if (selectedTag.value) {
      filtered = filtered.filter(note =>
        (note.tags || []).includes(selectedTag.value!)
      )
    }

    return filtered
  })

  async function fetchNotes() {
    const { $api } = useNuxtApp()
    loading.value = true
    fetchError.value = null
    try {
      const params = new URLSearchParams()
      if (searchQuery.value) params.append('search', searchQuery.value)
      if (selectedTag.value) params.append('tag', selectedTag.value)

      const queryString = params.toString()
      const url = queryString ? `notes?${queryString}` : 'notes'
      const fetchedNotes = await $api<NotesResponse>(url)

      if (fetchedNotes) {
        notes.value = fetchedNotes.notes ?? []
      }
    } catch (error) {
      fetchError.value = error as Error
      addToast('Failed to fetch notes', 'error')
      console.error('Error fetching notes:', error)
    } finally {
      loading.value = false
    }
  }

  async function createNote(payload: NoteFormData) {
    const optimisticNote: Note = {
      id: `temp-${Date.now()}`,
      user_id: '',
      title: payload.title,
      content: payload.content,
      tags: payload.tags,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }

    notes.value.unshift(optimisticNote)
    selectedNoteId.value = optimisticNote.id

    try {
      const { $api } = useNuxtApp()
      const newNote = await $api<Note>('notes', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      const index = notes.value.findIndex(n => n.id === optimisticNote.id)
      if (index !== -1) {
        notes.value[index] = newNote
        selectedNoteId.value = newNote.id
      }

      addToast('Note created successfully', 'success')
      return newNote
    } catch (err) {
      notes.value = notes.value.filter(n => n.id !== optimisticNote.id)
      if (selectedNoteId.value === optimisticNote.id) {
        selectedNoteId.value = notes.value[0]?.id || null
      }
      addToast('Failed to create note', 'error')
      console.error('Error creating note:', err)
      return null
    }
  }

  async function updateNote(id: string, payload: Partial<NoteFormData>) {
    const originalIndex = notes.value.findIndex(n => n.id === id)
    const originalNote = originalIndex !== -1 ? { ...notes.value[originalIndex] } : null

    if (originalIndex !== -1) {
      notes.value[originalIndex] = { ...notes.value[originalIndex], ...payload, updated_at: new Date().toISOString() }
    }

    try {
      const { $api } = useNuxtApp()
      const updatedNote = await $api<Note>(`notes/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      if (originalIndex !== -1 && updatedNote) {
        notes.value[originalIndex] = updatedNote
      }

      addToast('Note saved', 'success')
    } catch (err) {
      if (originalNote && originalIndex !== -1) {
        notes.value[originalIndex] = originalNote
      }
      addToast('Failed to save note', 'error')
      console.error('Error updating note:', err)
    }
  }

  async function deleteNote(id: string) {
    const noteToDelete = notes.value.find(n => n.id === id)
    if (!noteToDelete) {
      addToast('Note not found', 'error')
      return
    }

    notes.value = notes.value.filter(n => n.id !== id)

    if (selectedNoteId.value === id) {
      selectedNoteId.value = notes.value[0]?.id || null
    }

    try {
      const { $api } = useNuxtApp()
      await $api(`notes/${id}`, {
        method: 'DELETE'
      })

      addToast('Note deleted', 'success')
    } catch (err) {
      notes.value.push(noteToDelete)
      if (selectedNoteId.value === null && notes.value.length > 0) {
        selectedNoteId.value = notes.value[0].id
      }
      addToast('Failed to delete note', 'error')
      console.error('Error deleting note:', err)
    }
  }

  function selectNote(id: string | null) {
    selectedNoteId.value = id
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function setSelectedTag(tag: string | null) {
    selectedTag.value = tag
  }

  function reset() {
    notes.value = []
    selectedNoteId.value = null
    searchQuery.value = ''
    selectedTag.value = null
  }

  return {
    notes,
    selectedNoteId,
    selectedNote,
    searchQuery,
    selectedTag,
    allTags,
    filteredNotes,
    loading,
    fetchError,
    fetchNotes,
    createNote,
    updateNote,
    deleteNote,
    selectNote,
    setSearchQuery,
    setSelectedTag,
    reset
  }
})
