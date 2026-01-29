import type { Meta, StoryObj } from '@storybook/vue3-vite';
 
import Tabs from '~/components/ui/Tabs.vue';
 
const meta = {
  component: Tabs,
} satisfies Meta<typeof Tabs>;
 
export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    primary: true,
    label: 'Tabs',
  },
};

export const Secondary: Story = {
  args: {
    primary: true,
    label: 'Tabs Secondary',
  },
};
export const Large: Story = {
  args: {
    primary: true,
    label: 'Tabs Large',
  },
};
export const Small: Story = {
  args: {
    primary: true,
    label: 'Tabs Small',
  },
};
