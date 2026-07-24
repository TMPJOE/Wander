import { ref } from 'vue'

export interface Toast {
  id: number
  type: 'success' | 'error' | 'info'
  message: string
  duration?: number
}

const toasts = ref<Toast[]>([])
let nextId = 1

export function useToast() {
  const addToast = (type: 'success' | 'error' | 'info', message: string, duration = 3500) => {
    const id = nextId++
    const toast: Toast = { id, type, message, duration }
    toasts.value.push(toast)

    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }
  }

  const removeToast = (id: number) => {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  const success = (message: string, duration?: number) => addToast('success', message, duration)
  const error = (message: string, duration?: number) => addToast('error', message, duration)
  const info = (message: string, duration?: number) => addToast('info', message, duration)

  return {
    toasts,
    removeToast,
    success,
    error,
    info,
  }
}
