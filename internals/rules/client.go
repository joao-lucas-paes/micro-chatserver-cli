package rules

import (
	"net"
)

type Client struct {
	conn net.Conn
	nick string
}