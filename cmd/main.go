package main

import (
	"fmt"

	"github.com/OlafStar/jobboard-api/internal/api"
	"github.com/OlafStar/jobboard-api/internal/cache"
	"github.com/OlafStar/jobboard-api/internal/database"
	"github.com/OlafStar/jobboard-api/internal/queue"

	_ "github.com/go-sql-driver/mysql"
)

// "github.com/olafstar/salejobs-api/internal/api"
// "github.com/olafstar/salejobs-api/internal/s3"

func main(){
	db := database.InitDatabase()
	// s3 := s3.InitS3Client()

	rqm := queue.NewRequestQueueManager(10, 10)
	defer rqm.Shutdown()

	c := cache.NewCache()

	server := api.NewAPIServer(":4200", db, rqm, c, 
	// s3
	)
	fmt.Printf("Address of store: %p\n", &db)
	fmt.Printf("Address of requestQueueManager: %p\n", &rqm)
	fmt.Printf("Address of cache: %p\n", &c)
	server.Run()
}