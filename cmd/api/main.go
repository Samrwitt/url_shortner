package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)


type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortCode string `json:"short_code"`
	ShortUrl string `json:"short_url"`
}



func main() {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	router.Post("/api/shorten", func(w http.ResponseWriter, req *http.Request) {
		var input ShortenRequest

		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		response := ShortenResponse{
			OriginalUrl: input.URL,
			ShortCode:   "abc123",
			ShortUrl:    "http://localhost:8080/abc123",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	})


	log.Println("server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}