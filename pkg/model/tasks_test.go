package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasks_ValidateRequestBody(t *testing.T) {
	tcs := []struct {
		description string
		req         Task
		expErr      error
	}{
		{
			description: "Positive case: valid request body",
			req: Task{
				Method: "GET",
				URL:    "https://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
			},
		},
		{
			description: "Negative case: GET method and data passed",
			req: Task{
				Method: "GET",
				URL:    "https://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
				Data: map[string]interface{}{
					"name": "fistName",
				},
			},
			expErr: errors.New("Invalid request: GET/DELETE method doesn't accept data attribute"),
		},
		{
			description: "Negative case: Write call with data attribute not passed",
			req: Task{
				Method: "POST",
				URL:    "https://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
			},
			expErr: errors.New("Invalid request: data attribute is required for POST, PUT, or PATCH requests"),
		},
		{
			description: "Negative case: Empty method",
			req: Task{
				Method: "",
				URL:    "https://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
			},
			expErr: errors.New("Invalid request: method cannot be empty"),
		},
		{
			description: "Negative case: Empty method",
			req: Task{
				Method: "PERTH",
				URL:    "https://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
			},
			expErr: errors.New("Invalid request: only the following methods are supported: [GET, PATCH, POST, PUT, DELETE]"),
		},
		{
			description: "Negative case: Invalid URL",
			req: Task{
				Method: "GET",
				URL:    "ftp://www.getyourtasks.com/task",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
			},
			expErr: errors.New("Invalid URL: only the following schemes are supported: [http, https]"),
		},
		{
			description: "Negative case: Empty URL",
			req: Task{
				Method: "GET",
				URL:    "",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
			},
			expErr: errors.New("Invalid request: url cannot be empty"),
		},
		{
			description: "Negative case: Invalid content type for patch method",
			req: Task{
				Method: "PATCH",
				URL:    "https://www.getyourtasks.com/tas",
				Headers: map[string]interface{}{
					ContentType: "application/json",
				},
				Data: map[string]interface{}{
					"name": "kelly",
				},
			},
			expErr: errors.New("Invalid Request: Content-type with value application/merge-patch+json is required for PATCH"),
		},
		{
			description: "Negative case: missing content type for post method",
			req: Task{
				Method: "POST",
				URL:    "https://www.getyourtasks.com/tas",
				Data: map[string]interface{}{
					"name": "kelly",
				},
			},
			expErr: errors.New("Invalid Request: Content-type header is required for POST/PUT/PATCH"),
		},
		{
			description: "Negative case: missing content type for post method",
			req: Task{
				Method: "PATCH",
				URL:    "https://www.getyourtasks.com/tas",
				Headers: map[string]interface{}{
					ContentType: "json",
				},
				Data: map[string]interface{}{
					"name": "kelly",
				},
			},
			expErr: errors.New("Invalid Request: Content-type with value application/merge-patch+json is required for PATCH"),
		},
		{
			description: "Negative case: missing content type for post method",
			req: Task{
				Method: "PATCH",
				URL:    "https://www.getyourtasks.com/tas",
				Headers: map[string]interface{}{
					"auth": "abcd",
				},
				Data: map[string]interface{}{
					"name": "kelly",
				},
			},
			expErr: errors.New("Invalid Request: Content-type header is required for POST/PUT/PATCH"),
		},
		{
			description: "Negative case: missing content type for post method",
			req: Task{
				Method: "POST",
				URL:    "https://www.getyourtasks.com/tas",
				Headers: map[string]interface{}{
					ContentType: "json",
				},
				Data: map[string]interface{}{
					"name": "kelly",
				},
			},
			expErr: errors.New("Invalid Request: Content-type with value application/json is required for POST/PUT"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			err := ValidateRequestBody(tc.req)

			assert.Equal(t, tc.expErr, err)
		})
	}
}
