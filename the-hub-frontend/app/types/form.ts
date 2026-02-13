export interface FormField {
  name: string
  label: string
  type: 'text' | 'textarea' | 'select' | 'combobox' | 'date' | 'datetime-local' | 'color' | 'number' | 'email' | 'password'
  placeholder?: string
  required?: boolean
  options?: { value: any; label: string }[]
  categories?: { budget_category_id: string; name: string }[]
  allowCreate?: boolean
  rows?: number
  min?: number
  max?: number
  step?: number
}

export interface FormUIProps {
  title: string
  fields: FormField[]
  submitLabel: string
  cancelLabel?: string
  validationSchema?: Record<string, any>
  initialData?: Record<string, any>
  showForm?: boolean
  teleportTarget?: string
  size?: 'sm' | 'md' | 'lg'
}
