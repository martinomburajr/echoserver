package main

import (
	"flag"
	"fmt"
	"github.com/martinomburajr/echoserver/server"
	"log"
	"os"
	"strings"
)

// echo --listener --network="tcp" --address="localhost:8080" || echo -l -n="tcp" -a="localhost:8080"
// Creates a TCP listener on the supplied address

// echo --dialer --network="tcp" --address="localhost:8080" || echo -d -n="tcp" -a="localhost:8080,localhost:8081"
// Creates a TCP dialer on the supplied address

// echo --network="tcp" --listenerAddress="" --dialerAddress=""

func cliParser() (dialer *server.EchoDialer, listener *server.EchoListener, err error) {
	network := flag.String("network", "tcp", "Select the network either 'tcp' or 'udp'")
	isListener := flag.Bool("listener", false, "Select whether the server is a listener")
	isDialer := flag.Bool("dialer", false, "Select whether the server is a dialer")
	listenerAddress := flag.String("listenerAddress", "", "The address to listen to, "+
		"this should include the port e.g. 127.0.0.1:8080")
	dialerAddresses := flag.String("dialerAddress", "", "The addresses to dial to as a comma-separated list")
	dialInterval := flag.Int("dialInterval", 3000, "The dialer dial interval in milliseconds, "+
		"value should be greater than 10 and less than 60000")

	flag.Parse()
	if *network == "" {
		return nil, nil, fmt.Errorf("parser error: network flag cannot be empty")
	}
	if !*isListener && !*isDialer {
		return nil, nil, fmt.Errorf("parser error: please use either --listener or --dialer")
	}
	if *isListener && *isDialer {
		return nil, nil, fmt.Errorf("parser error: please use either --listener or --dialer")
	}
	if *isListener && !*isDialer {
		if *listenerAddress == "" {
			return nil,nil, fmt.Errorf("parser error: listener address cannot be empty")
		}
	}
	if !*isListener && *isDialer {
		if *dialerAddresses == "" {
			return nil,nil, fmt.Errorf("parser error: dialer addresses cannot be empty, " +
				"need at least one valid address e.g. 'localhost:8080', or multiple separated by comma 'localhost:8080," +
				"127.0.0.2:7655'")
		}
	}
	if *listenerAddress == "" && *dialerAddresses == "" {
		return nil, nil,fmt.Errorf("parser error: listener and dialer addresses cannot be empty")
	}
	if *dialInterval < 1 {
		return nil, nil,fmt.Errorf("dialInterval must be greater than 1 and less than 60000")
	}
	if *dialInterval > 60000 {
		return nil,nil, fmt.Errorf("dialInterval must be greater than 1 and less than 60000")
	}

	logger := log.New(os.Stdout, "", log.Lshortfile)
	errChan := make(chan error, 10)

	if *isDialer {
		splitDialerAddresses := strings.Split(*dialerAddresses, ",")
		dialer = &server.EchoDialer{
			Name:      "dialer",
			ErrChan:   errChan,
			Network:   *network,
			Addresses: splitDialerAddresses,
			Interval:  *dialInterval,
			Logger:    logger,
		}

		return dialer, nil, nil
	}
	if *isListener {
		listener = &server.EchoListener{
			Name:    "listener",
			Address: *listenerAddress,
			Network: *network,
			ErrChan: errChan,
			Logger:  logger,
		}

		return nil, listener, nil
	}

	return nil, nil, nil
}

func main() {
	// Parsing CLI
	dialer, listener, err := cliParser()
	if err != nil {
		log.Fatal(err)
	}
	if dialer == nil && listener == nil {
		log.Fatal(fmt.Errorf("parser error: dialer and listener cannot be both unselected"))
	}

	doneChan := make(chan bool, 1)
	if dialer != nil {
		go dialer.Setup()
		log.Println(<-dialer.ErrChan)
	}
	if listener != nil {
		go listener.Setup()
		log.Println(<-listener.ErrChan)
	}
	<-doneChan
}

