package model

import (
	"time"
)

type User struct {
	ID           string     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName    string     `gorm:"type:varchar(255);not null" json:"first_name" validate:"required,min=2,max=255"`
	LastName     string     `gorm:"type:varchar(255);not null" json:"last_name" validate:"required,min=2,max=255"`
	Email        string     `gorm:"type:citext;unique;not null" json:"email" validate:"required,email"`
	Password     string     `gorm:"type:text;not null" json:"password" validate:"required,strong_password"` // bcrypt password; strong password
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"-"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"-"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`                  // soft delete
	IsDeleted    bool       `gorm:"default:false" json:"is_deleted"` // soft delete timestamp
	IsVerified   bool       `gorm:"default:false" json:"is_verified"`
	IsBanned     bool       `gorm:"default:false" json:"is_banned"`
	BannedAt     *time.Time `json:"banned_at,omitempty"`                                                      // only if banned
	Avatar       string     `gorm:"type:varchar(255)" json:"avatar,omitempty" validate:"omitempty,url"`       // URL to avatar image
	PhoneNumber  string     `gorm:"type:varchar(20)" json:"phone_number,omitempty" validate:"omitempty,e164"` // E.164 phone format
	YearOfBirth  string     `gorm:"type:char(4)" json:"year_of_birth,omitempty" validate:"omitempty,len=4,numeric"`
	MonthOfBirth string     `gorm:"type:char(2)" json:"month_of_birth,omitempty" validate:"omitempty,len=2,numeric"`
	DateOfBirth  string     `gorm:"type:char(2)" json:"date_of_birth,omitempty" validate:"omitempty,len=2,numeric"`
	Gender       *string    `gorm:"type:varchar(20)" json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	UserTokens []UserToken `gorm:"foreignKey:UserID" json:"-"`
}

// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	if u.ID == uuid.Nil {
// 		u.ID = uuid.New()
// 	}
// 	return
// }
