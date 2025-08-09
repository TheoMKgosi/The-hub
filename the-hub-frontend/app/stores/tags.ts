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
    try {
      const { $api } = useNuxtApp()
      const { data, error } = await $api('tags', {
        method: 'POST',
        body: payload
      })
      fetchTags()
      addToast("Tag added succesfully", "success")

    } catch (err) {
      addToast("Tag not added", "error")
    }
  }

  async function editTag(tag: Tag) {
    const { $api } = useNuxtApp()
    const { error } = await $api(`tags/${tag.tag_id}`).patch(tag).json()

    if (!error.value) {
      const index = tags.value.findIndex(t => t.tag_id === tag.tag_id)
      if (index !== -1) {
        tags.value[index] = { ...tags.value[index], ...tag }
        addToast("Edited tag succesfully", "success")
      } else {
        addToast("Editing tag failed", "error")
      }
    } else {
      addToast("Editing tag failed", "error")
    }
  }

  async function deleteTag(id: number) {
    const { $api } = useNuxtApp()
    await $api(`tags/${id}`).delete().json()
    tags.value = tags.value.filter((t) => t.tag_id !== id)
    addToast("Tag deleted succesfully", "success")
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
