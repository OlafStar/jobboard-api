package jwt

import (
	"time"

	"github.com/OlafStar/jobboard-api/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewJWTCompany(registerUser types.RegisterCompany) (types.JWTCompany, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return types.JWTCompany{}, err
	}

	return types.JWTCompany{
		Email:         registerUser.Email,
		Name:          registerUser.Name,
		PasswordHash:  string(hashedPassword),
		Country:       registerUser.Country,
		City:          registerUser.City,
		Address:       registerUser.Address,
	}, nil
}

func CreateCompanyToken(company types.JWTCompany) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"email":   company.Email,
		"name":    company.Name,
		"country": company.Country,
		"city": company.City,
		"address": company.Address,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "secret" //TO CHANGE TO .ENV

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return tokenString
}

func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
