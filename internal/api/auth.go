package api

import (
	"net/http"

	"github.com/OlafStar/jobboard-api/internal/jwt"
	"github.com/OlafStar/jobboard-api/internal/types"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func (s *APIServer) HandleLoginCompany(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.JWTLoginHandler(w, r)
	}
	
	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) HandleRegisterCompany(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.registerCompany(w, r)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) JWTLoginHandler(w http.ResponseWriter, r *http.Request) error {
	var company types.LoginCompany

	if err := DecodeJSONBody(r, &company); err != nil {
		httpErr, ok := err.(*HTTPError)
		if !ok {
				httpErr = &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Internal server error"}
		}
		return &HTTPError{StatusCode: httpErr.StatusCode, Message: httpErr.Message}
	}

	jwtCompany, err := s.store.GetJWTCompany(company.Email)
	if err != nil {
			return &HTTPError{StatusCode: http.StatusUnauthorized, Message: "Authorization failed"}
	}

	if !jwt.ValidatePassword(jwtCompany.PasswordHash, company.Password) {
			return &HTTPError{StatusCode: http.StatusUnauthorized, Message: "invalid credentials"}
	}

	tokenString := jwt.CreateCompanyToken(*jwtCompany)
	if tokenString == "" {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Internal server error"}
	}

	return WriteJSON(w, http.StatusOK, LoginResponse{
		Token: tokenString,
	})
}

func (s *APIServer) registerCompany(w http.ResponseWriter, r *http.Request) error {
	var body types.RegisterCompany
	
	if err := DecodeJSONBody(r, &body); err != nil {
		httpErr, ok := err.(*HTTPError)
		if !ok {
				httpErr = &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Internal server error"}
		}
		return &HTTPError{StatusCode: httpErr.StatusCode, Message: httpErr.Message}
	}


	isExist, err := s.store.DoesCompanyExist(body.Email)

	if isExist {
		return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Username currently exists"}
	}

	if err != nil {
		return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Error checking username existence"}
	}

	jwtCompany, err := jwt.NewJWTCompany(body)
	
	if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to create company profile"}
	}

	err = s.store.InsertCompany(jwtCompany)

	if err != nil {
		return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to create company profile"}
	}

	return WriteJSON(w, http.StatusOK, "Success")
}