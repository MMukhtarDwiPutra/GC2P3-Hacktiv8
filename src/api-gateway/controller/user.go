package controller

import (
	"api-gateway/dto"
	"api-gateway/helpers"
	userpb "api-gateway/pb/userpb"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validator global yang digunakan untuk memvalidasi struktur request body.
var validate = validator.New()

type UserController struct {
	Client userpb.UserServiceClient
}

func NewUserController(client userpb.UserServiceClient) UserController {
	return UserController{
		Client: client,
	}
}

// @Summary     Register a new user
// @Description Register a user with username and password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body dto.RegisterRequest true "User registration details"
// @Success     201 {object} dto.RegisterResponse
// @Failure     400 {object} dto.ErrResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /users/register [post]
func (u UserController) Register(c echo.Context) error {
	var userRequest dto.RegisterRequest

	// Bind request body
	if err := c.Bind(&userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	// Validate request
	if err := validate.Struct(userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid request parameters: %v", err)})
	}

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Map to gRPC request
	grpcRequest := &userpb.RegisterRequest{
		Username: userRequest.Username,
		Password: userRequest.Password,
	}

	// Call gRPC service
	responseGrpc, err := u.Client.Register(ctx, grpcRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return c.JSON(http.StatusConflict, map[string]string{"message": "User already exists"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

// @Summary     User login
// @Description Authenticate a user with username and password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body dto.LoginRequest true "Login details"
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} dto.ErrResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /users/login [post]
func (u UserController) Login(c echo.Context) error {
	var userRequest dto.LoginRequest

	// Bind request body
	if err := c.Bind(&userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	// Validate request
	if err := validate.Struct(userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid request parameters: %v", err)})
	}

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Map to gRPC request
	grpcRequest := &userpb.LoginRequest{
		Username: userRequest.Username,
		Password: userRequest.Password,
	}

	// Call gRPC service
	responseGrpc, err := u.Client.Login(ctx, grpcRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	log.Println(data)

	return c.JSON(data.Status, data)
}
