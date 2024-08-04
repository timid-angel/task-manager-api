package models

type Library struct {
	Books   map[int]Book
	Members map[int]Member
}

type LibError struct {
	message string
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
	if lib.Books[bookID].Status == "borrowed" {
		return LibError{message: "Book not available: It has already been borrowed."}
	}

	return nil
}
