import { ref, onMounted, onUnmounted } from 'vue'

export interface TouchGestureOptions {
  onSwipeLeft?: () => void
  onSwipeRight?: () => void
  onSwipeUp?: () => void
  onSwipeDown?: () => void
  onTap?: () => void
  threshold?: number
}

export function useTouchGestures(element: HTMLElement | null, options: TouchGestureOptions) {
  const touchStartX = ref(0)
  const touchStartY = ref(0)
  const touchEndX = ref(0)
  const touchEndY = ref(0)

  const threshold = options.threshold || 50

  const handleTouchStart = (e: TouchEvent) => {
    touchStartX.value = e.touches[0].clientX
    touchStartY.value = e.touches[0].clientY
  }

  const handleTouchEnd = (e: TouchEvent) => {
    touchEndX.value = e.changedTouches[0].clientX
    touchEndY.value = e.changedTouches[0].clientY

    const deltaX = touchEndX.value - touchStartX.value
    const deltaY = touchEndY.value - touchStartY.value

    const absDeltaX = Math.abs(deltaX)
    const absDeltaY = Math.abs(deltaY)

    // Determine if it's a swipe or tap
    if (Math.max(absDeltaX, absDeltaY) < threshold) {
      // It's a tap
      options.onTap?.()
      return
    }

    // Determine swipe direction
    if (absDeltaX > absDeltaY) {
      // Horizontal swipe
      if (deltaX > threshold) {
        options.onSwipeRight?.()
      } else if (deltaX < -threshold) {
        options.onSwipeLeft?.()
      }
    } else {
      // Vertical swipe
      if (deltaY > threshold) {
        options.onSwipeDown?.()
      } else if (deltaY < -threshold) {
        options.onSwipeUp?.()
      }
    }
  }

  onMounted(() => {
    if (element) {
      element.addEventListener('touchstart', handleTouchStart, { passive: true })
      element.addEventListener('touchend', handleTouchEnd, { passive: true })
    }
  })

  onUnmounted(() => {
    if (element) {
      element.removeEventListener('touchstart', handleTouchStart)
      element.removeEventListener('touchend', handleTouchEnd)
    }
  })

  return {
    touchStartX: readonly(touchStartX),
    touchStartY: readonly(touchStartY),
    touchEndX: readonly(touchEndX),
    touchEndY: readonly(touchEndY),
  }
}

// Composable for task swipe actions
export function useTaskSwipeActions(taskId: string) {
  const element = ref<HTMLElement | null>(null)

  const { onSwipeLeft, onSwipeRight } = useTouchGestures(element.value, {
    onSwipeLeft: () => {
      // Mark as complete
      console.log(`Mark task ${taskId} as complete`)
    },
    onSwipeRight: () => {
      // Mark as pending
      console.log(`Mark task ${taskId} as pending`)
    },
    threshold: 75,
  })

  return {
    element,
    onSwipeLeft,
    onSwipeRight,
  }
}