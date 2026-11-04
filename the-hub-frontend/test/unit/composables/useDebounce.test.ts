import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { debounce } from '../../../app/composables/useDebounce'

describe('debounce', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('delays invocation until the wait time has elapsed', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 300)

    debounced('a')

    vi.advanceTimersByTime(299)
    expect(fn).not.toHaveBeenCalled()

    vi.advanceTimersByTime(1)
    expect(fn).toHaveBeenCalledTimes(1)
    expect(fn).toHaveBeenCalledWith('a')
  })

  it('collapses rapid successive calls into one invocation with the latest args', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 300)

    debounced('first')
    vi.advanceTimersByTime(100)
    debounced('second')
    vi.advanceTimersByTime(100)
    debounced('third')

    vi.advanceTimersByTime(300)

    expect(fn).toHaveBeenCalledTimes(1)
    expect(fn).toHaveBeenCalledWith('third')
  })

  it('defaults to a 300ms wait', () => {
    const fn = vi.fn()
    const debounced = debounce(fn)

    debounced()
    vi.advanceTimersByTime(299)
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(1)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('fires immediately on the leading edge when immediate is true', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 300, true)

    debounced('first')
    expect(fn).toHaveBeenCalledTimes(1)
    expect(fn).toHaveBeenCalledWith('first')

    // Calls within the wait window are swallowed (no trailing call)
    debounced('second')
    vi.advanceTimersByTime(300)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('allows a new immediate call once the wait window has passed', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 300, true)

    debounced('first')
    vi.advanceTimersByTime(400)
    debounced('second')

    expect(fn).toHaveBeenCalledTimes(2)
    expect(fn).toHaveBeenNthCalledWith(2, 'second')
  })

  it('cancel prevents the pending invocation', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 300)

    debounced('a')
    debounced.cancel()
    vi.advanceTimersByTime(1000)

    expect(fn).not.toHaveBeenCalled()
  })

  it('keeps working after cancel', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 300)

    debounced('a')
    debounced.cancel()
    debounced('b')
    vi.advanceTimersByTime(300)

    expect(fn).toHaveBeenCalledTimes(1)
    expect(fn).toHaveBeenCalledWith('b')
  })
})
