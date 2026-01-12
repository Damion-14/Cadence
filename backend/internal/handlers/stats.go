package handlers

import (
	"net/http"
	"strconv"

	"github.com/damion-14/cadence/backend/internal/middleware"
	"github.com/damion-14/cadence/backend/internal/models"
	"github.com/damion-14/cadence/backend/internal/services"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

func (h *StatsHandler) GetPRs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	prs, err := h.statsService.GetPersonalRecords(r.Context(), userID)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.PRsResponse{
		PRs: prs,
	})
}

func (h *StatsHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	workouts, total, err := h.statsService.GetWorkoutHistory(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.HistoryResponse{
		Workouts: workouts,
		Total:    total,
	})
}

func (h *StatsHandler) GetWeeklySummary(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	week := r.URL.Query().Get("week")

	summary, err := h.statsService.GetWeeklySummary(r.Context(), userID, week)
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", err.Error(), 400))
		return
	}

	respondJSON(w, http.StatusOK, summary)
}

func (h *StatsHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	exerciseName := r.PathValue("exerciseName")
	if exerciseName == "" {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Exercise name is required", 400))
		return
	}

	period := r.URL.Query().Get("period")

	dataPoints, err := h.statsService.GetExerciseProgress(r.Context(), userID, exerciseName, period)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.ProgressResponse{
		ExerciseName: exerciseName,
		DataPoints:   dataPoints,
	})
}
