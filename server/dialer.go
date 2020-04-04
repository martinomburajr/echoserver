package server

import (
	"bufio"
	"fmt"
	"github.com/martinomburajr/echoserver/utils"
	"log"
	"net"
	"strings"
	"time"
)

type EchoDialer struct {
	Name      string
	ErrChan   chan error
	Network   string
	Addresses []string
	Interval  int
	Logger    *log.Logger
}

func (e *EchoDialer) Setup() {
	if e.Name == "" {
		e.Name = "dialer"
	}
	prefixBuilder := utils.SetupLoggerPrefix(e.Name)
	e.Logger.SetPrefix(prefixBuilder.String())

	e.Dial()
}

func (e* EchoDialer) SetupLoggerPrefix() strings.Builder {
	dialerId := utils.RandString(5)
	builder := strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(e.Name)
	builder.WriteString("]:[")
	builder.WriteString(dialerId)
	builder.WriteString("]\t - ")
	return builder
}

// dialer dials a set of addresses on the given Network. If isStdin is true,
// it will only send out data input as stdin on the terminal.
func (e *EchoDialer) Dial() {
	for _, address := range e.Addresses {
		go dial(e.Network, address, e.Interval, e.Logger, e.ErrChan)
	}
}

func dial(network, address string, interval int, logger *log.Logger, errChan chan error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		errChan <- err
	}
	logger.Printf("Dialing %s on Network: %s\n", address, network)
	counter := 1
	for {
		time.Sleep(time.Millisecond * time.Duration(interval))
		randomStr := utils.RandString(6)

		if conn != nil {
			logger.Printf("Sending '%s' ==> %s", randomStr, conn.RemoteAddr())
			dialerMsg := fmt.Sprintf("%s\n", randomStr)

			connectionWriter := bufio.NewWriter(conn)
			connectionWriter.Write([]byte(dialerMsg))
			err := connectionWriter.Flush()
			if err != nil {
				errChan <- err
			}
		}
		counter++
	}
}
