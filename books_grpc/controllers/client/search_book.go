package client

import (
	"books_grpc/service"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

//Search Books / Bi-Directional Streaming
func SearchBooks(client service.BookServiceClient, ctx context.Context) {

	//Initializing stream
	var stream service.BookService_SearchBooksClient
	//Input
	fmt.Println(`Search Books from Library: 
	Chose an option...
	1. Search By Title
	2. Search By Author`)

	fmt.Printf("Your Decition : ")
	input, _ := reader.ReadString('\n')
	decition, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	if err != nil {
		log.Fatal("Failed to convert string into int")
	}

	switch decition {
	case 1:
		fmt.Print("Title: ")
		title, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Wrong input: Title", err)
		}
		title = strings.TrimSpace(title)

		// Call SearchBook that returns a stream
		stream, err = client.SearchBooks(ctx, &service.SearchBooksReq{
			Search: &service.SearchBooksReq_Title{Title: title},
		})
		// Check for errors
		if err != nil {
			fmt.Println("SearchBook did not return stream: ", err)
		}

	case 2:

		fmt.Print("Author: ")
		author, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Wrong input: Author", err)
		}
		author = strings.TrimSpace(author)

		// Call ListBook that returns a stream
		stream, err = client.SearchBooks(ctx, &service.SearchBooksReq{
			Search: &service.SearchBooksReq_Author{Author: author},
		})

		// Check for errors
		if err != nil {
			fmt.Println("SearchBook did not return stream: ", err)
		}

	default:
		fmt.Println("Invalid Option")
	}

	//Print Result
	fmt.Println("Books we have found...")

	// Start iterating
	for i := 0; i < 3; i++ {

		// stream.Recv returns a pointer to a book in a current iteration
		responseStream, err := stream.Recv()
		// If end of stream, break the loop
		if err == io.EOF {
			break
		}
		// if err, print error
		if err != nil {
			log.Fatal("Stream error: ", err)
		}

		// If everything went well use the generated getter to print the Book Details
		fmt.Println(responseStream.GetBook())

	}
}
