<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import { useCategoryStore } from "@/stores/finance";

const categoryStore = useCategoryStore()

const searchQuery = ref('')
const formData = reactive({
  name: ''
})

const filteredTasks = computed(() => {
  let result = categoryStore.categories

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(task =>
      task.title.toLowerCase().includes(query) ||
      task.description.toLowerCase().includes(query)
    )
  }

  return result
})

const submitForm = async () => {
  const dataToSend = { ...formData }
  categoryStore.submitForm(dataToSend)
  Object.assign(formData, {
    name: ''
  })
}

onMounted(() => {
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})
</script>

<template>
  <div>
    <!-- Filters -->
    <!-- Filters + Search -->
    <div class="shadow-sm p-3 bg-white/20 backdrop-blur-md rounded-lg mt-2">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <input v-model="searchQuery" placeholder="Search tasks..." class="flex-grow shadow-sm  bg-gradient-to-r from-gray-50 to-gray-100 px-3 py-2 rounded
          w-full sm:w-0" />
      </div>
    </div>

    <div class="grid grid-cols-3">
      <!-- Adding card for Category -->
      <div>
        <form @submit.prevent="submitForm">
          <label for="name">Name</label>
          <input type="text" id="name" v-model="formData.name">
          <button type="submit">Create Category</button>
        </form>
      </div>


      <!-- Category cards -->
      <div v-if="categoryStore.categories.length === 0">The are no categories</div>
      <div v-for="category in categoryStore.categories">
        <h2>{{ category.name }}</h2>
      </div>
    </div>
  </div>
</template>
