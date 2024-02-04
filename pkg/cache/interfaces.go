package cache

import (
	"context"
	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache interface {
	StoreTask(ctx context.Context, taskId string, taskObj *model.TasksObject) error
	GetTask(ctx context.Context, taskID string) (*model.TasksObject, error)
}

// Client interface for mocking redis client
type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}
