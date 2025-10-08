package api

import (
	"bufio"
	"chatServer/internals/logger"
	"chatServer/internals/rules"
	"chatServer/internals/rules/syncdto"
)

func broadcastMessage(l *logger.Logger, channel *rules.Channel, msg *rules.Msg) {
	channel.ForEachClient(func(client rules.Client) {
		if client.Nick == msg.User().Nick {
			return
		}
		_, err := client.Conn.Write([]byte(msg.Body()))
		if err != nil {
			l.Errorf("Error to send msg to client " + client.Nick + ": " + err.Error())
			return
		}
	})
}

func connectionSend(l *logger.Logger, channel *rules.Channel) {
	for msg := range channel.GetMsgs() {
		broadcastMessage(l, channel, &msg)
	}
}

func watcherUser(l *logger.Logger, client rules.Client, channel *rules.Channel) {
	reader := bufio.NewReader(client.Conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			l.Errorf("Error to read msg from client " + client.Nick + ": " + err.Error())
			channel.RemoveClient(client)
			return
		}
		msg := rules.NewMsg(client, input)
		
		if (!channel.TrySend(*msg)) {
			l.Errorf("Channel " + channel.Id + " is full, dropping msg from " + client.Nick + ": " + input)
			client.Conn.Write([]byte("<error>Buffer is full, try again later</error>\n"))
			return
		}
		_, errSend := client.Conn.Write([]byte("<ok>Message sent</ok>\n"))
		if errSend != nil {
			l.Errorf("Error to send confirmation to client " + client.Nick + ": " + errSend.Error())
			channel.RemoveClient(client)
			return
		}
		l.Infof("Client " + client.Nick + " sent msg: " + input)
	}
}

func ConnectionRead(l *logger.Logger, clients *syncdto.SafeList[rules.Client], channels *syncdto.SafeMap[rules.Channel]) {
	for {
		for i := 0; i < clients.GetSize(); i++ {
			newClient := clients.PopItem(0)
			nick, channelString, isLogged := LoginTalk(l, newClient, channels)

			if isLogged {
				newClient.Conn.Write([]byte("ok"))
				if !channels.ThereIsKey(channelString) {
					newChannel := rules.NewChannel(channelString)
					channels.PushMap(channelString, newChannel)
					go connectionSend(l, &newChannel) // start channel message dispatcher
					l.Infof("Channel %s created", channelString)
					continue
				}
				newClient.Nick = nick
				existingChannel, _ := channels.GetMap(channelString)
				existingChannel.AddClient(newClient)
				channels.PushMap(channelString, existingChannel)
				watcherUser(l, newClient, &existingChannel) // watch user messages and manage buffer
				l.Infof("Client %s logged in channel %s", nick, channelString)
			} else {
				l.Errorf("Error to login: " + newClient.Nick)
				clients.PushItem(newClient)
				newClient.Conn.Write([]byte("error"))
			}
		}
	}
}