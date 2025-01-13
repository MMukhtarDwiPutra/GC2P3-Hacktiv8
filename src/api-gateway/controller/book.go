package controller

import (
	"api-gateway/dto"
	"api-gateway/helpers"
	bookpb "api-gateway/pb/bookpb"
	"api-gateway/pb/borrowpb"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookController struct {
	Client       bookpb.BookServiceClient
	BorrowClient borrowpb.BorrowServiceClient
}

func NewBookController(client bookpb.BookServiceClient) BookController {
	conn, err := grpc.Dial(os.Getenv("BORROW_SERVICE_URI"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to book-service: %v", err)
	}

	borrowClient := borrowpb.NewBorrowServiceClient(conn)
	return BookController{
		Client:       client,
		BorrowClient: borrowClient,
	}
}

// @Summary     Create a new book
// @Description Register a new book with the role 'Book'
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       request body dto.BookRequest true "Book details"
// @Success     201 {object} dto.WebResponse
// @Failure     400 {object} dto.ErrResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /books [post]
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

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

// @Summary     Get a book by ID
// @Description Retrieve a book by its unique ID
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       id path string true "Book ID"
// @Success     200 {object} dto.BookResponse
// @Failure     404 {object} dto.ErrResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /books/{id} [get]
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

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

// @Summary     List all books
// @Description Retrieve all books
// @Tags        books
// @Accept      json
// @Produce     json
// @Success     200 {array} dto.BookResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /books [get]
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

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

// @Summary     Update a book
// @Description Update a book's information by ID
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       id path string true "Book ID"
// @Param       request body dto.BookRequest true "Updated book details"
// @Success     200 {object} dto.BookResponse
// @Failure     400 {object} dto.ErrResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /books/{id} [put]
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

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

// @Summary     Delete a book
// @Description Delete a book by its unique ID
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       id path string true "Book ID"
// @Success     200 {object} dto.SuccessResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /books/{id} [delete]
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
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

func (b BookController) UpdateStatusNotReturnedYet() {
	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		log.Printf("Error making book context: %v", err)
		return
	}
	defer cancel()

	emptyRequest := &borrowpb.Empty{}
	responseGrpc, err := b.BorrowClient.GetAllNotReturnedBook(ctx, emptyRequest)
	if err != nil {
		log.Printf("Error getting all not returned books: %v", err)
		return
	}

	// Parse the gRPC response data
	webResponseData := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	var books []dto.BorrowedBook

	// If data is in map form, manually populate the fields of books
	if data, ok := webResponseData.Data["data"].(map[string]interface{}); ok {
		books, ok = data["data"].([]dto.BorrowedBook)
		if !ok {
			log.Println("Failed to assert books to []dto.BorrowedBook")
			return
		}
	}

	for _, book := range books {
		// Map to gRPC request
		grpcRequest := &bookpb.UpdateStatusBookRequest{
			BookId: book.BookID,
			Status: "",
		}

		// Call gRPC service to update the status
		_, err := b.Client.UpdateStatusBookById(ctx, grpcRequest)
		if err != nil {
			log.Printf("Failed to update book status for book ID %v: %v", book.BookID, err)
		}
	}

	log.Printf("%d books have been updated to null status", len(books))
}
