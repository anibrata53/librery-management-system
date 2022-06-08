package server

import (
	"fmt"

	"log"
	"net"

	"books_grpc/database"
	"books_grpc/service"

	"google.golang.org/grpc"
)

const port = ":50051"

func main() {

	fmt.Println("\nWelcome to CRUD GRPC")

	//Creatiung Database Connection
	database.InitDB()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen database: ", err)
	}

	//Initialize new Server
	server := grpc.NewServer()

	//Regester the server as a new grpc service
	service.RegisterBookServiceServer(server, &service.BookServiceServer{})
	log.Println("server listening at ", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}

}
