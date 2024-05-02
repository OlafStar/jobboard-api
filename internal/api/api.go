package api

import (
	"fmt"
	"net/http"

	"github.com/OlafStar/jobboard-api/internal/cache"
	"github.com/OlafStar/jobboard-api/internal/database"
	"github.com/OlafStar/jobboard-api/internal/queue"
	// "github.com/olafstar/salejobs-api/internal/s3"
)

type APIServer struct {
	listenAddr string
	store database.Store
	requestQueueManager *queue.RequestQueueManager
	cache *cache.AllCache
	// s3 *s3.S3Client
}

func NewAPIServer (listenAddr string, db database.Store, rqm *queue.RequestQueueManager, cache *cache.AllCache, 
	) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: db,
		requestQueueManager: rqm,
		cache: cache,
	}
}

func (s *APIServer) Run() {
	fmt.Printf("Memory addresses in APIServer struct:\n")
	fmt.Printf("Address of store: %p\n", &s.store)
	fmt.Printf("Address of requestQueueManager: %p\n", &s.requestQueueManager)
	fmt.Printf("Address of cache: %p\n", &s.cache)


	mux := http.NewServeMux()
	s.SetupCompanyAPI(mux)

	fmt.Printf("Server listening on http://localhost%s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}