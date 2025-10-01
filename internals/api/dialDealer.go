package api

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"chatServer/internals/logger"
)

const (
	ctrlAddr           = "127.0.0.1:6000"
	bufferReadTimeout  = 30 * time.Second
	initialPort        = 32000
	maxPortAttempts    = 10000
	protocolHello      = "hello"
	protocolConfirmOK  = "ok"
	msgTryAgain        = "WAS NOT POSSIBLE TO RESERVE SOME PORT TO YOU, TRYING AGAIN\n"
	msgBadCommunication = "BAD COMMUNICATION THIS PORT ONLY SHOULD BE USED TO DEALER SERVER\n"
)

func dealer(l logger.Logger, listConn *list.List) error {
	ctrlLn, err := net.Listen("tcp", ctrlAddr)
	if err != nil {
		l.Errorf("Failed to listen control addr %s: %v", ctrlAddr, err)
		return err
	}
	l.Infof("Dealer is starting on %s", ctrlAddr)
	defer ctrlLn.Close()

	var portMu sync.Mutex
	lastPort := initialPort

	var listMu sync.Mutex

	for {
		conn, err := ctrlLn.Accept()
		if err != nil {
			l.Errorf("Accept error: %v", err)
			continue
		}
		l.Infof("Control connection from %s", conn.RemoteAddr())

		go func(c net.Conn) {
			defer c.Close()
			_ = c.SetDeadline(time.Now().Add(bufferReadTimeout))
			r := bufio.NewReader(c)

			msg, err := r.ReadString('\n')
			if err != nil {
				l.Errorf("Error reading initial msg: %v", err)
				return
			}
			msg = strings.TrimSpace(msg)
			l.Infof("Received control msg: %s", msg)

			if msg != protocolHello {
				l.Errorf("Bad communication: %s", msg)
				_, _ = c.Write([]byte(msgBadCommunication))
				return
			}

			attempts := 0
			for {
				if attempts > maxPortAttempts {
					l.Errorf("Exceeded max port attempts")
					return
				}
				attempts++

				portMu.Lock()
				port := lastPort
				lastPort++
				portMu.Unlock()

				addr := fmt.Sprintf("127.0.0.1:%d", port)
				ln, listenErr := net.Listen("tcp", addr)
				if listenErr != nil {
					l.Errorf("Port %d not available: %v", port, listenErr)
					_, _ = c.Write([]byte(msgTryAgain))
					continue
				}

				_, writeErr := c.Write([]byte(fmt.Sprintf("%s\n", addr)))
				if writeErr != nil {
					l.Errorf("Error writing port to control connection: %v", writeErr)
					ln.Close() // ta com vazamento, o close garante funcionamento
					return
				}

				// espera confirmação
				_ = c.SetDeadline(time.Now().Add(bufferReadTimeout))
				confirm, err := r.ReadString('\n')
				if err != nil {
					l.Errorf("Error reading confirmation: %v", err)
					ln.Close()
					return
				}
				confirm = strings.TrimSpace(confirm)
				if confirm == protocolConfirmOK {
					// salvo o listener protegido por mutex
					listMu.Lock()
					listConn.PushBack(ln)
					listMu.Unlock()
					l.Infof("Reserved and stored listener on %s", addr)
					break
				} else {
					l.Errorf("Client did not confirm (got %q), closing listener %s", confirm, addr)
					ln.Close()
					// tenta dnv
					continue
				}
			}
		}(conn)
	}
}