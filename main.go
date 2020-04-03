package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
	"unsafe"
)

func main() {
	fmt.Println("EchoServer - Set the DIALER_NETWORK, DIALER_ADDRESS, LISTENER_NETWORK, " +
		"LISTENER_ADDRESS environment variables use this application")

	dialerNetwork := os.Getenv("DIALER_NETWORK")
	if dialerNetwork == "" {
		log.Fatal("DIALER_NETWORK must be set")
	}
	dialerAddress := os.Getenv("DIALER_ADDRESS")
	if dialerAddress == "" {
		log.Fatal("DIALER_ADDRESS must be set")
	}
	listenerNetwork := os.Getenv("LISTENER_NETWORK")
	if listenerNetwork == "" {
		log.Fatal("LISTENER_NETWORK must be set")
	}
	listenerAddress := os.Getenv("LISTENER_ADDRESS")
	if listenerAddress == "" {
		log.Fatal("LISTENER_ADDRESS must be set")
	}

	errChan := make(chan error, 1)
	go dialer(errChan, dialerNetwork, dialerAddress)
	go listener(errChan, listenerNetwork, listenerAddress)

	log.Fatal(<-errChan)
}

func dialer (errChan chan error, network, address string) {
	const interval = 3
	conn, err := net.Dial(network, address)
	if err != nil {
		errChan <- err
	}
	log.Printf("[dialer]\t - Dialing %s on network: %s\n", address, network)
	counter := 1
	for {
		time.Sleep(time.Second * interval)
		randomStr := RandString(6)

		log.Printf("[dialer]\t - %s - Sending %s ==> %s", time.Now().Format(time.RFC3339), randomStr, conn.RemoteAddr())
		dialerMsg := fmt.Sprintf("%s\n", randomStr)

		connectionWriter := bufio.NewWriter(conn)
		connectionWriter.Write([]byte(dialerMsg))
		err := connectionWriter.Flush()
		if err != nil {
			errChan <- err
		}
		counter++
	}
}

func listener (errChan chan error, network, address string) {
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

// Rand String generator from https://stackoverflow.com/a/31832326
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}