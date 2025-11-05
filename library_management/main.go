package main

import (
	"fmt"
	"library_management/controllers"
	"library_management/models"
)

func main() {
	library := controllers.Library{
		Books:         make(map[int]models.Book),
		Members:       make(map[int]models.Member),
		BorrowedBooks: make(map[int]models.BorrowedBook),
	}

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1 Add member")
		fmt.Println("2 Add book")
		fmt.Println("3 Remove book")
		fmt.Println("4 Borrow book")
		fmt.Println("5 Return book")
		fmt.Println("6 List available books")
		fmt.Println("7 List borrowed books by member")
		fmt.Println("8 Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var id int
			var name string
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&id)
			fmt.Print("Enter name: ")
			fmt.Scanln(&name)
			library.AddMember(models.Member{ID: id, Name: name})

		case 2:
			var id int
			var title, author string
			fmt.Print("Enter book ID: ")
			fmt.Scanln(&id)
			fmt.Print("Enter title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter author: ")
			fmt.Scanln(&author)
			library.AddBook(models.Book{ID: id, Title: title, Author: author, Status: "Available"})

		case 3:
			var id int
			fmt.Print("Enter book ID to remove: ")
			fmt.Scanln(&id)
			library.RemoveBook(id)

		case 4:
			var bookID, memberID int
			fmt.Print("Enter book ID to borrow: ")
			fmt.Scanln(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&memberID)
			if err := library.BorrowBook(bookID, memberID); err != nil {
				fmt.Println("Error:", err)
			}

		case 5:
			var bookID, memberID int
			fmt.Print("Enter book ID to return: ")
			fmt.Scanln(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&memberID)
			if err := library.ReturnBook(bookID, memberID); err != nil {
				fmt.Println("Error:", err)
			}

		case 6:
			library.ListAvailableBooks()

		case 7:
			var memberID int
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&memberID)
			library.ListBorrowedBooks(memberID)

		case 8:
			fmt.Println("👋 Goodbye!")
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}
