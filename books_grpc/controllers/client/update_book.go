package client

import (
	"books_grpc/service"
	"context"
	"fmt"
	"log"
	"strings"
)

//Update Book
func UpdateBook(client service.BookServiceClient, ctx context.Context) {

	//Input
	//reader := bufio.NewReader(os.Stdin)
	fmt.Println("Insert Details to update: ")

	fmt.Print("Id: ")
	id, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Wrong input: Title", err)
	}
	id = strings.TrimSpace(id)

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

	if id == "" || title == "" || author == "" {
		log.Fatal("Empty input")
	}

	//Creating Request
	updateBook := &service.Book{
		Id:     id,
		Title:  title,
		Author: author,
	}
	//Call UpdateBook that returns a Book as response
	response, err := client.UpdateBook(ctx, &service.UpdateBookReq{Book: updateBook})
	if err != nil {
		log.Fatal("Could not update book: \n", err)
	}
	//print
	log.Printf(`Book Updated:
	Book Id: %s
	Title: %s
	Author: %s`, response.Book.Id, response.Book.Title, response.Book.Author)
}
