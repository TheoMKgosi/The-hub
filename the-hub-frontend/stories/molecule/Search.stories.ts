import type { Meta, StoryObj } from '@nuxtjs/storybook';
 
import Search from '~/components/ui/Search.vue';
 
const meta = {
  component: Search,
} satisfies Meta<typeof Search>;
 
export default meta 
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    placeholder: "Search Filter"
  }
}
