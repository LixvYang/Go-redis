package database

import "Go-redis/interface/redis"

type CmdLine = [][]byte

type DB interface {
	Exec(client redis.Connection, cmdLine [][]byte) redis.Reply
	AfterClientClose(c redis.Connection)
	Close()
}



// DataEntity stores data bound to a key, including a string, list, hash, set and so on
type DataEntity struct {
	Data interface{}
}
