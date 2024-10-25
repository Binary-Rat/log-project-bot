package fsm

import "context"

type Interface interface {
	GetState(ctx context.Context, userID string) string
	SetState(ctx context.Context, userID string, state string) error
}
