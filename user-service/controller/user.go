package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"user-service/models"
	"user-service/pb/userpb" // Ensure this is the correct import path
	"user-service/service"

	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate = validator.New()

type userController struct {
	userpb.UnimplementedUserServiceServer // Correctly reference UnimplementedUserServiceServer
	userService                           service.UserService
}

func NewUserController(userService service.UserService) userController {
	return userController{
		userService: userService,
	}
}

func (h userController) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.WebResponse, error) {
	// Convert the gRPC RegisterRequest to your internal RegisterRequest model
	userRequest := models.RegisterRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	// Validate request
	err := validate.Struct(userRequest)
	if err != nil {
		return nil, fmt.Errorf("invalid request parameters: %v", err)
	}

	// Call your service to register the user
	status, webResponse := h.userService.RegisterUser(userRequest)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &userpb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

// LoginUser godoc
// @Summary User login
// @Description Authenticates a user by email and password, and returns a JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string "User successfully logged in"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func (h userController) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.WebResponse, error) {
	// Map the gRPC request to your `models.LoginRequest`
	loginRequest := models.LoginRequest{
		Username: req.Username, // Assuming `Username` exists in your proto definition
		Password: req.Password, // Assuming `Password` exists in your proto definition
	}

	// Validate the login request
	if err := validate.Struct(loginRequest); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Call your service to register the user
	status, webResponse := h.userService.LoginUser(loginRequest)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &userpb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}
