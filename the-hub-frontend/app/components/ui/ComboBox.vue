<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue";

interface Category {
  budget_category_id: string;
  name: string;
}

interface Props {
  modelValue?: string;
  categories: Category[];
  placeholder?: string;
  disabled?: boolean;
  loading?: boolean;
  allowCreate?: boolean;
}

interface Emits {
  (e: "update:modelValue", value: string): void;
  (e: "select", category: Category): void;
  (e: "create", name: string): void;
  (e: "focus"): void;
  (e: "blur"): void;
  (e: "clear"): void;
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: "Select or type...",
  disabled: false,
  loading: false,
  allowCreate: true,
});

const emit = defineEmits<Emits>();

const searchQuery = ref("");
const isOpen = ref(false);
const highlightedIndex = ref(-1);
const inputRef = ref<HTMLInputElement>();
const blurTimeout = ref<NodeJS.Timeout | null>(null);

// Computed properties
const filteredCategories = computed(() => {
  if (!searchQuery.value.trim()) {
    return props.categories;
  }

  const query = searchQuery.value.toLowerCase();
  return props.categories.filter((category) =>
    category.name.toLowerCase().includes(query),
  );
});

const showCreateOption = computed(() => {
  return (
    props.allowCreate &&
    searchQuery.value.trim() &&
    !filteredCategories.value.some(
      (cat) => cat.name.toLowerCase() === searchQuery.value.toLowerCase(),
    )
  );
});

const displayValue = computed(() => {
  if (props.modelValue) {
    const category = props.categories.find(
      (cat) => cat.budget_category_id === props.modelValue,
    );
    return category?.name || "";
  }
  return searchQuery.value;
});

// Methods
const openDropdown = () => {
  if (props.disabled) return;
  isOpen.value = true;
  highlightedIndex.value = -1;
  emit("focus");
};

const closeDropdown = () => {
  isOpen.value = false;
  highlightedIndex.value = -1;
  emit("blur");
};

const selectCategory = (category: Category) => {
  if (blurTimeout.value) {
    clearTimeout(blurTimeout.value);
    blurTimeout.value = null;
  }
  searchQuery.value = category.name;
  emit("update:modelValue", category.budget_category_id);
  emit("select", category);
  closeDropdown();
};

const createCategory = () => {
  if (blurTimeout.value) {
    clearTimeout(blurTimeout.value);
    blurTimeout.value = null;
  }
  const name = searchQuery.value.trim();
  if (name) {
    emit("create", name);
    closeDropdown();
  }
};

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  searchQuery.value = target.value;
  emit("update:modelValue", "");
  openDropdown();
};

const handleKeydown = (event: KeyboardEvent) => {
  if (!isOpen.value) {
    if (event.key === "Enter" || event.key === "ArrowDown") {
      event.preventDefault();
      openDropdown();
      return;
    }
  }

  switch (event.key) {
    case "ArrowDown":
      event.preventDefault();
      highlightedIndex.value = Math.min(
        highlightedIndex.value + 1,
        filteredCategories.value.length + (showCreateOption.value ? 0 : -1),
      );
      break;

    case "ArrowUp":
      event.preventDefault();
      highlightedIndex.value = Math.max(highlightedIndex.value - 1, -1);
      break;

    case "Enter":
      event.preventDefault();
      if (highlightedIndex.value === -1 && showCreateOption.value) {
        createCategory();
      } else if (
        highlightedIndex.value >= 0 &&
        highlightedIndex.value < filteredCategories.value.length
      ) {
        selectCategory(filteredCategories.value[highlightedIndex.value]);
      }
      break;

    case "Escape":
      event.preventDefault();
      closeDropdown();
      break;

    case "Tab":
      closeDropdown();
      break;
  }
};

const handleFocus = () => {
  if (!props.disabled) {
    openDropdown();
  }
};

const handleClickOutside = () => {
  // Delay closing to allow click events on dropdown items to fire first
  blurTimeout.value = setTimeout(() => {
    closeDropdown();
  }, 200);
};

const clearSelection = () => {
  searchQuery.value = "";
  emit("update:modelValue", "");
  emit("clear");
  closeDropdown();
};

// Watch for external model value changes
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue) {
      const category = props.categories.find(
        (cat) => cat.budget_category_id === newValue,
      );
      if (category) {
        searchQuery.value = category.name;
      }
    } else {
      searchQuery.value = "";
    }
  },
);

// Focus input when dropdown opens
watch(isOpen, (open) => {
  if (open) {
    nextTick(() => {
      inputRef.value?.focus();
    });
  }
});
</script>

<template>
  <div class="relative">
    <input ref="inputRef" v-model="displayValue" :placeholder="placeholder" :disabled="disabled" :class="[
      'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary',
      'border-surface-light dark:border-surface-dark',
      'bg-surface-light dark:bg-surface-dark',
      'text-text-light dark:text-text-dark',
      'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
      disabled && 'opacity-50 cursor-not-allowed',
    ]" @input="handleInput" @keydown="handleKeydown" @focus="handleFocus" @blur="handleClickOutside"
      autocomplete="off" />

    <!-- Loading indicator -->
    <div v-if="loading" class="absolute right-3 top-1/2 transform -translate-y-1/2">
      <div class="animate-spin h-4 w-4 border-2 border-primary border-t-transparent rounded-full"></div>
    </div>

    <!-- Clear button -->
    <button v-if="modelValue && !loading" 
      type="button"
      @click="clearSelection"
      class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 p-1 rounded">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
      </svg>
    </button>

    <!-- Dropdown -->
    <div v-if="isOpen && (filteredCategories.length > 0 || showCreateOption)"
      class="absolute z-[9999] w-full mt-1 bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark rounded-md shadow-lg max-h-60 overflow-auto">
      <!-- Existing categories -->
      <div v-for="(category, index) in filteredCategories" :key="category.budget_category_id" :class="[
        'px-3 py-2 cursor-pointer transition-colors',
        'hover:bg-primary/10 dark:hover:bg-primary/20',
        highlightedIndex === index && 'bg-primary/10 dark:bg-primary/20',
      ]" @click.stop="selectCategory(category)" @mousedown.stop @mouseenter="highlightedIndex = index">
        {{ category.name }}
      </div>

      <!-- Create new option -->
      <div v-if="showCreateOption" :class="[
        'px-3 py-2 cursor-pointer transition-colors border-t',
        'border-surface-light dark:border-surface-dark',
        'hover:bg-green-50 dark:hover:bg-green-900/20',
        'text-green-700 dark:text-green-300',
        highlightedIndex === filteredCategories.length &&
        'bg-green-50 dark:bg-green-900/20',
      ]" @click.stop="createCategory" @mousedown.stop @mouseenter="highlightedIndex = filteredCategories.length">
        Create: "{{ searchQuery }}"
      </div>
    </div>

    <!-- No results -->
    <div v-if="isOpen && filteredCategories.length === 0 && !showCreateOption"
      class="absolute z-50 w-full mt-1 bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark rounded-md shadow-lg p-3 text-center text-text-light dark:text-text-dark/60">
      No categories found
    </div>
  </div>
</template>
