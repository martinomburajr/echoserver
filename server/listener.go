package server

import (
	"bufio"
	"github.com/martinomburajr/echoserver/utils"
	"log"
	"net"
	"strings"
)

type EchoListener struct {
	Name    string
	Address string
	Network string
	ErrChan chan error
	Logger  *log.Logger
}

func (e* EchoListener) Setup() {
	if e.Name == "" {
		e.Name = "listener"
	}
	prefixBuilder := utils.SetupLoggerPrefix(e.Name)
	e.Logger.SetPrefix(prefixBuilder.String())

	listener(e.ErrChan, e.Network, e.Address)
}

func (e* EchoListener) SetupLoggerPrefix() strings.Builder {
	dialerId := utils.RandString(5)
	builder := strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(e.Name)
	builder.WriteString("]:[")
	builder.WriteString(dialerId)
	builder.WriteString("]\t - ")
	return builder
}

func listener(errChan chan error, network, address string) {
	listen, err := net.Listen(network, address)
	if err != nil {
		errChan <- err
	}
	for {
		conn, err := listen.Accept()
		log.Printf("[listener]\t - Accepted Connection at %s from %s", conn.LocalAddr(), conn.RemoteAddr())
		go func() {
			for {

				if err != nil {
					errChan <- err
				}
				connReader := bufio.NewReader(conn)
				line, _, _ := connReader.ReadLine()

				log.Printf("[listener]\t - Read %s", line)
			}
		}()
	}
}
