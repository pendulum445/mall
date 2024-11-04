package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_key")

type AuthService struct {
	UnimplementedAuthServiceServer
}

func (s *AuthService) DeliverTokenByRPC(ctx context.Context, req *DeliverTokenReq) (*DeliveryResp, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   string(req.UserId),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}
	return &DeliveryResp{Token: tokenString}, nil
}

func (s *AuthService) VerifyTokenByRPC(ctx context.Context, req *VerifyTokenReq) (*VerifyResp, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &VerifyResp{Res: true}, nil
}
