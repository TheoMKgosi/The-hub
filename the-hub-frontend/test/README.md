# Frontend Testing

This directory contains the test setup and configuration for the frontend application.

## Setup

The testing environment is configured with:
- **Vitest**: Fast testing framework
- **Vue Test Utils**: Vue component testing utilities
- **jsdom**: DOM environment for testing
- **TypeScript**: Full type support

## Running Tests

```bash
# Run all tests once
npm run test:run

# Run tests in watch mode
npm run test

# Run tests with UI
npm run test:ui
```

## Test Structure

```
test/
├── setup.ts              # Global test setup and configuration
├── components/           # Component tests
│   └── ui/
│       └── Button.test.ts
└── composables/          # Composable tests
    └── useToast.test.ts
```

## Writing Tests

### Component Tests
Use Vue Test Utils to mount and test components:

```typescript
import { mount } from '@vue/test-utils'
import MyComponent from '~/components/MyComponent.vue'

describe('MyComponent', () => {
  it('renders correctly', () => {
    const wrapper = mount(MyComponent, {
      props: { /* props */ },
      slots: { /* slots */ }
    })

    expect(wrapper.text()).toBe('Expected text')
  })
})
```

### Composable Tests
Test composables by calling them directly:

```typescript
import { useMyComposable } from '~/composables/useMyComposable'

describe('useMyComposable', () => {
  it('works correctly', () => {
    const { result } = useMyComposable()

    expect(result.value).toBe('expected value')
  })
})
```

## Configuration

- Tests run in jsdom environment
- Global Vue APIs are available
- Auto-cleanup is handled by Vitest
- Path aliases (`~` and `@`) are configured