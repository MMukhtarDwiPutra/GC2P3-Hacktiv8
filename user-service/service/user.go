package service

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"user-service/models"
	"user-service/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(request models.RegisterRequest) (int, map[string]interface{})
	LoginUser(request models.LoginRequest) (int, map[string]interface{})
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository}
}

func (s *userService) RegisterUser(request models.RegisterRequest) (int, map[string]interface{}) {
	findUser, err := s.userRepository.GetUserByUsername(request.Username)
	if findUser != nil {
		return http.StatusConflict, map[string]interface{}{"message": "Tidak berhasil! Email sudah terdaftar!"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	user := models.User{
		Username: request.Username,
		Password: string(hashedPassword),
	}

	userResult, err := s.userRepository.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	userResponse := models.UserResponse{
		UserID:   *userResult,
		Username: request.Username,
	}

	return http.StatusCreated, map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "User response created successfully",
		"data":    userResponse,
	}
}

func (s *userService) LoginUser(request models.LoginRequest) (int, map[string]interface{}) {
	secretKey := os.Getenv("LOGIN_SECRET_KEY")

	// Mengecek apakah username siswa ada di database
	user, err := s.userRepository.GetUserByUsername(request.Username)
	if user == nil {
		return http.StatusNotFound, map[string]interface{}{"message": "Invalid Email or Password"}
	}
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return http.StatusNotFound, map[string]interface{}{"message": "Invalid Email or Password"}
		}

		return http.StatusInternalServerError, map[string]interface{}{
			"message": "Error querying the database.",
		}
	}

	// Memverifikasi password siswa
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"status":  http.StatusNotFound,
			"message": "Invalid Email or Password",
		}
	}

	// Membuat token JWT untuk siswa yang berhasil login
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// Menandatangani token dengan secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "error signed jwt token",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Login successfuly",
		"token":   tokenString,
	}
}
