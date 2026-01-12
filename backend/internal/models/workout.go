package models

import "time"

type WorkoutSession struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Exercises   []Exercise `json:"exercises,omitempty"`
}

type Exercise struct {
	ID               int       `json:"id"`
	WorkoutSessionID int       `json:"workout_session_id"`
	Name             string    `json:"name"`
	OrderIndex       int       `json:"order_index"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Sets             []Set     `json:"sets,omitempty"`
}

type Set struct {
	ID           int       `json:"id"`
	ExerciseID   int       `json:"exercise_id"`
	SetNumber    int       `json:"set_number"`
	Reps         int       `json:"reps"`
	Weight       *float64  `json:"weight,omitempty"`
	IsBodyweight bool      `json:"is_bodyweight"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateWorkoutRequest struct {
	Name string `json:"name"`
}

type CreateWorkoutResponse struct {
	Workout WorkoutSession `json:"workout"`
}

type GetWorkoutResponse struct {
	Workout WorkoutSession `json:"workout"`
}

type CreateExerciseRequest struct {
	Name string    `json:"name"`
	Sets []SetInput `json:"sets"`
}

type SetInput struct {
	Reps         int      `json:"reps"`
	Weight       *float64 `json:"weight,omitempty"`
	IsBodyweight bool     `json:"is_bodyweight"`
}

type UpdateExerciseRequest struct {
	Name string     `json:"name,omitempty"`
	Sets []SetInput `json:"sets,omitempty"`
}

type CreateExerciseResponse struct {
	Exercise Exercise `json:"exercise"`
}

type UpdateExerciseResponse struct {
	Exercise Exercise `json:"exercise"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}
