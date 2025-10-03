package rules

import (
	"chatServer/internals/rules/syncdto"
	"container/list"
)

type Channel struct {
	Clients []list.List
	MsgList *syncdto.SafeList[Msg]
	Id string
}