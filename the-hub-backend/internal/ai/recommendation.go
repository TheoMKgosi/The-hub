package ai

import (
	"strconv"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/google/uuid"
)

// EnergyProfile represents a user's energy levels throughout the day
type EnergyProfile struct {
	UserID      uuid.UUID              `json:"user_id"`
	TimeSlots   map[string]int         `json:"time_slots"` // hour -> energy level (1-10)
	Workload    map[string]int         `json:"workload"`   // day -> workload score (1-10)
	Preferences map[string]interface{} `json:"preferences"`
	// New configurable fields
	PreferredStartHour int      `json:"preferred_start_hour"` // Preferred start hour (0-23)
	PreferredEndHour   int      `json:"preferred_end_hour"`   // Preferred end hour (0-23)
	WorkDays           []string `json:"work_days"`            // Days of the week user works
	BreakDuration      int      `json:"break_duration"`       // Preferred break duration in minutes
}

// GenerateScheduleSuggestions creates AI-powered schedule suggestions
func GenerateScheduleSuggestions(userID uuid.UUID, tasks []models.Task, existingEvents []models.ScheduledTask) ([]models.ScheduledTask, error) {
	var suggestions []models.ScheduledTask

	// Get user's energy profile (for now, use default if not found)
	energyProfile := getDefaultEnergyProfile(userID)

	// Get available time slots for the next 7 days
	availableSlots := findAvailableTimeSlots(userID, existingEvents, 7)

	// Sort tasks by priority and deadline
	prioritizedTasks := prioritizeTasks(tasks)

	for _, task := range prioritizedTasks {
		if len(availableSlots) == 0 {
			break
		}

		// Find best time slot for this task
		bestSlot := findBestTimeSlot(task, availableSlots, energyProfile)

		if bestSlot != nil {
			suggestion := models.ScheduledTask{
				Title:       task.Title,
				Start:       bestSlot.Start,
				End:         bestSlot.End,
				UserID:      userID,
				TaskID:      &task.ID,
				CreatedByAI: true,
			}

			suggestions = append(suggestions, suggestion)

			// Remove this slot from available slots
			availableSlots = removeTimeSlot(availableSlots, bestSlot)
		}
	}

	return suggestions, nil
}

// getDefaultEnergyProfile returns a default energy profile
func getDefaultEnergyProfile(userID uuid.UUID) *EnergyProfile {
	return &EnergyProfile{
		UserID: userID,
		TimeSlots: map[string]int{
			"6":  3, // Early morning - low energy
			"9":  8, // Morning peak
			"12": 6, // Lunch time - moderate
			"15": 7, // Afternoon
			"18": 5, // Evening
			"21": 2, // Night - low energy
		},
		Workload: map[string]int{
			"monday":    7,
			"tuesday":   8,
			"wednesday": 6,
			"thursday":  8,
			"friday":    5,
			"saturday":  3,
			"sunday":    2,
		},
		PreferredStartHour: 9,
		PreferredEndHour:   17,
		WorkDays:           []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
		BreakDuration:      15,
	}
}

// TimeSlot represents an available time slot
type TimeSlot struct {
	Start time.Time
	End   time.Time
}

// findAvailableTimeSlots finds available time slots for scheduling
func findAvailableTimeSlots(userID uuid.UUID, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var availableSlots []TimeSlot

	now := time.Now()
	endDate := now.AddDate(0, 0, days)

	// Create time slots from 9 AM to 6 PM each day
	for d := now; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dayStart := time.Date(d.Year(), d.Month(), d.Day(), 9, 0, 0, 0, d.Location())
		dayEnd := time.Date(d.Year(), d.Month(), d.Day(), 18, 0, 0, 0, d.Location())

		current := dayStart
		for current.Before(dayEnd) {
			slotEnd := current.Add(time.Hour)

			// Check if this slot conflicts with existing events
			if !hasConflict(current, slotEnd, existingEvents) {
				availableSlots = append(availableSlots, TimeSlot{
					Start: current,
					End:   slotEnd,
				})
			}

			current = slotEnd
		}
	}

	return availableSlots
}

// hasConflict checks if a time slot conflicts with existing events
func hasConflict(start, end time.Time, events []models.ScheduledTask) bool {
	for _, event := range events {
		if (start.Before(event.End) && end.After(event.Start)) ||
			(event.Start.Before(end) && event.End.After(start)) {
			return true
		}
	}
	return false
}

// prioritizeTasks sorts tasks by priority, deadline, and other factors
func prioritizeTasks(tasks []models.Task) []models.Task {
	if len(tasks) <= 1 {
		return tasks
	}

	// Create a copy to avoid modifying the original slice
	prioritized := make([]models.Task, len(tasks))
	copy(prioritized, tasks)

	// Sort tasks by a scoring system
	for i := 0; i < len(prioritized)-1; i++ {
		for j := i + 1; j < len(prioritized); j++ {
			if taskScore(prioritized[j]) > taskScore(prioritized[i]) {
				prioritized[i], prioritized[j] = prioritized[j], prioritized[i]
			}
		}
	}

	return prioritized
}

// taskScore calculates a priority score for a task
func taskScore(task models.Task) int {
	score := 0

	// Priority score (higher priority = higher score)
	if task.Priority != nil {
		score += *task.Priority * 10
	}

	// Deadline proximity score
	if task.DueDate != nil {
		hoursUntilDue := time.Until(*task.DueDate).Hours()
		if hoursUntilDue < 0 {
			// Overdue tasks get highest priority
			score += 100
		} else if hoursUntilDue < 24 {
			score += 50
		} else if hoursUntilDue < 72 {
			score += 25
		}
	}

	// Time estimate consideration (shorter tasks get slight preference for scheduling)
	if task.TimeEstimate != nil && *task.TimeEstimate <= 30 {
		score += 5
	}

	// Recurring tasks get slight boost
	if task.IsRecurring {
		score += 3
	}

	return score
}

// findBestTimeSlot finds the best time slot for a task
func findBestTimeSlot(task models.Task, availableSlots []TimeSlot, profile *EnergyProfile) *TimeSlot {
	if len(availableSlots) == 0 {
		return nil
	}

	bestSlot := &availableSlots[0]
	bestScore := calculateSlotScore(availableSlots[0], task, profile)

	for _, slot := range availableSlots[1:] {
		score := calculateSlotScore(slot, task, profile)
		if score > bestScore {
			bestScore = score
			bestSlot = &slot
		}
	}

	return bestSlot
}

// calculateSlotScore calculates how good a time slot is for a task
func calculateSlotScore(slot TimeSlot, task models.Task, profile *EnergyProfile) int {
	score := 0

	// Energy level score (0-10)
	hour := slot.Start.Hour()
	hourStr := strconv.Itoa(hour)
	energyLevel := profile.TimeSlots[hourStr]
	if energyLevel == 0 {
		energyLevel = 5 // Default moderate energy if not specified
	}
	score += energyLevel * 3

	// Workload score (prefer days with lower workload)
	day := strings.ToLower(slot.Start.Weekday().String())
	workload := profile.Workload[day]
	if workload == 0 {
		workload = 5 // Default moderate workload
	}
	score += (10 - workload) * 2

	// Preferred hours bonus
	if hour >= profile.PreferredStartHour && hour <= profile.PreferredEndHour {
		score += 10
	}

	// Work day preference
	isWorkDay := false
	for _, workDay := range profile.WorkDays {
		if strings.ToLower(workDay) == day {
			isWorkDay = true
			break
		}
	}

	if isWorkDay {
		score += 5
	}

	// Task-specific scoring
	if task.TimeEstimate != nil {
		slotDuration := slot.End.Sub(slot.Start).Minutes()
		// Prefer slots that match or are larger than task estimate
		if slotDuration >= float64(*task.TimeEstimate) {
			score += 8
		} else if slotDuration >= float64(*task.TimeEstimate)*0.8 {
			score += 5
		}
	}

	// Avoid scheduling during typical break times
	if hour >= 12 && hour <= 13 {
		score -= 5 // Lunch time penalty
	}

	return score
}

// removeTimeSlot removes a time slot from the list
func removeTimeSlot(slots []TimeSlot, toRemove *TimeSlot) []TimeSlot {
	var result []TimeSlot
	for _, slot := range slots {
		if slot.Start != toRemove.Start || slot.End != toRemove.End {
			result = append(result, slot)
		}
	}
	return result
}

// GetAISuggestions is the main entry point for AI scheduling
func GetAISuggestions(userID uuid.UUID) ([]models.ScheduledTask, error) {
	// Get user's tasks
	var tasks []models.Task
	if err := config.GetDB().Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}

	// Get existing scheduled events
	var existingEvents []models.ScheduledTask
	if err := config.GetDB().Where("user_id = ?", userID).Find(&existingEvents).Error; err != nil {
		return nil, err
	}

	return GenerateScheduleSuggestions(userID, tasks, existingEvents)
}
