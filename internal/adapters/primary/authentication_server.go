package primary

import (
	"context"
	"errors"
	"log"

	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	AuthService *services.AuthService
}

func (s *AuthServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	userID, email, token, err := s.AuthService.RegisterUser(req.Email, req.Password)
	if err != nil {
		log.Println("Registration error:", err)
		return nil, errors.New("registration failed")
	}

	return &proto.RegisterResponse{
		UserId: userID,
		Email:  email,
		Token:  token,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	userID, email, token, err := s.AuthService.LoginUser(req.Email, req.Password)
	if err != nil {
		log.Println("Login error:", err)
		return nil, errors.New("login failed")
	}

	return &proto.LoginResponse{
		UserId: userID,
		Email:  email,
		Token:  token,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	err := s.AuthService.LogoutUser(req.Token)
	if err != nil {
		log.Println("Logout error:", err)
		return nil, errors.New("logout failed")
	}

	return &proto.LogoutResponse{
		Message: "Logout successful",
	}, nil
}
