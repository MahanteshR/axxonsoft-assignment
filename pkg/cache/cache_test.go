package cache

import (
	"context"
	"testing"

	"github.com/axxonsoft-assignment/pkg/model"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func NewCache() *cache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &cache{
		client: client,
	}
}

func TestCache(t *testing.T) {
	c := NewCache()
	ctx := context.Background()
	data := &model.TasksObject{
		ID:     "2321",
		Status: "200",
	}

	err := c.StoreTask(ctx, "2321", data)

	assert.Equal(t, err, nil)

	resp, err := c.GetTask(ctx, "2321")
	assert.Equal(t, err, nil)
	assert.Equal(t, resp, data)

	// Task details do not exist
	_, err = c.GetTask(ctx, "224321")
	assert.Equal(t, err, nil)
	assert.Equal(t, nil, nil)
}
