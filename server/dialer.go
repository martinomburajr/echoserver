package server

import (
	"bufio"
	"fmt"
	"github.com/martinomburajr/echoserver/utils"
	"log"
	"net"
	"time"
)

type EchoDialer struct {
	Name      string
	ErrChan   chan error
	Network   string
	Addresses []string
	Interval  int
	Logger    *log.Logger
	Id        string
}

func (e *EchoDialer) Setup() {
	if e.Name == "" {
		e.Name = "dialer"
	}
	prefixBuilder, id := utils.SetupLoggerPrefix(e.Name)
	e.Logger.SetPrefix(prefixBuilder.String())
	e.Id = id
	e.Dial()
}

//func (e *EchoDialer) SetupLoggerPrefix() (strings.Builder, string) {
//	dialerId := utils.RandString(5)
//	builder := strings.Builder{}
//	builder.WriteString("[")
//	builder.WriteString(e.Name)
//	builder.WriteString("]:[")
//	builder.WriteString(dialerId)
//	builder.WriteString("]\t - ")
//	return builder, dialerId
//}

// dialer dials a set of addresses on the given Network. If isStdin is true,
// it will only send out data input as stdin on the terminal.
func (e *EchoDialer) Dial() {
	for _, address := range e.Addresses {
		go dial(e.Network, address, e.Id, e.Interval, e.Logger, e.ErrChan)
	}
}

func dial(network, address string, dialerId string, interval int, logger *log.Logger, errChan chan error) {
	shouldRetry := true
	logger.Printf("Attempting to dial...")

	for shouldRetry {
		conn, err := net.Dial(network, address)
		if err != nil {
			errChan <- err
		} else {
			shouldRetry = false
			logger.Printf("Dialing %s on Network: %s\n", address, network)
			counter := 1
			for {
				time.Sleep(time.Millisecond * time.Duration(interval))
				randomStr := utils.RandString(6)

				if conn != nil {
					logger.Printf("Sending %d:'%s' ==> %s", counter, randomStr, conn.RemoteAddr())
					dialerMsg := fmt.Sprintf("[id:%s] - %s\n", dialerId, randomStr)

					connectionWriter := bufio.NewWriter(conn)
					_, err = connectionWriter.Write([]byte(dialerMsg))
					err = connectionWriter.Flush()
					if err != nil {
						errChan <- err
						shouldRetry = true
						break
					}
				} else {
					shouldRetry = true
					break
				}
				counter++
			}
		}
		logger.Printf("Retrying connection to %s://%s ...", network, address)
		time.Sleep(time.Millisecond * 500)
		if conn != nil {
			logger.Printf("Closing connection ...")
			conn.Close()
		}
	}
	logger.Println("Dialer exited")
}
