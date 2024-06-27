package models

import "time"

// import "github.com/google/uuid"

// import

// User represents a user in the system
// type User struct {
// 	ID       string `gorm:"primaryKey"`
// 	Email    string
// 	Name     string
// 	Password string
// 	Phone    string `binding:"required" gorm:"unique"`
// }

type OTP struct {
	ID        string `gorm:"primaryKey"`
	Phone     string
	HashedOTP string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// type UserData struct {
// 	Id                string `json:"id" gorm:"primaryKey"`
// 	UserId            string `json:"user_id"`
// 	Latitude          string `json:"latitude"`
// 	Longitude         string `json:"longitude"`
// 	Address           string `json:"address"`
// 	Pincode           string `json:"pincode"`
// 	City              string `json:"city"`
// 	State             string `json:"state"`
// 	Country           string `json:"country"`
// 	NotificationToken string `json:"notification_token"`
// }

type User struct {
	ID                string  `gorm:"primaryKey" json:"id"`
	Email             *string `json:"email,omitempty"`
	Name              *string `json:"name,omitempty"`
	Password          *string `json:"password,omitempty"`
	Phone             string  `json:"phone" gorm:"size:15"`
	Latitude          *string `json:"latitude,omitempty"`
	Longitude         *string `json:"longitude,omitempty"`
	Address           *string `json:"address,omitempty"`
	Pincode           *string `json:"pincode,omitempty"`
	City              *string `json:"city,omitempty"`
	State             *string `json:"state,omitempty"`
	Country           *string `json:"country,omitempty"`
	NotificationToken *string `json:"notification_token,omitempty"`
	IsVerified        bool    `json:"is_verified"`
	CreatedAt         time.Time
}

// type LoggedInUser struct {
// 	UserID    string `gorm:"primaryKey"`
// 	Preferences map[string]interface{}
// 	// Add other fields specific to logged-in users
// }

type Customer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	UserDataId string `json:"user_data_id" gorm:"size:36"`
	SaleId     string `json:"sale_id" gorm:"size:36"`
}
