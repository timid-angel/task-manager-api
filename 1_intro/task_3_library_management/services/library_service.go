package services

import (
	"fmt"
	"library_management/models"
)

type MyLibrary models.Library

type LibError struct {
	message string
}

func (err LibError) Error() string {
	return err.message
}

var nextInt = 0
var LIB = MyLibrary{Books: make(map[int]models.Book), Members: map[int]models.Member{
	0: {ID: 0, Name: "Shear", BorrowedBooks: []models.Book{}},
	1: {ID: 1, Name: "Aben", BorrowedBooks: []models.Book{}},
	2: {ID: 2, Name: "Zyri", BorrowedBooks: []models.Book{}},
}}

func GetBooks(lib *MyLibrary, status string) []models.Book {
	availableBooks := make([]models.Book, 0)
	for _, v := range lib.Books {
		if v.Status == status {
			availableBooks = append(availableBooks, v)
		}
	}

	return availableBooks
}

func (lib *MyLibrary) AddBook(newBook models.Book) {
	lib.Books[newBook.ID] = newBook
}

func (lib *MyLibrary) RemoveBook(removedID int) {
	delete(lib.Books, removedID)
}

func (lib *MyLibrary) BorrowBook(bookID int, memberID int) error {
	if lib.Books[bookID].Status == "borrowed" {
		return LibError{message: "Book not available: It has already been borrowed."}
	}

	book := lib.Books[bookID]
	book.Status = "borrowed"
	lib.Books[bookID] = book
	member := lib.Members[memberID]
	member.BorrowedBooks = append(member.BorrowedBooks, lib.Books[bookID])
	lib.Members[memberID] = member

	return nil
}

func (lib *MyLibrary) ReturnBook(bookID int, memberID int) error {
	book := lib.Books[bookID]
	if book.Status == "available" {
		return LibError{message: "Book has not been borrowed"}
	}

	member := lib.Members[memberID]
	found := false
	newBorrowedBooks := make([]models.Book, 0)
	for _, v := range member.BorrowedBooks {
		if v.ID == bookID {
			found = true
			book.Status = "available"
			lib.Books[bookID] = book
			continue
		}

		newBorrowedBooks = append(newBorrowedBooks, v)
	}

	if !found {
		return LibError{message: fmt.Sprintf("Book has not been borrowed by member with ID: %v", memberID)}
	}

	member.BorrowedBooks = newBorrowedBooks
	return nil
}

func (lib *MyLibrary) ListAvailableBooks() []models.Book {
	return GetBooks(lib, "available")
}

func (lib *MyLibrary) ListBorrowedBooks() []models.Book {
	return GetBooks(lib, "borrowed")
}

func HandleAddBook(title string, author string) {
	newBook := models.Book{ID: nextInt, Title: title, Status: "available", Author: author}
	nextInt++
	LIB.AddBook(newBook)
}

func HandleRemoveBook(bookID int) {
	LIB.RemoveBook(bookID)
}

func HandleBorrow(bookID int, memberID int) error {
	return LIB.BorrowBook(bookID, memberID)
}

func HandleReturn(bookID int, memberID int) error {
	return LIB.ReturnBook(bookID, memberID)
}

func HandleListAvailable() []models.Book {
	return LIB.ListAvailableBooks()
}

func HandleListBorrowed() []models.Book {
	return LIB.ListBorrowedBooks()
}
