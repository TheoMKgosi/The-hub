import type { Meta, StoryObj } from '@nuxtjs/storybook'
import TaskTemplate from '~/components/task/TaskTemplate.vue';
import Nav from '~/components/ui/Nav.vue';
import Task from '~/types/task'

const mockTasks: Omit<Task, 'time_spent_minutes' | 'is_recurring'>[] = [
  {
    task_id: '1',
    title: 'Complete project documentation',
    description: 'Write comprehensive documentation for the new feature',
    status: 'in_progress',
    due_date: new Date('2023-12-15'),
    priority: 1,
    time_estimate_minutes: 120
  },
  {
    task_id: '2',
    title: 'Fix login page bug',
    description: 'Investigate and fix the authentication error on mobile devices',
    status: 'todo',
    due_date: new Date('2023-12-10'),
    priority: 3,
    time_estimate_minutes: 60
  },
  {
    task_id: '3',
    title: 'Review PR #42',
    description: 'Review the pull request for the new API endpoints',
    status: 'done',
    due_date: new Date('2023-12-05'),
    priority: 5,
    time_estimate_minutes: 30
  }
];

const meta = {
  component: TaskTemplate
} satisfies Meta<typeof TaskTemplate>

type Story = StoryObj<typeof meta>

export default meta

export const Default: Story = {
  args: {
    taskList: mockTasks
  }
}

export const WithNav: Story = {
  render: () => ({
    components: { TaskTemplate, Nav },
    setup() {
      return { mockTasks }
    },
    template: ' <div><Nav /><div class="m-64"><TaskTemplate :taskList="mockTasks"/></div> </div> '
  }),
}
