<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { createFetch } from '@vueuse/core'
import { VueCal } from 'vue-cal'
import type { Event } from 'vue-cal'
import 'vue-cal/style'

const useMyFetch = createFetch({ baseUrl: import.meta.env.VITE_BASE_URL })
const events = ref<Event[]>([])

async function fetchEvents(start?: Date, end?: Date) {
  const params: string[] = []
  if (start) params.push(`start=${start.toISOString()}`)
  if (end) params.push(`end=${end.toISOString()}`)
  const url = `calendar${params.length ? `?${params.join('&')}` : ''}`
  const { data } = await useMyFetch(url).json<Event[]>()
  events.value = (data.value || []).map((e) => ({
    ...e,
    start: new Date(e.start),
    end: new Date(e.end),
  }))
}

onMounted(() => {
  fetchEvents()
})

function onViewChange(viewMeta: { start: Date; end: Date }) {
  fetchEvents(viewMeta.start, viewMeta.end)
}

async function onEventDropped(meta: { event: Event }) {
  console.log(meta.event)
  await useMyFetch(`calendar/${meta.event.task_id}`).patch({ start: meta.event.start }).json()
  fetchEvents()
}

async function onEventDelete(meta: { event: Event }) {
  await useMyFetch(`calendar/${meta.event.id}`).delete().json()
  fetchEvents()
}
</script>

<template>
  <h2>Calendar</h2>
  <VueCal
    :events="events"
    :editable-events="true"
    @view-change="onViewChange"
    @event-drop="onEventDropped"
    @event-delete="onEventDelete"
  />
</template>
