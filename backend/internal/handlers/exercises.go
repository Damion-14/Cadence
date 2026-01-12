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

type ExerciseHandler struct {
	workoutService *services.WorkoutService
}

func NewExerciseHandler(workoutService *services.WorkoutService) *ExerciseHandler {
	return &ExerciseHandler{
		workoutService: workoutService,
	}
}

func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workoutID, err := strconv.Atoi(r.PathValue("workoutId"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid workout ID", 400))
		return
	}

	var req models.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid request body", 400))
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Exercise name is required", 400))
		return
	}

	if len(req.Sets) == 0 {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "At least one set is required", 400))
		return
	}

	for i, set := range req.Sets {
		if set.Reps <= 0 {
			respondError(w, r, models.NewAppError("INVALID_INPUT", "Reps must be greater than 0", 400))
			return
		}
		if !set.IsBodyweight && (set.Weight == nil || *set.Weight <= 0) {
			respondError(w, r, models.NewAppError("INVALID_INPUT", "Weight must be greater than 0 for non-bodyweight sets", 400))
			return
		}
		if set.IsBodyweight {
			req.Sets[i].Weight = nil
		}
	}

	exercise, err := h.workoutService.AddExercise(r.Context(), userID, workoutID, req.Name, req.Sets)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, r, models.ErrNotFound)
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			respondError(w, r, models.ErrForbidden)
			return
		}
		if strings.Contains(err.Error(), "not active") {
			respondError(w, r, models.NewAppError("INVALID_INPUT", "Workout is not active", 400))
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusCreated, models.CreateExerciseResponse{
		Exercise: *exercise,
	})
}

func (h *ExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workoutID, err := strconv.Atoi(r.PathValue("workoutId"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid workout ID", 400))
		return
	}

	exerciseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid exercise ID", 400))
		return
	}

	var req models.UpdateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid request body", 400))
		return
	}

	var namePtr *string
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)
		namePtr = &req.Name
	}

	if len(req.Sets) > 0 {
		for i, set := range req.Sets {
			if set.Reps <= 0 {
				respondError(w, r, models.NewAppError("INVALID_INPUT", "Reps must be greater than 0", 400))
				return
			}
			if !set.IsBodyweight && (set.Weight == nil || *set.Weight <= 0) {
				respondError(w, r, models.NewAppError("INVALID_INPUT", "Weight must be greater than 0 for non-bodyweight sets", 400))
				return
			}
			if set.IsBodyweight {
				req.Sets[i].Weight = nil
			}
		}
	}

	exercise, err := h.workoutService.UpdateExercise(r.Context(), userID, workoutID, exerciseID, namePtr, req.Sets)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, r, models.ErrNotFound)
			return
		}
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "does not belong") {
			respondError(w, r, models.ErrForbidden)
			return
		}
		if strings.Contains(err.Error(), "not active") {
			respondError(w, r, models.NewAppError("INVALID_INPUT", "Workout is not active", 400))
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.UpdateExerciseResponse{
		Exercise: *exercise,
	})
}

func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		respondError(w, r, models.ErrUnauthorized)
		return
	}

	workoutID, err := strconv.Atoi(r.PathValue("workoutId"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid workout ID", 400))
		return
	}

	exerciseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid exercise ID", 400))
		return
	}

	if err := h.workoutService.DeleteExercise(r.Context(), userID, workoutID, exerciseID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, r, models.ErrNotFound)
			return
		}
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "does not belong") {
			respondError(w, r, models.ErrForbidden)
			return
		}
		if strings.Contains(err.Error(), "not active") {
			respondError(w, r, models.NewAppError("INVALID_INPUT", "Workout is not active", 400))
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.DeleteResponse{
		Message: "Exercise deleted successfully",
	})
}
