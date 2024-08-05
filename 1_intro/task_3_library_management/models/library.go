package models

import "fmt"

type Library struct {
	Books   map[int]Book
	Members map[int]Member
}

type LibError struct {
	message string
}

func GetBooks(lib *Library, status string) []Book {
	availableBooks := make([]Book, 0)
	for _, v := range lib.Books {
		if v.Status == status {
			availableBooks = append(availableBooks, v)
		}
	}

	return availableBooks
}

func (err LibError) Error() string {
	return err.message
}

func (lib *Library) AddBook(newBook Book) {
	lib.Books[newBook.ID] = newBook
}

func (lib *Library) RemoveBook(removedID int) {
	delete(lib.Books, removedID)
}

func (lib *Library) BorrowBook(bookID int, memberID int) error {
	if lib.Books[bookID].Status == "borrowed" {
		return LibError{message: "Book not available: It has already been borrowed."}
	}

	member := lib.Members[memberID]
	member.BorrowedBooks = append(member.BorrowedBooks, lib.Books[bookID])
	lib.Members[memberID] = member

	return nil
}

func (lib *Library) ReturnBook(bookID int, memberID int) error {
	book := lib.Books[bookID]
	if book.Status == "available" {
		return LibError{message: "Book has not been borrowed"}
	}

	member := lib.Members[memberID]
	found := false
	newBorrowedBooks := make([]Book, 0)
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

func (lib *Library) ListAvailableBooks() []Book {
	return GetBooks(lib, "available")
}

func (lib *Library) ListBorrowedBooks() []Book {
	return GetBooks(lib, "borrowed")
}
