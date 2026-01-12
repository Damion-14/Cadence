package models

import "time"

type PersonalRecord struct {
	ExerciseName string    `json:"exercise_name"`
	MaxWeight    *float64  `json:"max_weight,omitempty"`
	MaxReps      int       `json:"max_reps"`
	MaxVolume    *float64  `json:"max_volume,omitempty"`
	AchievedAt   time.Time `json:"achieved_at"`
}

type WorkoutSummary struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	CompletedAt   time.Time `json:"completed_at"`
	ExerciseCount int       `json:"exercise_count"`
	TotalSets     int       `json:"total_sets"`
	TotalVolume   float64   `json:"total_volume"`
}

type WeeklySummary struct {
	Week          string           `json:"week"`
	TotalWorkouts int              `json:"total_workouts"`
	TotalExercises int             `json:"total_exercises"`
	TotalVolume   float64          `json:"total_volume"`
	Workouts      []WorkoutInWeek  `json:"workouts"`
}

type WorkoutInWeek struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CompletedAt time.Time `json:"completed_at"`
	DayOfWeek   int       `json:"day_of_week"`
}

type ProgressDataPoint struct {
	Date      time.Time `json:"date"`
	MaxWeight *float64  `json:"max_weight,omitempty"`
	MaxReps   int       `json:"max_reps"`
	Volume    float64   `json:"volume"`
}

type ProgressResponse struct {
	ExerciseName string              `json:"exercise_name"`
	DataPoints   []ProgressDataPoint `json:"data_points"`
}

type HistoryResponse struct {
	Workouts []WorkoutSummary `json:"workouts"`
	Total    int              `json:"total"`
}

type PRsResponse struct {
	PRs []PersonalRecord `json:"prs"`
}
