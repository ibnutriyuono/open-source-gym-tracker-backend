package main

import (
	"caloria-backend/internal/controller/health"
	"caloria-backend/internal/controller/permission"
	"caloria-backend/internal/controller/user"
	customMiddleware "caloria-backend/internal/middleware"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

type application struct {
	config config
	db     *gorm.DB
	userController   *user.UserController
	healthController *health.HealthController
	permissionController *permission.PermissionController
}

type config struct {
	address string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		defer func() {
	// 			if err := recover(); err != nil {
	// 				log.Printf("RECOVERED in custom middleware: %v", err)
	// 				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 			}
	// 		}()
	// 		next.ServeHTTP(w, r)
	// 	})
	// })

	r.Use(middleware.Timeout(60 * time.Second))
	// userController := &user.UserController{
	// 	DB: app.db,
	// }

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthController.HealthCheck)
		r.Route("/users", func(r chi.Router) {
			r.Use(customMiddleware.Authentication(app.db))
			r.Get("/", app.userController.FindAll)
			r.Put("/{id}", app.userController.Update)
			r.Delete("/{id}", app.userController.Delete)
			r.Get("/{id}", app.userController.FindById)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.userController.Create)
			r.Get("/login", app.userController.Login)
			r.Get("/refresh", app.userController.RefreshToken)
		})
		r.Route("/permissions", func(r chi.Router) {
			r.Use(customMiddleware.Authentication(app.db))
			r.Get("/", app.permissionController.FindAll)
			r.Post("/", app.permissionController.Create)
			r.Put("/{id}", app.permissionController.Update)
			r.Delete("/{id}", app.permissionController.Delete)
			r.Get("/{id}", app.permissionController.FindById)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	server := &http.Server{
		Addr:         app.config.address,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server started at %s", app.config.address)

	return server.ListenAndServe()
}
