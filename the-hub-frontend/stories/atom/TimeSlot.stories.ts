import { type Meta, type StoryObj } from '@storybook/vue3-vite';
import timeSlot from '~/components/task/TimeSlot.vue';

const meta = {
  component: timeSlot,
  // The 'render' function belongs at the top level of the object, not inside 'args'
} satisfies Meta<typeof timeSlot>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    time: "1:00"
  }
};
