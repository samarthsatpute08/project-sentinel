package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"source":"primary","path":"%s"}`, r.URL.Path)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	http.ListenAndServe(":"+port, nil)
}