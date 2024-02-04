package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/axxonsoft-assignment/pkg/service"
	"github.com/gorilla/mux"
)

type Task struct {
	tasksService service.Tasks
}

func New(tasksService service.Tasks) Tasks {
	return Task{tasksService: tasksService}
}

// CreateTask handles incoming create HTTP requests, process it and returns a taskID or error respectively
func (t Task) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := context.Background()

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)

		return
	}

	var taskData model.Task

	err = json.Unmarshal(reqBody, &taskData)
	if err != nil {
		http.Error(w, "Error in unmarshalling JSON", http.StatusBadRequest)

		return
	}

	resp, err := t.tasksService.TasksCreate(ctx, taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error in marshalling response", http.StatusBadRequest)

		return
	}

	w.Header().Set(model.ContentType, "application/json")

	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, "Error sending JSON response", http.StatusInternalServerError)

		return
	}
}

// GetTask handles incoming get HTTP requests, and returns the data present for that taskID.
func (t Task) GetTask(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := context.Background()

	// get the path param
	vars := mux.Vars(r)
	taskID := vars["taskID"]
	if taskID == "" {
		http.Error(w, "Missing value for the parameter: taskID", http.StatusBadRequest)

		return
	}

	resp, err := t.tasksService.TasksGet(ctx, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error in marshalling response", http.StatusBadRequest)

		return
	}

	w.Header().Set(model.ContentType, "application/json")

	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, "Error sending JSON response", http.StatusInternalServerError)

		return
	}
}
