<script setup lang="ts">
import type { ShortcutAction } from '~/composables/useKeyboardShortcuts'

const { showHelp, closeHelp, shortcuts } = useKeyboardShortcuts()

interface ShortcutRow {
  description: string
  keys: string
}

const grouped = computed<Record<string, ShortcutRow[]>>(() => {
  const groups: Record<string, ShortcutRow[]> = {}
  shortcuts.value.forEach((shortcut: ShortcutAction) => {
    if (!groups[shortcut.category]) groups[shortcut.category] = []
    groups[shortcut.category].push({
      description: shortcut.description,
      keys: formatKeys(shortcut)
    })
  })
  return groups
})

const formatKeys = (shortcut: ShortcutAction) => {
  const parts: string[] = []
  if (shortcut.ctrl) parts.push('Ctrl')
  if (shortcut.alt) parts.push('Alt')
  if (shortcut.shift) parts.push('Shift')
  if (shortcut.meta) parts.push('⌘')
  parts.push(shortcut.key === ' ' ? 'Space' : shortcut.key.length === 1 ? shortcut.key.toUpperCase() : shortcut.key)
  return parts.join(' + ')
}
</script>

<template>
  <UModal v-model:open="showHelp" title="Keyboard Shortcuts">
    <template #body>
      <div class="space-y-5">
        <div v-for="(rows, category) in grouped" :key="category">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-text-light/60 dark:text-text-dark/60 mb-2">
            {{ category }}
          </h3>
          <div class="divide-y divide-surface-light/20 dark:divide-surface-dark/20">
            <div v-for="row in rows" :key="row.keys + row.description"
              class="flex items-center justify-between py-2 gap-4">
              <span class="text-sm text-text-light dark:text-text-dark">{{ row.description }}</span>
              <kbd
                class="px-2 py-0.5 bg-surface-light/50 dark:bg-surface-dark/50 rounded text-xs text-text-light/70 dark:text-text-dark/70 whitespace-nowrap">
                {{ row.keys }}
              </kbd>
            </div>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>
