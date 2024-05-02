package api

import (
	"net/http"
)

func (s *APIServer) SetupCompanyAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/company/auth/login", s.makeHTTPHandleFunc(s.HandleLoginCompany))
	mux.HandleFunc("/api/company/auth/register", s.makeHTTPHandleFunc(s.HandleRegisterCompany))
	// mux.HandleFunc("/api/iam", s.makeHTTPHandleFunc(ProtectedRequest(s.Iam)))
}