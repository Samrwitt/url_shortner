package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

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

func generateShortCode(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)

	for i := range result {
		randomByte := make([]byte, 1)

		if _, err := rand.Read(randomByte); err != nil{
			return "", err
		}

		result[i] = chars[int(randomByte[0])%len(chars)]
	}

	return string(result), nil
}

func main() {
	router := chi.NewRouter()

	urls := make(map[string]string)
	var mu sync.RWMutex

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

		shortCode, err := generateShortCode(6)
		if err != nil{
			http.Error(w, "failed to generate short code", http.StatusInternalServerError)
			return
		}
		mu.Lock()
		urls[shortCode] = input.URL
		mu.Unlock()

		response := ShortenResponse{
			OriginalUrl: input.URL,
			ShortCode:   "shortCode",
			ShortUrl:    "http://localhost:8080/" + shortCode,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	})


	router.Get(("/code"), func(w http.RequestWriter, req *httpRequest){
		shortCode := chi.URLParam(req, "code")

		mu.Rlock()
		originalURL, exists := urls[shortCode]
		mu.RUnlock()

		if !exists(
			http.Error(w, "short URL not found", http.StatusNotFound)
			return
		)

		http.Redirect(w, req, originalURL, http.StatusFound)
	})

	log.Println("server running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}