package queries

import (
	"context"
	"database/sql"
	"time"

	"github.com/damion-14/cadence/backend/internal/models"
)

type StatsQueries struct {
	db *sql.DB
}

func NewStatsQueries(db *sql.DB) *StatsQueries {
	return &StatsQueries{db: db}
}

func (q *StatsQueries) GetPersonalRecords(ctx context.Context, userID int) ([]models.PersonalRecord, error) {
	query := `
		WITH exercise_prs AS (
			SELECT
				e.name AS exercise_name,
				MAX(s.weight) AS max_weight,
				MAX(s.reps) AS max_reps,
				MAX(COALESCE(s.weight, 0) * s.reps) AS max_volume,
				ws.completed_at
			FROM exercises e
			JOIN workout_sessions ws ON e.workout_session_id = ws.id
			JOIN sets s ON s.exercise_id = e.id
			WHERE ws.user_id = $1 AND ws.status = 'completed'
			GROUP BY e.name, ws.completed_at
		),
		ranked_prs AS (
			SELECT
				exercise_name,
				max_weight,
				max_reps,
				max_volume,
				completed_at,
				ROW_NUMBER() OVER (PARTITION BY exercise_name ORDER BY max_volume DESC, completed_at DESC) as rn
			FROM exercise_prs
		)
		SELECT
			exercise_name,
			max_weight,
			max_reps,
			max_volume,
			completed_at
		FROM ranked_prs
		WHERE rn = 1
		ORDER BY exercise_name ASC
	`

	rows, err := q.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := []models.PersonalRecord{}
	for rows.Next() {
		var pr models.PersonalRecord
		err := rows.Scan(
			&pr.ExerciseName,
			&pr.MaxWeight,
			&pr.MaxReps,
			&pr.MaxVolume,
			&pr.AchievedAt,
		)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (q *StatsQueries) GetWorkoutHistory(ctx context.Context, userID int, limit, offset int) ([]models.WorkoutSummary, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM workout_sessions
		WHERE user_id = $1 AND status = 'completed'
	`

	var total int
	err := q.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			ws.id,
			ws.name,
			ws.completed_at,
			COUNT(DISTINCT e.id) AS exercise_count,
			COUNT(s.id) AS total_sets,
			COALESCE(SUM(COALESCE(s.weight, 0) * s.reps), 0) AS total_volume
		FROM workout_sessions ws
		LEFT JOIN exercises e ON e.workout_session_id = ws.id
		LEFT JOIN sets s ON s.exercise_id = e.id
		WHERE ws.user_id = $1 AND ws.status = 'completed'
		GROUP BY ws.id, ws.name, ws.completed_at
		ORDER BY ws.completed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := q.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	workouts := []models.WorkoutSummary{}
	for rows.Next() {
		var workout models.WorkoutSummary
		err := rows.Scan(
			&workout.ID,
			&workout.Name,
			&workout.CompletedAt,
			&workout.ExerciseCount,
			&workout.TotalSets,
			&workout.TotalVolume,
		)
		if err != nil {
			return nil, 0, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, total, nil
}

func (q *StatsQueries) GetWeeklySummary(ctx context.Context, userID int, year int, week int) (*models.WeeklySummary, error) {
	startOfWeek, endOfWeek := getWeekBounds(year, week)

	summaryQuery := `
		SELECT
			COUNT(DISTINCT ws.id) AS total_workouts,
			COUNT(DISTINCT e.id) AS total_exercises,
			COALESCE(SUM(COALESCE(s.weight, 0) * s.reps), 0) AS total_volume
		FROM workout_sessions ws
		LEFT JOIN exercises e ON e.workout_session_id = ws.id
		LEFT JOIN sets s ON s.exercise_id = e.id
		WHERE ws.user_id = $1
			AND ws.status = 'completed'
			AND ws.completed_at >= $2
			AND ws.completed_at < $3
	`

	var summary models.WeeklySummary
	err := q.db.QueryRowContext(ctx, summaryQuery, userID, startOfWeek, endOfWeek).Scan(
		&summary.TotalWorkouts,
		&summary.TotalExercises,
		&summary.TotalVolume,
	)
	if err != nil {
		return nil, err
	}

	workoutsQuery := `
		SELECT
			ws.id,
			ws.name,
			ws.completed_at,
			EXTRACT(DOW FROM ws.completed_at)::INTEGER AS day_of_week
		FROM workout_sessions ws
		WHERE ws.user_id = $1
			AND ws.status = 'completed'
			AND ws.completed_at >= $2
			AND ws.completed_at < $3
		ORDER BY ws.completed_at ASC
	`

	rows, err := q.db.QueryContext(ctx, workoutsQuery, userID, startOfWeek, endOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []models.WorkoutInWeek{}
	for rows.Next() {
		var workout models.WorkoutInWeek
		err := rows.Scan(
			&workout.ID,
			&workout.Name,
			&workout.CompletedAt,
			&workout.DayOfWeek,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	summary.Week = formatWeek(year, week)
	summary.Workouts = workouts

	return &summary, nil
}

func (q *StatsQueries) GetExerciseProgress(ctx context.Context, userID int, exerciseName string, days int) ([]models.ProgressDataPoint, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT
			DATE(ws.completed_at) AS workout_date,
			MAX(s.weight) AS max_weight,
			MAX(s.reps) AS max_reps,
			SUM(COALESCE(s.weight, 0) * s.reps) AS volume
		FROM exercises e
		JOIN workout_sessions ws ON e.workout_session_id = ws.id
		JOIN sets s ON s.exercise_id = e.id
		WHERE ws.user_id = $1
			AND ws.status = 'completed'
			AND e.name = $2
			AND ws.completed_at >= $3
		GROUP BY DATE(ws.completed_at)
		ORDER BY workout_date ASC
	`

	rows, err := q.db.QueryContext(ctx, query, userID, exerciseName, cutoffDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataPoints := []models.ProgressDataPoint{}
	for rows.Next() {
		var point models.ProgressDataPoint
		err := rows.Scan(
			&point.Date,
			&point.MaxWeight,
			&point.MaxReps,
			&point.Volume,
		)
		if err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, point)
	}

	return dataPoints, nil
}

func getWeekBounds(year, week int) (time.Time, time.Time) {
	jan1 := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	daysToMonday := int(time.Monday - jan1.Weekday())
	if daysToMonday > 0 {
		daysToMonday -= 7
	}

	firstMonday := jan1.AddDate(0, 0, daysToMonday)

	startOfWeek := firstMonday.AddDate(0, 0, (week-1)*7)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return startOfWeek, endOfWeek
}

func formatWeek(year, week int) string {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, (week-1)*7).
		Format("2006-W02")
}
