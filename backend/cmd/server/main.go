package main

import (
	"log"
	"net/http"

	"github.com/prostiate/bank-statement-viewer/internal/handler"
	"github.com/prostiate/bank-statement-viewer/internal/repository"
	"github.com/prostiate/bank-statement-viewer/internal/service"
)

func main() {
	repo := repository.NewTransactionRepository()
	svc := service.NewTransactionService(repo)
	h := handler.NewTransactionHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /upload", h.Upload)
	mux.HandleFunc("GET /balance", h.GetBalance)
	mux.HandleFunc("GET /issues", h.GetIssues)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
