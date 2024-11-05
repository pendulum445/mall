package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	authpb "mall/auth"
	userpb "mall/user"
)

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
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, &authpb.AuthService{})
	userpb.RegisterUserServiceServer(grpcServer, &userpb.UserService{
		Db: db,
	})
	log.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
