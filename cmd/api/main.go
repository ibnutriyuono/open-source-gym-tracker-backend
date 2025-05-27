package main

import (
	"caloria-backend/internal/controller/permission"
	"caloria-backend/internal/controller/user"
	"caloria-backend/internal/env"
	"caloria-backend/internal/model"
	"expvar"
	"fmt"
	"log"
	"runtime"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &config{
		address: env.GetString("PORT", ":8080"),
	}

	dbHost := env.GetString("DB_HOST", "")
	dbUser := env.GetString("DB_USER", "")
	dbPassword := env.GetString("DB_PASSWORD", "")
	dbName := env.GetString("DB_NAME", "")
	dbPort := env.GetString("DB_PORT", "")
	dbSslMode := env.GetString("DB_SSL_MODE", "")
	dbTimeZone := env.GetString("DB_TIMEZONE", "")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbSslMode, dbTimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.UserToken{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.Exercise{},
		&model.WorkoutRoutine{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	app := &application{
		config: *config,
		db:     db,
		userController: &user.UserController{
			DB: db,
		},
		permissionController: &permission.PermissionController{
			DB: db,
		},
	}

	mux := app.mount()
	version := env.GetString("VERSION", "1.0.0")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm DB: %v", err)
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		stats := sqlDB.Stats()
		return map[string]any{
			"OpenConnections":   stats.OpenConnections,
			"InUse":             stats.InUse,
			"Idle":              stats.Idle,
			"WaitCount":         stats.WaitCount,
			"WaitDuration":      stats.WaitDuration.String(),
			"MaxIdleClosed":     stats.MaxIdleClosed,
			"MaxLifetimeClosed": stats.MaxLifetimeClosed,
		}
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("mem_stats", expvar.Func(func() any {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return map[string]any{
			"Alloc":      m.Alloc,
			"TotalAlloc": m.TotalAlloc,
			"Sys":        m.Sys,
			"NumGC":      m.NumGC,
		}
	}))

	log.Fatal(app.run(mux))
}
