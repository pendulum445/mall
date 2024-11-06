package user

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log" // 使用 log 包
)

type UserService struct {
	UnimplementedUserServiceServer
	Db *sql.DB
}

func (s *UserService) Register(ctx context.Context, req *RegisterReq) (*RegisterResp, error) {
	// 使用 log 记录收到的请求信息
	log.Printf("Register request received: Email=%s, Password=%s, ConfirmPassword=%s", req.Email, req.Password, req.ConfirmPassword)

	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		log.Println("Error: invalid parameters") // 使用 log 记录错误
		return nil, errors.New("invalid params")
	}
	if req.Password != req.ConfirmPassword {
		log.Println("Error: passwords do not match") // 使用 log 记录错误
		return nil, errors.New("password not match")
	}

	var existingID int
	err := s.Db.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err == nil {
		log.Println("Error: email already registered") // 使用 log 记录错误
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Error querying database: %v", err) // 使用 log 记录数据库查询错误
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error: failed to hash password") // 使用 log 记录错误
		return nil, errors.New("failed to hash password")
	}

	result, err := s.Db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", req.Email, hashedPassword)
	if err != nil {
		log.Printf("Error inserting into database: %v", err) // 使用 log 记录数据库插入错误
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err) // 使用 log 记录获取插入 ID 错误
		return nil, err
	}

	log.Printf("User registered successfully: UserID=%d", userID) // 使用 log 记录成功注册
	return &RegisterResp{UserId: int32(userID)}, nil
}

func (s *UserService) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	// 使用 log 记录收到的登录请求信息
	log.Printf("Login request received: Email=%s, Password=%s", req.Email, req.Password)

	if req.Email == "" || req.Password == "" {
		log.Println("Error: invalid parameters") // 使用 log 记录错误
		return nil, errors.New("invalid params")
	}

	var userID int
	var password []byte
	err := s.Db.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email).Scan(&userID, &password)
	if err != nil {
		log.Printf("Error querying user: %v", err) // 使用 log 记录数据库查询错误
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(password, []byte(req.Password))
	if err != nil {
		log.Println("Error: invalid password") // 使用 log 记录错误
		return nil, errors.New("invalid password")
	}

	log.Printf("User logged in successfully: UserID=%d", userID) // 使用 log 记录成功登录
	return &LoginResp{UserId: int32(userID)}, nil
}
