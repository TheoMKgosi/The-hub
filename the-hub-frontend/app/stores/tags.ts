import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useToast } from '@/composables/useToast'



interface Tag {
  tag_id: number
  name: string
  color: string
}

interface TagFormData {
  name: string
  color: string
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
    loading.value = true
    const { data, error } = await useMyFetch('tags').json<TagResponse>()

    if (data.value) tags.value = data.value.tags
    fetchError.value = error.value

    loading.value = false
  }


  async function submitForm(formData: TagFormData) {
    const { data, error } = await useMyFetch('tags').post(formData).json()
    if (!data.value.tag_id) {
      data.value.tag_id = Date.now() // fallback if backend didn’t return ID
    }
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Tag not added", "error")
    } else {
      tags.value.push(data.value)
      addToast("Tag added succesfully", "success")
    }
  }

  async function editTag(tag: Tag) {
    const { error } = await useMyFetch(`tags/${tag.tag_id}`).patch(tag).json()

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
    await useMyFetch(`tags/${id}`).delete().json()
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
