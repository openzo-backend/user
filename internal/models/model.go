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

type OrderStatus string

const (
	OrderPlaced    OrderStatus = "placed"
	OrderAccepted  OrderStatus = "accepted"
	OrderRejected  OrderStatus = "rejected"
	OrderOutForDel OrderStatus = "out_for_delivery"
	OrderCancelled OrderStatus = "cancelled"
	OrderDelivered OrderStatus = "delivered"
)

type OnlineOrder struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	OrderItems  []OnlineOrderItem `json:"order_items"`
	StoreID     string            `json:"store_id"`
	Customer    OnlineCustomer    `json:"customer"`
	OrderTime   string            `json:"order_time"`
	OrderStatus OrderStatus       `json:"order_status"`
	TotalAmount float64           `json:"total_amount"`
}

type OnlineOrderItem struct {
	ID            int    `json:"id"`
	ProductID     string `json:"product_id"`
	OnlineOrderId string `json:"sale_id"`
	Quantity      int    `json:"quantity"`
}

type OnlineCustomer struct {
	ID            int    `json:"id"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	UserDataId    string `json:"user_data_id"`
	OnlineOrderID string `json:"sale_id"`
}
