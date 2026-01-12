package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/damion-14/cadence/backend/internal/middleware"
	"github.com/damion-14/cadence/backend/internal/models"
	"github.com/damion-14/cadence/backend/internal/services"
)

type WorkoutHandler struct {
	workoutService *services.WorkoutService
}

func NewWorkoutHandler(workoutService *services.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: workoutService,
	}
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	var req models.CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid request body", 400))
		return
	}

	workout, err := h.workoutService.CreateWorkout(r.Context(), userID, req.Name)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusCreated, models.CreateWorkoutResponse{
		Workout: *workout,
	})
}

func (h *WorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workoutID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid workout ID", 400))
		return
	}

	workout, err := h.workoutService.GetWorkout(r.Context(), workoutID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, r, models.ErrNotFound)
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	if workout.UserID != userID {
		respondError(w, r, models.ErrForbidden)
		return
	}

	respondJSON(w, http.StatusOK, models.GetWorkoutResponse{
		Workout: *workout,
	})
}

func (h *WorkoutHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workout, err := h.workoutService.GetActiveWorkout(r.Context(), userID)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	if workout == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"workout": nil,
		})
		return
	}

	respondJSON(w, http.StatusOK, models.GetWorkoutResponse{
		Workout: *workout,
	})
}

func (h *WorkoutHandler) Complete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workoutID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid workout ID", 400))
		return
	}

	workout, err := h.workoutService.CompleteWorkout(r.Context(), userID, workoutID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "already completed") {
			respondError(w, r, models.NewAppError("INVALID_INPUT", err.Error(), 400))
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.GetWorkoutResponse{
		Workout: *workout,
	})
}

func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workoutID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid workout ID", 400))
		return
	}

	if err := h.workoutService.DeleteWorkout(r.Context(), userID, workoutID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, r, models.ErrNotFound)
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			respondError(w, r, models.ErrForbidden)
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.DeleteResponse{
		Message: "Workout deleted successfully",
	})
}
