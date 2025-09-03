package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecurrenceRule defines how an event repeats
type RecurrenceRule struct {
	ID          uuid.UUID  `json:"recurrence_rule_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Frequency   string     `json:"frequency" gorm:"not null"` // daily, weekly, monthly, yearly
	Interval    int        `json:"interval" gorm:"default:1"` // every N frequency units
	ByDay       string     `json:"by_day"`                    // e.g., "MO,TU,WE" for weekly (0=sunday, 1=monday, etc.)
	ByMonthDay  *int       `json:"by_month_day"`              // day of month (1-31)
	ByMonth     *int       `json:"by_month"`                  // month for yearly (1-12)
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Count       *int       `json:"count"` // number of occurrences
	// Template for recurring tasks
	TitleTemplate       string         `json:"title_template"`
	DescriptionTemplate string         `json:"description_template"`
	Priority            *int           `json:"priority"`
	TimeEstimate        *int           `json:"time_estimate_minutes"`
	DueDateOffset       *int           `json:"due_date_offset_days"` // Days from occurrence date
	User                User           `json:"-" gorm:"foreignKey:UserID"`
	Tasks               []Task         `json:"-" gorm:"foreignKey:RecurrenceRuleID"`
	CreatedAt           time.Time      `json:"-"`
	UpdatedAt           time.Time      `json:"-"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// GenerateOccurrences generates the next occurrence dates for this recurrence rule
func (rr *RecurrenceRule) GenerateOccurrences(fromDate time.Time, limit int) []time.Time {
	var occurrences []time.Time
	current := fromDate

	if rr.StartDate != nil && rr.StartDate.After(fromDate) {
		current = *rr.StartDate
	}

	for len(occurrences) < limit {
		nextDate := rr.getNextOccurrence(current)
		if nextDate.IsZero() {
			break
		}

		// Check end conditions
		if rr.EndDate != nil && nextDate.After(*rr.EndDate) {
			break
		}

		if rr.Count != nil && len(occurrences) >= *rr.Count {
			break
		}

		occurrences = append(occurrences, nextDate)
		current = nextDate.AddDate(0, 0, 1) // Move to next day to avoid infinite loop
	}

	return occurrences
}

// getNextOccurrence calculates the next occurrence date based on the recurrence pattern
func (rr *RecurrenceRule) getNextOccurrence(fromDate time.Time) time.Time {
	switch rr.Frequency {
	case "daily":
		return rr.getNextDailyOccurrence(fromDate)
	case "weekly":
		return rr.getNextWeeklyOccurrence(fromDate)
	case "monthly":
		return rr.getNextMonthlyOccurrence(fromDate)
	case "yearly":
		return rr.getNextYearlyOccurrence(fromDate)
	default:
		return time.Time{}
	}
}

func (rr *RecurrenceRule) getNextDailyOccurrence(fromDate time.Time) time.Time {
	days := rr.Interval
	if days <= 0 {
		days = 1
	}
	return fromDate.AddDate(0, 0, days)
}

func (rr *RecurrenceRule) getNextWeeklyOccurrence(fromDate time.Time) time.Time {
	if rr.ByDay == "" {
		// Default to same day of week
		days := rr.Interval * 7
		return fromDate.AddDate(0, 0, days)
	}

	// Parse by_day (comma-separated day numbers: 0=sunday, 1=monday, etc.)
	dayStrings := strings.Split(rr.ByDay, ",")
	var targetDays []int

	for _, dayStr := range dayStrings {
		if day, err := strconv.Atoi(strings.TrimSpace(dayStr)); err == nil {
			targetDays = append(targetDays, day)
		}
	}

	if len(targetDays) == 0 {
		return fromDate.AddDate(0, 0, rr.Interval*7)
	}

	// Find next occurrence on target days
	current := fromDate
	weeks := rr.Interval
	if weeks <= 0 {
		weeks = 1
	}

	for i := 0; i < 7*weeks; i++ {
		current = current.AddDate(0, 0, 1)
		currentWeekday := int(current.Weekday())

		for _, targetDay := range targetDays {
			if currentWeekday == targetDay {
				return current
			}
		}
	}

	return time.Time{}
}

func (rr *RecurrenceRule) getNextMonthlyOccurrence(fromDate time.Time) time.Time {
	if rr.ByMonthDay == nil {
		// Default to same day of month
		return fromDate.AddDate(0, rr.Interval, 0)
	}

	targetDay := *rr.ByMonthDay
	months := rr.Interval
	if months <= 0 {
		months = 1
	}

	current := fromDate.AddDate(0, months, 0)

	// Set to target day of month
	if targetDay > 0 && targetDay <= 31 {
		current = time.Date(current.Year(), current.Month(), targetDay, 0, 0, 0, 0, current.Location())

		// If the target day doesn't exist in this month (e.g., Feb 31), use last day of month
		if current.Month() != fromDate.AddDate(0, months, 0).Month() {
			current = time.Date(current.Year(), current.Month()+1, 0, 0, 0, 0, 0, current.Location())
		}
	}

	return current
}

func (rr *RecurrenceRule) getNextYearlyOccurrence(fromDate time.Time) time.Time {
	years := rr.Interval
	if years <= 0 {
		years = 1
	}

	current := fromDate.AddDate(years, 0, 0)

	if rr.ByMonth != nil && rr.ByMonthDay != nil {
		month := *rr.ByMonth
		day := *rr.ByMonthDay

		if month >= 1 && month <= 12 && day >= 1 && day <= 31 {
			targetDate := time.Date(current.Year(), time.Month(month), day, 0, 0, 0, 0, current.Location())

			// Handle invalid dates (e.g., Feb 31)
			if targetDate.Month() != time.Month(month) {
				targetDate = time.Date(current.Year(), time.Month(month)+1, 0, 0, 0, 0, 0, current.Location())
			}

			return targetDate
		}
	}

	return current
}

// CreateTaskFromRule creates a new task instance from this recurrence rule
func (rr *RecurrenceRule) CreateTaskFromRule(userID uuid.UUID, occurrenceDate time.Time) *Task {
	task := &Task{
		UserID:           userID,
		Title:            rr.TitleTemplate,
		Description:      rr.DescriptionTemplate,
		Priority:         rr.Priority,
		TimeEstimate:     rr.TimeEstimate,
		IsRecurring:      true,
		RecurrenceRuleID: &rr.ID,
	}

	// Set due date based on offset
	if rr.DueDateOffset != nil {
		dueDate := occurrenceDate.AddDate(0, 0, *rr.DueDateOffset)
		task.DueDate = &dueDate
	}

	return task
}
