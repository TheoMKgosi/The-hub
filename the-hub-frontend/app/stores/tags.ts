import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useToast } from '@/composables/useToast'



interface Tag {
  tag_id: number
  name: string
  color: string
}

interface TagFormData {
}

export interface TagResponse {
  tags: Tag[]
}

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchTags() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedTags = await $api<TagResponse>('tags')
    if (fetchedTags) tags.value = fetchedTags.tags
    loading.value = false
  }


  async function submitForm(payload: { name: string; color: string; }) {
    // Create optimistic tag
    const optimisticTag: Tag = {
      tag_id: -Date.now(), // Use negative ID to avoid conflicts
      name: payload.name,
      color: payload.color
    }

    // Optimistically add to local state
    tags.value.push(optimisticTag)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Tag>('tags', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic tag with real data
      const optimisticIndex = tags.value.findIndex(t => t.tag_id === optimisticTag.tag_id)
      if (optimisticIndex !== -1) {
        tags.value[optimisticIndex] = data
      }

      addToast("Tag added succesfully", "success")

    } catch (err) {
      // Remove optimistic tag on error
      tags.value = tags.value.filter(t => t.tag_id !== optimisticTag.tag_id)
      addToast("Tag not added", "error")
    }
  }

  async function editTag(tag: Tag) {
    // Store original tag for potential rollback
    const originalTagIndex = tags.value.findIndex(t => t.tag_id === tag.tag_id)
    const originalTag = originalTagIndex !== -1 ? { ...tags.value[originalTagIndex] } : null

    // Optimistically update the tag
    if (originalTagIndex !== -1) {
      tags.value[originalTagIndex] = { ...tags.value[originalTagIndex], ...tag }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Tag>(`tags/${tag.tag_id}`, {
        method: 'PATCH',
        body: JSON.stringify(tag)
      })

      // Update with server response to ensure consistency
      if (originalTagIndex !== -1 && data) {
        tags.value[originalTagIndex] = data
      }

      addToast("Edited tag succesfully", "success")

    } catch (err) {
      // Revert optimistic update on error
      if (originalTag && originalTagIndex !== -1) {
        tags.value[originalTagIndex] = originalTag
      }
      addToast("Editing tag failed", "error")
    }
  }

  async function deleteTag(id: number) {
    // Store the tag for potential rollback
    const tagToDelete = tags.value.find(t => t.tag_id === id)
    if (!tagToDelete) {
      addToast("Tag not found", "error")
      return
    }

    // Optimistically remove from local state
    tags.value = tags.value.filter((t) => t.tag_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`tags/${id}`, {
        method: 'DELETE'
      })

      addToast("Tag deleted succesfully", "success")

    } catch (err) {
      // Restore the tag on error
      tags.value.push(tagToDelete)
      addToast("Tag deletion failed", "error")
    }
  }

  function reset() {
    tags.value = []
  }

  return {
    tags,
    loading,
    fetchError,
    fetchTags,
    editTag,
    deleteTag,
    submitForm,
    reset,
  }
})
