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

//All Books / Bi-Directional Streaming
func ListAllBooks(client service.BookServiceClient, ctx context.Context) {

	fmt.Print("Page no: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	pageNo, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		log.Fatal("Wrong input: Title", err)
	}

	fmt.Print("How many books you want in a page? ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	noOfItems, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		log.Fatal("Wrong input: Title", err)
	}

	dontPrintUpto := (pageNo - 1) * noOfItems

	// Call ListAllBook that returns a stream
	stream, err := client.ListAllBooks(ctx, &service.CreateNewBookReq{})
	// Check for errors
	if err != nil {
		log.Fatal("ListAllBooks did not return stream: ", err)
	}

	fmt.Println("List of Books: ")

	// Start iterating
	for i := 0; i < int(pageNo*noOfItems); i++ {

		// stream.Recv returns a pointer to a book in a current iteration
		response, err := stream.Recv()
		// If end of stream, break the loop
		if err == io.EOF {
			break
		}
		// if err, print error
		if err != nil {
			log.Fatal("Stream error: ", err)
		}

		// If everything went well use the generated getter to print the Book Details
		if i >= int(dontPrintUpto) {

			fmt.Printf("%d %v\n", i+1, response.GetBook())

		}

	}
}
