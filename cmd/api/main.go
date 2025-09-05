package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/EduardoMark/gastro-api/internal/auth"
	"github.com/EduardoMark/gastro-api/internal/config"
	"github.com/EduardoMark/gastro-api/internal/database"
	"github.com/EduardoMark/gastro-api/internal/dishes"
	appmw "github.com/EduardoMark/gastro-api/internal/middleware"
	"github.com/EduardoMark/gastro-api/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	env := config.Load()

	db, err := database.New(env)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to running migrate: %v", err)
	}

	authService := auth.NewAuthJWTService(env)
	jwtMiddleware := appmw.NewJWTMiddleware(authService)

	userRepo := users.NewUserRepo(db)
	userService := users.NewUserService(userRepo)
	userHandler := users.NerUserHandler(userService, jwtMiddleware, authService)

	dishRepo := dishes.NewDishRepository(db)
	dishService := dishes.NewDishService(dishRepo)
	dishHandler := dishes.NewDishHandler(dishService, jwtMiddleware)

	router := chi.NewMux()
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		userHandler.UserRoutes(r)
		dishHandler.DishRoutes(r)
	})

	server := http.Server{
		Addr:              ":3000",
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
