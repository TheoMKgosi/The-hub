interface Toast {
  id: number
  message: string
  type: string
  duration: number
}

const toasts: Ref<Toast[]> = ref([])

export function useToast() {
  const addToast = (message: string, type: string = 'info', duration: number = 3000): void => {
    const toast: Toast = {
      id: Date.now(),
      message,
      type,
      duration
    }
    toasts.value.push(toast)

    if (duration > 0) {
      setTimeout(() => removeToast(toast.id), duration)
    }
  }

  const removeToast = (id: number): void => {
    const index = toasts.value.findIndex(toast => toast.id === id)
    if (index > -1) toasts.value.splice(index, 1)
  }

  return {
    toasts,
    addToast,
    removeToast
  }
}

