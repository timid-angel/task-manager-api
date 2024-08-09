# Getting Started

To run the applicaiton, run the following command in the root directory of the project:
```bash
go run .
```

## Commands

| Operation | Command | Description |
| - | - | - |
| Add book  |   `add`   | Adds a book in the in-memory db |
| Remove book | `remove` | Removes a book from the library after prompting for the bookID |
| Borrow book | `borrow` | Lends a book specified by the bookID to the member specified by the memberID |
| Return book | `return` | Returns a previously borrowed book specified by the bookID by the member specified by the memberID |
| List available books | `la` | List of all the available (non-borrowed) books from the list of books |
| List borrowed books | `lb` | List of all the borrowed books from the list of books |

	"Remove book":          "remove",
	"Borrow book":          "borrow",
	"Return book":          "return",
	"List available books": "la",
	"List borrowed books":  "lb",