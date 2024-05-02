package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/OlafStar/jobboard-api/internal/middleware"
	"github.com/OlafStar/jobboard-api/internal/queue"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

type ApiError struct {
	Error string `json:"message"`
}

func (s *APIServer) makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: []string{"https://example.com", "*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}
	return middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		errc := make(chan error, 1)

		job := queue.Job{
			Fn: func() error {
				err := f(w, r)
				return err
			},
			Errc: errc,
		}

		s.requestQueueManager.EnqueueJob(job)

		err := <-errc
		if err != nil {
			fmt.Println(err)
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				WriteJSON(w, httpErr.StatusCode, ApiError{Error: httpErr.Message})
			} else {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "Internal server error"})
			}
		}
	}, middleware.Logging(), middleware.CORS(corsConfig))
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func DecodeJSONBody(r *http.Request, dst interface{}) error {
	defer r.Body.Close()

	var dataMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&dataMap); err != nil {
			if err == io.EOF {
					return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Request body must not be empty"}
			}
			return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Invalid JSON format"}
	}

	validKeys := make(map[string]bool)
	t := reflect.TypeOf(dst).Elem() 
	for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" {
					validKeys[jsonTag] = true
			}
	}

	for key := range dataMap {
			if !validKeys[key] {
					return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Unexpected field in JSON: " + key}
			}
	}

	jsonData, err := json.Marshal(dataMap)
	if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Error processing JSON"}
	}
	if err := json.Unmarshal(jsonData, dst); err != nil {
			return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Invalid request body"}
	}

	return nil
}