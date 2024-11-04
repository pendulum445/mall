package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"google.golang.org/grpc"
	authpb "mall/auth"
	userpb "mall/user"
)

func main() {
	// 初始化数据库连接
	dsn := "root:082425lph@tcp(localhost:3306)/mall" // 替换为你的实际数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping the database: %v", err)
	}

	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// 注册 AuthService 和 UserService
	authpb.RegisterAuthServiceServer(grpcServer, &authpb.AuthService{})
	userpb.RegisterUserServiceServer(grpcServer, &userpb.UserService{
		// 将数据库连接传递给 UserService
		Db: db,
	})

	log.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
