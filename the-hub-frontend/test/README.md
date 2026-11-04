# Frontend Testing

This directory contains the test setup and configuration for the frontend application.

## Setup

The testing environment is configured with:
- **Vitest**: Fast testing framework (3 projects: `unit`, `e2e`, `nuxt`)
- **@nuxt/test-utils** + **Vue Test Utils**: Component and page testing (`mountSuspended`)
- **jsdom** / Nuxt environment: DOM for component and page tests
- **TypeScript**: Full type support

## Running Tests

```bash
# Run all tests once
npx vitest run

# Run a single project
npx vitest run --project unit
npx vitest run --project nuxt

# Run a single file
npx vitest run test/unit/stores/tags.test.ts

# Watch mode / UI
npx vitest
npx vitest --ui
```

## Test Structure

```
test/
├── setup.ts                  # Global mocks (matchMedia, ResizeObserver)
├── unit/                     # Node environment (no Nuxt context)
│   ├── composables/          # Pure-logic composable tests
│   │   ├── useValidation.test.ts
│   │   ├── useDebounce.test.ts
│   │   ├── useMarkdown.test.ts
│   │   └── useTaskFilters.test.ts
│   └── stores/               # Pinia store tests
│       ├── tags.test.ts
│       ├── income.test.ts
│       └── focus-sessions.test.ts
└── nuxt/                     # Nuxt environment (auto-imports, Nuxt UI, router)
    ├── components/           # Component tests (mountSuspended)
    │   ├── Tabs.test.ts
    │   ├── ConfirmDialog.test.ts
    │   ├── NavLink.test.ts
    │   └── FilterChips.test.ts
    └── pages/                # Page tests (mountSuspended)
        ├── login.test.ts
        ├── offline.test.ts
        └── static-pages.test.ts
```

## Conventions

### Composables and stores → `test/unit/`
Pure-logic composables run in the plain node environment. Composables/stores that
rely on Nuxt auto-imports (`useRoute`, `useRouter`, `useNuxtApp`, `useToast`,
`defineStore`, `ref`, `computed`…) stub them with `vi.stubGlobal`, passing the
real implementations for Vue APIs:

```typescript
import { ref, reactive, computed, watch } from 'vue'
import { defineStore } from 'pinia'

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('defineStore', defineStore)

// If the module references stubbed globals at module top level, import it
// dynamically after the stubs are in place:
const { useFocusSessionStore } = await import('~/stores/focus-sessions')
```

Stores need a fresh Pinia per test:
```typescript
beforeEach(() => setActivePinia(createPinia()))
```

### Components and pages → `test/nuxt/`
Use `mountSuspended` for the full Nuxt context (Nuxt UI components, `NuxtLink`,
auto-imports):

```typescript
import { mountSuspended } from '@nuxt/test-utils/runtime'

const wrapper = await mountSuspended(MyComponent, {
  props: { /* props */ },
  slots: { /* slots */ }
})
expect(wrapper.emitted('update:modelValue')).toEqual([['value']])
```

Gotchas (documented in existing tests):
- `wrapper.text()` collapses whitespace between sibling spans — assert text
  without spaces, or per element
- `mountSuspended` does not simulate a parent `v-model`; to test two-way
  binding end-to-end, mount a small `defineComponent` parent harness
- To mock a store used via auto-import (e.g. `useAuthStore`), mock its module:
  `vi.mock('~/stores/auth', () => ({ useAuthStore: () => ({ login: loginMock }) }))`
  and declare the mock with `vi.hoisted`
