package redis_test

import (
	"log-proj/pkg/fsm/redis"
	"testing"
)

func TestMain(m *testing.M) {
	redis.New()

	m.Run()
}

func Test(t *testing.T) {

}
