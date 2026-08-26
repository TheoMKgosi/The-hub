<script setup lang="ts">
defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  close: []
  startFocus: []
}>()

const focusStore = useFocusSessionStore()
const todayMinutes = computed(() => focusStore.getTodayMinutes)
</script>

<template>
  <div>
    <!-- Desktop slide-in panel -->
    <div class="hidden md:block fixed right-0 top-0 bottom-0 z-30 w-80 bg-surface-light dark:bg-surface-dark border-l border-surface-light/20 dark:border-surface-dark/20 shadow-2xl overflow-y-auto custom-scrollbar transition-transform duration-300"
      :class="show ? 'translate-x-0' : 'translate-x-full'">
      <div class="p-4">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-lg bg-primary/10 dark:bg-primary/20 flex items-center justify-center">
              <UIcon name="i-lucide-timer" class="w-4 h-4 text-primary" />
            </div>
            <h3 class="text-sm font-semibold text-text-light dark:text-text-dark">Focus</h3>
          </div>
          <button @click="emit('close')" class="p-1.5 rounded-lg hover:bg-surface-light/10 dark:hover:bg-surface-dark/10 text-text-light/50 dark:text-text-dark/50 hover:text-text-light dark:hover:text-text-dark transition-colors">
            <UIcon name="i-lucide-chevron-right" class="w-4 h-4" />
          </button>
        </div>

        <UButton block label="Start Focus Session" icon="i-lucide-play" color="primary" size="lg" class="mb-4" @click="emit('startFocus')" />

        <FocusSessionStats />
      </div>
    </div>

    <!-- Desktop toggle tab (visible when panel is closed) -->
    <div class="hidden md:flex fixed right-0 top-1/2 -translate-y-1/2 z-20 transition-opacity duration-300"
      :class="show ? 'opacity-0 pointer-events-none' : 'opacity-100'">
      <button @click="emit('startFocus')" title="Start focus session"
        class="flex flex-col items-center gap-1 px-2 py-4 bg-surface-light dark:bg-surface-dark border border-surface-light/20 dark:border-surface-dark/20 rounded-l-xl shadow-lg hover:shadow-xl transition-all hover:bg-primary/5 dark:hover:bg-primary/10 text-text-light dark:text-text-dark group">
        <div class="w-8 h-8 rounded-lg bg-primary/10 dark:bg-primary/20 flex items-center justify-center group-hover:bg-primary/20 dark:group-hover:bg-primary/30 transition-colors">
          <UIcon name="i-lucide-timer" class="w-4 h-4 text-primary" />
        </div>
        <div v-if="todayMinutes > 0" class="text-[10px] font-semibold text-primary leading-none">{{ Math.floor(todayMinutes / 60) }}h {{ todayMinutes % 60 }}m</div>
      </button>
    </div>

    <!-- Mobile bottom drawer -->
    <Teleport to="body">
      <div v-if="show" class="md:hidden fixed inset-0 z-50">
        <div class="absolute inset-0 bg-black/50" @click="emit('close')" />
        <div class="absolute bottom-0 left-0 right-0 bg-surface-light dark:bg-surface-dark rounded-t-2xl shadow-2xl max-h-[85vh] overflow-y-auto custom-scrollbar transition-transform duration-300 translate-y-0">
          <div class="sticky top-0 bg-surface-light dark:bg-surface-dark rounded-t-2xl pt-3 pb-2 px-4 border-b border-surface-light/20 dark:border-surface-dark/20">
            <div class="w-10 h-1 bg-text-light/20 dark:bg-text-dark/20 rounded-full mx-auto mb-3" />
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="w-8 h-8 rounded-lg bg-primary/10 dark:bg-primary/20 flex items-center justify-center">
                  <UIcon name="i-lucide-timer" class="w-4 h-4 text-primary" />
                </div>
                <h3 class="text-sm font-semibold text-text-light dark:text-text-dark">Focus</h3>
              </div>
              <button @click="emit('close')" class="p-1.5 rounded-lg hover:bg-surface-light/10 dark:hover:bg-surface-dark/10 text-text-light/50 dark:text-text-dark/50">
                <UIcon name="i-lucide-x" class="w-4 h-4" />
              </button>
            </div>
          </div>
          <div class="p-4">
            <UButton block label="Start Focus Session" icon="i-lucide-play" color="primary" size="lg" class="mb-4" @click="emit('startFocus'); emit('close')" />
            <FocusSessionStats />
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
