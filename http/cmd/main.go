package main

import (
	"net/http"

	"github.com/izveigor/p2p-words/http/pkg/server"
)

func main() {
	router := server.GetRouter()

	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}
}
