package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/damion-14/cadence/backend/internal/cache"
	"github.com/damion-14/cadence/backend/internal/database/queries"
	"github.com/damion-14/cadence/backend/internal/models"
)

type StatsService struct {
	statsQueries *queries.StatsQueries
	cache        *cache.Cache
}

func NewStatsService(db *sql.DB, cacheClient *cache.Cache) *StatsService {
	return &StatsService{
		statsQueries: queries.NewStatsQueries(db),
		cache:        cacheClient,
	}
}

func (s *StatsService) GetPersonalRecords(ctx context.Context, userID int) ([]models.PersonalRecord, error) {
	cacheKey := cache.GetUserPRsKey(userID)

	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var prs []models.PersonalRecord
		if err := json.Unmarshal([]byte(cachedData), &prs); err == nil {
			return prs, nil
		}
	}

	prs, err := s.statsQueries.GetPersonalRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(prs); err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.TTLUserPRs)
	}

	return prs, nil
}

func (s *StatsService) GetWorkoutHistory(ctx context.Context, userID int, limit, offset int) ([]models.WorkoutSummary, int, error) {
	return s.statsQueries.GetWorkoutHistory(ctx, userID, limit, offset)
}

func (s *StatsService) GetWeeklySummary(ctx context.Context, userID int, weekStr string) (*models.WeeklySummary, error) {
	var year, week int
	if weekStr == "" {
		year, week = time.Now().ISOWeek()
	} else {
		parts := strings.Split(weekStr, "-W")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid week format, expected YYYY-WNN")
		}
		var err error
		year, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid year in week format")
		}
		week, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid week number in week format")
		}
	}

	cacheKey := cache.GetWeeklySummaryKey(userID, fmt.Sprintf("%d-W%02d", year, week))

	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var summary models.WeeklySummary
		if err := json.Unmarshal([]byte(cachedData), &summary); err == nil {
			return &summary, nil
		}
	}

	summary, err := s.statsQueries.GetWeeklySummary(ctx, userID, year, week)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(summary); err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.TTLWeeklySummary)
	}

	return summary, nil
}

func (s *StatsService) GetExerciseProgress(ctx context.Context, userID int, exerciseName string, period string) ([]models.ProgressDataPoint, error) {
	days := 30
	if period != "" {
		if strings.HasSuffix(period, "d") {
			d, err := strconv.Atoi(strings.TrimSuffix(period, "d"))
			if err == nil && d > 0 && d <= 365 {
				days = d
			}
		}
	}

	cacheKey := cache.GetExerciseProgressKey(userID, exerciseName)

	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var dataPoints []models.ProgressDataPoint
		if err := json.Unmarshal([]byte(cachedData), &dataPoints); err == nil {
			return dataPoints, nil
		}
	}

	dataPoints, err := s.statsQueries.GetExerciseProgress(ctx, userID, exerciseName, days)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(dataPoints); err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.TTLExerciseProgress)
	}

	return dataPoints, nil
}
