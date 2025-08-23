import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useToast } from '~/composables/useToast'

// Mock timer functions
beforeEach(() => {
  vi.useFakeTimers()
  // Clear toasts before each test to ensure isolation
  const { clearToasts } = useToast()
  clearToasts()
})

describe('useToast', () => {
  it('should initialize with empty toasts array', () => {
    const { toasts } = useToast()
    expect(toasts.value).toEqual([])
  })

  it('should add a toast with default values', () => {
    const { toasts, addToast } = useToast()

    addToast('Test message')

    expect(toasts.value).toHaveLength(1)
    expect(toasts.value[0]).toMatchObject({
      message: 'Test message',
      type: 'info',
      duration: 3000
    })
    expect(typeof toasts.value[0].id).toBe('number')
  })

  it('should add a toast with custom values', () => {
    const { toasts, addToast } = useToast()

    addToast('Custom message', 'success', 5000)

    expect(toasts.value).toHaveLength(1)
    expect(toasts.value[0]).toMatchObject({
      message: 'Custom message',
      type: 'success',
      duration: 5000
    })
  })

  it('should remove toast after duration', () => {
    const { toasts, addToast } = useToast()

    addToast('Test message', 'info', 1000)

    expect(toasts.value).toHaveLength(1)

    // Fast-forward time
    vi.advanceTimersByTime(1000)

    expect(toasts.value).toHaveLength(0)
  })

  it('should not auto-remove toast with duration 0', () => {
    const { toasts, addToast } = useToast()

    addToast('Persistent message', 'info', 0)

    expect(toasts.value).toHaveLength(1)

    // Fast-forward time
    vi.advanceTimersByTime(5000)

    expect(toasts.value).toHaveLength(1)
  })

  it('should remove toast manually', () => {
    const { toasts, addToast, removeToast } = useToast()

    addToast('Test message')
    const toastId = toasts.value[0].id

    expect(toasts.value).toHaveLength(1)

    removeToast(toastId)

    expect(toasts.value).toHaveLength(0)
  })

  it('should handle removing non-existent toast', () => {
    const { toasts, addToast, removeToast } = useToast()

    addToast('Test message')
    expect(toasts.value).toHaveLength(1)

    // Try to remove a non-existent toast
    removeToast(999)

    // Should not affect existing toasts
    expect(toasts.value).toHaveLength(1)
  })

  it('should generate unique IDs for multiple toasts', () => {
    const { toasts, addToast } = useToast()

    // Advance time between each toast to ensure unique IDs
    addToast('First message')
    vi.advanceTimersByTime(1)
    addToast('Second message')
    vi.advanceTimersByTime(1)
    addToast('Third message')

    expect(toasts.value).toHaveLength(3)

    const ids = toasts.value.map(toast => toast.id)
    const uniqueIds = new Set(ids)

    expect(uniqueIds.size).toBe(3)
  })

  it('should handle multiple toasts with different durations', () => {
    const { toasts, addToast } = useToast()

    addToast('Short toast', 'info', 1000)
    addToast('Long toast', 'info', 3000)

    expect(toasts.value).toHaveLength(2)

    // Fast-forward 1 second
    vi.advanceTimersByTime(1000)

    expect(toasts.value).toHaveLength(1)
    expect(toasts.value[0].message).toBe('Long toast')

    // Fast-forward another 2 seconds
    vi.advanceTimersByTime(2000)

    expect(toasts.value).toHaveLength(0)
  })
})