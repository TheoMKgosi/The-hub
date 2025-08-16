<script setup lang="ts">

const categoryStore = useCategoryStore()

const searchQuery = ref('')
const categoryID = ref(0)
const categoryName = ref('')
const formData = reactive({
  name: ''
})

const showDialog = ref(false)

const filteredCategory = computed(() => {
  let result = categoryStore.categories

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(category => category.name.toLowerCase().includes(query))
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

const deleteCategory = (id: number) => {
  categoryStore.deleteCategory(id)
}

onMounted(() => {
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})
</script>
<template>
  <div class="p-4 space-y-4">
    <!-- Filters + Search -->
    <div class="shadow-sm p-4 bg-white/30 backdrop-blur-md rounded-lg flex flex-wrap gap-2 items-center">
      <input v-model="searchQuery" placeholder="Search tasks..."
        class="flex-grow shadow-inner bg-gradient-to-r from-gray-50 to-gray-100 px-3 py-2 rounded border border-gray-200 focus:ring-2 focus:ring-blue-400 focus:outline-none" />
    </div>

    <!-- Main Grid -->
    <div class="grid gap-4 lg:grid-cols-4">
      <!-- Create Category Card -->
      <div class="p-4 bg-white rounded-lg shadow space-y-3 lg:col-span-1 max-h-72 md:max-h-56">
        <h3 class="font-semibold text-lg">Create Category</h3>
        <form @submit.prevent="submitForm" class="space-y-2">
          <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
          <input type="text" id="name" v-model="formData.name"
            class="w-full px-3 py-2 rounded border border-gray-300 focus:ring-2 focus:ring-blue-400 focus:outline-none" />
          <button type="submit" class="w-full py-2 px-4 rounded bg-blue-500 hover:bg-blue-600 text-white font-semibold">
            Create
          </button>
        </form>
      </div>

      <!-- Category Grid -->
      <div class="lg:col-span-3">
        <div v-if="categoryStore.categories.length === 0"
          class="p-4 text-gray-500 text-center bg-gray-50 rounded border border-gray-200">
          There are no categories
        </div>
        <div v-else class="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3">
          <div v-for="category in categoryStore.categories" :key="category.budget_category_id"
            class="p-4 bg-white rounded-lg shadow hover:shadow-md transition"
              @dblclick="showDialog = true; categoryID = category.budget_category_id; categoryName = category.name; ">
            <h2 class="font-semibold text-lg">{{ category.name }}</h2>
            <ConfirmDialog v-model:show="showDialog" :message="'Delete category ' + categoryName" 
            @confirm="deleteCategory(categoryID)"/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
