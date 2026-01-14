import { type Meta, type StoryObj } from '@storybook/vue3-vite';
import dateSlots from '~/components/task/DateSlots.vue';

const meta = {
  component: dateSlots,
  // The 'render' function belongs at the top level of the object, not inside 'args'
} satisfies Meta<typeof dateSlots>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    label: "Today"
  }
};
