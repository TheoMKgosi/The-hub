export interface ValidationResult {
  isValid: boolean
  errors: Record<string, string>
}

export interface ValidationRule {
  required?: boolean
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  custom?: (value: any) => string | null
}

export function useValidation() {
  const validateField = (value: any, rules: ValidationRule, fieldName: string): string | null => {
    if (rules.required && (!value || (typeof value === 'string' && !value.trim()))) {
      return `${fieldName} is required`
    }

    if (value && typeof value === 'string') {
      if (rules.minLength && value.length < rules.minLength) {
        return `${fieldName} must be at least ${rules.minLength} characters`
      }

      if (rules.maxLength && value.length > rules.maxLength) {
        return `${fieldName} must be no more than ${rules.maxLength} characters`
      }

      if (rules.pattern && !rules.pattern.test(value)) {
        return `${fieldName} format is invalid`
      }
    }

    if (rules.custom && value) {
      const customError = rules.custom(value)
      if (customError) return customError
    }

    return null
  }

  const validateObject = (data: Record<string, any>, schema: Record<string, ValidationRule>): ValidationResult => {
    const errors: Record<string, string> = {}

    for (const [field, rules] of Object.entries(schema)) {
      const error = validateField(data[field], rules, field)
      if (error) {
        errors[field] = error
      }
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors
    }
  }

  // Predefined validation schemas
  const schemas = {
    auth: {
      login: {
        email: { required: true, pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/ },
        password: { required: true, minLength: 6 }
      },
      register: {
        name: { required: true, minLength: 2, maxLength: 50 },
        email: { required: true, pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/ },
        password: { required: true, minLength: 8 }
      },
      forgotPassword: {
        email: { required: true, pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/ }
      },
      resetPassword: {
        token: { required: true, minLength: 32 },
        password: { required: true, minLength: 8 }
      }
    },
    task: {
      create: {
        title: { required: true, minLength: 1, maxLength: 200 },
        description: { maxLength: 1000 },
        priority: {
          custom: (value: number) => {
            if (value && (value < 1 || value > 5)) {
              return 'Priority must be between 1 and 5'
            }
            return null
          }
        },
        due_date: {
          custom: (value: string) => {
            if (value) {
              const date = new Date(value)
              if (isNaN(date.getTime())) {
                return 'Invalid date format'
              }
              if (date < new Date()) {
                return 'Due date cannot be in the past'
              }
            }
            return null
          }
        }
      },
      naturalLanguage: {
        natural_language: { required: true, minLength: 3, maxLength: 500 }
      }
    },
    goal: {
      create: {
        title: { required: true, minLength: 1, maxLength: 200 },
        description: { maxLength: 1000 }
      }
    },
    category: {
      create: {
        name: { required: true, minLength: 1, maxLength: 100 }
      }
    }
  }

  return {
    validateField,
    validateObject,
    schemas
  }
}