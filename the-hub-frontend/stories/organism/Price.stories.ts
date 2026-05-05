import type { Meta, StoryObj } from '@nuxtjs/storybook'

import Price from '~/components/Pricing'

const meta = {
  component: Price,
} satisfies Meta<typeof Price>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
}
