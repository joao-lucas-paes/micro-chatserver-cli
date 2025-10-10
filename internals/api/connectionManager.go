package api

import (
	"bufio"
	"chatServer/internals/logger"
	"chatServer/internals/rules"
	"chatServer/internals/rules/syncdto"
	"fmt"
	"strings"
)

func broadcastMessage(l *logger.Logger, channel *rules.Channel, msg *rules.Msg) {
	for idx := range channel.Clients {
		client := channel.Clients[idx]
		if client.Nick == msg.User().Nick {
			return
		}
		_, err := client.Conn.Write([]byte(msg.User().Nick+":"+msg.Body()))
		if err != nil {
			l.Errorf("Error to send msg to client " + client.Nick + ": " + err.Error())
			return
		}
		l.Infof("Message sent to " + client.Nick + ": " + msg.Body())
	}
}

func connectionSend(l *logger.Logger, channel *rules.Channel) {
	for {
		for msg := range channel.GetMsgs() {
			channel.Mu.Lock()
			broadcastMessage(l, channel, &msg)
			channel.Mu.Unlock()
		}
	}
}

func watcherUser(l *logger.Logger, client *rules.Client, channel *rules.Channel) {
	defer func() {
		if r := recover(); r != nil {
			l.Errorf("panic em watcherUser para %s: %v", client.Nick, r)
		}
		channel.RemoveClient(*client)
		l.Infof("watcherUser finalizado para %s", client.Nick)
	}()

	l.Infof("watcherUser iniciado para %s", client.Nick)

	reader := bufio.NewReader(client.Conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			l.Errorf("Error reading from client %s: %v", client.Nick, err)
			channel.RemoveClient(*client)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		msg := rules.NewMsg(*client, input)

		if !channel.TrySend(*msg) {
			l.Errorf("Channel %s is full, dropping msg from %s: %s", channel.Id, client.Nick, input)
			_, errW := client.Conn.Write([]byte("<error>Buffer is full, try again later</error>\n"))
			if errW != nil {
				l.Errorf("Erro ao enviar erro pro cliente %s: %v", client.Nick, errW)
				channel.RemoveClient(*client)
				return
			}
			continue
		}

		l.Infof("Client %s sent msg: %s", client.Nick, input)
	}
}

func ConnectionRead(l *logger.Logger, clients *syncdto.SafeList[rules.Client], channels *syncdto.SafeMap[rules.Channel]) {
	for {
		for i := 0; i < clients.GetSize(); i++ {
			newClient := clients.PopItem(0)
			nick, channelString, isLogged := LoginTalk(l, newClient, channels)
			l.Infof(fmt.Sprintf("User %s is trying to connect in %s", nick, channelString))

			if isLogged {
				if !channels.ThereIsKey(channelString) {
					newChannel := rules.NewChannel(channelString)
					channels.PushMap(channelString, &newChannel)
					go connectionSend(l, &newChannel) // start channel message dispatcher
					l.Infof("Channel %s created", channelString)
				}
				newClient.Nick = nick
				existingChannel, _ := channels.GetMap(channelString)
				existingChannel.AddClient(newClient)
				channels.PushMap(channelString, existingChannel)
				l.Infof("Client %s logged in channel %s", nick, channelString)
				go watcherUser(l, &newClient, existingChannel) // watch user messages and manage buffer
			} else {
				l.Errorf("Error to login: " + newClient.Nick)
				clients.PushItem(newClient)
				newClient.Conn.Write([]byte("error\n"))
			}
		}
	}
}