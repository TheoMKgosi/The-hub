import { type Meta, type StoryObj } from "@storybook/vue3-vite";
import timeSlot from "~/components/task/TimeSlot.vue";

const meta = {
  component: timeSlot,
  argTypes: {
    complete: {
      // 'boolean' is default, but 'select' or 'radio' handles null better
      control: { type: 'select' },
      options: [true, false, null],
      labels: {
        true: 'Completed',
        false: 'In Progress',
        null: 'Not Started (Null)'
      },
      description: 'The status of the task',
    },
    position: { control: 'text'}

  }
} satisfies Meta<typeof timeSlot>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    complete: null,
    time: 1,
  },
};

export const Title: Story = {
  args: {
    time: 1,
    title: "Task to be done",
  },
};

export const TitleMuted: Story = {
  args: {
    time: 1,
    title: "Task to be done",
    mute: true,
  },
};

export const WithColor: Story = {
  args: {
    time: 1,
    title: "Task to be done",
    color: "#c73030",
  },
};
