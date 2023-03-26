package client

import (
	"fmt"
	"github.com/donovanhubbard/memcache-go/utils"
	"testing"
)

const NAMESPACE = "testing.client"
const EXPIRATION = 3

func generateClient() Client {
	return Client{Host: "localhost", Port: 11211}
}

func TestSet(t *testing.T) {
	utils.InitializeLogger()
	c := generateClient()
	testKey := fmt.Sprintf("%s.%s", NAMESPACE, "testing.foo1")
	testValue := "bar"
	err := c.ExecuteSet(testKey, 0, EXPIRATION, testValue)

	if err != nil {
		t.Fatalf("Failed to set key to memcached. %s", err)
	}
}

func TestGet(t *testing.T) {
	var value, retreivedValue string
	var err error
	utils.InitializeLogger()
	c := generateClient()
	testKey := fmt.Sprintf("%s.%s", NAMESPACE, "testing.foo2")
	value = "bar"

	retreivedValue, err = c.ExecuteGet(testKey)

	if err.Error() != "specified key not found" {
		t.Fatalf("Failed to execute the first get. %s", err)
	}

	if retreivedValue != "" {
		t.Fatalf("Got wrong value from a get. Should be empty string but got %s", retreivedValue)
	}

	err = c.ExecuteSet(testKey, 0, EXPIRATION, value)

	if err != nil {
		t.Fatalf("Failed to set key to memcached. %s", err)
	}

	retreivedValue, err = c.ExecuteGet(testKey)

	if err != nil {
		t.Fatalf("Failed to execute the second get. %s", err)
	}

	if value != retreivedValue {
		t.Fatalf("Got wrong value from a get. Should be %s but got %s", value, retreivedValue)
	}
}
