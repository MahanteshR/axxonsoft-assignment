package service

import (
	"context"
	"github.com/axxonsoft-assignment/pkg/model"
)

type Tasks interface {
	TasksCreate(ctx context.Context, body model.Task) (*model.TasksResponse, error)
	TasksGet(ctx context.Context, taskID string) (*model.TasksObject, error)
}
