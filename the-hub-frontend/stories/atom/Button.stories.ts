import type { Meta, StoryObj } from '@storybook/vue3-vite'
// import { expect, fn } from 'storybook/test';
import BaseButton from '~/components/ui/BaseButton.vue'

const meta = {
  component: BaseButton,
  argTypes: {
    text: { control: 'text' },
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'danger', 'default'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
    type: {
      control: 'select',
      options: ['button', 'submit', 'reset'],
    },
    disabled: { control: 'boolean' },
  },
} satisfies Meta<typeof BaseButton>

export default meta
type Story = StoryObj<typeof meta>

export const Primary: Story = {
  args: {
    text: 'Primary',
    variant: 'primary',
  },
}

export const Danger: Story = {
  args: {
    variant: 'danger',
  },
}
