import { type Meta, type StoryObj } from '@storybook/vue3-vite';
import Toast from '~/components/ui/Toast.vue';
// Use a relative path instead of #imports for Storybook
import { useToast } from '@/composables/useToast';

const meta = {
  component: Toast,
  // The 'render' function belongs at the top level of the object, not inside 'args'
  render: (args) => ({
    components: { Toast },
    setup() {
      const { addToast } = useToast();

      onMounted(() => {
        addToast('Accessibility Test success', 'success')
        addToast('Accessibility Test error', 'error')
        addToast('Accessibility Test warning', 'warning')
        addToast('Accessibility Test info', 'info')
      })

      const triggerToast = (type: 'success' | 'error' | 'warning' | 'info') => {
        // Updated to pass the object as we discussed previously
        addToast(`This is a ${type} toast!`, type)
      }

      return { triggerToast, args };
    },
    template: `
      <div class="p-6">
        <div class="flex gap-2 mb-8">
          <button @click="triggerToast('success')" class="px-4 py-2 bg-green-600 text-white rounded shadow">Success</button>
          <button @click="triggerToast('error')" class="px-4 py-2 bg-red-600 text-white rounded shadow">Error</button>
        </div>
        
        <Toast v-bind="args" />
      </div>
    `,
  }),
} satisfies Meta<typeof Toast>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
