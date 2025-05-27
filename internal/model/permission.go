package model

type Permission struct {
	ID   string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"uniqueIndex"`
}

type RolePermission struct {
	ID           string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoleID       string
	PermissionID string
	Role         Role
	Permission   Permission
}
