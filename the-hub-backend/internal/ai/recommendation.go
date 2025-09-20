package ai

import (
	"encoding/json"
	"fmt"
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

// GenerateScheduleSuggestions creates algorithmic schedule suggestions
func GenerateScheduleSuggestions(userID uuid.UUID, tasks []models.Task, existingEvents []models.ScheduledTask) ([]models.ScheduledTask, error) {
	// Use algorithmic approach for scheduling suggestions
	return generateEnhancedRuleBasedSuggestions(userID, tasks, existingEvents)
}

// generateEnhancedRuleBasedSuggestions provides enhanced rule-based scheduling with hybrid zone + non-zone support
func generateEnhancedRuleBasedSuggestions(userID uuid.UUID, tasks []models.Task, existingEvents []models.ScheduledTask) ([]models.ScheduledTask, error) {
	var suggestions []models.ScheduledTask

	// Get user's energy profile (for now, use default if not found)
	energyProfile := getDefaultEnergyProfile(userID)

	// Get user's calendar zones
	zones, err := getUserCalendarZones(userID)
	if err != nil {
		// Log error but continue with default behavior
		zones = []models.CalendarZone{}
	}

	// Use the new hybrid scheduling system
	availableSlots := findAvailableTimeSlotsHybrid(userID, tasks, zones, existingEvents, 7)

	// Sort tasks by priority and deadline
	prioritizedTasks := prioritizeTasks(tasks)

	for _, task := range prioritizedTasks {
		if len(availableSlots) == 0 {
			break
		}

		// Find best time slot for this task, considering zones
		bestSlot := findBestTimeSlotWithZones(task, availableSlots, energyProfile, zones)

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

// generateRuleBasedSuggestions provides fallback rule-based scheduling (legacy function)
func generateRuleBasedSuggestions(userID uuid.UUID, tasks []models.Task, existingEvents []models.ScheduledTask) ([]models.ScheduledTask, error) {
	// Use the enhanced version for backward compatibility
	return generateEnhancedRuleBasedSuggestions(userID, tasks, existingEvents)
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

// getUserCalendarZones fetches calendar zones for a user
func getUserCalendarZones(userID uuid.UUID) ([]models.CalendarZone, error) {
	var zones []models.CalendarZone
	err := config.GetDB().Where("user_id = ? AND is_active = ?", userID, true).Find(&zones).Error
	return zones, err
}

// findAvailableTimeSlotsWithZones finds available time slots considering calendar zones
func findAvailableTimeSlotsWithZones(userID uuid.UUID, existingEvents []models.ScheduledTask, zones []models.CalendarZone, days int) []TimeSlot {
	var availableSlots []TimeSlot

	now := time.Now()
	endDate := now.AddDate(0, 0, days)

	// If no zones are defined, fall back to default behavior
	if len(zones) == 0 {
		return findAvailableTimeSlots(userID, existingEvents, days)
	}

	// Create time slots based on zones
	for d := now; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		for _, zone := range zones {
			if !zone.AllowScheduling {
				continue
			}

			// Check if zone applies to this day
			if !zone.IsTimeInZone(d) {
				continue
			}

			// Get available slots within this zone
			zoneSlots := zone.GetAvailableTimeSlots(d, existingEvents, time.Hour)
			for _, slot := range zoneSlots {
				availableSlots = append(availableSlots, TimeSlot{
					Start: slot,
					End:   slot.Add(time.Hour),
				})
			}
		}
	}

	return availableSlots
}

// findBestTimeSlotWithZones finds the best time slot considering calendar zones
func findBestTimeSlotWithZones(task models.Task, availableSlots []TimeSlot, profile *EnergyProfile, zones []models.CalendarZone) *TimeSlot {
	if len(availableSlots) == 0 {
		return nil
	}

	bestSlot := &availableSlots[0]
	bestScore := calculateSlotScoreWithZones(availableSlots[0], task, profile, zones)

	for _, slot := range availableSlots[1:] {
		score := calculateSlotScoreWithZones(slot, task, profile, zones)
		if score > bestScore {
			bestScore = score
			bestSlot = &slot
		}
	}

	return bestSlot
}

// calculateSlotScoreWithZones calculates slot score considering calendar zones
func calculateSlotScoreWithZones(slot TimeSlot, task models.Task, profile *EnergyProfile, zones []models.CalendarZone) int {
	score := 0

	// Base energy and workload score
	score += calculateSlotScore(slot, task, profile)

	// Zone-based scoring
	for _, zone := range zones {
		if zone.IsTimeInZone(slot.Start) {
			// Zone preference bonus
			score += zone.GetZoneScore()

			// Category matching bonus
			if zone.Category == "work" && task.Priority != nil && *task.Priority >= 4 {
				score += 10 // High priority tasks in work zones
			}
			if zone.Category == "study" && strings.Contains(strings.ToLower(task.Title), "learn") {
				score += 8 // Learning tasks in study zones
			}
			if zone.Category == "personal" && task.Priority != nil && *task.Priority <= 2 {
				score += 6 // Low priority tasks in personal zones
			}

			break // Found matching zone, no need to check others
		}
	}

	return score
}

// GoalTaskRecommendation represents an AI-recommended task for a goal
type GoalTaskRecommendation struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Priority       int    `json:"priority"`
	EstimatedHours int    `json:"estimated_hours"`
	Reasoning      string `json:"reasoning"`
}

// GenerateGoalTaskRecommendations creates AI-powered task recommendations for a goal
func GenerateGoalTaskRecommendations(goalID uuid.UUID, userID uuid.UUID) ([]GoalTaskRecommendation, error) {
	// Get the goal details
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		return nil, fmt.Errorf("goal not found: %w", err)
	}

	// Get existing tasks for this goal
	var existingTasks []models.Task
	if err := config.GetDB().Where("goal_id = ? AND user_id = ?", goalID, userID).Find(&existingTasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get existing tasks: %w", err)
	}

	// Convert existing tasks to string format
	existingTaskStrings := make([]string, len(existingTasks))
	for i, task := range existingTasks {
		existingTaskStrings[i] = task.Title
		if task.Description != "" {
			existingTaskStrings[i] += ": " + task.Description
		}
	}

	// Initialize OpenRouter client
	client, err := NewOpenRouterClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI client: %w", err)
	}

	// Get AI recommendations from OpenRouter
	aiResponse, err := client.GenerateGoalTaskRecommendations(goal.Title, goal.Description, existingTaskStrings)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI recommendations: %w", err)
	}

	// Parse AI response
	recommendations, err := parseGoalTaskRecommendations(aiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI recommendations: %w", err)
	}

	return recommendations, nil
}

// parseGoalTaskRecommendations parses the AI response into structured recommendations
func parseGoalTaskRecommendations(aiResponse string) ([]GoalTaskRecommendation, error) {
	var recommendations []GoalTaskRecommendation

	// Try to unmarshal as JSON array
	if err := json.Unmarshal([]byte(aiResponse), &recommendations); err != nil {
		// If direct unmarshal fails, try to extract JSON from the response
		start := strings.Index(aiResponse, "[")
		end := strings.LastIndex(aiResponse, "]")
		if start == -1 || end == -1 || start >= end {
			return nil, fmt.Errorf("failed to parse AI response: no JSON array found")
		}

		jsonStr := aiResponse[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &recommendations); err != nil {
			return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
		}
	}

	// Validate and clean up the recommendations
	for i := range recommendations {
		// Ensure priority is within valid range
		if recommendations[i].Priority < 1 {
			recommendations[i].Priority = 1
		} else if recommendations[i].Priority > 5 {
			recommendations[i].Priority = 5
		}

		// Ensure estimated hours is reasonable
		if recommendations[i].EstimatedHours < 1 {
			recommendations[i].EstimatedHours = 1
		} else if recommendations[i].EstimatedHours > 40 {
			recommendations[i].EstimatedHours = 8 // Default to 8 hours
		}
	}

	return recommendations, nil
}

// hasTaskWithKeyword checks if any existing task contains a specific keyword
func hasTaskWithKeyword(tasks []models.Task, keyword string) bool {
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), keyword) ||
			strings.Contains(strings.ToLower(task.Description), keyword) {
			return true
		}
	}
	return false
}

// GetAISuggestions is the main entry point for AI scheduling (deprecated - use GenerateGoalTaskRecommendations instead)
func GetAISuggestions(userID uuid.UUID) ([]models.ScheduledTask, error) {
	// Get user's tasks (exclude tasks that already have due dates)
	var tasks []models.Task
	if err := config.GetDB().Where("user_id = ? AND due_date IS NULL", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}

	// Get existing scheduled events
	var existingEvents []models.ScheduledTask
	if err := config.GetDB().Where("user_id = ?", userID).Find(&existingEvents).Error; err != nil {
		return nil, err
	}

	return GenerateScheduleSuggestions(userID, tasks, existingEvents)
}

// === ENHANCED HYBRID SCHEDULING SYSTEM ===

// NonZonePreferences represents non-zone scheduling preferences
type NonZonePreferences struct {
	StartTime  time.Time
	EndTime    time.Time
	DaysOfWeek []string
}

// findAvailableTimeSlotsHybrid implements the new hybrid scheduling system
func findAvailableTimeSlotsHybrid(userID uuid.UUID, tasks []models.Task, zones []models.CalendarZone, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var allAvailableSlots []TimeSlot

	// Phase 1: Try to schedule in zones with smart filtering
	zoneSlots := findAvailableTimeSlotsWithSmartZones(userID, tasks, zones, existingEvents, days)
	allAvailableSlots = append(allAvailableSlots, zoneSlots...)

	// Phase 2: If zones didn't provide enough slots or no zones exist, add non-zone slots
	if len(allAvailableSlots) == 0 || shouldAddNonZoneSlots(tasks, zoneSlots, zones) {
		nonZoneSlots := findAvailableNonZoneSlots(userID, tasks, zones, existingEvents, days)
		allAvailableSlots = append(allAvailableSlots, nonZoneSlots...)
	}

	return allAvailableSlots
}

// findAvailableTimeSlotsWithSmartZones processes zones with enhanced task filtering
func findAvailableTimeSlotsWithSmartZones(userID uuid.UUID, tasks []models.Task, zones []models.CalendarZone, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var availableSlots []TimeSlot

	// If no zones exist, return empty (let non-zone handle it)
	if len(zones) == 0 {
		return []TimeSlot{}
	}

	// Process each zone based on its scheduling mode
	for _, zone := range zones {
		if !zone.IsActive {
			continue
		}

		switch zone.SchedulingMode {
		case "none":
			continue // Skip this zone entirely
		case "non_zone":
			continue // Handle in non-zone phase
		case "whitelist":
			// Only allow explicitly listed task categories/types
			zoneSlots := processZoneWithWhitelist(zone, tasks, existingEvents, days)
			availableSlots = append(availableSlots, zoneSlots...)
		case "blacklist":
			// Allow all except explicitly blocked categories/types
			zoneSlots := processZoneWithBlacklist(zone, tasks, existingEvents, days)
			availableSlots = append(availableSlots, zoneSlots...)
		default:
			// Fallback to current behavior for backward compatibility
			if zone.AllowScheduling {
				zoneSlots := zone.GetAvailableTimeSlots(time.Now(), existingEvents, time.Hour)
				for _, slot := range zoneSlots {
					availableSlots = append(availableSlots, TimeSlot{
						Start: slot,
						End:   slot.Add(time.Hour),
					})
				}
			}
		}
	}

	return availableSlots
}

// processZoneWithWhitelist processes a zone using whitelist filtering
func processZoneWithWhitelist(zone models.CalendarZone, tasks []models.Task, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var availableSlots []TimeSlot

	// Find tasks that are allowed in this zone
	allowedTasks := filterTasksForZoneWhitelist(tasks, zone)
	if len(allowedTasks) == 0 {
		return availableSlots // No tasks can be scheduled in this zone
	}

	// Get available slots within this zone
	zoneSlots := zone.GetAvailableTimeSlots(time.Now(), existingEvents, time.Hour)
	for _, slot := range zoneSlots {
		availableSlots = append(availableSlots, TimeSlot{
			Start: slot,
			End:   slot.Add(time.Hour),
		})
	}

	return availableSlots
}

// processZoneWithBlacklist processes a zone using blacklist filtering
func processZoneWithBlacklist(zone models.CalendarZone, tasks []models.Task, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var availableSlots []TimeSlot

	// Find tasks that are not blocked in this zone
	allowedTasks := filterTasksForZoneBlacklist(tasks, zone)
	if len(allowedTasks) == 0 {
		return availableSlots // No tasks can be scheduled in this zone
	}

	// Get available slots within this zone
	zoneSlots := zone.GetAvailableTimeSlots(time.Now(), existingEvents, time.Hour)
	for _, slot := range zoneSlots {
		availableSlots = append(availableSlots, TimeSlot{
			Start: slot,
			End:   slot.Add(time.Hour),
		})
	}

	return availableSlots
}

// findAvailableNonZoneSlots generates slots outside of zone times
func findAvailableNonZoneSlots(userID uuid.UUID, tasks []models.Task, zones []models.CalendarZone, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var availableSlots []TimeSlot

	now := time.Now()
	endDate := now.AddDate(0, 0, days)

	// Extract non-zone preferences from zones or use defaults
	nonZonePrefs := extractNonZonePreferences(zones)

	for d := now; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		// Check if this day should be included
		if !shouldScheduleOnDay(d, nonZonePrefs.DaysOfWeek) {
			continue
		}

		// Generate slots outside of zone times
		daySlots := generateNonZoneSlotsForDay(d, nonZonePrefs.StartTime, nonZonePrefs.EndTime, existingEvents, zones)
		availableSlots = append(availableSlots, daySlots...)
	}

	return availableSlots
}

// extractNonZonePreferences extracts non-zone scheduling preferences from zones
func extractNonZonePreferences(zones []models.CalendarZone) NonZonePreferences {
	// Default preferences
	prefs := NonZonePreferences{
		StartTime:  time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),  // 9 AM
		EndTime:    time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC), // 6 PM
		DaysOfWeek: []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
	}

	// Look for a non-zone configuration zone
	for _, zone := range zones {
		if zone.SchedulingMode == "non_zone" && zone.IsActive {
			prefs.StartTime = zone.NonZoneStartTime
			prefs.EndTime = zone.NonZoneEndTime
			prefs.DaysOfWeek = zone.NonZoneDaysOfWeek
			break
		}
	}

	return prefs
}

// shouldScheduleOnDay checks if scheduling should occur on a given day
func shouldScheduleOnDay(date time.Time, allowedDays []string) bool {
	if len(allowedDays) == 0 {
		return true // No restrictions
	}

	dayName := strings.ToLower(date.Weekday().String())
	for _, allowedDay := range allowedDays {
		if strings.ToLower(allowedDay) == dayName {
			return true
		}
	}
	return false
}

// generateNonZoneSlotsForDay generates time slots outside of zone coverage
func generateNonZoneSlotsForDay(date time.Time, startTime, endTime time.Time, existingEvents []models.ScheduledTask, zones []models.CalendarZone) []TimeSlot {
	var availableSlots []TimeSlot

	// Create the day's time range
	dayStart := time.Date(date.Year(), date.Month(), date.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, date.Location())
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, date.Location())

	current := dayStart
	for current.Before(dayEnd) {
		slotEnd := current.Add(time.Hour)

		// Check if this slot is outside all active zones (true non-zone)
		if !isSlotInAnyActiveZone(current, slotEnd, zones) {
			// Check for conflicts with existing events
			if !hasConflict(current, slotEnd, existingEvents) {
				availableSlots = append(availableSlots, TimeSlot{
					Start: current,
					End:   slotEnd,
				})
			}
		}

		current = slotEnd
	}

	return availableSlots
}

// isSlotInAnyActiveZone checks if a time slot falls within any active zone
func isSlotInAnyActiveZone(start, end time.Time, zones []models.CalendarZone) bool {
	for _, zone := range zones {
		if zone.IsActive && zone.IsTimeInZone(start) {
			return true
		}
	}
	return false
}

// shouldAddNonZoneSlots determines if non-zone slots should be added
func shouldAddNonZoneSlots(tasks []models.Task, zoneSlots []TimeSlot, zones []models.CalendarZone) bool {
	// Always add non-zone slots if no zones allow scheduling
	hasSchedulingZone := false
	for _, zone := range zones {
		if zone.AllowScheduling || zone.SchedulingMode == "whitelist" || zone.SchedulingMode == "blacklist" {
			hasSchedulingZone = true
			break
		}
	}

	if !hasSchedulingZone {
		return true
	}

	// Add non-zone slots if we have tasks but insufficient zone slots
	minSlotsNeeded := len(tasks)
	return len(zoneSlots) < minSlotsNeeded
}

// filterTasksForZoneWhitelist filters tasks that are allowed in a zone (whitelist mode)
func filterTasksForZoneWhitelist(tasks []models.Task, zone models.CalendarZone) []models.Task {
	var filteredTasks []models.Task

	for _, task := range tasks {
		if isTaskAllowedByZoneWhitelist(task, zone) {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}

// filterTasksForZoneBlacklist filters tasks that are not blocked in a zone (blacklist mode)
func filterTasksForZoneBlacklist(tasks []models.Task, zone models.CalendarZone) []models.Task {
	var filteredTasks []models.Task

	for _, task := range tasks {
		if !isTaskBlockedByZoneBlacklist(task, zone) {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}

// isTaskAllowedByZoneWhitelist checks if a task is allowed in a zone (whitelist mode)
func isTaskAllowedByZoneWhitelist(task models.Task, zone models.CalendarZone) bool {
	// Check category whitelist
	if len(zone.AllowedTaskCategories) > 0 {
		if !containsString(zone.AllowedTaskCategories, task.Category) {
			return false
		}
	}

	// Check type whitelist
	if len(zone.AllowedTaskTypes) > 0 {
		if !containsString(zone.AllowedTaskTypes, task.TaskType) {
			return false
		}
	}

	return true
}

// isTaskBlockedByZoneBlacklist checks if a task is blocked in a zone (blacklist mode)
func isTaskBlockedByZoneBlacklist(task models.Task, zone models.CalendarZone) bool {
	// Check category blacklist
	if containsString(zone.BlockedTaskCategories, task.Category) {
		return true
	}

	// Check type blacklist
	if containsString(zone.BlockedTaskTypes, task.TaskType) {
		return true
	}

	return false
}

// containsString checks if a string slice contains a specific string
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// === TASK-ZONE COMPATIBILITY FUNCTIONS ===

// buildTaskZoneMapping creates a mapping of tasks to compatible zones
func buildTaskZoneMapping(tasks []models.Task, zones []models.CalendarZone) map[uuid.UUID][]models.CalendarZone {
	taskZoneMap := make(map[uuid.UUID][]models.CalendarZone)

	for _, task := range tasks {
		var compatibleZones []models.CalendarZone

		for _, zone := range zones {
			if isTaskCompatibleWithZone(task, zone) {
				compatibleZones = append(compatibleZones, zone)
			}
		}

		if len(compatibleZones) > 0 {
			taskZoneMap[task.ID] = compatibleZones
		}
	}

	return taskZoneMap
}

// isTaskCompatibleWithZone checks if a task can be scheduled in a zone
func isTaskCompatibleWithZone(task models.Task, zone models.CalendarZone) bool {
	if !zone.IsActive {
		return false
	}

	switch zone.SchedulingMode {
	case "none":
		return false
	case "non_zone":
		return false // Non-zone zones don't schedule tasks directly
	case "whitelist":
		return isTaskAllowedByZoneWhitelist(task, zone)
	case "blacklist":
		return !isTaskBlockedByZoneBlacklist(task, zone)
	default:
		// Backward compatibility with old AllowScheduling boolean
		return zone.AllowScheduling
	}
}

// findSlotsForTaskInZones finds available slots for a task within its compatible zones
func findSlotsForTaskInZones(task models.Task, compatibleZones []models.CalendarZone, existingEvents []models.ScheduledTask, days int) []TimeSlot {
	var availableSlots []TimeSlot

	now := time.Now()
	endDate := now.AddDate(0, 0, days)

	for d := now; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		for _, zone := range compatibleZones {
			if !zone.IsTimeInZone(d) {
				continue
			}

			// Get available slots within this zone
			zoneSlots := zone.GetAvailableTimeSlots(d, existingEvents, time.Hour)
			for _, slot := range zoneSlots {
				availableSlots = append(availableSlots, TimeSlot{
					Start: slot,
					End:   slot.Add(time.Hour),
				})
			}
		}
	}

	return availableSlots
}
