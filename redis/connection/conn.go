package connection

import "net"

const (
	// NormalCli is client with user
	NormalCli = iota
	// ReplicationRecvCli is fake client with replication master
	ReplicationRecvCli
)

type Connection struct {
	conn net.Conn

	// waiting until protocol finished
	waiting wait.wait

}