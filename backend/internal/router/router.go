package router

import (
	"database/sql"
	"net/http"

	"github.com/damion-14/cadence/backend/internal/auth"
	"github.com/damion-14/cadence/backend/internal/cache"
	"github.com/damion-14/cadence/backend/internal/config"
	"github.com/damion-14/cadence/backend/internal/handlers"
	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	DB              *sql.DB
	Redis           *redis.Client
	Cache           *cache.Cache
	Config          *config.Config
	AuthHandler     *handlers.AuthHandler
	WorkoutHandler  *handlers.WorkoutHandler
	ExerciseHandler *handlers.ExerciseHandler
	StatsHandler    *handlers.StatsHandler
}

func NewRouter(deps *Dependencies) *http.ServeMux {
	mux := http.NewServeMux()

	authMiddleware := auth.Middleware(deps.Config.JWT.Secret)

	mux.HandleFunc("POST /api/v1/auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", deps.AuthHandler.Login)

	mux.Handle("POST /api/v1/workouts", authMiddleware(http.HandlerFunc(deps.WorkoutHandler.Create)))
	mux.Handle("GET /api/v1/workouts/active", authMiddleware(http.HandlerFunc(deps.WorkoutHandler.GetActive)))
	mux.Handle("GET /api/v1/workouts/{id}", authMiddleware(http.HandlerFunc(deps.WorkoutHandler.GetByID)))
	mux.Handle("POST /api/v1/workouts/{id}/complete", authMiddleware(http.HandlerFunc(deps.WorkoutHandler.Complete)))
	mux.Handle("DELETE /api/v1/workouts/{id}", authMiddleware(http.HandlerFunc(deps.WorkoutHandler.Delete)))

	mux.Handle("POST /api/v1/workouts/{workoutId}/exercises", authMiddleware(http.HandlerFunc(deps.ExerciseHandler.Create)))
	mux.Handle("PUT /api/v1/workouts/{workoutId}/exercises/{id}", authMiddleware(http.HandlerFunc(deps.ExerciseHandler.Update)))
	mux.Handle("DELETE /api/v1/workouts/{workoutId}/exercises/{id}", authMiddleware(http.HandlerFunc(deps.ExerciseHandler.Delete)))

	mux.Handle("GET /api/v1/history", authMiddleware(http.HandlerFunc(deps.StatsHandler.GetHistory)))
	mux.Handle("GET /api/v1/stats/prs", authMiddleware(http.HandlerFunc(deps.StatsHandler.GetPRs)))
	mux.Handle("GET /api/v1/stats/weekly", authMiddleware(http.HandlerFunc(deps.StatsHandler.GetWeeklySummary)))
	mux.Handle("GET /api/v1/stats/progress/{exerciseName}", authMiddleware(http.HandlerFunc(deps.StatsHandler.GetProgress)))

	mux.HandleFunc("GET /health", healthCheck)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
