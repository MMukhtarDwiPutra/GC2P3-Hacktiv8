package controller

import (
	"book-service/models"  // Ensure this is the correct import path
	"book-service/pb"      // Ensure this is the correct import path
	"book-service/service" // Ensure this is the correct import path
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

var validate = validator.New()

type bookController struct {
	pb.UnimplementedBookServiceServer // Correctly reference UnimplementedBookServiceServer
	bookService                       service.BookService
}

func NewBookController(bookService service.BookService) bookController {
	return bookController{
		bookService: bookService,
	}
}

func (h bookController) CreateBook(ctx context.Context, req *pb.BookRequest) (*pb.WebResponse, error) {
	// Parse the string date to time.Time
	publishedDate, err := time.Parse("2006-01-02", req.GetPublishedDate()) // Sesuaikan format sesuai input tanggal Anda

	if err != nil {
		return nil, fmt.Errorf("Error parsing PublishedDate: %v", err)
	}

	// Create the BookRequest struct
	bookRequest := models.BookRequest{
		Title:         req.GetTitle(),
		Author:        req.GetAuthor(),
		PublishedDate: publishedDate, // Use the parsed time.Time value
		Status:        req.GetStatus(),
		UserID:        "Belum ada peminjaman",
	}

	// Validate request
	err = validate.Struct(bookRequest)
	if err != nil {
		return nil, fmt.Errorf("invalid request parameters bookRequest: %v", err)
	}

	// Call your service to register the book
	status, webResponse := h.bookService.CreateBook(bookRequest)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

func (h bookController) GetBookById(ctx context.Context, req *pb.BookId) (*pb.WebResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error())
	}

	// Call your service to register the book
	status, webResponse := h.bookService.GetBookById(id)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

func (h bookController) GetAllBooks(ctx context.Context, param *emptypb.Empty) (*pb.WebResponse, error) {
	// Call your service to register the book
	status, webResponse := h.bookService.GetAllBooks()

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

func (h bookController) UpdateBookById(ctx context.Context, req *pb.UpdateBookRequest) (*pb.WebResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetBookId())
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error())
	}

	// Call your service to register the book
	status, webResponse := h.bookService.GetBookById(id)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	if status != 200 {
		// Map the response to WebResponse
		return &pb.WebResponse{
			Status: fmt.Sprintf("%d", status),
			Data:   string(webResponseJSON), // Convert to string
		}, nil
	}

	// Parse the string date to time.Time
	publishedDate, err := time.Parse("2006-01-02", req.GetPublishedDate()) // Sesuaikan format sesuai input tanggal Anda

	if err != nil {
		return nil, fmt.Errorf("Error parsing PublishedDate: %v", err)
	}

	// Create the BookRequest struct
	bookRequest := models.BookRequest{
		Title:         req.GetTitle(),
		Author:        req.GetAuthor(),
		PublishedDate: publishedDate, // Use the parsed time.Time value
		Status:        req.GetStatus(),
		UserID:        req.GetUserId(),
	}

	// Call your service to register the book
	status, webResponse = h.bookService.UpdateBookById(id, bookRequest)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err = json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

func (h bookController) DeleteBookById(ctx context.Context, req *pb.BookId) (*pb.WebResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error())
	}

	// Call your service to register the book
	status, webResponse := h.bookService.GetBookById(id)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	if status != 200 {
		// Map the response to WebResponse
		return &pb.WebResponse{
			Status: fmt.Sprintf("%d", status),
			Data:   string(webResponseJSON), // Convert to string
		}, nil
	}

	// Call your service to register the book
	status, webResponse = h.bookService.DeleteBookById(id)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err = json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}
