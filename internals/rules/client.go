package rules

import (
	"net"
)

type Client struct {
	idConnection uint64
	conn net.Conn
	nick string
	ch   chan string
}