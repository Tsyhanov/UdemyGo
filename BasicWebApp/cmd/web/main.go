package main

import (
	"net/http"
	"udemygo/basicwebapp/pkg/handlers"
)

func main() {

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
