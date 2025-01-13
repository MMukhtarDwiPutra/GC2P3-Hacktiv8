package routes

import (
	"api-gateway/controller"
	_ "api-gateway/docs" // Ensure this is enabled for Swagger
	"api-gateway/middlewares"

	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Template represents the template renderer for Echo
type Template struct {
	templates *template.Template
}

// NewRouter initializes the routes for the application
func NewRouter(
	userController controller.UserController,
	bookController controller.BookController,
	borrowController controller.BorrowController,
) *echo.Echo {
	// Create a new Echo instance
	e := echo.New()

	// Middleware setup
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)),
	)

	// Swagger documentation route
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// User routes
	e.POST("/users/register", userController.Register)
	e.POST("/users/login", userController.Login)

	// Book routes
	e.POST("/books", bookController.Create, middlewares.RequireAuth)
	e.GET("/books", bookController.GetAllBooks, middlewares.RequireAuth)
	e.PUT("/books/:id", bookController.UpdateBookById, middlewares.RequireAuth)
	e.GET("/books/:id", bookController.GetBookById, middlewares.RequireAuth)
	e.DELETE("/books/:id", bookController.DeleteBookById, middlewares.RequireAuth)

	// Borrow routes
	e.POST("/borrows", borrowController.BorrowABook, middlewares.RequireAuth)
	e.POST("/borrows/return", borrowController.ReturnBook, middlewares.RequireAuth)
	e.GET("/borrows", borrowController.GetAllBorrows, middlewares.RequireAuth)

	return e
}
