package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tanush-128/openzo_backend/user/config"
	handlers "github.com/tanush-128/openzo_backend/user/internal/api"
	"github.com/tanush-128/openzo_backend/user/internal/middlewares"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	userpb "github.com/tanush-128/openzo_backend/user/internal/pb"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/service"
)

type Server struct {
	userpb.UserServiceServer
	userRepository repository.UserRepository
	userService    service.UserService
}

// Define Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	userRepository := repository.NewUserRepository(db)

	otpRepository := repository.NewOTPRepository(db)

	userService := service.NewUserService(userRepository)

	otpService := service.NewOTPService(otpRepository, userRepository)

	conf := ReadConfig()
	p, _ := kafka.NewProducer(&conf)
	// topic := "notification"

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	go consumeKafka(userRepository, p)

	go service.GrpcServer(cfg, &service.Server{UserRepository: userRepository, UserService: userService})

	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&userService)
	otp_handler := handlers.NewOTPHandler(&otpService)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Define routes
	router.GET("ping", measureMetrics("ping", "GET", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong v3 trying to deploy",
		})
	}))
	router.POST("/", measureMetrics("/", "POST", handler.CreateUser))
	router.GET("/:id", measureMetrics("/:id", "GET", handler.GetUserByID))
	router.PUT("/", measureMetrics("/", "PUT", handler.UpdateUser))
	router.GET("/email/:email", measureMetrics("/email/:email", "GET", handler.GetUserByEmail))
	router.POST("/signin", measureMetrics("/signin", "POST", handler.UserSignIn))

	router.POST("/otp", measureMetrics("/otp", "POST", otp_handler.GenerateOTP))
	router.POST("/otp/verify", measureMetrics("/otp/verify", "POST", otp_handler.VerifyOTP))

	router.Use(middlewares.JwtMiddleware)
	router.GET("/jwt", measureMetrics("/jwt", "GET", handler.GetUserWithJWT))

	// Start server
	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}

func measureMetrics(path string, method string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(path, method))
		defer timer.ObserveDuration()

		// Increment counter
		httpRequestsTotal.WithLabelValues(path, method).Inc()

		// Call the handler
		handlerFunc(c)
	}
}

type Notification struct {
	Message  string `json:"message"`
	FCMToken string `json:"fcm_token"`
	Data     string `json:"data,omitempty"`
	Topic    string `json:"topic,omitempty"`
}

type data struct {
	// OrderStatus string `json:"order_status"`
	ID          string          `json:"id"`
	StoreID     string          `json:"store_id"`
	Customer    models.Customer `json:"customer"`
	OrderStatus string          `json:"status"`
	Type        string          `json:"type"`
}

func consumeKafka(userRepo repository.UserRepository, notificationProducer *kafka.Producer) {
	// Read the Kafka configuration
	conf := ReadConfig()
	topic := "sales"

	// Set the consumer group ID and offset
	conf["group.id"] = "UserGroup"
	conf["auto.offset.reset"] = "earliest"

	// Create a new consumer and subscribe to the topic
	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	var order data
	run := true

	for run {
		// Poll for new messages
		event := consumer.Poll(1000)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			log.Printf("Message on %s: %s\n", e.TopicPartition, string(e.Value))

			// Unmarshal the message into the order struct
			err := json.Unmarshal(e.Value, &order)
			if err != nil {
				log.Printf("Error unmarshalling JSON: %v", err)
				continue
			}

			log.Printf("Order received: %+v", order)

			// Validate the order data
			if order.Type != "online_order" || order.OrderStatus == "" || order.Customer.UserDataId == "" || order.OrderStatus == "not_placed" {
				continue
			}

			// Fetch the user data
			userData, err := userRepo.GetUserByID(order.Customer.UserDataId)
			if err != nil {
				log.Printf("Error getting user data: %v", err)
				continue
			}

			if userData.NotificationToken == nil {
				log.Printf("User does not have an FCM token")
				continue
			}

			// Construct the notification message
			var notificationMessage string
			switch order.OrderStatus {
			case "accepted":
				notificationMessage = "Your order has been accepted"
			case "cancelled":
				notificationMessage = "Your order has been cancelled"
			case "out_for_delivery":
				notificationMessage = "Your order is out for delivery"
			case "delivered":
				notificationMessage = "Your order has been delivered"
			case "rejected":
				notificationMessage = "Your order has been rejected"
			default:
				notificationMessage = "Your order has been placed"
			}

			notification := Notification{
				Message:  notificationMessage,
				FCMToken: *userData.NotificationToken,
				Data:     fmt.Sprintf(`{"order_id": "%s", "status": "%s"}`, order.ID, order.OrderStatus),
			}

			notificationMsg, err := json.Marshal(notification)
			if err != nil {
				log.Printf("Error marshalling notification: %v", err)
				continue
			}

			notificationTopic := "notification"

			// Send the notification
			err = notificationProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &notificationTopic, Partition: kafka.PartitionAny},
				Value:          notificationMsg,
			}, nil)
			if err != nil {
				log.Printf("Error producing notification: %v", err)
				continue
			}

			notificationProducer.Flush(15000)

		case kafka.Error:
			log.Printf("Kafka error: %v", e)
			run = false
		default:
			log.Printf("Ignored event: %v", e)
		}
	}

	// Closing the consumer connection
	err = consumer.Close()
	if err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}

// func consumeKafka(userRepo repository.UserRepository, notificationProducer *kafka.Producer) {
// 	conf := ReadConfig()
// 	topic := "sales"

// 	// sets the consumer group ID and offset
// 	conf["group.id"] = "UserGroup"
// 	conf["auto.offset.reset"] = "earliest"

// 	// creates a new consumer and subscribes to your topic
// 	consumer, _ := kafka.NewConsumer(&conf)
// 	consumer.SubscribeTopics([]string{topic}, nil)

// 	var order data

// 	run := true
// 	for run {
// 		// consumes messages from the subscribed topic and prints them to the console
// 		e := consumer.Poll(1000)
// 		switch ev := e.(type) {
// 		case *kafka.Message:
// 			// application-specific processing
// 			log.Printf("Message on %s: %s\n", ev.TopicPartition, string(ev.Value))
// 			err := json.Unmarshal(ev.Value, &order)
// 			if err != nil {
// 				fmt.Println("Error unmarshalling JSON: ", err)
// 				continue
// 			}

// 			fmt.Printf("Order received: %+v ", order)

// 			if order.Type != "online_order" || order.OrderStatus == "" || order.Customer.UserDataId == "" || order.OrderStatus == "not_placed" {

// 				continue
// 			}

// 			userData, err := userRepo.GetUserByID(order.Customer.UserDataId)
// 			fcm := userData.NotificationToken
// 			if err != nil || fcm == nil {
// 				fmt.Println("Error getting FCM token: ", err)
// 				continue
// 			}

// 			fmt.Println("FCM token: ", fcm)

// 			notificationMessage := ""

// 			if order.OrderStatus == "accepted" {
// 				notificationMessage = "Your order has been accepted"
// 			} else if order.OrderStatus == "cancelled" {
// 				notificationMessage = "Your order has been cancelled"
// 			} else if order.OrderStatus == "out_for_delivery" {
// 				notificationMessage = "Your order is out for delivery"
// 			} else if order.OrderStatus == "delivered" {
// 				notificationMessage = "Your order has been delivered"
// 			} else if order.OrderStatus == "rejected" {
// 				notificationMessage = "Your order has been rejected"
// 			} else {
// 				notificationMessage = "Your order has been placed"
// 			}

// 			notificationMsg, _ := json.Marshal(Notification{
// 				Message:  notificationMessage,
// 				FCMToken: *fcm,

// 				Data: fmt.Sprintf(`{"order_id": "%s", "status": "%s"}`, order.ID, order.OrderStatus),
// 			})

// 			notificationTopic := "notification"

// 			// send a notification to the store
// 			notificationProducer.Produce(&kafka.Message{
// 				TopicPartition: kafka.TopicPartition{Topic: &notificationTopic, Partition: kafka.PartitionAny},
// 				Value:          notificationMsg,
// 			}, nil)

// 			notificationProducer.Flush(15 * 1000)

// 		case kafka.Error:
// 			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
// 			run = false
// 		}
// 	}

// 	// closes the consumer connection
// 	consumer.Close()

// }
