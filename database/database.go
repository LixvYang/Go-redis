package database

import (
	"Go-redis/aof"
	"Go-redis/config"
	"Go-redis/interface/database"
	"Go-redis/interface/redis"
	"Go-redis/lib/logger"
	"Go-redis/pubsub"
	"Go-redis/redis/protocol"
	"fmt"
	"runtime/debug"
	"strings"
	"sync/atomic"
)

// MultiDB is a set of multiple database set
type MultiDB struct {
	dbSet []*atomic.Value // *DB

	// handle publish/subscribe
	hub *pubsub.Hub
	// handle aof persistence
	aofHandler *aof.Handler

	// store master node address
	slaveOf     string
	role        int32
	replication *replicationStatus
}

// NewStandaloneServer creates a standalone redis server, with multi database and all other funtions
func NewStandloneServer() *MultiDB {
	mdb := &MultiDB{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	mdb.dbSet = make([]*atomic.Value, config.Properties.Databases)
	for i := range mdb.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		holder := &atomic.Value{}
		holder.Store(singleDB)
		mdb.dbSet[i] = holder
	}
	mdb.hub = pubsub.MakeHub()
	validAof := false
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAOFHandler(mdb, func() database.EmbedDB {
			return MakeBasicMultiDB()
		})
		if err != nil {
			panic(err)
		}
		mdb.aofHandler = aofHandler
		for _, db := range mdb.dbSet {
			singleDB := db.Load().(*DB)
			singleDB.addAof = func(line CmdLine) {
				mdb.aofHandler.AddAof(singleDB.index, line)
			}
		}
		validAof = true
	}
	mdb.replication = initReplStatus()
	mdb.role = masterRole // The initialization process does not require atomicity
	return mdb
}

// MakeBasicMultiDB create a MultiDB only with basic abilities for aof rewrite and other usages
func MakeBasicMultiDB() *MultiDB {
	mdb := &MultiDB{}
	mdb.dbSet = make([]*atomic.Value, config.Properties.Databases)
	for i := range mdb.dbSet {
		holder := &atomic.Value{}
		holder.Store(makeBasicDB())
		mdb.dbSet[i] = holder
	}
	return mdb
}

func (mdb *MultiDB) Exec(c redis.Connection, cmdLine [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occur: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(cmdLine[0]))
	// authenticate
	if cmdName == "auth" {
		return Auth(c, cmdLine[1:])
	}
	if !isAuthenticated(c) {
		return protocol.MakeErrReply("NOAUTH Authentication required")
	}

	// special commands which cannot execute within transaction
	if cmdName == "subscribe" {
		if len(cmdLine) < 2 {
			return protocol.MakeArgNumErrReply("subscribe")
		}
		return pubsub.Subscribe(mdb.hub, c, cmdLine[1:])
	} else if cmdName == "publish" {
		return pubsub.Publish(mdb.hub, cmdLine[1:])
	} else if cmdName == "unsubscribe" {
		return pubsub.UnSubscribe(mdb.hub, c, cmdLine[1:])
	} else if cmdName == "bgrewriteaof" {
		// aof.go imports router.go, router.go cannot import BGRewriteAOF from aof.go
		return BGRewriteAOF(mdb, cmdLine[1:])
	}

}

// BGRewriteAOF asynchronously rewrites Append-Only-File
func BGRewriteAOF(db *MultiDB, args [][]byte) redis.Reply {
	go db.aofHandler.Rewrite()
	return protocol.MakeStatusReply("Background append only file rewriting started")
}