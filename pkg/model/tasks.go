package model

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// TasksObject represents the structure for returning the task details
type TasksObject struct {
	ID             string      `json:"id,omitempty"`
	Status         string      `json:"status,omitempty"`
	HTTPStatusCode *int        `json:"httpStatusCode,omitempty"`
	Headers        http.Header `json:"headers,omitempty"`
	Length         *int64      `json:"length,omitempty"`
}

// TasksResponse represents the structure of POST response
type TasksResponse struct {
	ID string `json:"id"`
}

// Task represents the structure of a task, used while unmarshalling the request body
type Task struct {
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
	Data    map[string]interface{} `json:"data"`
}

// ValidateRequestBody provides basic validations on the request body like validating the method and url passed in the request body
func ValidateRequestBody(task Task) error {
	if task.Method == "" {
		return errors.New("Invalid request: method cannot be empty")
	}

	// check if the method attribute's value is one of [GET, PATCH, POST, PUT, DELETE]
	task.Method = strings.ToUpper(task.Method)
	if task.Method != http.MethodDelete && task.Method != http.MethodGet && task.Method != http.MethodPost && task.Method != http.MethodPut && task.Method != http.MethodPatch {
		return errors.New("Invalid request: only the following methods are supported: [GET, PATCH, POST, PUT, DELETE]")
	}

	if err := validateURL(task.URL); err != nil {
		return err
	}

	// check to validate that data attribute is passed only for POST/PUT/PATCH
	if checkMethodForBody(task.Method) && task.Data == nil {
		return errors.New("Invalid request: data attribute is required for POST, PUT, or PATCH requests")
	}

	if !checkMethodForBody(task.Method) && task.Data != nil {
		return errors.New("Invalid request: GET/DELETE method doesn't accept data attribute")
	}

	if err := validateHeaders(task); err != nil {
		return err
	}

	return nil
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("Invalid request: url cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.New("Invalid URL in request body: " + err.Error())
	}

	if parsedURL.Scheme == "" || !isValidScheme(parsedURL.Scheme) {
		return errors.New("Invalid URL: only the following schemes are supported: [http, https]")
	}

	return nil
}

func validateHeaders(task Task) error {
	if task.Headers == nil && checkMethodForBody(task.Method) {
		return errors.New("Invalid Request: Content-type header is required for POST/PUT/PATCH")
	}

	if task.Headers != nil && checkMethodForBody(task.Method) {
		contentType, contentTypePresent := task.Headers[ContentType]

		if !contentTypePresent {
			return errors.New("Invalid Request: Content-type header is required for POST/PUT/PATCH")
		}

		if (task.Method == "POST" || task.Method == "PUT") && contentType != "application/json" {
			return errors.New("Invalid Request: Content-type with value application/json is required for POST/PUT")
		}

		if task.Method == "PATCH" && contentType != "application/merge-patch+json" {
			return errors.New("Invalid Request: Content-type with value application/merge-patch+json is required for PATCH")
		}
	}

	return nil
}

func isValidScheme(scheme string) bool {
	return scheme == "http" || scheme == "https"
}

func checkMethodForBody(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}
