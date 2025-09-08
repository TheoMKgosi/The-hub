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

const deleteCategory = (id: string) => {
  categoryStore.deleteCategory(id)
}

onMounted(() => {
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})
</script>
<template>
  <div class="p-4 max-w-6xl mx-auto space-y-6">
    <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Category Management</h2>

    <!-- Filters + Search -->
    <div class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search categories..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <!-- Main Grid -->
    <div class="grid gap-6 lg:grid-cols-4">
      <!-- Create Category Card -->
      <div class="p-6 bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark lg:col-span-1">
        <h3 class="font-semibold text-lg mb-4 text-text-light dark:text-text-dark">Create Category</h3>
        <form @submit.prevent="submitForm" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Category Name</label>
            <input type="text" id="name" v-model="formData.name" placeholder="Enter category name"
              class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <UiButton type="submit" variant="primary" size="md" class="w-full">
            Create Category
          </UiButton>
        </form>
      </div>

      <!-- Category Grid -->
      <div class="lg:col-span-3">
        <div v-if="categoryStore.categories.length === 0"
          class="p-8 text-center text-text-light dark:text-text-dark/60 bg-surface-light dark:bg-surface-dark rounded-lg border border-surface-light dark:border-surface-dark">
          <p class="text-lg mb-2">No categories yet</p>
          <p class="text-sm">Create your first category to get started</p>
        </div>
        <div v-else>
          <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Your Categories</h3>
          <div class="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3">
            <div v-for="category in filteredCategory" :key="category.budget_category_id"
              class="p-4 bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark hover:shadow-lg transition-all duration-200 cursor-pointer group"
              @dblclick="showDialog = true; categoryID = category.budget_category_id; categoryName = category.name;">
              <div class="flex items-center justify-between">
                <h4 class="font-semibold text-lg text-text-light dark:text-text-dark group-hover:text-primary dark:group-hover:text-primary transition-colors">
                  {{ category.name }}
                </h4>
                <UiButton variant="danger" size="sm" class="opacity-0 group-hover:opacity-100 transition-opacity"
                  @click.stop="showDialog = true; categoryID = category.budget_category_id; categoryName = category.name;">
                  Delete
                </UiButton>
              </div>
              <p class="text-sm text-text-light dark:text-text-dark/60 mt-2">Double-click to edit</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog v-model:show="showDialog" :message="`Are you sure you want to delete the category '${categoryName}'?`"
      @confirm="deleteCategory(categoryID)"/>
  </div>
</template>
