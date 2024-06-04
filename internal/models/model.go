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

type UserData struct {
	Id                string `json:"id" gorm:"primaryKey"`
	UserId            string `json:"user_id"`
	Latitude          string `json:"latitude"`
	Longitude         string `json:"longitude"`
	Address           string `json:"address"`
	Pincode           string `json:"pincode"`
	City              string `json:"city"`
	State             string `json:"state"`
	Country           string `json:"country"`
	NotificationToken string `json:"notification_token"`
}

type Customer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	UserDataId string `json:"user_data_id" gorm:"size:36"`
	SaleId     string `json:"sale_id" gorm:"size:36"`
}
