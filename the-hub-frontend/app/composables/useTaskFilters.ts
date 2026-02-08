import type { Task } from '~/types/task'

export interface TaskFilters {
  search?: string
  status?: string
  priority?: string // "1,2" for multiple
  dueDate?: 'today' | 'tomorrow' | 'this_week' | 'overdue' | 'none' | 'custom'
  dueDateStart?: string // ISO date for custom range
  dueDateEnd?: string // ISO date for custom range
  linked?: 'true' | 'false'
  hasStartTime?: 'true' | 'false'
  duration?: '0-30' | '30-60' | '60-120' | '120+'
  created?: '7d' | '30d' | 'older'
  sortBy?: 'due_date' | 'priority' | 'created' | 'title'
  sortOrder?: 'asc' | 'desc'
}

const defaultFilters: TaskFilters = {
  search: '',
  status: '',
  priority: '',
  dueDate: '',
  dueDateStart: '',
  dueDateEnd: '',
  linked: '',
  hasStartTime: '',
  duration: '',
  created: '',
  sortBy: 'due_date',
  sortOrder: 'asc'
}

export const useTaskFilters = () => {
  const route = useRoute()
  const router = useRouter()
  
  // Filter state
  const filters = reactive<TaskFilters>({ ...defaultFilters })
  
  // Advanced filters expanded state
  const isAdvancedOpen = ref(false)
  
  // Filter bar collapsed state
  const isFilterBarOpen = ref(true)
  
  // Load filters from URL on mount
  const loadFiltersFromURL = () => {
    const query = route.query
    
    Object.keys(defaultFilters).forEach((key) => {
      const value = query[key]
      if (value && typeof value === 'string') {
        filters[key as keyof TaskFilters] = value as any
      }
    })
  }
  
  // Sync filters to URL
  const syncFiltersToURL = () => {
    const query: Record<string, string> = {}
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value !== '' && value !== 'all') {
        query[key] = String(value)
      }
    })
    
    router.replace({ query })
  }
  
  // Watch for filter changes and sync to URL
  watch(filters, syncFiltersToURL, { deep: true })
  
  // Helper: Check if same day
  const isSameDay = (date1: Date, date2: Date): boolean => {
    return date1.getFullYear() === date2.getFullYear() &&
           date1.getMonth() === date2.getMonth() &&
           date1.getDate() === date2.getDate()
  }
  
  // Helper: Check if date is this week
  const isThisWeek = (date: Date): boolean => {
    const now = new Date()
    const startOfWeek = new Date(now)
    startOfWeek.setDate(now.getDate() - now.getDay())
    startOfWeek.setHours(0, 0, 0, 0)
    
    const endOfWeek = new Date(startOfWeek)
    endOfWeek.setDate(startOfWeek.getDate() + 7)
    
    return date >= startOfWeek && date < endOfWeek
  }
  
  // Helper: Match due date
  const matchesDueDate = (dueDate: Date | null | undefined, filter: string): boolean => {
    if (!dueDate) return filter === 'none'
    if (filter === 'none') return false
    
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    
    const due = new Date(dueDate)
    due.setHours(0, 0, 0, 0)
    
    switch (filter) {
      case 'today':
        return isSameDay(due, today)
      case 'tomorrow': {
        const tomorrow = new Date(today)
        tomorrow.setDate(tomorrow.getDate() + 1)
        return isSameDay(due, tomorrow)
      }
      case 'this_week':
        return isThisWeek(due)
      case 'overdue':
        return due < today
      default:
        return false
    }
  }
  
  // Helper: Match duration
  const matchesDuration = (minutes: number | undefined, range: string): boolean => {
    if (!minutes) return false
    
    switch (range) {
      case '0-30': return minutes >= 0 && minutes <= 30
      case '30-60': return minutes > 30 && minutes <= 60
      case '60-120': return minutes > 60 && minutes <= 120
      case '120+': return minutes > 120
      default: return false
    }
  }
  
  // Helper: Match created date
  const matchesCreatedDate = (createdAt: Date | undefined, filter: string): boolean => {
    if (!createdAt) return false
    
    const created = new Date(createdAt)
    const now = new Date()
    const daysDiff = Math.floor((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24))
    
    switch (filter) {
      case '7d': return daysDiff <= 7
      case '30d': return daysDiff <= 30
      case 'older': return daysDiff > 30
      default: return false
    }
  }
  
  // Main filter function
  const matchFilter = (task: Task): boolean => {
    // Search text
    if (filters.search) {
      const search = filters.search.toLowerCase()
      const titleMatch = task.title?.toLowerCase().includes(search)
      const descMatch = task.description?.toLowerCase().includes(search)
      if (!titleMatch && !descMatch) return false
    }
    
    // Status
    if (filters.status && task.status !== filters.status) return false
    
    // Priority (supports multiple: "4,5")
    if (filters.priority) {
      const priorities = filters.priority.split(',').map(Number)
      if (!priorities.includes(task.priority || 0)) return false
    }
    
    // Due date
    if (filters.dueDate && !matchesDueDate(task.due_date, filters.dueDate)) return false
    
    // Custom date range
    if (filters.dueDateStart && filters.dueDateEnd && task.due_date) {
      const start = new Date(filters.dueDateStart)
      const end = new Date(filters.dueDateEnd)
      const due = new Date(task.due_date)
      if (due < start || due > end) return false
    }
    
    // Goal linked
    if (filters.linked === 'true' && !task.goal_id) return false
    if (filters.linked === 'false' && task.goal_id) return false
    
    // Has start time
    if (filters.hasStartTime === 'true' && !task.start_time) return false
    if (filters.hasStartTime === 'false' && task.start_time) return false
    
    // Duration
    if (filters.duration && !matchesDuration(task.time_estimate_minutes, filters.duration)) return false
    
    // Created date
    if (filters.created && !matchesCreatedDate(task.created_at, filters.created)) return false
    
    return true
  }
  
  // Sort function
  const sortTasks = (tasks: Task[]): Task[] => {
    return [...tasks].sort((a, b) => {
      let comparison = 0
      
      switch (filters.sortBy) {
        case 'due_date': {
          const aDate = a.due_date ? new Date(a.due_date).getTime() : Infinity
          const bDate = b.due_date ? new Date(b.due_date).getTime() : Infinity
          comparison = aDate - bDate
          break
        }
        case 'priority':
          comparison = (b.priority || 0) - (a.priority || 0)
          break
        case 'created': {
          const aCreated = a.created_at ? new Date(a.created_at).getTime() : 0
          const bCreated = b.created_at ? new Date(b.created_at).getTime() : 0
          comparison = bCreated - aCreated
          break
        }
        case 'title':
          comparison = (a.title || '').localeCompare(b.title || '')
          break
      }
      
      return filters.sortOrder === 'desc' ? -comparison : comparison
    })
  }
  
  // Clear all filters
  const clearAllFilters = () => {
    Object.assign(filters, defaultFilters)
  }
  
  // Clear specific filter
  const clearFilter = (key: keyof TaskFilters) => {
    filters[key] = defaultFilters[key]
  }
  
  // Get active filter count
  const activeFilterCount = computed(() => {
    return Object.entries(filters).filter(([key, value]) => {
      if (key === 'sortBy' || key === 'sortOrder') return false
      return value && value !== ''
    }).length
  })
  
  // Get active filters as chips
  const activeFilterChips = computed(() => {
    const chips: { key: keyof TaskFilters; label: string; value: string }[] = []
    
    const labels: Record<string, string> = {
      status: 'Status',
      priority: 'Priority',
      dueDate: 'Due Date',
      linked: 'Linked',
      hasStartTime: 'Schedule',
      duration: 'Duration',
      created: 'Created',
      search: 'Search'
    }
    
    const formatValue = (key: string, value: string): string => {
      switch (key) {
        case 'dueDate':
          return value.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
        case 'linked':
          return value === 'true' ? 'Yes' : 'No'
        case 'hasStartTime':
          return value === 'true' ? 'Scheduled' : 'Not Scheduled'
        case 'duration':
          return value.replace('-', '-').replace('+', '+')
        case 'priority':
          return value.split(',').map(p => `P${p}`).join(', ')
        default:
          return value.charAt(0).toUpperCase() + value.slice(1)
      }
    }
    
    Object.entries(filters).forEach(([key, value]) => {
      if (key === 'sortBy' || key === 'sortOrder') return
      if (value && value !== '') {
        chips.push({
          key: key as keyof TaskFilters,
          label: labels[key] || key,
          value: formatValue(key, value)
        })
      }
    })
    
    return chips
  })
  
  return {
    filters,
    isAdvancedOpen,
    isFilterBarOpen,
    loadFiltersFromURL,
    matchFilter,
    sortTasks,
    clearAllFilters,
    clearFilter,
    activeFilterCount,
    activeFilterChips
  }
}
