package fsm

import (
	"context"

	"github.com/Binary-Rat/atisu"
)

type Interface interface {
	AllUserData(ctx context.Context, userID string) (map[string]string, error)
	GetState(ctx context.Context, userID string) string
	SetState(ctx context.Context, userID string, state string) error
	SetLoadW(ctx context.Context, userID string, loadW float64) error
	SetLoadV(ctx context.Context, userID string, loadV float64) error
	GetLoad(ctx context.Context, userID string) (loadV float64, loadW float64)
	SetFilter(ctx context.Context, userID string, filter []byte) error
	GetFilter(ctx context.Context, userID string) atisu.Filter
	SetCityTO(ctx context.Context, userID string, cityTo string) error
	SetCityFrom(ctx context.Context, userID string, cityFrom string) error
	GetRoadCities(ctx context.Context, userID string) []string
}
