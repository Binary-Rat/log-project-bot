package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type fsm struct {
	client *redis.Client
}

//In Redis the data of user will be stored as hash table with userID as key
//and state, loadV, loadW as values

func New() *fsm {
	return &fsm{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

// SetState sets the state associated with the given userID in Redis.
// It returns an error if something fails while executing the command in Redis.
func (f *fsm) SetState(ctx context.Context, userID string, state string) error {
	return f.client.HSet(ctx, userID, "state", state).Err()
}

// GetState retrieves the state associated with the given userID from Redis.
// It returns the state as a string.
func (f *fsm) GetState(ctx context.Context, userID string) string {
	return f.client.HGet(ctx, userID, "state").Val()
}

// SetLoadV sets the loadV associated with the given userID in Redis.
// It returns an error if something fails while executing the command in Redis.
func (f *fsm) SetLoadV(ctx context.Context, userID string, loadV float32) error {
	return f.client.HSet(ctx, userID, "loadV", loadV).Err()
}

// SetLoadW sets the loadW associated with the given userID in Redis.
// It returns an error if something fails while executing the command in Redis.
func (f *fsm) SetLoadW(ctx context.Context, userID string, loadW float32) error {
	return f.client.HSet(ctx, userID, "loadW", loadW).Err()
}

func (f *fsm) GetLoad(ctx context.Context, userID string) (float32, float32) {
	vals := f.client.HMGet(ctx, userID, "loadV", "loadW").Val()

	loadV, _ := vals[0].(float32)
	loadW, _ := vals[1].(float32)
	return loadV, loadW
}
