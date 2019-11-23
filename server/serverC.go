package main

import (
	"fmt"
	"net/http"
)

var (
	port = "8085"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":"+port, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from port %s!", port)
}
