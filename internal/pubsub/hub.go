package pubsub

import (
	"github.com/lixvyang/Go-redis/internal/datastruct/dict"
	"github.com/lixvyang/Go-redis/internal/datastruct/lock"
)

// Hub stores all subscribe relations
type Hub struct {
	// channel -> list(*Client)
	subs dict.Dict
	// lock channel
	subsLocker *lock.Locks
}

// MakeHub creates new hub
func MakeHub() *Hub {
	return &Hub{
		subs:       dict.MakeConcurrent(4),
		subsLocker: lock.Make(16),
	}
}
