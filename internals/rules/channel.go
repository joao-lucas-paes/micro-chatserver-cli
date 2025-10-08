package rules

type Channel struct {
	Clients []Client
	msg 		chan Msg
	Id     	string
}


func (c *Channel) ForEachClient(f func(client Client)) {
	for _, client := range c.Clients {
		f(client)
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
		Clients: make([]Client, 0),
		msg:     make(chan Msg, 64),
		Id:      id,
	}
}

func (c *Channel) AddClient(client Client) {
	c.Clients = append(c.Clients, client)
}

func (c *Channel) RemoveClient(client Client) {
	for i := range c.Clients {
		if c.Clients[i].Nick == client.Nick { 
			c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
			return
		}
	}
}