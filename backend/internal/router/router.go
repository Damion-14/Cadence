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
	DB         *sql.DB
	Redis      *redis.Client
	Cache      *cache.Cache
	Config     *config.Config
	AuthHandler *handlers.AuthHandler
}

func NewRouter(deps *Dependencies) *http.ServeMux {
	mux := http.NewServeMux()

	authMiddleware := auth.Middleware(deps.Config.JWT.Secret)

	mux.HandleFunc("POST /api/v1/auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", deps.AuthHandler.Login)

	mux.HandleFunc("GET /health", healthCheck)

	_ = authMiddleware

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
