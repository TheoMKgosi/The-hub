import type { Meta, StoryObj } from '@nuxtjs/storybook';
 
import Task from '~/components/task/Task.vue';
 
const meta = {
  component: Task,
} satisfies Meta<typeof Task>;
 
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    title: 'Title'
  }
}
