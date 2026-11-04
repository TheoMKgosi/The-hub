import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useNoteStore, type Note } from '../../../app/stores/notes'

// The note store reaches for Nuxt auto-imports (useNuxtApp, useToast).
// Stub them all; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

// Plain (non-mock) DOM stubs for exportNote: it appends an <a> and clicks it.
const createdLinks: Array<{ href: string; download: string; clicks: number }> = []

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  // exportNote reads the raw response through $api.raw — alias it to the mock
  apiMock.raw = apiMock
  vi.stubGlobal('document', {
    createElement: () => {
      const link = { href: '', download: '', clicks: 0, click() { this.clicks++ } }
      createdLinks.push(link)
      return link
    },
    body: {
      appendChild: () => {},
      removeChild: () => {}
    }
  })
})

beforeEach(() => {
  // Fresh pinia per test gives fresh store state
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  createdLinks.length = 0
  vi.spyOn(console, 'error').mockImplementation(() => {})
  vi.spyOn(console, 'log').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const note = (id: string, overrides: Partial<Note> = {}): Note => ({
  id,
  user_id: 'u1',
  title: `note-${id}`,
  content: `content-${id}`,
  tags: [],
  created_at: '2026-01-01T00:00:00.000Z',
  updated_at: '2026-01-01T00:00:00.000Z',
  ...overrides
})

describe('useNoteStore', () => {
  it('starts with empty notes, no selection and no filters', () => {
    const store = useNoteStore()

    expect(store.notes).toEqual([])
    expect(store.selectedNoteId).toBeNull()
    expect(store.selectedNote).toBeNull()
    expect(store.searchQuery).toBe('')
    expect(store.selectedTag).toBeNull()
    expect(store.allTags).toEqual([])
    expect(store.filteredNotes).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.fetchError).toBeNull()
  })

  describe('computed getters', () => {
    it('selectedNote returns the note matching selectedNoteId', () => {
      const store = useNoteStore()
      store.notes.push(note('n1'), note('n2'))

      store.selectNote('n2')

      expect(store.selectedNote).toEqual(note('n2'))
    })

    it('selectedNote returns null for an unknown id', () => {
      const store = useNoteStore()
      store.notes.push(note('n1'))

      store.selectNote('missing')

      expect(store.selectedNote).toBeNull()
    })

    it('allTags collects a sorted list of unique tags', () => {
      const store = useNoteStore()
      store.notes.push(note('n1', { tags: ['b', 'a'] }), note('n2', { tags: ['a', 'c'] }))

      expect(store.allTags).toEqual(['a', 'b', 'c'])
    })

    it('filteredNotes matches the search query against title and content case-insensitively', () => {
      const store = useNoteStore()
      store.notes.push(
        note('n1', { title: 'Vue tips', content: 'composables' }),
        note('n2', { title: 'Go patterns', content: 'VUE router' }),
        note('n3', { title: 'Unrelated', content: 'nothing' })
      )

      store.setSearchQuery('vue')

      expect(store.filteredNotes.map(n => n.id)).toEqual(['n1', 'n2'])
    })

    it('filteredNotes narrows by the selected tag', () => {
      const store = useNoteStore()
      store.notes.push(note('n1', { tags: ['work'] }), note('n2', { tags: ['home'] }))

      store.setSelectedTag('home')

      expect(store.filteredNotes.map(n => n.id)).toEqual(['n2'])
    })
  })

  describe('fetchNotes', () => {
    it('requests a bare path when no filters are set', async () => {
      const store = useNoteStore()
      apiMock.mockResolvedValue({ notes: [note('n1')] })

      await store.fetchNotes()

      expect(apiMock).toHaveBeenCalledWith('notes')
      expect(store.notes).toEqual([note('n1')])
      expect(store.loading).toBe(false)
    })

    it('appends search and tag filters to the query string', async () => {
      const store = useNoteStore()
      store.setSearchQuery('vue')
      store.setSelectedTag('work')
      apiMock.mockResolvedValue({ notes: [] })

      await store.fetchNotes()

      expect(apiMock).toHaveBeenCalledWith('notes?search=vue&tag=work')
      expect(store.notes).toEqual([])
    })

    it('falls back to an empty list when the API omits notes', async () => {
      const store = useNoteStore()
      store.notes.push(note('n1'))
      apiMock.mockResolvedValue({})

      await store.fetchNotes()

      expect(store.notes).toEqual([])
      expect(store.loading).toBe(false)
    })

    it('sets fetchError, shows a toast and clears loading on failure', async () => {
      const store = useNoteStore()
      const failure = new Error('network down')
      apiMock.mockRejectedValue(failure)

      await store.fetchNotes()

      expect(store.fetchError).toBe(failure)
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch notes',
        color: 'error'
      })
    })
  })

  describe('createNote', () => {
    it('prepends the note optimistically and selects it, then swaps in the server note', async () => {
      const store = useNoteStore()
      const payload = { title: 'New note', content: 'body', tags: ['work'] }
      apiMock.mockResolvedValue(note('n9', { title: 'New note' }))

      const result = await store.createNote(payload)

      expect(apiMock).toHaveBeenCalledWith('notes', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.notes).toHaveLength(1)
      expect(store.notes[0]).toEqual(note('n9', { title: 'New note' }))
      expect(store.selectedNoteId).toBe('n9')
      expect(result).toEqual(note('n9', { title: 'New note' }))
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Notes',
        description: 'Note created successfully',
        color: 'success'
      })
    })

    it('rolls back the optimistic note and clears the selection on failure', async () => {
      const store = useNoteStore()
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.createNote({ title: 'New note', content: 'body', tags: [] })

      expect(result).toBeNull()
      expect(store.notes).toHaveLength(0)
      expect(store.selectedNoteId).toBeNull()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to create note',
        color: 'error'
      })
    })
  })

  describe('updateNote', () => {
    it('applies the update optimistically, then reconciles with the server response', async () => {
      const store = useNoteStore()
      store.notes.push(note('n1'))
      apiMock.mockResolvedValue(note('n1', { title: 'Updated' }))

      await store.updateNote('n1', { title: 'Updated' })

      expect(apiMock).toHaveBeenCalledWith('notes/n1', {
        method: 'PATCH',
        body: JSON.stringify({ title: 'Updated' })
      })
      expect(store.notes[0]).toEqual(note('n1', { title: 'Updated' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Note', description: 'Note saved', color: 'success' })
    })

    it('bumps updated_at optimistically before the server responds', async () => {
      const store = useNoteStore()
      store.notes.push(note('n1'))
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.updateNote('n1', { title: 'Updated' })
      // While the request is in flight, the optimistic update already applied
      expect(store.notes[0].title).toBe('Updated')
      expect(store.notes[0].updated_at).not.toBe('2026-01-01T00:00:00.000Z')

      resolveApi(note('n1', { title: 'Updated' }))
      await pending
    })

    it('reverts to the original note on failure', async () => {
      const store = useNoteStore()
      store.notes.push(note('n1'))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.updateNote('n1', { title: 'Updated' })

      expect(store.notes[0]).toEqual(note('n1'))
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to save note',
        color: 'error'
      })
    })

    it('still reports success for an unknown note id', async () => {
      const store = useNoteStore()
      apiMock.mockResolvedValue(note('n99'))

      await store.updateNote('n99', { title: 'Updated' })

      expect(store.notes).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Note', description: 'Note saved', color: 'success' })
    })
  })

  describe('deleteNote', () => {
    it('removes the note optimistically and moves the selection to the first remaining note', async () => {
      const store = useNoteStore()
      store.notes.push(note('n1'), note('n2'))
      store.selectNote('n1')
      apiMock.mockResolvedValue(undefined)

      await store.deleteNote('n1')

      expect(apiMock).toHaveBeenCalledWith('notes/n1', { method: 'DELETE' })
      expect(store.notes.map(n => n.id)).toEqual(['n2'])
      expect(store.selectedNoteId).toBe('n2')
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Note', description: 'Note deleted', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useNoteStore()

      await store.deleteNote('n42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Note not found',
        color: 'error'
      })
    })

    it('NOTE: documents current behavior — on failure the note is re-appended at the end and the selection moves to the first note', async () => {
      const store = useNoteStore()
      store.notes.push(note('n1'), note('n2'))
      store.selectNote('n1')
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteNote('n1')

      // Rollback pushes the note back to the end of the list
      expect(store.notes.map(n => n.id)).toEqual(['n2', 'n1'])
      // The selection was already moved to n2 before the API call, and stays there
      expect(store.selectedNoteId).toBe('n2')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to delete note',
        color: 'error'
      })
    })
  })

  describe('exportNote', () => {
    it('exports the note payload as a download and shows a success toast', async () => {
      const store = useNoteStore()
      apiMock.mockResolvedValue({
        headers: { get: (name: string) => (name === 'Content-Disposition' ? 'attachment; filename="note.json"' : null) },
        _data: { id: 'n1', title: 'note-n1' }
      })

      await store.exportNote('n1')

      expect(apiMock).toHaveBeenCalledWith('notes/n1/export')
      expect(createdLinks).toHaveLength(1)
      expect(createdLinks[0].download).toBe('note.json')
      expect(createdLinks[0].clicks).toBe(1)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Note', description: 'Note exported', color: 'success' })
    })

    it('uses the default filename without a Content-Disposition header', async () => {
      const store = useNoteStore()
      apiMock.mockResolvedValue({
        headers: { get: () => null },
        _data: { id: 'n1' }
      })

      await store.exportNote('n1')

      expect(createdLinks[0].download).toBe('note_export.json')
    })

    it('shows an error toast instead of throwing when the export fails', async () => {
      const store = useNoteStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.exportNote('n1')

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Note not exported',
        color: 'error'
      })
    })
  })

  describe('filters and selection helpers', () => {
    it('selectNote, setSearchQuery and setSelectedTag update their state', () => {
      const store = useNoteStore()

      store.selectNote('n1')
      store.setSearchQuery('vue')
      store.setSelectedTag('work')

      expect(store.selectedNoteId).toBe('n1')
      expect(store.searchQuery).toBe('vue')
      expect(store.selectedTag).toBe('work')
    })
  })

  describe('reset', () => {
    it('clears notes, selection and filters', () => {
      const store = useNoteStore()
      store.notes.push(note('n1'))
      store.selectNote('n1')
      store.setSearchQuery('vue')
      store.setSelectedTag('work')

      store.reset()

      expect(store.notes).toEqual([])
      expect(store.selectedNoteId).toBeNull()
      expect(store.searchQuery).toBe('')
      expect(store.selectedTag).toBeNull()
    })
  })
})
