import type { Meta, StoryObj } from '@storybook/vue3-vite';
 
import ComboBox from '~/components/ui/BaseComboBox.vue';
 
const meta = {
  component: ComboBox,
} satisfies Meta<typeof ComboBox>;
 
export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    primary: true,
    label: 'ComboBox Primary',
  },
};

export const Secondary: Story = {
  args: {
    primary: true,
    label: 'ComboBox Secondary',
  },
};

export const Large: Story = {
  args: {
    primary: true,
    label: 'ComboBox Large',
  },
};

export const Small: Story = {
  args: {
    primary: true,
    label: 'ComboBox Small',
  },
};
