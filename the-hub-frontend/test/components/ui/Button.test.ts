import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { computed } from 'vue'
import Button from '~/components/ui/Button.vue'

describe('Button', () => {
  it('renders with default props', () => {
    const wrapper = mount(Button)
    expect(wrapper.element.tagName).toBe('BUTTON')
    expect(wrapper.classes()).toContain('inline-flex')
    expect(wrapper.classes()).toContain('items-center')
    expect(wrapper.classes()).toContain('justify-center')
  })

  it('renders slot content', () => {
    const wrapper = mount(Button, {
      slots: {
        default: 'Click me'
      }
    })
    expect(wrapper.text()).toBe('Click me')
  })

  it('applies correct variant classes', () => {
    const variants = ['primary', 'secondary', 'danger', 'default'] as const

    variants.forEach(variant => {
      const wrapper = mount(Button, {
        props: { variant }
      })

      if (variant === 'primary') {
        expect(wrapper.classes()).toContain('bg-primary')
        expect(wrapper.classes()).toContain('text-white')
      } else if (variant === 'secondary') {
        expect(wrapper.classes()).toContain('bg-secondary')
        expect(wrapper.classes()).toContain('text-white')
      } else if (variant === 'danger') {
        expect(wrapper.classes()).toContain('bg-red-500')
        expect(wrapper.classes()).toContain('text-white')
      } else if (variant === 'default') {
        expect(wrapper.classes()).toContain('bg-gray-200')
        expect(wrapper.classes()).toContain('text-gray-900')
      }
    })
  })

  it('applies correct size classes', () => {
    const sizes = ['sm', 'md', 'lg'] as const

    sizes.forEach(size => {
      const wrapper = mount(Button, {
        props: { size }
      })

      if (size === 'sm') {
        expect(wrapper.classes()).toContain('px-3')
        expect(wrapper.classes()).toContain('py-1.5')
        expect(wrapper.classes()).toContain('text-sm')
      } else if (size === 'md') {
        expect(wrapper.classes()).toContain('px-4')
        expect(wrapper.classes()).toContain('py-2')
        expect(wrapper.classes()).toContain('text-base')
      } else if (size === 'lg') {
        expect(wrapper.classes()).toContain('px-6')
        expect(wrapper.classes()).toContain('py-3')
        expect(wrapper.classes()).toContain('text-lg')
      }
    })
  })

  it('handles disabled state', () => {
    const wrapper = mount(Button, {
      props: { disabled: true }
    })

    expect(wrapper.attributes('disabled')).toBeDefined()
    expect(wrapper.classes()).toContain('disabled:opacity-50')
    expect(wrapper.classes()).toContain('disabled:cursor-not-allowed')
  })

  it('sets correct button type', () => {
    const types = ['button', 'submit', 'reset'] as const

    types.forEach(type => {
      const wrapper = mount(Button, {
        props: { type }
      })
      expect(wrapper.attributes('type')).toBe(type)
    })
  })

  it('emits click events', async () => {
    const wrapper = mount(Button)
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('does not emit click events when disabled', async () => {
    const wrapper = mount(Button, {
      props: { disabled: true }
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })
})