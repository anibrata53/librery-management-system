package client

import (
	"books_grpc/service"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const address = ":50051"

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

func main() {
	//dial conn to grpc server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Dial Failed:", err)
	}

	defer conn.Close()
	//create new client
	client := service.NewBookServiceClient(conn)

	// init context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Input option

	for {
		fmt.Println(`Welcome to book management system...
		
		**MENU**
		1. Upload a New Book
		2. View All Books
		3. Search a Book
		4. Update Book
		5. Delete Book
		choose other to exit`)

		fmt.Println("Choose Option: ")

		input, _ := reader.ReadString('\n')
		option, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if err != nil {
			log.Fatal("Failed to convert string into int")
		}

		switch option {
		case 1:
			CreateBook(client, ctx)
			continue
		case 2:
			ListAllBooks(client, ctx)
			continue
		case 3:
			SearchBooks(client, ctx)
			continue
		case 4:
			UpdateBook(client, ctx)
			continue
		case 5:
			DeleteBook(client, ctx)
			continue
		default:
			os.Exit(0)
		}
	}
}
