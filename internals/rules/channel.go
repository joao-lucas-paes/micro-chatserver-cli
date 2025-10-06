package rules

import (
	"container/list"
)

type Channel struct {
	Clients []list.List
	msg 		chan Msg
	Id     	string
}


func (c *Channel) ForEachClient(f func(client Client)) {
	for _, clientList := range c.Clients {
		for e := clientList.Front(); e != nil; e = e.Next() {
			f(e.Value.(Client))
		}
	}
}

func (c *Channel) TrySend(msg Msg) bool {
    select {
    case c.msg <- msg:
        return true
    default:
        return false
    }
}

func (c *Channel) GetMsgs() chan Msg {
	return c.msg
}


func NewChannel(id string) Channel {
	return Channel{
		Clients: make([]list.List, 0),
		msg:     make(chan Msg, 64),
		Id:      id,
	}
}

func (c *Channel) AddClient(client Client, index int) {
	if index >= len(c.Clients) {
		c.Clients = append(c.Clients, list.List{})
	}
	c.Clients[index].PushBack(client)
}

func (c *Channel) RemoveClient(client Client) {
	for i := range c.Clients {
		for e := c.Clients[i].Front(); e != nil; e = e.Next() {
			if e.Value.(Client).Nick == client.Nick { 
				c.Clients[i].Remove(e)
				return
			}
		}
	}
}