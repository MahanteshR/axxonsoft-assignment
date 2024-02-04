package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/axxonsoft-assignment/pkg/cache"
	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/google/uuid"
)

type tasks struct {
	cache  cache.Cache
	client http.Client
}

func New(cache cache.Cache) Tasks {
	return &tasks{cache: cache}
}

// TasksCreate takes the request body, makes the call to third party service and updates the cache respectively.
func (t tasks) TasksCreate(ctx context.Context, taskDetails model.Task) (*model.TasksResponse, error) {
	// validate request body
	if err := model.ValidateRequestBody(taskDetails); err != nil {
		return nil, err
	}

	taskID := uuid.New().String()

	// when a new task is created, its status is "new"
	taskObj := &model.TasksObject{
		ID:     taskID,
		Status: model.New,
	}

	// store the new task details into the cache
	if err := t.cache.StoreTask(ctx, taskID, taskObj); err != nil {
		return nil, err
	}

	// launch a go routine by pass task details to call the 3rd party service
	go func(taskObj *model.TasksObject) {
		var (
			taskBytes []byte
		)

		// create a new http request instance, if failed update the task's status to "error" in the cache
		request, er := http.NewRequest(taskDetails.Method, taskDetails.URL, nil)
		if er != nil {
			log.Printf("Error creating request: %v", er)

			taskObj.Status = model.Error
			if er = t.cache.StoreTask(ctx, taskID, taskObj); er != nil {
				return
			}

			return
		}

		// set headers from the task details
		for key, value := range taskDetails.Headers {
			request.Header.Set(key, fmt.Sprintf("%v", value))
		}

		// in-case the method is POST/PUT/PATCH, fetch the request body
		if taskDetails.Data != nil {
			taskBytes, er = json.Marshal(taskDetails.Data)
			if er != nil {
				log.Printf("Error marshaling JSON")

				taskObj.Status = model.Error
				if er = t.cache.StoreTask(ctx, taskID, taskObj); er != nil {
					return
				}

				return
			}

			body := bytes.NewBuffer(taskBytes)
			request.Body = io.NopCloser(body)
		}

		// make the http call, if failed update the task's status in the cache
		response, er := t.client.Do(request)
		if er != nil {
			log.Printf("Error while calling the 3rd party servicce: %v", er)

			taskObj.Status = model.Error
			if er = t.cache.StoreTask(ctx, taskID, taskObj); er != nil {
				return
			}

			return
		}
		defer response.Body.Close()

		// when the call to 3rd party service is successfully made, update the task's status to "in_process".
		taskObj.Status = model.InProcess
		if er = t.cache.StoreTask(ctx, taskID, taskObj); er != nil {
			return
		}

		_, e := io.ReadAll(response.Body)
		if e != nil {
			log.Printf("Error reading response body: %v", e)

			taskObj.Status = model.Error
			if er = t.cache.StoreTask(ctx, taskID, taskObj); er != nil {
				return
			}

			return
		}

		// update the task's status as per the status of the 3rd party service's call and update the cache respectively
		if er != nil || response.StatusCode != http.StatusOK {
			taskObj.Status = model.Error
		} else {
			taskObj.Status = model.Done
		}

		taskObj.HTTPStatusCode = &response.StatusCode
		taskObj.Length = &response.ContentLength
		taskObj.Headers = response.Header

		if er = t.cache.StoreTask(ctx, taskID, taskObj); er != nil {
			return
		}

		return
	}(taskObj)

	return &model.TasksResponse{ID: taskID}, nil
}

// TasksGet gives the complete task details given a taskID, return an empty object if not found.
func (t tasks) TasksGet(ctx context.Context, taskID string) (*model.TasksObject, error) {
	taskObj, err := t.cache.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return taskObj, nil
}
