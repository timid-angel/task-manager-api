package services

import "library_management/models"

var nextInt = 0
var LIB = models.Library{Books: make(map[int]models.Book), Members: map[int]models.Member{
	0: {ID: 0, Name: "Shear", BorrowedBooks: []models.Book{}},
	1: {ID: 1, Name: "Aben", BorrowedBooks: []models.Book{}},
	2: {ID: 2, Name: "Zyri", BorrowedBooks: []models.Book{}},
}}

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
