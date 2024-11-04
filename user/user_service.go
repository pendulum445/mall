package user

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UnimplementedUserServiceServer
	Db *sql.DB
}

func (s *UserService) Register(ctx context.Context, req *RegisterReq) (*RegisterResp, error) {
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return nil, errors.New("invalid params")
	}
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password not match")
	}
	var existingID int
	err := s.Db.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	result, err := s.Db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &RegisterResp{UserId: int32(userID)}, nil
}

func (s *UserService) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("invalid params")
	}
	var userID int
	var password []byte
	err := s.Db.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email).Scan(&userID, &password)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(password, []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}
	return &LoginResp{UserId: int32(userID)}, nil
}
