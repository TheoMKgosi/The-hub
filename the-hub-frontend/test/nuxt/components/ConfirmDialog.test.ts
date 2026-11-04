import { describe, it, expect } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import ConfirmDialog from '~/components/ui/ConfirmDialog.vue'

describe('ConfirmDialog', () => {
  it('renders nothing when show is false', async () => {
    const wrapper = await mountSuspended(ConfirmDialog, {
      props: { show: false }
    })

    expect(wrapper.text()).toBe('')
  })

  it('renders the message when show is true', async () => {
    const wrapper = await mountSuspended(ConfirmDialog, {
      props: { show: true, message: 'Delete this task?' }
    })

    expect(wrapper.text()).toContain('Delete this task?')
  })

  it('falls back to the default message', async () => {
    const wrapper = await mountSuspended(ConfirmDialog, {
      props: { show: true }
    })

    expect(wrapper.text()).toContain('Are you sure?')
  })

  it('emits update:show false when Cancel is clicked', async () => {
    const wrapper = await mountSuspended(ConfirmDialog, {
      props: { show: true }
    })

    const buttons = wrapper.findAll('button')
    await buttons.find(b => b.text() === 'Cancel')!.trigger('click')

    expect(wrapper.emitted('update:show')).toEqual([[false]])
    expect(wrapper.emitted('confirm')).toBeUndefined()
  })

  it('emits confirm and update:show false when Confirm is clicked', async () => {
    const wrapper = await mountSuspended(ConfirmDialog, {
      props: { show: true }
    })

    const buttons = wrapper.findAll('button')
    await buttons.find(b => b.text() === 'Confirm')!.trigger('click')

    expect(wrapper.emitted('confirm')).toHaveLength(1)
    expect(wrapper.emitted('update:show')).toEqual([[false]])
  })
})
