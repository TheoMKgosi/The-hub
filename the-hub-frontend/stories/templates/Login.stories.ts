import type { Meta, StoryObj } from '@nuxtjs/storybook'
import Login from '~/pages/login.vue';

const meta = {
  component: Login
} satisfies Meta<typeof Login>

type Story = StoryObj<typeof meta>

export default meta

export const Default: Story = {
}

