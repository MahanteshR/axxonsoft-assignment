package handlers

import "net/http"

type Tasks interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetTask(w http.ResponseWriter, r *http.Request)
}
