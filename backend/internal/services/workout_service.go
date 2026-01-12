package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/damion-14/cadence/backend/internal/cache"
	"github.com/damion-14/cadence/backend/internal/database/queries"
	"github.com/damion-14/cadence/backend/internal/models"
)

type WorkoutService struct {
	workoutQueries *queries.WorkoutQueries
	cache          *cache.Cache
}

func NewWorkoutService(db *sql.DB, cacheClient *cache.Cache) *WorkoutService {
	return &WorkoutService{
		workoutQueries: queries.NewWorkoutQueries(db),
		cache:          cacheClient,
	}
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, userID int, name string) (*models.WorkoutSession, error) {
	workout, err := s.workoutQueries.CreateWorkout(ctx, userID, name)
	if err != nil {
		return nil, err
	}

	if err := s.cacheActiveWorkout(ctx, userID, workout); err != nil {
		fmt.Printf("Failed to cache active workout: %v\n", err)
	}

	return workout, nil
}

func (s *WorkoutService) GetWorkout(ctx context.Context, workoutID int) (*models.WorkoutSession, error) {
	return s.workoutQueries.GetWorkoutByID(ctx, workoutID)
}

func (s *WorkoutService) GetActiveWorkout(ctx context.Context, userID int) (*models.WorkoutSession, error) {
	cacheKey := cache.GetActiveWorkoutKey(userID)

	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var workout models.WorkoutSession
		if err := json.Unmarshal([]byte(cachedData), &workout); err == nil {
			return &workout, nil
		}
	}

	workout, err := s.workoutQueries.GetActiveWorkout(ctx, userID)
	if err != nil {
		return nil, err
	}

	if workout != nil {
		if err := s.cacheActiveWorkout(ctx, userID, workout); err != nil {
			fmt.Printf("Failed to cache active workout: %v\n", err)
		}
	}

	return workout, nil
}

func (s *WorkoutService) CompleteWorkout(ctx context.Context, userID, workoutID int) (*models.WorkoutSession, error) {
	if err := s.workoutQueries.CompleteWorkout(ctx, workoutID); err != nil {
		return nil, err
	}

	if err := s.invalidateCachesOnComplete(ctx, userID); err != nil {
		fmt.Printf("Failed to invalidate caches: %v\n", err)
	}

	return s.workoutQueries.GetWorkoutByID(ctx, workoutID)
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, userID, workoutID int) error {
	workout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return err
	}

	if workout.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	if err := s.workoutQueries.DeleteWorkout(ctx, workoutID); err != nil {
		return err
	}

	if workout.Status == "active" {
		cacheKey := cache.GetActiveWorkoutKey(userID)
		s.cache.Delete(ctx, cacheKey)
	}

	return nil
}

func (s *WorkoutService) AddExercise(ctx context.Context, userID, workoutID int, name string, sets []models.SetInput) (*models.Exercise, error) {
	workout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	if workout.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	if workout.Status != "active" {
		return nil, fmt.Errorf("workout is not active")
	}

	orderIndex, err := s.workoutQueries.GetNextExerciseOrderIndex(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	exercise, err := s.workoutQueries.CreateExercise(ctx, workoutID, name, orderIndex)
	if err != nil {
		return nil, err
	}

	for i, setInput := range sets {
		set, err := s.workoutQueries.CreateSet(ctx, exercise.ID, i+1, setInput.Reps, setInput.Weight, setInput.IsBodyweight)
		if err != nil {
			return nil, err
		}
		exercise.Sets = append(exercise.Sets, *set)
	}

	if workout.Status == "active" {
		updatedWorkout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
		if err == nil {
			s.cacheActiveWorkout(ctx, userID, updatedWorkout)
		}
	}

	return exercise, nil
}

func (s *WorkoutService) UpdateExercise(ctx context.Context, userID, workoutID, exerciseID int, name *string, sets []models.SetInput) (*models.Exercise, error) {
	workout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	if workout.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	if workout.Status != "active" {
		return nil, fmt.Errorf("workout is not active")
	}

	exercise, err := s.workoutQueries.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	if exercise.WorkoutSessionID != workoutID {
		return nil, fmt.Errorf("exercise does not belong to this workout")
	}

	if name != nil && *name != "" {
		if err := s.workoutQueries.UpdateExerciseName(ctx, exerciseID, *name); err != nil {
			return nil, err
		}
	}

	if len(sets) > 0 {
		if err := s.workoutQueries.DeleteSetsByExerciseID(ctx, exerciseID); err != nil {
			return nil, err
		}

		for i, setInput := range sets {
			_, err := s.workoutQueries.CreateSet(ctx, exerciseID, i+1, setInput.Reps, setInput.Weight, setInput.IsBodyweight)
			if err != nil {
				return nil, err
			}
		}
	}

	updatedExercise, err := s.workoutQueries.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	if workout.Status == "active" {
		updatedWorkout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
		if err == nil {
			s.cacheActiveWorkout(ctx, userID, updatedWorkout)
		}
	}

	return updatedExercise, nil
}

func (s *WorkoutService) DeleteExercise(ctx context.Context, userID, workoutID, exerciseID int) error {
	workout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return err
	}

	if workout.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	if workout.Status != "active" {
		return fmt.Errorf("workout is not active")
	}

	exercise, err := s.workoutQueries.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		return err
	}

	if exercise.WorkoutSessionID != workoutID {
		return fmt.Errorf("exercise does not belong to this workout")
	}

	if err := s.workoutQueries.DeleteExercise(ctx, exerciseID); err != nil {
		return err
	}

	if workout.Status == "active" {
		updatedWorkout, err := s.workoutQueries.GetWorkoutByID(ctx, workoutID)
		if err == nil {
			s.cacheActiveWorkout(ctx, userID, updatedWorkout)
		}
	}

	return nil
}

func (s *WorkoutService) cacheActiveWorkout(ctx context.Context, userID int, workout *models.WorkoutSession) error {
	cacheKey := cache.GetActiveWorkoutKey(userID)
	data, err := json.Marshal(workout)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, cacheKey, data, cache.TTLActiveWorkout)
}

func (s *WorkoutService) invalidateCachesOnComplete(ctx context.Context, userID int) error {
	keys := []string{
		cache.GetActiveWorkoutKey(userID),
		cache.GetUserPRsKey(userID),
	}

	return s.cache.Delete(ctx, keys...)
}
