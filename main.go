package main

import (
	"context"
	"database/sql"
	"log"
	authpb "mall/auth"
	cartpb "mall/cart"
	productpb "mall/product"
	userpb "mall/user"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Received request from: %s", p.Addr.String())
	}
	log.Printf("Handling gRPC method: %s", info.FullMethod)
	return handler(ctx, req)
}

func main() {
	dsn := "root:082425lph@tcp(localhost:3306)/mall"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close the database connection: %v", err)
		}
	}(db)
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping the database: %v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	authpb.RegisterAuthServiceServer(grpcServer, &authpb.AuthService{})
	userpb.RegisterUserServiceServer(grpcServer, &userpb.UserService{Db: db})
	productpb.RegisterProductCatalogServiceServer(grpcServer, &productpb.ProductService{Db: db})
	cartpb.RegisterCartServiceServer(grpcServer, &cartpb.CartService{Db: db})
	reflection.Register(grpcServer)
	log.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
