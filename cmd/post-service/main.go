package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nucleussss/hikayat-forum/post/db"
	"github.com/Nucleussss/hikayat-forum/post/internal/delivery/grpc"
	"github.com/Nucleussss/hikayat-forum/post/internal/repository/postgres"
	"github.com/Nucleussss/hikayat-forum/post/internal/service"

	pb "github.com/Nucleussss/hikayat-forum/post/api/post/v1"
)

func main() {
	connString := db.ConnString()
	pool, err := db.InitDB(connString)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		println("Closing database connection")
		pool.Close()
	}()

	// initiate repo
	postRepo := postgres.NewPostRepo(pool)

	postService := service.NewPostService(postRepo)

	postHandler := grpc.NewPostHandler(postService)

	grpcServer := grpc.NewServer()

	pb.RegisterPostServiceServer(grpcServer, postHandler)

	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen on port %s : %v", os.Getenv("GRPC_PORT"), err)
	}
	log.Printf("Starting gRPC server at %s\n", os.Getenv("GRPC_PORT"))

	// gracefull shoutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	grpcStopped := make(chan struct{})
	go func() {
		defer close(grpcStopped)
		log.Printf("start grpc server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("grpc server stopped with error: %v", err)
		} else {
			log.Printf("grpc server stopped gracefully")
		}
	}()

	// wait for shutdown signal
	<-sigChan
	log.Printf("recieve shoutdown signal, Stopping gRPC server\n ")

	// stop grpc server
	grpcServer.GracefulStop()
	log.Printf("gRPC stopped gracefully")

	// wait for shutdown signal
	<-grpcStopped
	log.Println("post service stopped")

}
