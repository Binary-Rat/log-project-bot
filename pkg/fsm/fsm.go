package fsm

import "context"

type Interface interface {
	GetState(ctx context.Context, userID string) string
	SetState(ctx context.Context, userID string, state string) error
	SetLoadW(ctx context.Context, userID string, loadW float64) error
	SetLoadV(ctx context.Context, userID string, loadV float64) error
	GetLoad(ctx context.Context, userID string) (loadV float64, loadW float64)
}
