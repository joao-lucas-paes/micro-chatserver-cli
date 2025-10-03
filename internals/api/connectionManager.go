package api

import (
	"chatServer/internals/logger"
	"chatServer/internals/rules"
	"chatServer/internals/rules/syncdto"
)

func connectionSend() {}

func connectionRead(l logger.Logger, clients *syncdto.SafeList[rules.Client], channels *syncdto.SafeMap[rules.Channel]) {
	for {
		for i := 0; i < clients.GetSize(); i++ {
			newClient := clients.PopItem(0)
			go LoginTalk(l, newClient, channels)
		}
	}
}