package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

var commands = map[string]string{
	"Add book":             "add",
	"Remove book":          "remove",
	"Borrow book":          "borrow",
	"Return book":          "return",
	"List available books": "la",
	"List borrowed books":  "lb",
}

func GetInputStr(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("\t\t**Error: ", err.Error())
			continue
		}

		return strings.TrimSpace(input)
	}
}

func GetInputInt(reader *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("\t\t**Error: ", err.Error())
			continue
		}

		res, err := strconv.ParseInt(strings.TrimSpace(input), 10, 32)
		if err != nil {
			fmt.Print("\t\t**Error: ", err.Error())
			continue
		}

		return int(res)
	}
}

func LogBooks(books []models.Book, title string) {
	fmt.Printf("\n%v\n\t%20v\t\t%20v\t\t%20v\n", title, "ID", "Title", "Author")
	for _, book := range books {
		fmt.Printf("\t%20v\t\t%20v\t\t%20v\n", book.ID, book.Title, book.Author)
	}
}

func Run() {
	reader := bufio.NewReader(os.Stdin)
	active := true
	fmt.Print("\nCommands:\n")
	for k, v := range commands {
		fmt.Printf("\t- %-25v %v\n", k+":", v)
	}

	for active {
		cmd := GetInputStr(reader, "\n> ")
		switch cmd {
		case "add":
			title := GetInputStr(reader, "\tBook title: ")
			author := GetInputStr(reader, "\tBook author: ")
			services.HandleAddBook(title, author)
			fmt.Println("\t\tBook added successfully")

		case "remove":
			bookID := GetInputInt(reader, "\tBook ID: ")
			services.HandleRemoveBook(bookID)
			fmt.Printf("\t\tBook with ID %v removed successfully\n", bookID)

		case "borrow":
			bookID := GetInputInt(reader, "\tBook ID: ")
			memberID := GetInputInt(reader, "\tMember ID: ")
			err := services.HandleBorrow(bookID, memberID)
			if err != nil {
				fmt.Printf("\t\t**Error: %v", err.Error())
				continue
			}

			fmt.Printf("\t\tMember with ID %v has borrowed book #%v\n", memberID, bookID)

		case "return":
			bookID := GetInputInt(reader, "\tBook ID: ")
			memberID := GetInputInt(reader, "\tMember ID: ")
			err := services.HandleReturn(bookID, memberID)
			if err != nil {
				fmt.Printf("\t\t**Error: %v", err.Error())
				continue
			}

			fmt.Printf("\t\tBook with ID: %v returned successfully\n", bookID)

		case "la":
			availableBooks := services.HandleListAvailable()
			LogBooks(availableBooks, "\tAvailable Books:")

		case "lb":
			borrowedBooks := services.HandleListBorrowed()
			LogBooks(borrowedBooks, "\tBorrowed Books:")

		case "q":
			active = false

		default:
			fmt.Println("\t\t**Invalid command")
		}

	}
}
