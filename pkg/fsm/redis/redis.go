package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type fsm struct {
	client *redis.Client
}

const (
	stateField = "state"
	loadVField = "loadV"
	loadWField = "loadW"
)

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
	return f.client.HSet(ctx, userID, stateField, state).Err()
}

// GetState retrieves the state associated with the given userID from Redis.
// It returns the state as a string.
func (f *fsm) GetState(ctx context.Context, userID string) string {
	return f.client.HGet(ctx, userID, stateField).Val()
}

// SetLoadV sets the loadV associated with the given userID in Redis.
// It returns an error if something fails while executing the command in Redis.
func (f *fsm) SetLoadV(ctx context.Context, userID string, loadV float64) error {
	return f.client.HSet(ctx, userID, loadVField, loadV).Err()
}

// SetLoadW sets the loadW associated with the given userID in Redis.
// It returns an error if something fails while executing the command in Redis.
func (f *fsm) SetLoadW(ctx context.Context, userID string, loadW float64) error {
	return f.client.HSet(ctx, userID, loadWField, loadW).Err()
}

func (f *fsm) GetLoad(ctx context.Context, userID string) (loadV float64, loadW float64) {
	vals := f.client.HMGet(ctx, userID, loadVField, loadWField).Val()
	log.Println(vals)
	v, _ := vals[0].(string)
	w, _ := vals[1].(string)
	loadV, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Println("cannot convert loadV to float")
	}
	loadW, err = strconv.ParseFloat(w, 64)
	if err != nil {
		log.Println("cannot convert loadW to float")
	}
	return loadV, loadW
}
