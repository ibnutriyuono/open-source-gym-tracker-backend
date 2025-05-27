package model

import (
	"time"
)

type UserToken struct {
	ID           string     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       string     `gorm:"type:uuid;not null;index" json:"user_id"`
	AccessToken  string     `gorm:"type:text;not null" json:"access_token"`
	RefreshToken string     `gorm:"type:text;not null" json:"refresh_token"`
	UserAgent    string     `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress    string     `gorm:"type:text" json:"ip_address,omitempty"`
	ExpiresAt    *time.Time `gorm:"type:timestamptz" json:"expires_at,omitempty"`
	IsRevoked    bool       `gorm:"default:false" json:"is_revoked"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"-"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}
