package client

import (
	"net"
	"strconv"
	"fmt"
	"bufio"
	"strings"
	"errors"

	"github.com/donovanhubbard/memcache-go/utils"
)

type Client struct {
	Host string
	Port int
}

func (c Client) setupConnection() (net.Conn, error) {
	socketAddress := c.Host + ":" + strconv.Itoa(c.Port)
	utils.Sugar.Debug("Attempting to connect to ["+ socketAddress+"]")
	conn, err := net.Dial("tcp", socketAddress)

	if err != nil{
		utils.Sugar.Error("Connection failed. Error: [" + err.Error() + "]")
		return nil, err
	}

	utils.Sugar.Debug("Connection established.")
	return conn, nil
}

func (c Client) ExecuteSet(key string, flags int, expiry int, value string) error {
	utils.Sugar.Debugf("starting ExecuteSet: key:[%s] flags:[%d] expiry:[%d] value:[%s]\n",key,flags,expiry,value)
	conn, err := c.setupConnection()
	
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("set %s %d %d %d", key, flags, expiry, len(value))
	utils.Sugar.Infof("Memcached command to execute: [%s]\n",cmdStr)

	fmt.Fprintf(conn, "%s\r\n",cmdStr)

	utils.Sugar.Infof("Value: [%s]",value)

	fmt.Fprintf(conn, "%s\r\n",value)
	
	utils.Sugar.Debug("Reading from connection")
	status, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil{
		utils.Sugar.Error("Failed to read from connection. Error: [" + err.Error() + "]")
		return err
	}

	response := strings.TrimSpace(status)

	utils.Sugar.Info("Received the following response:["+response+"]")

	if response != "STORED" {
		err := errors.New("failed to set value. " + response)
		utils.Sugar.Error(err)
		return err
	}
	conn.Close()
	utils.Sugar.Debug("Exiting ExecuteSet")
	return nil
}

func (c Client) ExecuteGet(key string) (string,error) {
	utils.Sugar.Debugf("starting ExecuteGet: key:[%s]\n",key)
	conn, err := c.setupConnection()
	
	if err != nil {
		return "",err
	}

	cmdStr := fmt.Sprintf("get %s", key)
	utils.Sugar.Infof("Memcached command to execute: [%s]\n",cmdStr)

	fmt.Fprintf(conn, "%s\r\n", cmdStr)
	
	utils.Sugar.Debug("Reading from connection")
	response, err := bufio.NewReader(conn).ReadString('\n')
	response = strings.TrimSpace(response)

	if err != nil{
		utils.Sugar.Error("Failed to read from connection. Error: [" + err.Error() + "]")
		conn.Close()
		return "",err
	}

	utils.Sugar.Info("Received the following response:["+response+"]")

	if  response == "END" {
		utils.Sugar.Error("Failed to find the specified key")
		return "", errors.New("specified key not found")
	}
	conn.Close()
	utils.Sugar.Debugf("Exiting ExecuteGet. Returning [%s]", response)
	return response,nil
}
