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

func (e *EchoListener) Setup() {
	if e.Name == "" {
		e.Name = "listener"
	}
	prefixBuilder, _ := utils.SetupLoggerPrefix(e.Name)
	e.Logger.SetPrefix(prefixBuilder.String())

	listener(e.ErrChan, e.Network, e.Address, e.Logger)
}

func (e *EchoListener) SetupLoggerPrefix() strings.Builder {
	dialerId := utils.RandString(5)
	builder := strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(e.Name)
	builder.WriteString("]:[")
	builder.WriteString(dialerId)
	builder.WriteString("]\t - ")
	return builder
}

func listener(errChan chan error, network, address string, logger *log.Logger) {
	logger.Printf("Listening on %s://%s ...", network, address)
	listen, err := net.Listen(network, address)
	if err != nil {
		errChan <- err
	} else {
		for {
			conn, err := listen.Accept()
			logger.Printf("Accepted Connection at %s from %s", conn.LocalAddr(), conn.RemoteAddr())
			go func(conn net.Conn, logger *log.Logger) {
				logger.Printf("Reading from connection: %s", conn.RemoteAddr())
				defer conn.Close()
				for {
					if err != nil {
						errChan <- err
						break
					}
					if conn == nil {
						errChan <- err
						break
					}
					connReader := bufio.NewReader(conn)
					line, _, err := connReader.ReadLine()
					if err != nil {
						errChan <- err
						break
					}

					logger.Printf("Read %s | %s", line, conn.RemoteAddr())
				}
			}(conn, logger)
			logger.Println("Listener exited 1")
		}
		logger.Println("Listener exited 2")
	}
}
