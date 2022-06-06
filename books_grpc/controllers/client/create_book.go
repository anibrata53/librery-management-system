package client

import (
	"books_grpc/service"
	"context"
	"fmt"
	"log"
	"strings"
)

//Create Book
func CreateBook(client service.BookServiceClient, ctx context.Context) {

	//Input
	fmt.Println("Insert Details: ")

	fmt.Print("Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Wrong input: Title", err)
	}
	title = strings.TrimSpace(title)

	fmt.Print("Author: ")
	author, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Wrong input: Title", err)
	}
	author = strings.TrimSpace(author)

	if title == "" || author == "" {
		log.Fatal("Empty input")
	}

	//Creating book to send as request
	newBook := &service.Book{
		Title:  title,
		Author: author,
	}
	//Call CreateBook that returns a book as response
	response, err := client.CreateBook(ctx, &service.CreateNewBookReq{Book: newBook})
	if err != nil {
		log.Fatal("Could not Create Book: \n", err)
	}
	//print
	log.Printf(`New Book Uploaded:
	Book Id: %s
	Title: %s
	Author: %s`, response.Book.Id, response.Book.Title, response.Book.Author)
}
