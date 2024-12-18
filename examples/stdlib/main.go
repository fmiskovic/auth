package main

import (
	"log"
	"net/http"

	"github.com/softika/auth"
)

func main() {
	mux := http.NewServeMux()

	// init auth middleware
	a := auth.New(auth.Config{
		Secret: "your-secret",
	})

	// wrap the handler with the auth middleware
	mux.Handle("/protected", a.Handler(http.HandlerFunc(HandleProtectedInfo)))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func HandleProtectedInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Protected Info!"))
}
