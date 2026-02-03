export interface Task {
  task_id: string
  title: string
  description: string
  due_date?: Date | null
  priority?: number
  status: string
  order?: number
  goal_id?: string
  parent_task_id?: string
  subtasks?: Task[]
  time_estimate_minutes?: number
  time_spent_minutes: number
  is_recurring: boolean
}

export interface TaskUpdate {
  task_id: string
  title?: string
  description?: string
  due_date?: Date | null
  priority?: number
  status?: string
  order?: number
  goal_id?: string
  parent_task_id?: string
  subtasks?: Task[]
  time_estimate_minutes?: number
  time_spent_minutes?: number
  is_recurring?: boolean
}

export interface TimeEntry {
  time_entry_id: string
  task_id: string
  description: string
  start_time: string
  end_time?: string
  duration_minutes: number
  is_running: boolean
}

export interface TaskTemplate {
  template_id: string
  name: string
  description: string
  category: string
  title_template: string
  description_template: string
  priority?: number
  time_estimate_minutes?: number
  tags: string
  is_public: boolean
  usage_count: number
}

export interface RecurrenceRule {
  recurrence_rule_id: string
  name: string
  description: string
  frequency: string
  interval: number
  by_day?: string
  by_month_day?: number
  by_month?: number
  start_date?: string
  end_date?: string
  count?: number
  title_template: string
  description_template: string
  priority?: number
  time_estimate_minutes?: number
  due_date_offset_days?: number
}
