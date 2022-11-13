package database

import "Go-redis/interface/redis"

type CmdLine = [][]byte

type DB interface {
	Exec(client redis.Connection, cmdLine [][]byte) redis.Reply
	AfterClientClose(c redis.Connection)
	Close()
}
