package controllers

import (
	"fmt"
	"library_management/models"
)

type Library struct {
	Books         map[int]models.Book
	Members       map[int]models.Member
	BorrowedBooks map[int]models.BorrowedBook
}

func (l *Library) AddBook(book models.Book) {
	l.Books[book.ID] = book
	fmt.Println("Book added:", book.Title)
}

func (l *Library) RemoveBook(bookID int) {
	if _, exists := l.Books[bookID]; exists {
		delete(l.Books, bookID)
		fmt.Println("Book removed with ID:", bookID)
	} else {
		fmt.Println("Book not found with ID:", bookID)
	}
}

func (l *Library) AddMember(member models.Member) {
	l.Members[member.ID] = member
	fmt.Println("Member added:", member.Name)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return fmt.Errorf("book not found with ID %d", bookID)
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
		return fmt.Errorf("member not found with ID %d", memberID)
	}

	delete(l.Books, bookID)
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	l.BorrowedBooks[bookID] = models.BorrowedBook{BookID: bookID, MemberID: memberID, Title: book.Title}

	fmt.Printf(" Book '%s' borrowed by %s\n", book.Title, member.Name)
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	member, ok := l.Members[memberID]
	if !ok {
		return fmt.Errorf("member not found with ID %d", memberID)
	}

	borrowed, exists := l.BorrowedBooks[bookID]
	if !exists {
		return fmt.Errorf("book not borrowed or does not exist")
	}

	if borrowed.MemberID != memberID {
		return fmt.Errorf("book was borrowed by another member")
	}

	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			break
		}
	}

	l.Members[memberID] = member
	delete(l.BorrowedBooks, bookID)

	l.Books[bookID] = models.Book{
		ID:     borrowed.BookID,
		Title:  borrowed.Title,
		Status: "Available",
	}

	fmt.Printf("Book '%s' returned by %s\n", borrowed.Title, member.Name)
	return nil
}

func (l *Library) ListAvailableBooks() {
	fmt.Println(" Available books:")
	for _, b := range l.Books {
		fmt.Printf("- %s by %s\n", b.Title, b.Author)
	}
}

func (l *Library) ListBorrowedBooks(memberID int) {
	member, exists := l.Members[memberID]
	if !exists {
		fmt.Println(" Member not found.")
		return
	}
	fmt.Printf("Books borrowed by %s:\n", member.Name)
	for _, b := range member.BorrowedBooks {
		fmt.Printf("- %s by %s\n", b.Title, b.Author)
	}
}
