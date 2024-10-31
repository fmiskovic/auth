package stdlib

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// init auth middleware
	//a := auth.New(auth.Config{
	//	Secret: "",
	//})

	mux.HandleFunc("/protected", HandleProtectedInfo)

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func HandleProtectedInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Protected Info!"))
}
