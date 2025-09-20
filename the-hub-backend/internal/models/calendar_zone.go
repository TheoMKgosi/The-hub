package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CalendarZone represents predefined time zones for different activities
type CalendarZone struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"` // e.g., "Work Hours", "Class Time", "Personal Time"
	Description string    `json:"description"`
	Category    string    `json:"category" gorm:"not null"`       // work, study, personal, exercise, etc.
	Color       string    `json:"color" gorm:"default:'#3b82f6'"` // Hex color for visualization

	// Time constraints
	StartTime  time.Time `json:"start_time" gorm:"not null"` // Daily start time (e.g., 09:00:00)
	EndTime    time.Time `json:"end_time" gorm:"not null"`   // Daily end time (e.g., 17:00:00)
	DaysOfWeek string    `json:"days_of_week"`               // JSON array: ["monday", "tuesday", "wednesday"]

	// Scheduling preferences
	Priority        int  `json:"priority" gorm:"default:5"`             // 1-10, higher = more preferred for scheduling
	IsActive        bool `json:"is_active" gorm:"default:true"`         // Whether this zone is currently active
	AllowScheduling bool `json:"allow_scheduling" gorm:"default:false"` // Whether AI can schedule events in this zone
	MaxEventsPerDay *int `json:"max_events_per_day"`                    // Maximum events allowed per day in this zone

	// Enhanced scheduling controls
	SchedulingMode string `json:"scheduling_mode" gorm:"default:'none'"` // "none", "whitelist", "blacklist", "non_zone"

	// Task type filtering
	AllowedTaskCategories []string `json:"allowed_task_categories" gorm:"type:text[]"` // ["work", "meeting"]
	AllowedTaskTypes      []string `json:"allowed_task_types" gorm:"type:text[]"`      // ["development", "planning"]
	BlockedTaskCategories []string `json:"blocked_task_categories" gorm:"type:text[]"` // Explicit blocks
	BlockedTaskTypes      []string `json:"blocked_task_types" gorm:"type:text[]"`      // Explicit blocks

	// Non-zone scheduling preferences
	AllowNonZoneScheduling bool      `json:"allow_non_zone_scheduling" gorm:"default:true"`
	NonZoneStartTime       time.Time `json:"non_zone_start_time"`
	NonZoneEndTime         time.Time `json:"non_zone_end_time"`
	NonZoneDaysOfWeek      []string  `json:"non_zone_days_of_week" gorm:"type:text[]"`

	// Recurrence pattern (for recurring zones)
	IsRecurring     bool       `json:"is_recurring" gorm:"default:false"`
	RecurrenceStart *time.Time `json:"recurrence_start"`
	RecurrenceEnd   *time.Time `json:"recurrence_end"`

	User      User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ZoneCategory represents predefined categories for zones
type ZoneCategory struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	Color       string    `json:"color" gorm:"default:'#3b82f6'"`
	Icon        string    `json:"icon" gorm:"default:'calendar'"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// GetDefaultZoneCategories returns the default zone categories
func GetDefaultZoneCategories() []ZoneCategory {
	return []ZoneCategory{
		{Name: "work", Description: "Work and professional activities", Color: "#3b82f6", Icon: "briefcase", IsDefault: true},
		{Name: "study", Description: "Learning and educational activities", Color: "#10b981", Icon: "book", IsDefault: true},
		{Name: "personal", Description: "Personal and leisure activities", Color: "#f59e0b", Icon: "user", IsDefault: true},
		{Name: "exercise", Description: "Physical activities and workouts", Color: "#ef4444", Icon: "activity", IsDefault: true},
		{Name: "family", Description: "Family time and commitments", Color: "#8b5cf6", Icon: "users", IsDefault: true},
		{Name: "health", Description: "Medical appointments and health activities", Color: "#06b6d4", Icon: "heart", IsDefault: true},
		{Name: "creative", Description: "Creative and artistic activities", Color: "#ec4899", Icon: "palette", IsDefault: true},
		{Name: "social", Description: "Social events and gatherings", Color: "#84cc16", Icon: "users", IsDefault: true},
	}
}

// IsTimeInZone checks if a given time falls within this zone
func (cz *CalendarZone) IsTimeInZone(checkTime time.Time) bool {
	if !cz.IsActive {
		return false
	}

	// Check if the day of week matches
	if cz.DaysOfWeek != "" {
		// Parse days of week (assuming JSON array format)
		// For simplicity, we'll check if the day name is contained in the string
		dayName := checkTime.Weekday().String()
		if !strings.Contains(strings.ToLower(cz.DaysOfWeek), strings.ToLower(dayName)) {
			return false
		}
	}

	// Check if the time of day falls within the zone
	checkHour := checkTime.Hour()
	checkMinute := checkTime.Minute()

	zoneStartHour := cz.StartTime.Hour()
	zoneStartMinute := cz.StartTime.Minute()
	zoneEndHour := cz.EndTime.Hour()
	zoneEndMinute := cz.EndTime.Minute()

	checkTimeMinutes := checkHour*60 + checkMinute
	zoneStartMinutes := zoneStartHour*60 + zoneStartMinute
	zoneEndMinutes := zoneEndHour*60 + zoneEndMinute

	return checkTimeMinutes >= zoneStartMinutes && checkTimeMinutes <= zoneEndMinutes
}

// GetZoneScore returns a scheduling preference score for this zone
func (cz *CalendarZone) GetZoneScore() int {
	if !cz.AllowScheduling || !cz.IsActive {
		return 0 // Not available for scheduling
	}

	score := cz.Priority

	// Bonus for higher priority zones
	if cz.Priority >= 8 {
		score += 5
	} else if cz.Priority >= 6 {
		score += 3
	}

	return score
}

// GetAvailableTimeSlots returns available time slots within this zone for a given date
func (cz *CalendarZone) GetAvailableTimeSlots(date time.Time, existingEvents []ScheduledTask, slotDuration time.Duration) []time.Time {
	var availableSlots []time.Time

	if !cz.IsActive || !cz.AllowScheduling {
		return availableSlots
	}

	// Check if this zone applies to the given date
	dateWithZoneStart := time.Date(date.Year(), date.Month(), date.Day(),
		cz.StartTime.Hour(), cz.StartTime.Minute(), 0, 0, date.Location())
	dateWithZoneEnd := time.Date(date.Year(), date.Month(), date.Day(),
		cz.EndTime.Hour(), cz.EndTime.Minute(), 0, 0, date.Location())

	// If zone doesn't apply to this day, return empty
	if !cz.IsTimeInZone(dateWithZoneStart) {
		return availableSlots
	}

	current := dateWithZoneStart
	for current.Before(dateWithZoneEnd) {
		slotEnd := current.Add(slotDuration)

		// Check if this slot conflicts with existing events
		hasConflict := false
		for _, event := range existingEvents {
			if (current.Before(event.End) && slotEnd.After(event.Start)) ||
				(event.Start.Before(slotEnd) && event.End.After(current)) {
				hasConflict = true
				break
			}
		}

		if !hasConflict {
			availableSlots = append(availableSlots, current)
		}

		current = slotEnd
	}

	return availableSlots
}
