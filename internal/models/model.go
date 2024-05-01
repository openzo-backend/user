package models

// import "github.com/google/uuid"

// import

// User represents a user in the system
type User struct {
	ID       string `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Name     string
	Password string
}
