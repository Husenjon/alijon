package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/Husenjon/InkassBack/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt      = "%&jegGf@3"
	signinKey = "sd%&jegGf@3gG%&jegGf@%&jegGf@3F#%9/*"
	tokenTTL  = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
type AuthService struct {
	repo repository.Authoration
}

func NewAuthService(repo repository.Authoration) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) CreateUser(user inkassback.User) (inkassback.User, error) {
	user.Password = s.GenerateHash(user.Password)
	return s.repo.CreateUser(user)
}
func (s *AuthService) GenerateHash(secret string) string {
	hash := sha1.New()
	hash.Write([]byte(secret))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
func (s *AuthService) GetUser(username, password string) (inkassback.User, error) {
	user, err := s.repo.GetUser(username, s.GenerateHash(password))
	if err != nil {
		return user, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	token2, err := token.SignedString([]byte(signinKey))
	user, err = s.repo.UpdateToken(user.Id, token2)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (s *AuthService) ParseToken(acessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		acessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid sign method")
			}
			return []byte(signinKey), nil
		})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type * tokenClaims")
	}
	return claims.UserId, nil
}
