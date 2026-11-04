import { describe, it, expect } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import FilterChips from '~/components/filter/FilterChips.vue'

const chips = [
  { key: 'status', label: 'Status', value: 'Pending' },
  { key: 'priority', label: 'Priority', value: 'P4, P5' }
]

describe('FilterChips', () => {
  it('renders one chip per entry with label and value', async () => {
    const wrapper = await mountSuspended(FilterChips, {
      props: { chips }
    })

    // wrapper.text() collapses whitespace between the label and value spans
    expect(wrapper.text()).toContain('Status:Pending')
    expect(wrapper.text()).toContain('Priority:P4, P5')
  })

  it('emits remove with the chip key when its remove button is clicked', async () => {
    const wrapper = await mountSuspended(FilterChips, {
      props: { chips }
    })

    const removeButtons = wrapper.findAll('button[title="Remove filter"]')
    expect(removeButtons).toHaveLength(2)

    await removeButtons[0].trigger('click')
    expect(wrapper.emitted('remove')).toEqual([['status']])
  })

  it('shows the clear all button only when there is more than one chip', async () => {
    const multi = await mountSuspended(FilterChips, { props: { chips } })
    const single = await mountSuspended(FilterChips, {
      props: { chips: [chips[0]] }
    })

    expect(multi.text()).toContain('Clear all')
    expect(single.text()).not.toContain('Clear all')
  })

  it('emits clearAll when the clear all button is clicked', async () => {
    const wrapper = await mountSuspended(FilterChips, {
      props: { chips }
    })

    const clearAll = wrapper.findAll('button').find(b => b.text() === 'Clear all')!
    await clearAll.trigger('click')

    expect(wrapper.emitted('clearAll')).toHaveLength(1)
  })

  it('renders nothing extra when the chip list is empty', async () => {
    const wrapper = await mountSuspended(FilterChips, {
      props: { chips: [] }
    })

    expect(wrapper.findAll('button')).toHaveLength(0)
  })
})
