package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/go-redis/redis/v8"
)

// cache represents a client for interacting with a Redis cache.
type cache struct {
	client *redis.Client
}

func New(client *redis.Client) Cache {
	return &cache{client: client}
}

// StoreTask stores the task details into the cache with a TTL of 1 week
func (c cache) StoreTask(ctx context.Context, taskId string, taskObj *model.TasksObject) error {
	data, err := json.Marshal(taskObj)
	if err != nil {
		log.Printf("Error marshalling task object")

		return err
	}

	err = c.client.Set(ctx, taskId, data, 7*24*time.Hour).Err()
	if err != nil {
		log.Printf("Error updating cache for task:%s: %v", taskId, err)

		return err
	}

	return nil
}

// GetTask fetches the task details from redis cache using the taskID.
func (c cache) GetTask(ctx context.Context, taskID string) (*model.TasksObject, error) {
	data, err := c.client.Get(ctx, taskID).Result()
	if err != nil {
		// If taskID is not present, return an empty object
		if err == redis.Nil {
			return &model.TasksObject{}, nil
		}

		log.Printf("Error in fetching the task:%s details from cache: %v", taskID, err)

		return nil, err
	}

	taskObj := &model.TasksObject{}
	err = json.Unmarshal([]byte(data), taskObj)
	if err != nil {
		log.Printf("Error unmarshalling task object")

		return nil, err
	}

	return taskObj, nil
}
