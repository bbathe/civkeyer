package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

type Function struct {
	Label   string
	Message []byte
}

// configuration holds the application configuration
type configuration struct {
	Connection struct {
		Port string
		Baud int
	}
	Functions []Function
}

// implements the Unmarshaler interface of the yaml pkg
// this is so we only have to convert the function strings to bytes once
// and we can have a little friendlier error messages
func (c *configuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	missingSection := "%s section missing from configuration file"
	missingSetting := "%s setting missing from %s section in configuration file"

	// get yaml contents as map
	var cfg map[string]interface{}
	err := unmarshal(&cfg)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// handle connection section
	if cfg["connection"] == nil {
		err = fmt.Errorf(missingSection, "connection")
		log.Printf("%+v", err)
		return err
	}
	connection := cfg["connection"].(map[interface{}]interface{})

	if connection["port"] == nil {
		err = fmt.Errorf(missingSetting, "connection port", "connection")
		log.Printf("%+v", err)
		return err
	}
	c.Connection.Port = connection["port"].(string)

	if connection["baud"] == nil {
		err = fmt.Errorf(missingSetting, "connection baud", "connection")
		log.Printf("%+v", err)
		return err
	}
	c.Connection.Baud = connection["baud"].(int)

	// handle functions section
	if cfg["functions"] == nil {
		err = fmt.Errorf(missingSection, "functions")
		log.Printf("%+v", err)
		return err
	}
	functions := cfg["functions"].([]interface{})

	// max 12 functions, limited by number of function keys on a keyboard for now
	qty := len(functions)
	if qty > 12 {
		err = fmt.Errorf("functions quantity must be less than 12")
		log.Printf("%+v", err)
		return err
	}
	c.Functions = make([]Function, qty)

	// iterate over the defined functions
	for i := 0; i < qty; i++ {
		fn := functions[i].(map[interface{}]interface{})

		// get label
		if fn["label"] == nil {
			err = fmt.Errorf("function %d label must be defined", i+1)
			log.Printf("%+v", err)
			return err
		}
		c.Functions[i].Label = fn["label"].(string)

		// convert CIV message string to bytes
		if fn["message"] == nil {
			err = fmt.Errorf("function %d message must be defined", i+1)
			log.Printf("%+v", err)
			return err
		}
		b, err := hex.DecodeString(fn["message"].(string))
		if err != nil {
			log.Printf("%+v", err)
			return err
		}
		c.Functions[i].Message = b
	}

	return nil
}
