package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/Binary-Rat/atisu"
	"github.com/go-redis/redis/v8"
)

type fsm struct {
	client *redis.Client
}

const (
	stateField  = "state"
	loadVField  = "loadV"
	loadWField  = "loadW"
	filterField = "filter"
	cityTo      = "cityTo"
	cityFrom    = "cityFrom"
)

//In Redis the data of user will be stored as hash table with userID as key
//and state, loadV, loadW as values

func New() *fsm {
	fsm := &fsm{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
	if fsm.client.Ping(context.Background()).Err() != nil {
		log.Fatal("cannot connect to redis")
	}
	return fsm
}

func (f *fsm) AllUserData(ctx context.Context, userID string) (map[string]string, error) {
	return f.client.HGetAll(ctx, userID).Result()
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

func (f *fsm) SetFilter(ctx context.Context, userID string, filter []byte) error {
	return f.client.HSet(ctx, userID, filterField, filter).Err()
}

// Not implemented. Problem is how to store filter in reddis maybe not reddis?
func (f *fsm) GetFilter(ctx context.Context, userID string) atisu.Filter {
	return atisu.Filter{}
}

func (f *fsm) SetCityFrom(ctx context.Context, userID string, city string) error {
	return f.client.HSet(ctx, userID, cityFrom, city).Err()
}

func (f *fsm) SetCityTO(ctx context.Context, userID string, city string) error {
	return f.client.HSet(ctx, userID, cityTo, city).Err()
}

func (f *fsm) GetRoadCities(ctx context.Context, userID string) []string {
	vals := f.client.HMGet(ctx, userID, cityFrom, cityTo).Val()
	to, _ := vals[0].(string)
	from, _ := vals[1].(string)
	return []string{from, to}
}
