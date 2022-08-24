package server

import (
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, err error) {
	log.Print(err.Error())
	http.Error(w, "Internal Server Error", 500)
	panic(err)
}
