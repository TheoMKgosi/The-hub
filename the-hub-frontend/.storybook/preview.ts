import { setup, type Preview } from '@storybook-vue/nuxt'
import '../app/assets/css/main.css'
import { createPinia } from 'pinia';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { withThemeByDataAttribute } from '@storybook/addon-themes';

dayjs.extend(relativeTime);

setup((app) => {
  app.use(createPinia())
}
)

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export const decorators = [
  withThemeByDataAttribute({
    themes: {
      light: '',
      dark: 'dark',
    },
    defaultTheme: 'light',
    attributeName: 'class'
  }),
];


export default preview;
