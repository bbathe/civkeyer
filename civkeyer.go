package main

import (
	"bytes"
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

// readCIVMessageFromPort reads bytes from port and returns CIV message as string
func readCIVMessageFromPort(p *serial.Port) ([]byte, error) {
	var buf bytes.Buffer
	b := []byte{0}

	for {
		n, err := p.Read(b)
		if err != nil {
			log.Printf("%+v", err)
			return []byte{0}, err
		}

		if n > 0 {
			// accumulate message bytes
			buf.Write(b)

			// return CIV message
			if b[0] == 0xFD {
				return buf.Bytes(), nil
			}
		}
	}
}

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

		log.Printf("%+v %X", err, b)
		return err
	}

	// check response from radio
	for {
		r, err := readCIVMessageFromPort(port)
		if err != nil {
			log.Printf("%+v %X", err, b)
			return err
		}

		// valid response should be 6 bytes
		if len(r) < 6 {
			err = fmt.Errorf("invalid response from radio")
			log.Printf("%+v %X %X", err, b, r)
			return err
		}

		// message for us from radio?
		if r[2] == '\xE0' && r[3] == '\x94' {
			// check status returned from radio
			if r[4] != '\xFB' {
				err = fmt.Errorf("error response from radio")
				log.Printf("%+v %X %X", err, b, r)
				return err
			}
			break
		}
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
