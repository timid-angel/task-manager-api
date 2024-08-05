package controllers

import (
	"bufio"
	"fmt"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

func GetInputStr(reader *bufio.Reader, prompt string) string {
	for {
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

func Run() {
	reader := bufio.NewReader(os.Stdin)
	active := true
	fmt.Print("\nCommands:\n\tAdd book: 'add'\n\tRemove book: 'remove'\n\tBorrow book: 'borrow'\n\tReturn book: 'return'\n\tList available books: 'la'\n\tList borrowed books: 'lb'\n\tExit: 'q'\n")
	for active {
		cmd := GetInputStr(reader, "> ")

		switch cmd {
		case "add":
			title := GetInputStr(reader, "\tBook title: ")
			author := GetInputStr(reader, "\tBook author: ")
			services.HandleAddBook(title, author)

		case "remove":
			bookID := GetInputInt(reader, "\tBook ID: ")
			services.HandleRemoveBook(bookID)

		case "borrow":
			bookID := GetInputInt(reader, "\tBook ID: ")
			memberID := GetInputInt(reader, "\tMember ID: ")
			services.HandleBorrow(bookID, memberID)

		case "return":
			bookID := GetInputInt(reader, "\tBook ID: ")
			memberID := GetInputInt(reader, "\tMember ID: ")
			services.HandleReturn(bookID, memberID)

		case "la":
			services.HandleListAvailable()
		case "lb":
			services.HandleListBorrowed()

		case "q":
			active = false

		default:
			fmt.Println("\t\t**Invalid command")
		}

	}
}
