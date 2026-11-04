import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { toValue } from 'vue'
import type { ShortcutAction } from '../../../app/composables/useKeyboardShortcuts'

// useKeyboardShortcuts imports useRouter from 'vue-router' explicitly (mocked
// below) and reads useTaskShortcuts via Nuxt auto-imports, which inject a real
// module import — so the composable module is mocked directly. The module
// holds shared registry state, so it is re-imported per test.
const taskMock = {
  moveSelection: vi.fn(),
  toggleComplete: vi.fn(),
  editSelected: vi.fn(),
  deleteTask: vi.fn(),
  undoDelete: vi.fn(),
  newTask: vi.fn(),
}
const router = { push: vi.fn() }

vi.mock('vue-router', () => ({
  useRouter: () => router
}))

vi.mock('../../../app/composables/useTaskShortcuts', () => ({
  useTaskShortcuts: () => taskMock
}))

let useKeyboardShortcuts: typeof import('../../../app/composables/useKeyboardShortcuts')['useKeyboardShortcuts']

beforeEach(async () => {
  vi.clearAllMocks()
  // Fresh module state (shortcut registry, palette/help flags)
  vi.resetModules()
  ;({ useKeyboardShortcuts } = await import('../../../app/composables/useKeyboardShortcuts'))
})

const findShortcut = (
  shortcuts: readonly ShortcutAction[] | { value: readonly ShortcutAction[] },
  key: string,
  ctrl = false,
  shift = false
) =>
  toValue(shortcuts).find(s => s.key === key && !!s.ctrl === ctrl && !!s.shift === shift)

describe('useKeyboardShortcuts', () => {
  it('registers the task shortcuts', () => {
    const { shortcuts } = useKeyboardShortcuts()

    expect(findShortcut(shortcuts, 'j')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, 'k')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, ' ')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, 'e')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, 'n')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, 'd')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, 'u')).toMatchObject({ category: 'Tasks' })
    expect(findShortcut(shortcuts, '?', false, true)).toMatchObject({ category: 'Actions' })
  })

  it('wires the task shortcut actions to useTaskShortcuts', () => {
    const { shortcuts } = useKeyboardShortcuts()

    findShortcut(shortcuts, 'j')!.action()
    expect(taskMock.moveSelection).toHaveBeenCalledWith(1)

    findShortcut(shortcuts, 'k')!.action()
    expect(taskMock.moveSelection).toHaveBeenCalledWith(-1)

    findShortcut(shortcuts, ' ')!.action()
    expect(taskMock.toggleComplete).toHaveBeenCalledTimes(1)

    findShortcut(shortcuts, 'e')!.action()
    expect(taskMock.editSelected).toHaveBeenCalledTimes(1)

    findShortcut(shortcuts, 'd')!.action()
    expect(taskMock.deleteTask).toHaveBeenCalledTimes(1)

    findShortcut(shortcuts, 'u')!.action()
    expect(taskMock.undoDelete).toHaveBeenCalledTimes(1)

    findShortcut(shortcuts, 'n')!.action()
    expect(taskMock.newTask).toHaveBeenCalledTimes(1)
  })

  it('navigates with Ctrl+<number> shortcuts', () => {
    const { shortcuts } = useKeyboardShortcuts()

    findShortcut(shortcuts, '2', true)!.action()
    expect(router.push).toHaveBeenCalledWith('/plan')
  })

  it('opens and closes the help modal', () => {
    const { openHelp, closeHelp, showHelp } = useKeyboardShortcuts()

    expect(showHelp.value).toBe(false)
    openHelp()
    expect(showHelp.value).toBe(true)
    closeHelp()
    expect(showHelp.value).toBe(false)
  })

  it('toggles the help modal with the ? shortcut', () => {
    const { shortcuts, showHelp } = useKeyboardShortcuts()

    findShortcut(shortcuts, '?', false, true)!.action()
    expect(showHelp.value).toBe(true)
  })

  it('supports adding and removing custom shortcuts', () => {
    const { addShortcut, removeShortcut, shortcuts } = useKeyboardShortcuts()
    const custom: ShortcutAction = {
      key: 'x',
      action: () => {},
      description: 'Custom action',
      category: 'Actions'
    }

    addShortcut(custom)
    expect(findShortcut(shortcuts, 'x')).toBeDefined()

    removeShortcut('x')
    expect(findShortcut(shortcuts, 'x')).toBeUndefined()
  })
})
