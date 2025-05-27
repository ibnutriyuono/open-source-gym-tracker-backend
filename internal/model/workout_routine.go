package model

import "time"

type WorkoutRoutine struct {
	ID        string     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    string     `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	Exercises []Exercise `gorm:"many2many:workout_routine_exercises;" json:"exercises"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Slug      string     `gorm:"type:varchar(255);not null" json:"slug"`
}
