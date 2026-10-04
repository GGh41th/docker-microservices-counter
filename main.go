package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	rdb         *redis.Client
	containerID string
)

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "db-service"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	var err error
	containerID, err = os.Hostname()
	if err != nil {
		containerID = "unknown"
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	hits, err := rdb.Incr(ctx, "hits").Result()
	if err != nil {
		fmt.Fprintf(w, "Bonjour ! Cette page a été vue [erreur Redis: %v] fois. Je suis le conteneur %s.\n", err, containerID)
		return
	}
	fmt.Fprintf(w, "Bonjour ! Cette page a été vue %d fois. Je suis le conteneur %s.\n", hits, containerID)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server running on port 5000...")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
