package controller

import (
	"api-gateway/dto"
	"api-gateway/helpers"
	bookpb "api-gateway/pb/bookpb"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookController struct {
	Client bookpb.BookServiceClient
}

func NewBookController(client bookpb.BookServiceClient) BookController {
	return BookController{
		Client: client,
	}
}

// @Summary     Register a new book
// @Description Register a new book with the role 'Book'
// @Tags        customer
// @Accept      json
// @Produce     json
// @Param       request body dto.BookRegister true "Book registration details"
// @Success     201 {object} dto.SwaggerResponseRegister
// @Failure     400 {object} utils.ErrResponse
// @Failure     409 {object} utils.ErrResponse
// @Failure     500 {object} utils.ErrResponse
// @Router      /books/register/book [post]
func (b BookController) CreateBook(c echo.Context) error {
	return b.Create(c)
}

func (b BookController) Create(c echo.Context) error {
	var bookRequest dto.BookRequest

	// Bind request body
	if err := c.Bind(&bookRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	// Validate request
	if err := validate.Struct(bookRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid request parameters: %v", err)})
	}

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Map to gRPC request
	grpcRequest := &bookpb.BookRequest{
		Title:         bookRequest.Title,
		Author:        bookRequest.Author,
		PublishedDate: bookRequest.PublishedDate,
		Status:        bookRequest.Status,
	}

	// Call gRPC service
	responseGrpc, err := b.Client.CreateBook(ctx, grpcRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToMap(responseGrpc.Data)

	return c.JSON(http.StatusCreated, data)
}

func (b BookController) GetBookById(c echo.Context) error {
	bookId := c.Param("id")

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Call gRPC service
	responseGrpc, err := b.Client.GetBookById(ctx, &bookpb.BookId{Id: bookId})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToMap(responseGrpc.Data)

	return c.JSON(http.StatusOK, data)
}

func (b BookController) GetAllBooks(c echo.Context) error {
	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Call gRPC service
	responseGrpc, err := b.Client.GetAllBooks(ctx, &empty.Empty{})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToMap(responseGrpc.Data)

	return c.JSON(http.StatusOK, data)
}

func (b BookController) UpdateBookById(c echo.Context) error {
	bookId := c.Param("id")
	var bookRequest dto.BookRequest

	// Bind request body
	if err := c.Bind(&bookRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	// Validate request
	if err := validate.Struct(bookRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid request parameters: %v", err)})
	}

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Map to gRPC request
	grpcRequest := &bookpb.UpdateBookRequest{
		Title:         bookRequest.Title,
		Author:        bookRequest.Author,
		PublishedDate: bookRequest.PublishedDate,
		Status:        bookRequest.Status,
		BookId:        bookId,
	}

	// Call gRPC service
	responseGrpc, err := b.Client.UpdateBookById(ctx, grpcRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToMap(responseGrpc.Data)

	return c.JSON(http.StatusOK, data)
}

func (b BookController) DeleteBookById(c echo.Context) error {
	bookId := c.Param("id")

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Call gRPC service
	responseGrpc, err := b.Client.DeleteBookById(ctx, &bookpb.BookId{Id: bookId})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return c.JSON(http.StatusConflict, map[string]string{"message": "Book already exists"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToMap(responseGrpc.Data)

	return c.JSON(http.StatusOK, data)
}
