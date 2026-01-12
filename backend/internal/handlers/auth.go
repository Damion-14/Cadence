package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/damion-14/cadence/backend/internal/auth"
	"github.com/damion-14/cadence/backend/internal/config"
	"github.com/damion-14/cadence/backend/internal/database/queries"
	"github.com/damion-14/cadence/backend/internal/middleware"
	"github.com/damion-14/cadence/backend/internal/models"
)

type AuthHandler struct {
	userQueries *queries.UserQueries
	jwtConfig   config.JWTConfig
}

func NewAuthHandler(db *sql.DB, jwtConfig config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		userQueries: queries.NewUserQueries(db),
		jwtConfig:   jwtConfig,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid request body", 400))
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)

	if req.Email == "" || req.Password == "" || req.Username == "" {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Email, password, and username are required", 400))
		return
	}

	if len(req.Password) < 8 {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Password must be at least 8 characters", 400))
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	user, err := h.userQueries.CreateUser(r.Context(), req.Email, passwordHash, req.Username)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			respondError(w, r, models.NewAppError("CONFLICT", "Email already exists", 409))
			return
		}
		respondError(w, r, models.ErrInternalServer)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.jwtConfig.Secret, h.jwtConfig.ExpiryHours)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusCreated, models.AuthResponse{
		User:  *user,
		Token: token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Invalid request body", 400))
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		respondError(w, r, models.NewAppError("INVALID_INPUT", "Email and password are required", 400))
		return
	}

	user, err := h.userQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondError(w, r, models.NewAppError("UNAUTHORIZED", "Invalid email or password", 401))
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		respondError(w, r, models.NewAppError("UNAUTHORIZED", "Invalid email or password", 401))
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.jwtConfig.Secret, h.jwtConfig.ExpiryHours)
	if err != nil {
		respondError(w, r, models.ErrInternalServer)
		return
	}

	respondJSON(w, http.StatusOK, models.AuthResponse{
		User:  *user,
		Token: token,
	})
}

func respondError(w http.ResponseWriter, r *http.Request, appErr *models.AppError) {
	requestID := middleware.GetRequestID(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":       appErr.Code,
			"message":    appErr.Message,
			"details":    appErr.Details,
			"request_id": requestID,
		},
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
