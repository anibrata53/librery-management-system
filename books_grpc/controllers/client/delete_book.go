package client

import (
	"books_grpc/service"
	"context"
	"fmt"
	"log"
	"strings"
)

//Delete Book
func DeleteBook(client service.BookServiceClient, ctx context.Context) {
	//Input
	//fmt.Println("Input Any details to delete a book")
	fmt.Printf("Book Id: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	id := strings.TrimSpace(input)
	//Call DeleteBook
	_, err = client.DeleteBook(ctx, &service.DeleteBookReq{Id: id})
	if err != nil {
		log.Fatal(err)
	}

	//Print Result
	fmt.Print("\nDeleted book with id: ", id)
}
