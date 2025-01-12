package routes

import (
	"api-gateway/controller"
	"api-gateway/helpers"
	"api-gateway/middlewares"

	// _ "api-gateway/docs"

	"html/template"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Template struct {
	templates *template.Template
}

func NewRouter(userController controller.UserController, bookController controller.BookController) *echo.Echo {
	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger(), middleware.Recover(), middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Swagger
	e.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// User Routes
	e.POST("/users/register", userController.RegisterUser)
	e.POST("/users/login", userController.LoginUser)

	e.POST("/books", bookController.CreateBook, middlewares.RequireAuth)
	e.GET("/books", bookController.GetAllBooks, middlewares.RequireAuth)
	e.PUT("/books/:id", bookController.UpdateBookById, middlewares.RequireAuth)
	e.GET("/books/:id", bookController.GetBookById, middlewares.RequireAuth)
	e.DELETE("/books/:id", bookController.DeleteBookById, middlewares.RequireAuth)

	return e
}
