import { type Meta, type StoryObj } from '@storybook/vue3-vite';
import dateSlots from '~/components/task/DateSlots.vue';

const meta = {
  component: dateSlots,
  // The 'render' function belongs at the top level of the object, not inside 'args'
} satisfies Meta<typeof dateSlots>;

export default meta;
type Story = StoryObj<typeof meta>;

const mockTask = [
  {
    start_time: 1,
    end_time: 2,
    title: "My favourite thing"
  },
  {
    start_time: 5,
    end_time: 5,
    title: "Stuff"
  },
  {
    start_time: 3,
    end_time: 3,
    title: "Theo"
  }
]

export const Default: Story = {
  args: {
    label: "Today",
    tasks: mockTask
  }
};
