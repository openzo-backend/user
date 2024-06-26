package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	userDataRepository := repository.NewUserDataRepository(db)
	otpRepository := repository.NewOTPRepository(db)

	userService := service.NewUserService(userRepository)
	userDataService := service.NewUserDataService(userDataRepository)
	otpService := service.NewOTPService(otpRepository)

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

	go consumeKafka(userDataRepository, p)

	go service.GrpcServer(cfg, &service.Server{UserRepository: userRepository, UserService: userService})
	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&userService, &userDataService)
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
	router.GET("/email/:email", measureMetrics("/email/:email", "GET", handler.GetUserByEmail))
	router.PUT("/", measureMetrics("/", "PUT", handler.UpdateUser))
	router.POST("/signin", measureMetrics("/signin", "POST", handler.UserSignIn))
	router.POST("/userdata", measureMetrics("/userdata", "POST", handler.CreateUserData))
	router.GET("/userdata/:id", measureMetrics("/userdata/:id", "GET", handler.GetUserDataByID))
	router.GET("/userdata/user/:id", measureMetrics("/userdata/user/:id", "GET", handler.GetUserDataByUserID))
	router.PUT("/userdata", measureMetrics("/userdata", "PUT", handler.UpdateUserData))
	router.DELETE("/userdata/:id", measureMetrics("/userdata/:id", "DELETE", handler.DeleteUserData))
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
}

func consumeKafka(userDataRepo repository.UserDataRepository, notificationProducer *kafka.Producer) {
	conf := ReadConfig()
	topic := "sales"

	// sets the consumer group ID and offset
	conf["group.id"] = "UserGroup"
	conf["auto.offset.reset"] = "latest"

	// creates a new consumer and subscribes to your topic
	consumer, _ := kafka.NewConsumer(&conf)
	consumer.SubscribeTopics([]string{topic}, nil)

	var order data

	run := true
	for run {
		// consumes messages from the subscribed topic and prints them to the console
		e := consumer.Poll(1000)
		switch ev := e.(type) {
		case *kafka.Message:
			// application-specific processing
			log.Printf("Message on %s: %s\n", ev.TopicPartition, string(ev.Value))
			err := json.Unmarshal(ev.Value, &order)
			if err != nil {
				fmt.Println("Error unmarshalling JSON: ", err)
			}
			fmt.Printf("Order received: %+v ", order)

			userData, err := userDataRepo.GetUserDataByID(order.Customer.UserDataId)
			if err != nil {
				fmt.Println("Error getting FCM token: ", err)
			}
			fcm := userData.NotificationToken

			fmt.Println("FCM token: ", fcm)

			notificationMessage := ""

			if order.OrderStatus == "accepted" {
				notificationMessage = "Your order has been accepted"
			} else if order.OrderStatus == "cancelled" {
				notificationMessage = "Your order has been cancelled"
			} else if order.OrderStatus == "out_for_delivery" {
				notificationMessage = "Your order is out for delivery"
			} else if order.OrderStatus == "delivered" {
				notificationMessage = "Your order has been delivered"
			} else if order.OrderStatus == "rejected" {
				notificationMessage = "Your order has been rejected"
			} else {
				notificationMessage = "Your order has been placed"
			}

			notificationMsg, _ := json.Marshal(Notification{
				Message:  notificationMessage,
				FCMToken: fcm,
				Data:     fmt.Sprintf(`{"order_id": "%s", "status": "%s"}`, order.ID, order.OrderStatus),
			})

			notificationTopic := "notification"

			// send a notification to the store
			notificationProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &notificationTopic, Partition: kafka.PartitionAny},
				Value:          notificationMsg,
			}, nil)

			notificationProducer.Flush(15 * 1000)

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
			run = false
		}
	}

	// closes the consumer connection
	consumer.Close()

}
