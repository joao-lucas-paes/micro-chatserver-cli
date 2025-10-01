package rules

import (

)

const (
	maxClients = 1024 // maximo de clientes em cada canal
)

type Channel struct {
	clients []Client
	id string
}