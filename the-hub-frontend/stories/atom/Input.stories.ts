import type { Meta, StoryObj } from '@storybook/vue3-vite';

import Input from '~/components/ui/BaseInput.vue';

const meta = {
  component: Input,
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    primary: true,
    label: 'Input Primary',
  },
};

export const Secondary: Story = {
  args: {
    primary: true,
    label: 'Input Secondary',
  },
};

export const Large: Story = {
  args: {
    primary: true,
    label: 'Input Large',
  },
};

export const Small: Story = {
  args: {
    primary: true,
    label: 'Input Small',
  },
};

export const Placeholder: Story = {
  args: {
    primary: true,
    label: "",
    placeholder: "Search filter",
    type: ""
  }
};
