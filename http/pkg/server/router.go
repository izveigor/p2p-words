package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("pkg/server/assets/"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Index(&HTTPServiceClient, w, r)
	})
	router.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		Search(&HTTPServiceClient, w, r)
	})
	router.HandleFunc("/about", About)

	return router
}
