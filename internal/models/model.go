package models

import "time"

type OTP struct {
	ID        string `gorm:"primaryKey"`
	Phone     string
	HashedOTP string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

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
	Role              string `json:"role" gorm:"default:'USER'"`
	DefaultAddress    string `json:"default_address"`
}

type Address struct {
	ID             string `gorm:"primaryKey" json:"id"`
	UserId         string `json:"user_id"`
	Name           string `json:"name"`
	PhoneNo        string `json:"phone_no"`
	Tag            string `json:"tag"`
	Area           string `json:"area"`
	Building       string `json:"building"`
	NearbyLandmark string `json:"nearby_landmark"`
	Address        string `json:"address"`
	Pincode        string `json:"pincode"`
	City           string `json:"city"`
	State          string `json:"state"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
}

type Customer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	UserDataId string `json:"user_data_id" gorm:"size:36"`
	SaleId     string `json:"sale_id" gorm:"size:36"`
}
