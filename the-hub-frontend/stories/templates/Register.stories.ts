import type { Meta, StoryObj } from '@nuxtjs/storybook'
import Register from '~/pages/register.vue';

const meta = {
  component: Register
} satisfies Meta<typeof Register>

type Story = StoryObj<typeof meta>

export default meta

export const Default: Story = {
}
