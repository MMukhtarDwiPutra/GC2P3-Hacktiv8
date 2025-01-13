package service

import (
	"book-service/models"
	"book-service/repository"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService interface {
	CreateBook(request models.BookRequest) (int, map[string]interface{})
	GetBookById(id primitive.ObjectID) (int, map[string]interface{})
	GetAllBooks() (int, map[string]interface{})
	UpdateBookById(id primitive.ObjectID, book models.BookRequest) (int, map[string]interface{})
	DeleteBookById(id primitive.ObjectID) (int, map[string]interface{})
	UpdateStatusBookById(id primitive.ObjectID, book models.UpdateStatusBookRequest) (int, map[string]interface{})
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) *bookService {
	return &bookService{bookRepository}
}

func (s *bookService) CreateBook(request models.BookRequest) (int, map[string]interface{}) {
	book := models.Book{
		Title:         request.Title,
		Author:        request.Author,
		PublishedDate: request.PublishedDate,
		Status:        request.Status,
		UserID:        "Belum ada peminjaman",
	}

	bookResult, err := s.bookRepository.CreateBook(book)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	book.BookID = *bookResult

	return http.StatusCreated, map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "Book response created successfully",
		"data":    book,
	}
}

func (s *bookService) GetBookById(id primitive.ObjectID) (int, map[string]interface{}) {
	bookResult, err := s.bookRepository.GetBookById(id)
	if bookResult == nil {
		return http.StatusNotFound, map[string]interface{}{
			"status":  http.StatusNotFound,
			"message": "Book ID not found!",
		}
	}

	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	bookResult.BookID = id

	return http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Get Book By Id Success",
		"data":    bookResult,
	}
}

func (s *bookService) GetAllBooks() (int, map[string]interface{}) {
	bookResult, err := s.bookRepository.GetAllBooks()
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Get All Books Success",
		"data":    bookResult,
	}
}

func (s *bookService) UpdateBookById(id primitive.ObjectID, book models.BookRequest) (int, map[string]interface{}) {
	_, err := s.bookRepository.UpdateBook(id, book)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error database: %v", err),
		}
	}

	bookUpdated := models.Book{
		BookID:        id,
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: book.PublishedDate,
		Status:        book.Status,
		UserID:        book.UserID,
	}

	webResponse := map[string]interface{}{
		"status": http.StatusOK,
		"data":   bookUpdated,
	}

	return http.StatusOK, webResponse
}

func (s *bookService) DeleteBookById(id primitive.ObjectID) (int, map[string]interface{}) {
	_, err := s.bookRepository.DeleteBookById(id)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error database: %v", err),
		}
	}

	webResponse := map[string]interface{}{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Book with ID %v has been deleted.", id),
	}

	return http.StatusOK, webResponse
}

func (s *bookService) UpdateStatusBookById(id primitive.ObjectID, book models.UpdateStatusBookRequest) (int, map[string]interface{}) {
	_, err := s.bookRepository.UpdateStatusBookById(id, book)
	if err != nil {
		return http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error database: %v", err),
		}
	}

	webResponse := map[string]interface{}{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Book with ID %v has been deleted.", id),
	}

	return http.StatusOK, webResponse
}
