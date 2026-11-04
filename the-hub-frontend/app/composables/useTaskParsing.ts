import * as chrono from 'chrono-node'

/**
 * Shared natural-language parsing helpers for task creation.
 * Used by both the quick-add input in the Command Center and FormTask.
 */
export function useTaskParsing() {
  // Parse dates from natural language input with enhanced patterns
  const parseDateFromText = (text: string): Date | null => {
    // First try chrono-node for standard date parsing
    const results = chrono.parse(text)
    if (results.length > 0) {
      return results[0].start.date()
    }

    // Enhanced fallback parsing for common patterns
    const lowerText = text.toLowerCase()
    const now = new Date()

    // Relative date patterns
    if (lowerText.includes('tomorrow')) {
      const tomorrow = new Date(now)
      tomorrow.setDate(now.getDate() + 1)
      return tomorrow
    }

    if (lowerText.includes('next week')) {
      const nextWeek = new Date(now)
      nextWeek.setDate(now.getDate() + 7)
      return nextWeek
    }

    if (lowerText.includes('next month')) {
      const nextMonth = new Date(now)
      nextMonth.setMonth(now.getMonth() + 1)
      return nextMonth
    }

    if (lowerText.includes('end of week') || lowerText.includes('this friday') || lowerText.includes('friday')) {
      const daysUntilFriday = (5 - now.getDay() + 7) % 7
      const friday = new Date(now)
      friday.setDate(now.getDate() + (daysUntilFriday === 0 ? 7 : daysUntilFriday))
      return friday
    }

    if (lowerText.includes('end of month')) {
      const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
      return endOfMonth
    }

    // Specific day patterns
    const dayPatterns = {
      'monday': 1, 'tuesday': 2, 'wednesday': 3, 'thursday': 4,
      'friday': 5, 'saturday': 6, 'sunday': 0
    }

    for (const [day, dayNum] of Object.entries(dayPatterns)) {
      if (lowerText.includes(`next ${day}`)) {
        const daysUntil = (dayNum - now.getDay() + 7) % 7
        const targetDate = new Date(now)
        targetDate.setDate(now.getDate() + (daysUntil === 0 ? 7 : daysUntil))
        return targetDate
      }
      if (lowerText.includes(`this ${day}`)) {
        const daysUntil = (dayNum - now.getDay() + 7) % 7
        const targetDate = new Date(now)
        targetDate.setDate(now.getDate() + daysUntil)
        return targetDate
      }
    }

    // Time parsing (e.g., "at 3pm", "at 14:30")
    const timeMatch = lowerText.match(/at (\d{1,2})(?::(\d{2}))?\s*(am|pm)?/)
    if (timeMatch) {
      let hour = parseInt(timeMatch[1])
      const minute = timeMatch[2] ? parseInt(timeMatch[2]) : 0
      const ampm = timeMatch[3]

      if (ampm === 'pm' && hour !== 12) hour += 12
      if (ampm === 'am' && hour === 12) hour = 0

      const timeDate = new Date(now)
      timeDate.setHours(hour, minute, 0, 0)
      return timeDate
    }

    return null
  }

  // Extract priority keywords from text with enhanced patterns
  const parsePriorityFromText = (text: string): number | null => {
    const lowerText = text.toLowerCase()

    // Priority 5 (Urgent/Critical/Emergency)
    if (lowerText.includes('urgent') || lowerText.includes('asap') ||
      lowerText.includes('critical') || lowerText.includes('emergency') ||
      lowerText.includes('immediately') || lowerText.includes('right now') ||
      lowerText.includes('deadline') || lowerText.includes('rush')) {
      return 5
    }

    // Priority 4 (High/Important)
    if (lowerText.includes('high priority') || lowerText.includes('high') ||
      lowerText.includes('important') || lowerText.includes('priority') ||
      lowerText.includes('crucial') || lowerText.includes('essential') ||
      lowerText.includes('key') || lowerText.includes('vital')) {
      return 4
    }

    // Priority 3 (Medium/Normal/Default)
    if (lowerText.includes('medium priority') || lowerText.includes('medium') ||
      lowerText.includes('normal') || lowerText.includes('standard') ||
      lowerText.includes('regular') || lowerText.includes('average')) {
      return 3
    }

    // Priority 2 (Low)
    if (lowerText.includes('low priority') || lowerText.includes('low') ||
      lowerText.includes('whenever') || lowerText.includes('sometime') ||
      lowerText.includes('eventually') || lowerText.includes('casual')) {
      return 2
    }

    // Priority 1 (Minor/Trivial)
    if (lowerText.includes('minor') || lowerText.includes('trivial') ||
      lowerText.includes('nice to have') || lowerText.includes('optional') ||
      lowerText.includes('if time') || lowerText.includes('low effort')) {
      return 1
    }

    return null
  }

  // Extract time estimate from text
  const parseTimeEstimateFromText = (text: string): number | null => {
    const lowerText = text.toLowerCase()

    // Match patterns like "30 minutes", "2 hours", "1h", "45m"
    const timePatterns = [
      /(\d+)\s*(?:hour|hr|h)(?:\w*)/,
      /(\d+)\s*(?:minute|min|m)(?:\w*)/,
      /(\d+):(\d+)/, // HH:MM format
      /(\d+)\s*(?:hr|hour|h)\s*(\d+)\s*(?:min|minute|m)?/ // "2h 30m" format
    ]

    for (const pattern of timePatterns) {
      const match = lowerText.match(pattern)
      if (match) {
        if (match[2]) { // HH:MM or H M format
          const hours = parseInt(match[1])
          const minutes = parseInt(match[2])
          return hours * 60 + minutes
        } else {
          const value = parseInt(match[1])
          if (lowerText.includes('hour') || lowerText.includes('hr') || lowerText.includes('h')) {
            return value * 60 // Convert hours to minutes
          } else if (lowerText.includes('minute') || lowerText.includes('min') || lowerText.includes('m')) {
            return value
          }
        }
      }
    }

    return null
  }

  return {
    parseDateFromText,
    parsePriorityFromText,
    parseTimeEstimateFromText
  }
}
