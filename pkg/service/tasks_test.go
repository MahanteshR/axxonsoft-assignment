package service

import (
	"context"
	"errors"
	"testing"

	"github.com/axxonsoft-assignment/pkg/cache"
	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func intPointer(i int64) *int64 {
	return &i
}

func TestTasks_TasksCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheMock := cache.NewMockCache(ctrl)

	task := New(cacheMock)

	tcs := []struct {
		description string
		taskID      string
		taskDetails model.Task
		resp        *model.TasksResponse
		mockCalls   []*gomock.Call
		expErr      error
	}{
		{
			description: "Positive case: valid request body",
			taskDetails: model.Task{Method: "GET", URL: "https://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					"Content-Type":  "application/json",
					"Authorization": "Basic AFsfr342FA",
				}},
			taskID: "2313",
			mockCalls: []*gomock.Call{
				cacheMock.EXPECT().StoreTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			},
			resp: &model.TasksResponse{ID: "2313"},
		},
		{
			description: "Positive case: valid request body; POST method",
			taskDetails: model.Task{Method: "POST", URL: "https://petstore.swagger.io/v2/pet",
				Data: map[string]interface{}{
					"id": 0,
					"category": map[string]interface{}{
						"id":   0,
						"name": "string",
					},
					"name":      "doggie",
					"photoUrls": []string{"string"},
					"tags": []map[string]interface{}{
						{
							"id":   0,
							"name": "string",
						},
					},
					"status": "available",
				},
			},
			taskID: "2313",
			mockCalls: []*gomock.Call{
				cacheMock.EXPECT().StoreTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			},
			resp: &model.TasksResponse{ID: "2313"},
		},
		{
			description: "Negative case: invalid request body; wrong scheme in url",
			taskDetails: model.Task{Method: "GET", URL: "ftp://www.getyourtasks.com/task"},
			taskID:      "2313",
			expErr:      errors.New("Invalid URL: only the following schemes are supported: [http, https]"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			resp, err := task.TasksCreate(context.TODO(), tc.taskDetails)
			if resp != nil {
				resp.ID = tc.taskID
			}

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}

func TestTasks_TasksGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheMock := cache.NewMockCache(ctrl)

	taskData := &model.TasksObject{
		ID:             "123122",
		Status:         "200",
		HTTPStatusCode: nil,
		Headers:        nil,
		Length:         intPointer(45),
	}

	tcs := []struct {
		description string
		taskID      string
		resp        *model.TasksObject
		mockCalls   []*gomock.Call
		expErr      error
	}{
		{
			description: "Positive case: valid taskID",
			mockCalls: []*gomock.Call{
				cacheMock.EXPECT().GetTask(gomock.Any(), "123122").Return(taskData, nil),
			},
			taskID: "123122",
			resp:   taskData,
		},
		{
			description: "Negative case: error from cache",
			mockCalls: []*gomock.Call{
				cacheMock.EXPECT().GetTask(gomock.Any(), "!@#!").Return(nil, errors.New("DB error")),
			},
			taskID: "!@#!",
			expErr: errors.New("DB error"),
		},
	}

	task := New(cacheMock)

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			resp, err := task.TasksGet(context.TODO(), tc.taskID)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}
