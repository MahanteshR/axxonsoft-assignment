package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/axxonsoft-assignment/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestTask_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServiceMock := service.NewMockTasks(ctrl)

	testCases := []struct {
		description            string
		reqBody                string
		mockCalls              []*gomock.Call
		expCode                int
		simulateMarshallingErr bool
	}{
		{
			description: "Positive case: valid request",
			reqBody:     `{"method":"GET","url":"https://httpstat.us/200"}`,
			mockCalls: []*gomock.Call{
				taskServiceMock.EXPECT().TasksCreate(gomock.Any(),
					model.Task{
						Method: "GET",
						URL:    "https://httpstat.us/200",
					}).
					Return(&model.TasksResponse{ID: "12323"}, nil),
			},
			expCode: http.StatusOK,
		},
		{
			description: "Negative case: error from service layer",
			reqBody:     `{"method":"PERTH","url":"https://httpstat.us/200"}`,
			mockCalls: []*gomock.Call{
				taskServiceMock.EXPECT().TasksCreate(gomock.Any(),
					model.Task{
						Method: "PERTH",
						URL:    "https://httpstat.us/200",
					}).
					Return(nil, errors.New("error from service layer")),
			},
			expCode: http.StatusBadRequest,
		},
		{
			description: "Negative case: invalid request body",
			reqBody:     `{`,
			expCode:     http.StatusBadRequest,
		},
		{
			description: "Negative case: invalid request body",
			reqBody:     "invalid-json",
			expCode:     http.StatusBadRequest,
		},
	}

	handler := New(taskServiceMock)

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(tc.reqBody))
			w := httptest.NewRecorder()

			handler.CreateTask(w, r)

			assert.Equal(t, tc.expCode, w.Code)
		})
	}
}

func TestTask_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServiceMock := service.NewMockTasks(ctrl)

	testCases := []struct {
		description string
		taskID      string
		mockCalls   []*gomock.Call
		expCode     int
	}{
		{
			description: "Positive case: valid request",
			taskID:      "12323",
			mockCalls: []*gomock.Call{
				taskServiceMock.EXPECT().TasksGet(gomock.Any(), "12323").
					Return(&model.TasksObject{ID: "12323"}, nil),
			},
			expCode: http.StatusOK,
		},
		{
			description: "Negative case: missing taskID",
			expCode:     http.StatusBadRequest,
		},
		{
			description: "Negative case: error from service layer",
			taskID:      "#!@!",
			mockCalls: []*gomock.Call{
				taskServiceMock.EXPECT().TasksGet(gomock.Any(), "#!@!").
					Return(nil, errors.New("error from service layer")),
			},
			expCode: http.StatusBadRequest,
		},
	}

	handler := New(taskServiceMock)

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/task/"+tc.taskID, nil)
			r = mux.SetURLVars(r, map[string]string{"taskID": tc.taskID})
			w := httptest.NewRecorder()

			handler.GetTask(w, r)

			assert.Equal(t, tc.expCode, w.Code)
		})
	}
}
