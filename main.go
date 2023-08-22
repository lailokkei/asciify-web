package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := "8080"
	publicDir := "public"
	fs := http.FileServer(http.Dir(publicDir))

	http.Handle("/", NoCache(fs))

	http.HandleFunc("/connect", connect)

	fmt.Println("serving on : http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
