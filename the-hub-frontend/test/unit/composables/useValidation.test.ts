import { describe, it, expect } from 'vitest'
import { useValidation } from '../../../app/composables/useValidation'

describe('useValidation', () => {
  describe('validateField', () => {
    it('returns an error for a required field that is empty or whitespace', () => {
      const { validateField } = useValidation()

      expect(validateField(null, { required: true }, 'Title')).toBe('Title is required')
      expect(validateField('', { required: true }, 'Title')).toBe('Title is required')
      expect(validateField('   ', { required: true }, 'Title')).toBe('Title is required')
    })

    it('returns null for a required field with a value', () => {
      const { validateField } = useValidation()

      expect(validateField('valid value', { required: true }, 'Title')).toBeNull()
    })

    it('does not require non-required fields even when empty', () => {
      const { validateField } = useValidation()

      expect(validateField('', { maxLength: 10 }, 'Description')).toBeNull()
    })

    it('enforces minLength', () => {
      const { validateField } = useValidation()

      expect(validateField('ab', { minLength: 3 }, 'Name')).toBe('Name must be at least 3 characters')
      expect(validateField('abc', { minLength: 3 }, 'Name')).toBeNull()
    })

    it('enforces maxLength', () => {
      const { validateField } = useValidation()

      expect(validateField('abcd', { maxLength: 3 }, 'Name')).toBe('Name must be no more than 3 characters')
      expect(validateField('abc', { maxLength: 3 }, 'Name')).toBeNull()
    })

    it('enforces pattern', () => {
      const { validateField } = useValidation()

      expect(validateField('not-an-email', { pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/ }, 'Email')).toBe('Email format is invalid')
      expect(validateField('a@b.com', { pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/ }, 'Email')).toBeNull()
    })

    it('enforces custom rules', () => {
      const { validateField } = useValidation()

      expect(validateField(9, { custom: (v: number) => (v > 5 ? 'Too high' : null) }, 'Priority')).toBe('Too high')
      expect(validateField(3, { custom: (v: number) => (v > 5 ? 'Too high' : null) }, 'Priority')).toBeNull()
    })

    it('skips custom rules when the value is empty', () => {
      const { validateField } = useValidation()

      expect(validateField('', { custom: () => 'should not run' }, 'Field')).toBeNull()
    })
  })

  describe('validateObject', () => {
    const schema = {
      name: { required: true, minLength: 2 },
      email: { required: true, pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/ },
      description: { maxLength: 10 }
    }

    it('returns isValid false and collects all errors', () => {
      const { validateObject } = useValidation()

      const result = validateObject({ name: 'a', email: 'bad', description: 'way too long description' }, schema)

      expect(result.isValid).toBe(false)
      expect(result.errors).toEqual({
        name: 'name must be at least 2 characters',
        email: 'email format is invalid',
        description: 'description must be no more than 10 characters'
      })
    })

    it('returns isValid true with no errors for valid data', () => {
      const { validateObject } = useValidation()

      const result = validateObject({ name: 'Alice', email: 'a@b.com', description: '' }, schema)

      expect(result.isValid).toBe(true)
      expect(result.errors).toEqual({})
    })
  })

  describe('predefined schemas', () => {
    it('auth.login validates email format', () => {
      const { validateObject, schemas } = useValidation()

      expect(validateObject({ email: 'user@example.com' }, schemas.auth.login).isValid).toBe(true)
      expect(validateObject({ email: 'nope' }, schemas.auth.login).isValid).toBe(false)
    })

    it('task.create rejects a due date in the past', () => {
      const { validateObject, schemas } = useValidation()
      const past = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString()

      const result = validateObject({ title: 'Write tests', due_date: past }, schemas.task.create)

      expect(result.isValid).toBe(false)
      expect(result.errors.due_date).toBe('Due date cannot be in the past')
    })

    it('task.create rejects an out-of-range priority', () => {
      const { validateObject, schemas } = useValidation()

      const result = validateObject({ title: 'Write tests', priority: 6 }, schemas.task.create)

      expect(result.isValid).toBe(false)
      expect(result.errors.priority).toBe('Priority must be between 1 and 5')
    })

    it('goal.create rejects an invalid hex color', () => {
      const { validateObject, schemas } = useValidation()

      const result = validateObject({ title: 'Learn Go', color: 'not-a-color' }, schemas.goal.create)

      expect(result.isValid).toBe(false)
      expect(result.errors.color).toBe('Color must be a valid hex color')
    })

    it('naturalLanguage requires a meaningful description', () => {
      const { validateObject, schemas } = useValidation()

      const tooShort = validateObject({ natural_language: 'ea' }, schemas.task.naturalLanguage)
      expect(tooShort.isValid).toBe(false)
      // minLength rule fires first and returns early, so the custom rule's
      // "Please provide at least a brief description" branch is unreachable
      expect(tooShort.errors.natural_language).toBe('natural_language must be at least 3 characters')

      const valid = validateObject({ natural_language: 'Finish the frontend tests tomorrow' }, schemas.task.naturalLanguage)
      expect(valid.isValid).toBe(true)
    })

    it('naturalLanguage suggests more detail for very short descriptions', () => {
      const { validateObject, schemas } = useValidation()

      const result = validateObject({ natural_language: 'eat' }, schemas.task.naturalLanguage)

      expect(result.isValid).toBe(false)
      expect(result.errors.natural_language).toBe('Consider adding more details like timing or priority for better task parsing')
    })
  })
})
