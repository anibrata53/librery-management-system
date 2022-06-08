package service

import (
	context "context"
	"fmt"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type BookDetails struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

//Implementation of proto interface
type BooksServiceServer struct {
	&service.UnimplementedBooksServiceServer
}

//Create Book
func (s *BookServiceServer) CreateBook(ctx context.Context, req *service.CreateBookReq) (*pb.CreateBookRes, error) {

	book := req.Book
	//convert message CreateBookRequest into a BookDetails type to convert into BSON
	data := BookDetails{
		Title:  book.Title,
		Author: book.Author,
	}

	//Insert Data
	inserted, err := database.Books.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Insert operation FAILED! \n %v", err))
	}

	//id assigned by mongodb
	id := inserted.InsertedID.(primitive.ObjectID)

	//convert to string Counterpart
	book.Id = id.Hex()

	return &pb.CreateBookRes{Book: book}, err

}

//List All Books
func (s *LmsServiceServer) ListAllBooks(req *pb.CreateBookReq, stream pb.LmsService_ListAllBooksServer) error {

	//data contains structure of book
	data := &BookDetails{}

	//Initializing Cursor
	cursor, err := database.Books.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	defer cursor.Close(context.Background())

	// Looping through Database
	for cursor.Next(context.Background()) {

		//whenever we decode, we pass on a referrence like, "if you decode use my structure to decode that"
		err := cursor.Decode(data)
		// check error
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}

		// If no error is found, send book over stream
		stream.Send(&pb.ListAllBooksRes{
			Book: &pb.Book{
				Id:     data.ID.Hex(),
				Title:  data.Title,
				Author: data.Author,
			},
		})

	}
	// Check if the cursor has any errors
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown cursor error : %v", err))
	}
	return nil
}

//Search Book
func (s *LmsServiceServer) SearchBooks(req *pb.SearchBooksReq, stream pb.LmsService_SearchBooksServer) error {

	//data contains structure of book
	data := &BookDetails{}

	var cursor *mongo.Cursor
	var err error

	title := req.GetTitle()
	author := req.GetAuthor()

	if title == "" && author == "" {
		return status.Errorf(codes.Internal, fmt.Sprintf("Empty Search Key: %v", err))
	}

	//Finding books Initializing Cursor
	if title != "" {
		cursor, err = database.Books.Find(context.Background(), bson.M{"title": title})
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Book not found, Check title: %v", err))
		}
	}

	if author != "" {
		cursor, err = database.Books.Find(context.Background(), bson.M{"author": author})
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Book not found, Check author: %v", err))
		}

	}

	defer cursor.Close(context.Background())

	// Looping through Database
	for cursor.Next(context.Background()) {

		//whenever we decode, we pass on a referrence like, "if you decode use my structure to decode that"
		err := cursor.Decode(data)
		// check error
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}

		// If no error is found, send book over stream
		stream.Send(&pb.SearchBooksRes{
			Book: &pb.Book{
				Id:     data.ID.Hex(),
				Title:  data.Title,
				Author: data.Author,
			},
		})

	}
	// Check if the cursor has any errors
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}

//Update Book
func (s *LmsServiceServer) UpdateBook(ctx context.Context, req *pb.UpdateBookReq) (*pb.UpdateBookRes, error) {

	// Get the Book data from the request
	book := req.GetBook()

	// Convert the Id string to a MongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(book.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert the supplied book id to a MongoDB ObjectId: %v", err))
	}

	// Convert the oid into an unordered bson document to search by id
	filter := bson.M{"_id": oid}

	// Convert the data to be updated into an unordered Bson document
	updates := bson.M{"$set": bson.M{"title": book.GetTitle(), "author": book.GetAuthor()}}

	// Result is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	result := database.Books.FindOneAndUpdate(ctx, filter, updates, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode result and write it to 'decoded'
	decoded := BookDetails{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find book with supplied ID: %v", err),
		)
	}
	return &pb.UpdateBookRes{Book: &pb.Book{Id: decoded.ID.Hex(), Title: decoded.Title, Author: decoded.Author}}, nil
}

//Delete Book
func (s *LmsServiceServer) DeleteBook(ctx context.Context, req *pb.DeleteBookReq) (*pb.DeleteBookRes, error) {

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	filter := bson.M{"_id": oid}
	_, err = database.Books.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf(`Delete operation FAILED!
		Could not find Book with id %s: %v`, req.GetId(), err))
	}

	return &pb.DeleteBookRes{Success: true}, nil
}
