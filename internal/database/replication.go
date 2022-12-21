package database

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/lixvyang/Go-redis/internal/redis/parser"
)

const (
	masterRole = iota
	slaveRole
)

type replicationStatus struct {
	mutex  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	// configVersion stands for the version of replication config. Any change of master host/port will cause configVersion increment
	// If configVersion change has been found during replication current replication procedure will stop.
	// It is designed to abort a running replication procedure
	configVersion int32

	masterHost string
	masterPort int

	masterConn   net.Conn
	masterChan   <-chan *parser.Payload
	replId       string
	replOffset   int64
	lastRecvTime time.Time
	running      sync.WaitGroup
}

var configChangedErr = errors.New("replication config changed")

func initReplStatus() *replicationStatus {
	repl := &replicationStatus{}
	// start cron
	return repl
}
