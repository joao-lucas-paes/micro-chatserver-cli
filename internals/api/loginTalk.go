package api

import (
	"bufio"
	"chatServer/internals/logger"
	"chatServer/internals/rules"
	"chatServer/internals/rules/syncdto"
	"strconv"
)

func LoginTalk(l logger.Logger, client rules.Client, channels *syncdto.SafeMap[rules.Channel]) {
		reader := bufio.NewReader(client.Conn)
    _, err := client.Conn.Write([]byte("Send your nick and the channel that you want to log in to\n"))
    if err != nil {
        l.Println("Error to send msg: " + err.Error())
        return
    }

    input, err := reader.ReadString('\n');
    if err != nil {
        l.Println("Error to read msg from client: " + err.Error())
        return
    }

    channel := -1
		isLogged := false
    input, channel, isLogged = rules.LoginMatch(input)

    l.Println("Client response: {msg:" + input + ", channel:" + strconv.Itoa(channel) + ", isLogged: " + strconv.FormatBool(isLogged) + "}")

}