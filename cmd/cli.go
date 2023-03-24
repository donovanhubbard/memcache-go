package main

import "fmt"
import "flag"
import "os"
import "strconv"
import "log"

import "github.com/donovanhubbard/memcache-go/client"

func main(){
	log.Println("Starting program")

	var host string
	var port int
	var command string
	var flags int
	var expiry int
	var key string
	var value string


	flag.StringVar(&host, "host", "localhost", "hostname or ip address of the memcached server")
	flag.IntVar(&port,"port",11211,"memcached port. Default: 11211")
	flag.StringVar(&command,"command","","memcached command to execute: set, get")
	flag.IntVar(&expiry, "expiry", 0, "How many seconds until the key is evicted. 0 never expires. unix timestamp required for values over 30 days")
	flag.IntVar(&flags, "flags", 0, "flags to set alongside the key")
	flag.StringVar(&key, "key", "", "Key to find the value again")
	flag.StringVar(&value, "value", "", "Value to be set")

	flag.Parse()

	if host == "" {
		os.Stderr.WriteString("Missing mandatory argument '-host'\n")
		os.Exit(1)
	}

	if port < 1 || port > 65535{
		os.Stderr.WriteString("Invalid port number '"+strconv.Itoa(port)+"'\n")
		os.Exit(1)
	}

	fmt.Println("host:["+host+"]")
	fmt.Printf("port:[%d]\n",port)

	c := client.Client {Host: host, Port: port}

	var err error 
	var returnedValue string

	if command == "set" {
		err = c.ExecuteSet(key, flags, expiry, value)
	}else if command == "get" {
		returnedValue, err = c.ExecuteGet(key)
	}

	if err != nil {
		log.Println("Command failed")
		os.Exit(2)
	}

	if returnedValue != ""{
		fmt.Println(returnedValue)
	}

	log.Println("Ending program")
}