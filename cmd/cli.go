package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/donovanhubbard/memcache-go/client"
	"github.com/donovanhubbard/memcache-go/utils"
)

func main() {
	utils.InitializeLogger()
	utils.Sugar.Info("Starting program")

	var host string
	var port int
	var command string
	var flags int
	var expiry int
	var key string
	var value string

	var err error

	flag.StringVar(&host, "host", "localhost", "hostname or ip address of the memcached server")
	flag.IntVar(&port, "port", 11211, "memcached port. Default: 11211")
	flag.StringVar(&command, "command", "", "memcached command to execute: set, get")
	flag.IntVar(&expiry, "expiry", 0, "How many seconds until the key is evicted. 0 never expires. unix timestamp required for values over 30 days")
	flag.IntVar(&flags, "flags", 0, "flags to set alongside the key")
	flag.StringVar(&key, "key", "", "Key to find the value again")
	flag.StringVar(&value, "value", "", "Value to be set")

	flag.Parse()

	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
		utils.Sugar.Debugf("Using environment variable HOST to set host=[%s]", host)
	}

	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			utils.Sugar.Error("Failed converting environment variable PORT to integer. " + err.Error())
			os.Stderr.WriteString("Invalid port number. Must be an integer")
			os.Exit(1)
		}
		utils.Sugar.Debugf("Using environment variable PORT to set port=[%d]", port)
	}

	utils.Sugar.Debugf("host=[%s]", host)
	utils.Sugar.Debugf("port=[%d]", port)

	if host == "" {
		os.Stderr.WriteString("Missing mandatory argument '-host'\n")
		os.Exit(1)
	}

	if port < 1 || port > 65535 {
		os.Stderr.WriteString("Invalid port number '" + strconv.Itoa(port) + "'\n")
		os.Exit(1)
	}

	c := client.Client{Host: host, Port: port}

	var returnedValue string

	if command == "set" {
		err = c.ExecuteSet(key, flags, expiry, value)
	} else if command == "get" {
		returnedValue, err = c.ExecuteGet(key)

		if err != nil {
			os.Stderr.WriteString("Failed to find key: [" + key + "]\n")
		}
	}

	if err != nil {
		utils.Sugar.Error("Command failed")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	if returnedValue != "" {
		fmt.Println(returnedValue)
	}

	utils.Sugar.Info("Ending program")
}
