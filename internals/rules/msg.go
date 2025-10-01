package rules

import (

)

type Msg struct {
	user Client
	body string
}

func (m *Msg) User() Client {
	return m.user
}
func (m *Msg) Body() string {
	return m.body
}

func NewMsg(user Client, body string) *Msg {
	return &Msg{
		user: user,
		body: body,
	}
}