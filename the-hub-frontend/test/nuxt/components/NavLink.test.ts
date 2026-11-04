import { describe, it, expect } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import NavLink from '~/components/ui/NavLink.vue'

describe('NavLink', () => {
  it('renders a NuxtLink anchor when `to` is provided', async () => {
    const wrapper = await mountSuspended(NavLink, {
      props: { to: '/dashboard' },
      slots: { default: 'Dashboard' }
    })

    expect(wrapper.find('a').exists()).toBe(true)
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('renders a plain anchor when `href` is provided', async () => {
    const wrapper = await mountSuspended(NavLink, {
      props: { href: 'https://example.com' },
      slots: { default: 'External' }
    })

    const anchor = wrapper.find('a')
    expect(anchor.exists()).toBe(true)
    expect(anchor.attributes('href')).toBe('https://example.com')
  })

  it('renders a button with click emit when neither to nor href is provided', async () => {
    const wrapper = await mountSuspended(NavLink, {
      slots: { default: 'Toggle' }
    })

    expect(wrapper.find('button').exists()).toBe(true)

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toHaveLength(1)
  })

  it('applies base transition classes in every variant', async () => {
    const link = await mountSuspended(NavLink, {
      props: { to: '/a' },
      slots: { default: 'A' }
    })
    const tab = await mountSuspended(NavLink, {
      props: { variant: 'tab' },
      slots: { default: 'A' }
    })

    expect(link.find('a').classes()).toContain('transition-all')
    expect(tab.find('button').classes()).toContain('transition-all')
  })

  it('applies active tab styling when active is true', async () => {
    const wrapper = await mountSuspended(NavLink, {
      props: { variant: 'tab', active: true },
      slots: { default: 'A' }
    })

    expect(wrapper.find('button').classes()).toContain('border-b-2')
    expect(wrapper.find('button').classes()).toContain('text-primary')
  })

  it('applies inactive tab styling when active is false', async () => {
    const wrapper = await mountSuspended(NavLink, {
      props: { variant: 'tab', active: false },
      slots: { default: 'A' }
    })

    expect(wrapper.find('button').classes()).not.toContain('border-b-2')
  })
})
