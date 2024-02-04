package routes

import (
	"net/http"

	"github.com/axxonsoft-assignment/pkg/http/handlers"
	"github.com/gorilla/mux"
)

func New(router *mux.Router, handler handlers.Tasks) {
	router.HandleFunc("/task", handler.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/task/{taskID}", handler.GetTask).Methods(http.MethodGet)
}
