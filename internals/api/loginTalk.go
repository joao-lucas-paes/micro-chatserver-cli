package api

import (
	"bufio"
	"chatServer/internals/logger"
	"chatServer/internals/rules"
	"chatServer/internals/rules/syncdto"
	"strconv"
)

func LoginTalk(l *logger.Logger, client rules.Client, channels *syncdto.SafeMap[rules.Channel]) (string, string, bool) {
  reader := bufio.NewReader(client.Conn)
  _, err := client.Conn.Write([]byte("Send your nick and the channel that you want to log in to\n"))
  if err != nil {
    l.Errorf("Error to send msg: " + err.Error())
    return "", "", false
  }

  input, err := reader.ReadString('\n');
  if err != nil {
    client.Conn.Write([]byte("<error>"))
    l.Errorf("Error to read msg from client: " + err.Error())
    return "", "", false
  }

  client.Conn.Write([]byte("<ok>"))
  channel := ""
  isLogged := false
  input, channel, isLogged = rules.LoginMatch(input)

  l.Infof("Client response: {msg:" + input + ", channel:" + channel + ", isLogged: " + strconv.FormatBool(isLogged) + "}")
  return input, channel, isLogged
}