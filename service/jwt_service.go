package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userId string, role string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDByToken(token string) (string, error)
	GetRoleByToken(token string) (string, error)
	GenerateAccessToken(userID string, role string) (string, int64, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(userId string, role string) string {
	claims := jwtCustomClaim{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return tx
}

func (j *jwtService) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
	}
	return []byte(j.secretKey), nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, j.parseToken)
}

func (j *jwtService) GetUserIDByToken(token string) (string, error) {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}

	claims := t_Token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id, nil
}
func (j *jwtService) GetRoleByToken(token string) (string, error) {
	t_, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}

	claims := t_.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	return role, nil
}

func (j *jwtService) GenerateAccessToken(userID string, role string) (string, int64, error) {
	expirationTime := time.Now().Add(time.Minute * 15)
	expiresIn := int64(expirationTime.Sub(time.Now()).Seconds())

	claims := &jwtCustomClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", 0, err
	}
	return t, expiresIn, nil
}
