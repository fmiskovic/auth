package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/softika/auth"
)

func main() {
	r := chi.NewRouter()

	// init auth configuration
	cfg := auth.Config{
		Secret: "your-secret",
	}

	// wrap the handler with the auth middleware
	r.Use(auth.Handle(cfg))
	r.Get("/protected", HandleProtectedInfo)

	log.Fatal(http.ListenAndServe(":3000", r))
}

func HandleProtectedInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Protected Info!"))
}
