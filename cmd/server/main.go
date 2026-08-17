package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/api"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/repository"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/service"
)

func main() {
	store := repository.NewMemoryStore()
	dispatch := service.NewDispatchService(store)
	handler := api.NewHandler(dispatch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("courier slot coordinator listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
