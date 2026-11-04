import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { ref, reactive, computed, watch } from 'vue'
import { useTaskFilters } from '../../../app/composables/useTaskFilters'
import type { Task } from '../../../app/types/task'

// useTaskFilters relies on Nuxt auto-imports (useRoute/useRouter) and Vue
// auto-imports (ref/reactive/computed/watch). Stub them all with the real
// Vue implementations plus lightweight router fakes so the composable works
// in the plain node unit environment.
const route = { query: {} as Record<string, string> }
const router = { replace: vi.fn() }

beforeAll(() => {
  vi.stubGlobal('useRoute', () => route)
  vi.stubGlobal('useRouter', () => router)
  vi.stubGlobal('ref', ref)
  vi.stubGlobal('reactive', reactive)
  vi.stubGlobal('computed', computed)
  vi.stubGlobal('watch', watch)
})

const DAY = 24 * 60 * 60 * 1000

const makeTask = (overrides: Partial<Task> = {}): Task => ({
  task_id: 't1',
  title: 'Write tests',
  description: 'Cover the filter logic',
  status: 'pending',
  start_time: new Date(),
  ...overrides
})

describe('useTaskFilters', () => {
  beforeEach(() => {
    route.query = {}
    router.replace.mockClear()
  })

  describe('matchFilter', () => {
    it('matches every task when no filters are set', () => {
      const { matchFilter, filters } = useTaskFilters()

      expect(matchFilter(makeTask())).toBe(true)
      expect(filters.status).toBe('')
    })

    it('filters by search across title and description, case-insensitively', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.search = 'WRITE'
      expect(matchFilter(makeTask({ title: 'Write tests' }))).toBe(true)
      expect(matchFilter(makeTask({ title: 'Other', description: 'please WRITE more' }))).toBe(true)
      expect(matchFilter(makeTask({ title: 'Go shopping', description: 'buy milk' }))).toBe(false)
    })

    it('filters by exact status', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.status = 'completed'
      expect(matchFilter(makeTask({ status: 'completed' }))).toBe(true)
      expect(matchFilter(makeTask({ status: 'pending' }))).toBe(false)
    })

    it('supports multiple priorities as a comma-separated list', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.priority = '4,5'
      expect(matchFilter(makeTask({ priority: 4 }))).toBe(true)
      expect(matchFilter(makeTask({ priority: 5 }))).toBe(true)
      expect(matchFilter(makeTask({ priority: 1 }))).toBe(false)
    })

    it('filters overdue tasks', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.dueDate = 'overdue'
      expect(matchFilter(makeTask({ due_date: new Date(Date.now() - DAY) }))).toBe(true)
      expect(matchFilter(makeTask({ due_date: new Date(Date.now() + DAY) }))).toBe(false)
      expect(matchFilter(makeTask({ due_date: null }))).toBe(false)
    })

    it('matches tasks with no due date when filtering for "none"', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.dueDate = 'none'
      expect(matchFilter(makeTask({ due_date: null }))).toBe(true)
      expect(matchFilter(makeTask({ due_date: new Date() }))).toBe(false)
    })

    it('filters tasks within a custom date range', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.dueDateStart = new Date(Date.now() - DAY).toISOString()
      filters.dueDateEnd = new Date(Date.now() + DAY).toISOString()

      expect(matchFilter(makeTask({ due_date: new Date() }))).toBe(true)
      expect(matchFilter(makeTask({ due_date: new Date(Date.now() + 10 * DAY) }))).toBe(false)
    })

    it('filters by goal linkage', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.linked = 'true'
      expect(matchFilter(makeTask({ goal_id: 'g1' }))).toBe(true)
      expect(matchFilter(makeTask())).toBe(false)

      filters.linked = 'false'
      expect(matchFilter(makeTask())).toBe(true)
      expect(matchFilter(makeTask({ goal_id: 'g1' }))).toBe(false)
    })

    it('filters by duration range', () => {
      const { matchFilter, filters } = useTaskFilters()

      filters.duration = '30-60'
      expect(matchFilter(makeTask({ time_estimate_minutes: 45 }))).toBe(true)
      expect(matchFilter(makeTask({ time_estimate_minutes: 30 }))).toBe(false)
      expect(matchFilter(makeTask({ time_estimate_minutes: 120 }))).toBe(false)
    })
  })

  describe('sortTasks', () => {
    it('sorts by due date ascending and keeps tasks without due dates last', () => {
      const { sortTasks } = useTaskFilters()

      const sorted = sortTasks([
        makeTask({ task_id: 'no-date' }),
        makeTask({ task_id: 'later', due_date: new Date(Date.now() + 2 * DAY) }),
        makeTask({ task_id: 'sooner', due_date: new Date(Date.now() + DAY) })
      ])

      expect(sorted.map(t => t.task_id)).toEqual(['sooner', 'later', 'no-date'])
    })

    it('sorts by priority with the highest priority first', () => {
      const { filters, sortTasks } = useTaskFilters()

      filters.sortBy = 'priority'
      const sorted = sortTasks([
        makeTask({ task_id: 'low', priority: 1 }),
        makeTask({ task_id: 'high', priority: 5 }),
        makeTask({ task_id: 'mid', priority: 3 })
      ])

      expect(sorted.map(t => t.task_id)).toEqual(['high', 'mid', 'low'])
    })

    it('sorts by title alphabetically', () => {
      const { filters, sortTasks } = useTaskFilters()

      filters.sortBy = 'title'
      const sorted = sortTasks([
        makeTask({ task_id: 'b', title: 'Banana' }),
        makeTask({ task_id: 'a', title: 'Apple' }),
        makeTask({ task_id: 'c', title: 'Cherry' })
      ])

      expect(sorted.map(t => t.task_id)).toEqual(['a', 'b', 'c'])
    })

    it('does not mutate the input array', () => {
      const { sortTasks } = useTaskFilters()

      const tasks = [
        makeTask({ task_id: 'b', due_date: new Date(Date.now() + DAY) }),
        makeTask({ task_id: 'a', due_date: new Date() })
      ]
      sortTasks(tasks)

      expect(tasks.map(t => t.task_id)).toEqual(['b', 'a'])
    })
  })

  describe('active filters', () => {
    it('counts active filters, excluding sort options', () => {
      const { filters, activeFilterCount } = useTaskFilters()

      expect(activeFilterCount.value).toBe(0)

      filters.status = 'pending'
      filters.search = 'tests'
      expect(activeFilterCount.value).toBe(2)

      filters.sortBy = 'title'
      expect(activeFilterCount.value).toBe(2)
    })

    it('produces labelled chips with formatted values', () => {
      const { filters, activeFilterChips } = useTaskFilters()

      filters.status = 'pending'
      filters.priority = '4,5'
      filters.dueDate = 'this_week'

      expect(activeFilterChips.value).toEqual([
        { key: 'status', label: 'Status', value: 'Pending' },
        { key: 'priority', label: 'Priority', value: 'P4, P5' },
        { key: 'dueDate', label: 'Due Date', value: 'This Week' }
      ])
    })

    it('clears a single filter and all filters', () => {
      const { filters, activeFilterCount, clearFilter, clearAllFilters } = useTaskFilters()

      filters.status = 'pending'
      filters.search = 'tests'
      clearFilter('status')

      expect(filters.status).toBe('')
      expect(activeFilterCount.value).toBe(1)

      clearAllFilters()
      expect(activeFilterCount.value).toBe(0)
      expect(filters.sortBy).toBe('due_date')
    })
  })

  describe('URL sync', () => {
    // NOTE: current behavior — syncFiltersToURL does not exclude sort options,
    // so default sortBy/sortOrder are included in the query.
    it('syncs non-empty filters to the URL', async () => {
      const { filters } = useTaskFilters()

      filters.status = 'pending'
      await vi.waitFor(() => {
        expect(router.replace).toHaveBeenCalledWith({ query: { status: 'pending', sortBy: 'due_date', sortOrder: 'asc' } })
      })
    })

    it('omits empty values from the URL query', async () => {
      const { filters } = useTaskFilters()

      filters.search = 'hello'
      await vi.waitFor(() => {
        expect(router.replace).toHaveBeenLastCalledWith({ query: { search: 'hello', sortBy: 'due_date', sortOrder: 'asc' } })
      })
    })

    it('loads filters from the URL query', () => {
      const { filters, loadFiltersFromURL } = useTaskFilters()

      route.query = { status: 'completed', priority: '3' }
      loadFiltersFromURL()

      expect(filters.status).toBe('completed')
      expect(filters.priority).toBe('3')
      expect(filters.search).toBe('')
    })
  })
})
