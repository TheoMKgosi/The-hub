import type { Meta, StoryObj } from '@nuxtjs/storybook'
import Settings from '~/pages/settings.vue';

const meta = {
  component: Settings
} satisfies Meta<typeof Settings>

type Story = StoryObj<typeof meta>

export default meta

export const Default: Story = { }

