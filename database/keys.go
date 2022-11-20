package database

import (
	"Go-redis/interface/redis"
	"Go-redis/lib/utils"
	"Go-redis/lib/wildcard"
	"Go-redis/redis/protocol"
	"strconv"
	"strings"
	"time"
)

// execKeys returns all keys matching the given pattern
func execKeys(db *DB, args [][]byte) redis.Reply {
	pattern, err := wildcard.CompilePattern(string(args[0]))
	if err != nil {
		return protocol.MakeErrReply("ERR illegal wildcard")
	}
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return protocol.MakeMultiBulkReply(result)
}

// execCopy usage: COPY source destination [DB destination-db] [REPLACE]
// This command copies the value stored at the source key to the destination key.
func execCopy(mdb *MultiDB, conn redis.Connection, args [][]byte) redis.Reply {
	dbIndex := conn.GetDBIndex()
	db := mdb.mustSelectDB(dbIndex) // Current DB
	replaceFlag := false
	srcKey := string(args[0])
	destKey := string(args[1])

	// Parse options
	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			arg := strings.ToLower(string(args[i]))
			if arg == "db" {
				if i+1 >= len(args) {
					return &protocol.SyntaxErrReply{}
				}
				idx, err := strconv.Atoi(string(args[i+1]))
				if err != nil {
					return &protocol.SyntaxErrReply{}
				}
				if idx >= len(mdb.dbSet) || idx < 0 {
					return protocol.MakeErrReply("ERR DB index is out of range")
				}
				dbIndex = idx
				i++
			} else if arg == "replace" {
				replaceFlag = true
			} else {
				return &protocol.SyntaxErrReply{}
			}
		}
	}

	if srcKey == destKey && dbIndex == conn.GetDBIndex() {
		return protocol.MakeErrReply("ERR source and destination objects are the same")
	}

	// source key does not exist
	src, exists := db.GetEntity(srcKey)
	if !exists {
		return protocol.MakeIntReply(0)
	}

	destDB := mdb.mustSelectDB(dbIndex)
	if _, exists = destDB.GetEntity(destKey); exists != false {
		// If destKey exists and there is no "replace" option
		if replaceFlag == false {
			return protocol.MakeIntReply(0)
		}
	}

	destDB.PutEntity(destKey, src)
	raw, exists := db.ttlMap.Get(srcKey)
	if exists {
		expire := raw.(time.Time)
		destDB.Expire(destKey, expire)
	}
	mdb.aofHandler.AddAof(conn.GetDBIndex(), utils.ToCmdLine3("copy", args...))
	return protocol.MakeIntReply(1)
}
