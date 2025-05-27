package model

import "time"

type Exercise struct {
	ID              string           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string           `gorm:"type:varchar(255);not null;unique" json:"name"`
	Description     string           `gorm:"type:text" json:"description,omitempty"`
	WorkoutRoutines []WorkoutRoutine `gorm:"many2many:workout_routine_exercises;" json:"-"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Slug            string           `gorm:"type:varchar(255);not null" json:"slug"`
}
