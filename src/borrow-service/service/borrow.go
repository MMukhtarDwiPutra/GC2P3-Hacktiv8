package service

import (
	"borrow-service/models"
	"borrow-service/repository"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BorrowService interface {
	CreateBorrow(request models.BorrowedBook) (int, map[string]interface{})
	GetBorrowById(id primitive.ObjectID) (int, map[string]interface{})
	GetAllBorrowsByUser(id primitive.ObjectID) (int, map[string]interface{})
	UpdateBorrowById(id primitive.ObjectID, borrow models.BorrowedBookRequest) (int, map[string]interface{})
	ReturnBook(id primitive.ObjectID, borrow models.UpdateBorrowedBook) (int, map[string]interface{})
	GetAllNotReturnedBook() (int, map[string]interface{})
}

type borrowService struct {
	borrowRepository repository.BorrowRepository
}

func NewBorrowService(borrowRepository repository.BorrowRepository) *borrowService {
	return &borrowService{borrowRepository}
}

func (s *borrowService) CreateBorrow(request models.BorrowedBook) (int, map[string]interface{}) {
	borrow := models.BorrowedBook{
		BookID:       request.BookID,
		UserID:       request.UserID,
		BorrowedDate: request.BorrowedDate,
		ReturnDate:   request.ReturnDate,
	}

	borrowResult, err := s.borrowRepository.CreateBorrow(borrow)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	borrow.ID = *borrowResult

	return http.StatusCreated, map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "Borrow response created successfully",
		"data":    borrow,
	}
}

func (s *borrowService) ReturnBook(id primitive.ObjectID, borrow models.UpdateBorrowedBook) (int, map[string]interface{}) {
	_, err := s.borrowRepository.UpdateReturnBorrow(id, borrow)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error database: %v", err),
		}
	}

	borrowUpdated := models.UpdateBorrowedBook{
		ReturnDate: borrow.ReturnDate,
	}

	webResponse := map[string]interface{}{
		"status": http.StatusOK,
		"data":   borrowUpdated,
	}

	return http.StatusOK, webResponse
}

func (s *borrowService) GetBorrowById(id primitive.ObjectID) (int, map[string]interface{}) {
	borrowResult, err := s.borrowRepository.GetBorrowById(id)
	if borrowResult == nil {
		return http.StatusNotFound, map[string]interface{}{
			"message": "Borrow ID not found!",
		}
	}
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	borrowResult.ID = id

	return http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Get Borrow By Id Success",
		"data":    borrowResult,
	}
}

func (s *borrowService) GetAllBorrowsByUser(id primitive.ObjectID) (int, map[string]interface{}) {
	borrowResult, err := s.borrowRepository.GetAllBorrowsByUser(id)
	log.Println(borrowResult)
	if borrowResult == nil {
		return http.StatusOK, map[string]interface{}{
			"message": "Belum ada data peminjaman",
		}
	}
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Get All Borrows Success",
		"data":    borrowResult,
	}
}

func (s *borrowService) UpdateBorrowById(id primitive.ObjectID, borrow models.BorrowedBookRequest) (int, map[string]interface{}) {
	_, err := s.borrowRepository.UpdateBorrow(id, borrow)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error database: %v", err),
		}
	}

	borrowUpdated := models.BorrowedBook{
		ID:           id,
		UserID:       borrow.UserID,
		BorrowedDate: borrow.BorrowedDate,
		ReturnDate:   borrow.ReturnDate,
	}

	webResponse := map[string]interface{}{
		"status": http.StatusOK,
		"data":   borrowUpdated,
	}

	return http.StatusOK, webResponse
}

func (s *borrowService) DeleteBorrowById(id primitive.ObjectID) (int, map[string]interface{}) {
	_, err := s.borrowRepository.DeleteBorrowById(id)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error database: %v", err),
		}
	}

	webResponse := map[string]interface{}{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Borrow with ID %v has been deleted.", id),
	}

	return http.StatusOK, webResponse
}

func (s *borrowService) GetAllNotReturnedBook() (int, map[string]interface{}) {
	borrowResult, err := s.borrowRepository.GetAllNotReturnedBook()

	if borrowResult == nil {
		return http.StatusOK, map[string]interface{}{
			"message": "Belum ada data peminjaman",
		}
	}
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Get All Not Returned Book Success",
		"data":    borrowResult,
	}
}
