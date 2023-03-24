package client

import "net"
import "strconv"
import "fmt"
import "bufio"
import "log"
import "strings"
import  "errors"

type Client struct {
	Host string
	Port int
}

func (c Client) setupConnection() (net.Conn, error) {
	socketAddress := c.Host + ":" + strconv.Itoa(c.Port)
	log.Println("Attempting to connect to ["+ socketAddress+"]")
	conn, err := net.Dial("tcp", socketAddress)

	if err != nil{
		log.Println("Connection failed. Error: [" + err.Error() + "]")
		return nil, err
	}

	log.Println("Connection established.")
	return conn, nil
}

func (c Client) ExecuteSet(key string, flags int, expiry int, value string) error {
	log.Printf("starting ExecuteSet: key:[%s] flags:[%d] expiry:[%d] value:[%s]\n",key,flags,expiry,value)
	conn, err := c.setupConnection()
	
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("set %s %d %d %d", key, flags, expiry, len(value))
	log.Printf("Memcached command to execute: [%s]\n",cmdStr)

	fmt.Fprintf(conn, "%s\r\n",cmdStr)

	log.Printf("Value: [%s]",value)

	fmt.Fprintf(conn, "%s\r\n",value)
	
	log.Println("Reading from connection")
	status, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil{
		log.Println("Failed to read from connection. Error: [" + err.Error() + "]")
		return err
	}

	log.Println("Received the following response:["+status+"]")

	if strings.TrimSpace(status) != "STORED" {
		return errors.New("failed to set value. "+strings.TrimSpace(status))
	}
	conn.Close()
	return nil
}

func (c Client) ExecuteGet(key string) (string,error) {
	log.Printf("starting ExecutGet: key:[%s]\n",key)
	conn, err := c.setupConnection()
	
	if err != nil {
		return "",err
	}

	cmdStr := fmt.Sprintf("get %s", key)
	log.Printf("Memcached command to execute: [%s]\n",cmdStr)

	fmt.Fprintf(conn, "%s\r\n", cmdStr)
	
	log.Println("Reading from connection")
	response, err := bufio.NewReader(conn).ReadString('\n')
	response = strings.TrimSpace(response)

	if err != nil{
		log.Println("Failed to read from connection. Error: [" + err.Error() + "]")
		conn.Close()
		return "",err
	}

	log.Println("Received the following response:["+response+"]")

	if  response == "END" {
		log.Println("Failed to find the specified key")
		return "", errors.New("specified key not found")
	}
	conn.Close()
	return response,nil
}
