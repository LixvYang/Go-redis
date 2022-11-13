package server

import (
	"Go-redis/interface/database"
	"Go-redis/lib/sync/atomic"
	"sync"
)

/**
*
 */

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type Handler struct {
	activeConn sync.Map
	db         database.DB
	closing    atomic.Boolean
}

func MakeHandle() *Handler {
	var db database.DB
	db = database2.NewStandaloneServer()
	return &Handler{
		db: db,
	}
}

func (h *Handler) closeClient(client *connection.Connection)
