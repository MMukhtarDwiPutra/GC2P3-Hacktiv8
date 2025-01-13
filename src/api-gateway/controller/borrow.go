package controller

import (
	"api-gateway/dto"
	"api-gateway/helpers"
	borrowpb "api-gateway/pb/borrowpb"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BorrowController struct {
	Client borrowpb.BorrowServiceClient
}

func NewBorrowController(client borrowpb.BorrowServiceClient) BorrowController {
	return BorrowController{
		Client: client,
	}
}

// @Summary     Borrow a book
// @Description Borrow a book using its ID
// @Tags        borrows
// @Accept      json
// @Produce     json
// @Param       request body dto.BorrowedBookRequest true "Borrow book details"
// @Success     200 {object} dto.SuccessResponse
// @Failure     400 {object} dto.ErrResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /borrows [post]
func (b BorrowController) BorrowABook(c echo.Context) error {
	var borrowRequest dto.BorrowedBookRequest

	// Bind request body
	if err := c.Bind(&borrowRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request parameters"})
	}

	// Validate request
	if err := validate.Struct(borrowRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid request parameters: %v", err)})
	}

	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Get the current date for BorrowedDate (assumes `helpers.GetCurrentDate` is implemented)
	borrowedDate := helpers.GetCurrentDate()

	// Extract user ID from JWT or session (this assumes `helpers.GetUserIDFromJWT` is implemented)
	user, err := helpers.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": fmt.Sprintf("Unauthorized: %v", err)})
	}

	log.Printf("User logged in: %v", user.UserID)

	// Map to gRPC request
	grpcRequest := &borrowpb.BorrowedBookRequest{
		BookId:       borrowRequest.BookID,
		UserId:       user.UserID,
		BorrowedDate: borrowedDate,
	}

	// Call gRPC service
	responseGrpc, err := b.Client.BorrowABook(ctx, grpcRequest)
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

// @Summary     Get all borrowed books
// @Description Retrieve all borrowed books for the current user
// @Tags        borrows
// @Accept      json
// @Produce     json
// @Success     200 {array} dto.BorrowResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /borrows [get]
func (b BorrowController) GetAllBorrows(c echo.Context) error {
	ctx, cancel, err := helpers.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	// Extract user ID from JWT or session (this assumes `helpers.GetUserIDFromJWT` is implemented)
	user, err := helpers.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": fmt.Sprintf("Unauthorized: %v", err)})
	}

	// Map to gRPC request
	grpcRequest := &borrowpb.UserId{
		Id: user.UserID,
	}

	// Call gRPC service
	responseGrpc, err := b.Client.GetAllBorrows(ctx, grpcRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error response: %v", err)})
	}

	data := helpers.UnmarshalJSONToWebResponse(responseGrpc.Data)

	return c.JSON(data.Status, data)
}

// @Summary     Return a borrowed book
// @Description Mark a borrowed book as returned
// @Tags        borrows
// @Accept      json
// @Produce     json
// @Param       id path string true "Borrow ID"
// @Success     200 {object} dto.SuccessResponse
// @Failure     500 {object} dto.ErrResponse
// @Router      /borrows/return [post]
func (b BorrowController) ReturnBook(c echo.Context) error {
	return c.JSON(200, nil)
}
