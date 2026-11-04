import { describe, it, expect } from 'vitest'
import { ref, defineComponent } from 'vue'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import Tabs from '~/components/ui/Tabs.vue'

describe('Tabs', () => {
  it('renders one button per tab', async () => {
    const wrapper = await mountSuspended(Tabs, {
      props: { tabs: ['Overview', 'Stats'] }
    })

    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(2)
    expect(buttons[0].text()).toBe('Overview')
    expect(buttons[1].text()).toBe('Stats')
  })

  it('defaults the active tab to the first tab', async () => {
    const wrapper = await mountSuspended(Tabs, {
      props: { tabs: ['Overview', 'Stats'] }
    })

    const buttons = wrapper.findAll('button')
    expect(buttons[0].classes()).toContain('border-b-2')
    expect(buttons[1].classes()).not.toContain('border-b-2')
  })

  it('marks the provided modelValue tab as active', async () => {
    const wrapper = await mountSuspended(Tabs, {
      props: { tabs: ['Overview', 'Stats'], modelValue: 'Stats' }
    })

    const buttons = wrapper.findAll('button')
    expect(buttons[1].classes()).toContain('border-b-2')
    expect(buttons[0].classes()).not.toContain('border-b-2')
  })

  it('emits update:modelValue when a tab is clicked', async () => {
    const wrapper = await mountSuspended(Tabs, {
      props: { tabs: ['Overview', 'Stats'], modelValue: 'Overview' }
    })

    await wrapper.findAll('button')[1].trigger('click')

    expect(wrapper.emitted('update:modelValue')).toEqual([['Stats']])
  })

  it('renders the slot content of the active tab', async () => {
    const wrapper = await mountSuspended(Tabs, {
      props: { tabs: ['Overview', 'Stats'] },
      slots: { Overview: 'Overview content', Stats: 'Stats content' }
    })

    expect(wrapper.text()).toContain('Overview content')
    expect(wrapper.text()).not.toContain('Stats content')
  })

  it('switches rendered slot content through a v-model parent', async () => {
    const Parent = defineComponent({
      components: { Tabs },
      setup() {
        const active = ref('Overview')
        return { active }
      },
      template: `
        <Tabs v-model="active" :tabs="['Overview', 'Stats']">
          <template #Overview>Overview content</template>
          <template #Stats>Stats content</template>
        </Tabs>
      `
    })

    const wrapper = await mountSuspended(Parent)
    expect(wrapper.text()).toContain('Overview content')
    expect(wrapper.text()).not.toContain('Stats content')

    await wrapper.findAll('button')[1].trigger('click')

    expect(wrapper.text()).toContain('Stats content')
    expect(wrapper.text()).not.toContain('Overview content')
  })
})
