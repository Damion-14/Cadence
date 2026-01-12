package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/damion-14/cadence/backend/internal/models"
)

type WorkoutQueries struct {
	db *sql.DB
}

func NewWorkoutQueries(db *sql.DB) *WorkoutQueries {
	return &WorkoutQueries{db: db}
}

func (q *WorkoutQueries) CreateWorkout(ctx context.Context, userID int, name string) (*models.WorkoutSession, error) {
	query := `
		INSERT INTO workout_sessions (user_id, name, status)
		VALUES ($1, $2, 'active')
		RETURNING id, user_id, name, status, started_at, completed_at, created_at, updated_at
	`

	var workout models.WorkoutSession
	err := q.db.QueryRowContext(ctx, query, userID, name).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.Status,
		&workout.StartedAt,
		&workout.CompletedAt,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	workout.Exercises = []models.Exercise{}
	return &workout, nil
}

func (q *WorkoutQueries) GetWorkoutByID(ctx context.Context, workoutID int) (*models.WorkoutSession, error) {
	query := `
		SELECT id, user_id, name, status, started_at, completed_at, created_at, updated_at
		FROM workout_sessions
		WHERE id = $1
	`

	var workout models.WorkoutSession
	err := q.db.QueryRowContext(ctx, query, workoutID).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.Status,
		&workout.StartedAt,
		&workout.CompletedAt,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workout not found")
	}

	if err != nil {
		return nil, err
	}

	exercises, err := q.GetExercisesByWorkoutID(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	workout.Exercises = exercises
	return &workout, nil
}

func (q *WorkoutQueries) GetActiveWorkout(ctx context.Context, userID int) (*models.WorkoutSession, error) {
	query := `
		SELECT id, user_id, name, status, started_at, completed_at, created_at, updated_at
		FROM workout_sessions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY started_at DESC
		LIMIT 1
	`

	var workout models.WorkoutSession
	err := q.db.QueryRowContext(ctx, query, userID).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.Status,
		&workout.StartedAt,
		&workout.CompletedAt,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	exercises, err := q.GetExercisesByWorkoutID(ctx, workout.ID)
	if err != nil {
		return nil, err
	}

	workout.Exercises = exercises
	return &workout, nil
}

func (q *WorkoutQueries) CompleteWorkout(ctx context.Context, workoutID int) error {
	query := `
		UPDATE workout_sessions
		SET status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND status = 'active'
	`

	result, err := q.db.ExecContext(ctx, query, workoutID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("workout not found or already completed")
	}

	return nil
}

func (q *WorkoutQueries) DeleteWorkout(ctx context.Context, workoutID int) error {
	query := `DELETE FROM workout_sessions WHERE id = $1`

	result, err := q.db.ExecContext(ctx, query, workoutID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("workout not found")
	}

	return nil
}

func (q *WorkoutQueries) GetExercisesByWorkoutID(ctx context.Context, workoutID int) ([]models.Exercise, error) {
	query := `
		SELECT id, workout_session_id, name, order_index, created_at, updated_at
		FROM exercises
		WHERE workout_session_id = $1
		ORDER BY order_index ASC
	`

	rows, err := q.db.QueryContext(ctx, query, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := []models.Exercise{}
	for rows.Next() {
		var exercise models.Exercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.WorkoutSessionID,
			&exercise.Name,
			&exercise.OrderIndex,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sets, err := q.GetSetsByExerciseID(ctx, exercise.ID)
		if err != nil {
			return nil, err
		}

		exercise.Sets = sets
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (q *WorkoutQueries) CreateExercise(ctx context.Context, workoutID int, name string, orderIndex int) (*models.Exercise, error) {
	query := `
		INSERT INTO exercises (workout_session_id, name, order_index)
		VALUES ($1, $2, $3)
		RETURNING id, workout_session_id, name, order_index, created_at, updated_at
	`

	var exercise models.Exercise
	err := q.db.QueryRowContext(ctx, query, workoutID, name, orderIndex).Scan(
		&exercise.ID,
		&exercise.WorkoutSessionID,
		&exercise.Name,
		&exercise.OrderIndex,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	exercise.Sets = []models.Set{}
	return &exercise, nil
}

func (q *WorkoutQueries) GetExerciseByID(ctx context.Context, exerciseID int) (*models.Exercise, error) {
	query := `
		SELECT id, workout_session_id, name, order_index, created_at, updated_at
		FROM exercises
		WHERE id = $1
	`

	var exercise models.Exercise
	err := q.db.QueryRowContext(ctx, query, exerciseID).Scan(
		&exercise.ID,
		&exercise.WorkoutSessionID,
		&exercise.Name,
		&exercise.OrderIndex,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("exercise not found")
	}

	if err != nil {
		return nil, err
	}

	sets, err := q.GetSetsByExerciseID(ctx, exercise.ID)
	if err != nil {
		return nil, err
	}

	exercise.Sets = sets
	return &exercise, nil
}

func (q *WorkoutQueries) UpdateExerciseName(ctx context.Context, exerciseID int, name string) error {
	query := `
		UPDATE exercises
		SET name = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := q.db.ExecContext(ctx, query, name, exerciseID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("exercise not found")
	}

	return nil
}

func (q *WorkoutQueries) DeleteExercise(ctx context.Context, exerciseID int) error {
	query := `DELETE FROM exercises WHERE id = $1`

	result, err := q.db.ExecContext(ctx, query, exerciseID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("exercise not found")
	}

	return nil
}

func (q *WorkoutQueries) GetNextExerciseOrderIndex(ctx context.Context, workoutID int) (int, error) {
	query := `
		SELECT COALESCE(MAX(order_index), -1) + 1
		FROM exercises
		WHERE workout_session_id = $1
	`

	var orderIndex int
	err := q.db.QueryRowContext(ctx, query, workoutID).Scan(&orderIndex)
	if err != nil {
		return 0, err
	}

	return orderIndex, nil
}

func (q *WorkoutQueries) GetSetsByExerciseID(ctx context.Context, exerciseID int) ([]models.Set, error) {
	query := `
		SELECT id, exercise_id, set_number, reps, weight, is_bodyweight, created_at, updated_at
		FROM sets
		WHERE exercise_id = $1
		ORDER BY set_number ASC
	`

	rows, err := q.db.QueryContext(ctx, query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sets := []models.Set{}
	for rows.Next() {
		var set models.Set
		err := rows.Scan(
			&set.ID,
			&set.ExerciseID,
			&set.SetNumber,
			&set.Reps,
			&set.Weight,
			&set.IsBodyweight,
			&set.CreatedAt,
			&set.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sets = append(sets, set)
	}

	return sets, nil
}

func (q *WorkoutQueries) CreateSet(ctx context.Context, exerciseID, setNumber, reps int, weight *float64, isBodyweight bool) (*models.Set, error) {
	query := `
		INSERT INTO sets (exercise_id, set_number, reps, weight, is_bodyweight)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, exercise_id, set_number, reps, weight, is_bodyweight, created_at, updated_at
	`

	var set models.Set
	err := q.db.QueryRowContext(ctx, query, exerciseID, setNumber, reps, weight, isBodyweight).Scan(
		&set.ID,
		&set.ExerciseID,
		&set.SetNumber,
		&set.Reps,
		&set.Weight,
		&set.IsBodyweight,
		&set.CreatedAt,
		&set.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &set, nil
}

func (q *WorkoutQueries) DeleteSetsByExerciseID(ctx context.Context, exerciseID int) error {
	query := `DELETE FROM sets WHERE exercise_id = $1`

	_, err := q.db.ExecContext(ctx, query, exerciseID)
	return err
}
