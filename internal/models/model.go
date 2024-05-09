package models

// import "github.com/google/uuid"

// import

// User represents a user in the system
type User struct {
	ID       string `gorm:"primaryKey"`
	Email    string
	Name     string
	Password string
	Phone    string `binding:"required" gorm:"unique"`
}
