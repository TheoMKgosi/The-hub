export type ThemeWeekday =
  | 'monday'
  | 'tuesday'
  | 'wednesday'
  | 'thursday'
  | 'friday'
  | 'saturday'
  | 'sunday'

export interface Theme {
  theme_id: string
  user_id: string
  name: string
  description: string
  color: string
  icon?: string
  days: ThemeWeekday[]
}

export interface ThemeSuggestion {
  name: string
  description: string
  color: string
  days: ThemeWeekday[]
  tasks: string[]
  reasoning: string
}

export interface ApplyTheme {
  name: string
  description: string
  color: string
  days: ThemeWeekday[]
}

export interface ThemeAssignment {
  task_id: string
  theme_name: string
}

export const WEEKDAYS: ThemeWeekday[] = [
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
  'sunday'
]

export const WEEKDAY_LABELS: Record<ThemeWeekday, string> = {
  monday: 'Monday',
  tuesday: 'Tuesday',
  wednesday: 'Wednesday',
  thursday: 'Thursday',
  friday: 'Friday',
  saturday: 'Saturday',
  sunday: 'Sunday'
}
