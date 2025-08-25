package ai

import (
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

// prioritizeTasks sorts tasks by priority and deadline
func prioritizeTasks(tasks []models.Task) []models.Task {
	// For now, return tasks as-is. In a real implementation,
	// you'd sort by priority, deadline, and other factors
	return tasks
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
	energyLevel := profile.TimeSlots[string(rune(hour+'0'))]
	score += energyLevel * 2

	// Workload score (prefer days with lower workload)
	day := slot.Start.Weekday().String()
	workload := profile.Workload[day]
	score += (10 - workload) * 1

	// Task-specific scoring could be added here
	// e.g., creative tasks in morning, physical tasks when energy is high

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
