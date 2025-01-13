package controller

import (
	"borrow-service/helpers"
	"borrow-service/models"           // Ensure this is the correct import path
	"borrow-service/pb"               // Ensure this is the correct import path
	bookpb "borrow-service/pb/bookpb" // Ensure this is the correct import path
	"borrow-service/service"          // Ensure this is the correct import path
	"log"
	"os"

	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

var validate = validator.New()

type borrowController struct {
	pb.UnimplementedBorrowServiceServer // Correctly reference UnimplementedBorrowServiceServer
	borrowService                       service.BorrowService
	bookServiceClient                   bookpb.BookServiceClient
}

func NewBorrowController(borrowService service.BorrowService) borrowController {
	conn, err := grpc.Dial(os.Getenv("BOOK_SERVICE_URI"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to book-service: %v", err)
	}

	bookClient := bookpb.NewBookServiceClient(conn)

	return borrowController{
		borrowService:     borrowService,
		bookServiceClient: bookClient, // gRPC client for book-service
	}
}

func (h borrowController) BorrowABook(ctx context.Context, req *pb.BorrowedBookRequest) (*pb.WebResponse, error) {
	// Parse the string date to time.Time
	borrowedDate, err := time.Parse("2006-01-02", req.GetBorrowedDate()) // Adjust format to match your input
	if err != nil {
		return nil, fmt.Errorf("Error parsing Borrowed Date: %v", err)
	}

	// Construct the borrow request model
	borrowRequest := models.BorrowedBook{
		BookID:       req.GetBookId(),
		UserID:       req.GetUserId(),
		BorrowedDate: borrowedDate,
		ReturnDate:   nil,
	}

	// Validate the request
	err = validate.Struct(borrowRequest)
	if err != nil {
		return nil, fmt.Errorf("Invalid request parameters for borrowRequest: %v", err)
	}

	ctxBook, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return nil, fmt.Errorf("Error making book context: %v", err)
	}
	defer cancel()

	// Call gRPC service to get book details
	responseGrpc, err := h.bookServiceClient.GetBookById(ctxBook, &bookpb.BookId{Id: req.GetBookId()})
	if err != nil {
		return nil, fmt.Errorf("Error calling book-service GetBookById: %v", err)
	}

	// Parse the gRPC response data
	webResponseData := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponseData)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal webResponse: %v", err)
	}

	if webResponseData.Status != 200 {
		return &pb.WebResponse{
			Status: fmt.Sprintf("%v", webResponseData.Status),
			Data:   string(webResponseJSON), // Convert to string
		}, nil
	}

	// Call your service to register the borrow
	status, webResponse := h.borrowService.CreateBorrow(borrowRequest)

	// Convert status to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err = json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal webResponse: %v", err)
	}

	// Return the mapped WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

func (h borrowController) ReturnBook(ctx context.Context, req *pb.UpdateBorrowedBook) (*pb.WebResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error())
	}

	// Call your service to register the borrow
	status, webResponse := h.borrowService.GetBorrowById(id)
	// Convert webResponse (map) to JSON string
	webResponseJSON, err := json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webResponse: %v", err)
	}

	if status != 200 {
		// Map the response to WebResponse
		return &pb.WebResponse{
			Status: fmt.Sprintf("%v", status),
			Data:   string(webResponseJSON), // Convert to string
		}, nil
	}

	// Parse the string date to time.Time
	returnDate, err := time.Parse("2006-01-02", req.GetReturnDate()) // Adjust format to match your input
	if err != nil {
		return nil, fmt.Errorf("Error parsing Borrowed Date: %v", err)
	}

	// Construct the borrow request model
	borrowRequest := models.UpdateBorrowedBook{
		ReturnDate: &returnDate,
	}

	status, webResponse = h.borrowService.ReturnBook(id, borrowRequest)

	// Convert status to string
	statusStr := fmt.Sprintf("%d", status)

	// Convert webResponse (map) to JSON string
	webResponseJSON, err = json.Marshal(webResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal webResponse: %v", err)
	}

	// Map the response to WebResponse
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(webResponseJSON), // Convert to string
	}, nil
}

func (h borrowController) GetAllBorrows(ctx context.Context, param *pb.UserId) (*pb.WebResponse, error) {
	// Validate and convert the userId from the request to primitive.ObjectID
	userId, err := primitive.ObjectIDFromHex(param.GetId())

	if err != nil {
		return nil, fmt.Errorf("invalid userId: %v", err)
	}

	// Call the service to retrieve all borrows by user
	status, borrows := h.borrowService.GetAllBorrowsByUser(userId)

	// Convert status (int) to string
	statusStr := fmt.Sprintf("%d", status)

	// Marshal the borrows data into JSON
	borrowsJSON, err := json.Marshal(borrows)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal borrows data: %v", err)
	}

	// Return the response
	return &pb.WebResponse{
		Status: statusStr,
		Data:   string(borrowsJSON),
	}, nil
}
