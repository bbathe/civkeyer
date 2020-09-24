package main

import (
	"fmt"
	"log"

	"github.com/tarm/serial"
)

var (
	// application configuration
	config configuration

	// civ connection
	port *serial.Port

	// user specified configuration file
	configFile string
)

// executeFunction writes the appropriate CIV message associated with function to the configured CIV port
func executeFunction(function int) error {
	var err error

	// connect to CIV port, if we aren't already
	if port == nil {
		c := &serial.Config{
			Name: config.Connection.Port,
			Baud: config.Connection.Baud,
		}

		port, err = serial.OpenPort(c)
		if err != nil {
			err = fmt.Errorf("%+v Port: %s, Baud: %d", err, c.Name, c.Baud)
			log.Printf("%+v", err)
			return err
		}
	}

	// write message bytes to port
	b := config.Functions[function].Message
	_, err = port.Write(b)
	if err != nil {
		// close the port to force a reconnect next time
		port.Close()
		port = nil

		err = fmt.Errorf("%+v Bytes: %X", err, b)
		log.Printf("%+v", err)
		return err
	}

	return nil
}

func main() {
	// show file & location, date & time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err := civkeyerWindow()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	port.Close()
}
