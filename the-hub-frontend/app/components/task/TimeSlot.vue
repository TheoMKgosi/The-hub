<script setup lang="ts">
interface Props {
  time: number
  title?: string | null
  mute?: boolean
  complete?: boolean | null
  position?: 'end' | 'start' | 'both' | ''
}

const props = withDefaults(defineProps<Props>(), {
  mute: false,
  complete: null,
  position: ''
}
)
const muted = props.mute ? 'text-gray-500' : ''
const withTitle = props.title == '' || props.title == null ? 'p-5 m-5' : 'p-3 mx-3'
const completed = computed(() => {
  if (props.complete == null) {
    return 'border-gray-500'
  } else if (props.complete === true) {
    return 'border-green-500'
  } else {
    return 'border-primary'
  }
})

const positionClass = computed(() => {
  if (props.position === '') {
    return ''
  } else if (props.position === 'end') {
    return 'rounded-bl-lg'
  } else if (props.position === 'start') {
    return 'rounded-tl-lg'
  } else {
    return 'rounded-l-lg'
  }
})
</script>
<template>
  <div class="dark:text-white border-l-2 " :class="completed, positionClass, withTitle">
    <div class="flex space-x-6">
      <p class="font-bold">{{ time }}:00</p>
      <p :class="muted">{{ title }}</p>
    </div>
  </div>
</template>
