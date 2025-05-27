package model

type Role struct {
	ID   string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"uniqueIndex"`
	RolePermissions []RolePermission
}

type UserRole struct {
	ID     string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID string
	RoleID string
	User   User
	Role   Role
}