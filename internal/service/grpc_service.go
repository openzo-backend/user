package service
import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tanush-128/openzo_backend/user/config"
	"github.com/tanush-128/openzo_backend/user/internal/middlewares"
	userpb "github.com/tanush-128/openzo_backend/user/internal/pb"

	// "github.com/tanush-128/openzo_backend/store/internal/pb"
	"github.com/tanush-128/openzo_backend/user/internal/repository"

	"google.golang.org/grpc"
)

type Server struct {
	userpb.UserServiceServer
	UserRepository repository.UserRepository
	UserService    UserService
}

func GrpcServer(
	cfg *config.Config,
	server *Server,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())
	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, server)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *Server) GetUserWithJWT(ctx context.Context, req *userpb.Token) (*userpb.User, error) {
	claims, err := middlewares.ValidateJwtToken(req.Token)
	if err != nil {
		return nil, err
	}
	user, err := s.UserRepository.GetUserByID(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	role := userpb.Role_USER
	if user.Role == "ADMIN" {
		role = userpb.Role_ADMIN
	}

	return &userpb.User{
		Id:         user.ID,
		Phone:      user.Phone,
		IsVerified: user.IsVerified,
		Role:       role,
	}, nil
}
// this is Tanush Agarwal from openzo backend