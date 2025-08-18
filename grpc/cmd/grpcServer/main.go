package main

import (
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-study/internal/database"
	"grpc-study/internal/pb"
	"grpc-study/internal/service"
	"net"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize the database connection and create the CategoryDB instance
	categoryDB := database.NewCategory(db)
	// Initialize the CategoryService with the CategoryDB instance
	categoryService := service.NewCategoryService(*categoryDB)

	// Create a new gRPC server and register the CategoryService
	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	// Start the gRPC server - open a listener on a specific port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}

}
