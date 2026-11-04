import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import OfflinePage from '~/pages/offline.vue'

const reloadMock = vi.fn()
const alertMock = vi.fn()

beforeEach(() => {
  vi.stubGlobal('alert', alertMock)
  Object.defineProperty(window, 'location', {
    configurable: true,
    value: { reload: reloadMock }
  })
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('offline page', () => {
  it('renders the offline message and retry button', async () => {
    const wrapper = await mountSuspended(OfflinePage)

    expect(wrapper.text()).toContain("You're Offline")
    expect(wrapper.text()).toContain("It looks like you're not connected to the internet")
    expect(wrapper.text()).toContain('Try Again')
  })

  it('reloads the page when retrying while online', async () => {
    vi.stubGlobal('navigator', { onLine: true })
    const wrapper = await mountSuspended(OfflinePage)

    await wrapper.findAll('button').find(b => b.text() === 'Try Again')!.trigger('click')

    expect(reloadMock).toHaveBeenCalledTimes(1)
    expect(alertMock).not.toHaveBeenCalled()
  })

  it('alerts instead of reloading when still offline', async () => {
    vi.stubGlobal('navigator', { onLine: false })
    const wrapper = await mountSuspended(OfflinePage)

    await wrapper.findAll('button').find(b => b.text() === 'Try Again')!.trigger('click')

    expect(reloadMock).not.toHaveBeenCalled()
    expect(alertMock).toHaveBeenCalledWith('You are still offline. Please check your internet connection.')
  })
})
