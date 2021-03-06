// +build !windows

package api_test

import (
	"os"
	"testing"
	"time"

	"github.com/baidu/openedge/api"
	"github.com/baidu/openedge/config"
	"github.com/baidu/openedge/module"
	"github.com/baidu/openedge/trans/http"
	"github.com/stretchr/testify/assert"
)

func TestAPIUnix(t *testing.T) {
	os.MkdirAll("./var/", 0755)
	defer os.RemoveAll("./var/")
	addr := "unix://./var/test.sock"
	s, err := api.NewServer(&mockEngine{pass: true}, http.ServerConfig{Address: addr, Timeout: time.Minute})
	assert.NoError(t, err)
	defer s.Close()
	err = s.Start()
	assert.NoError(t, err)
	c, err := api.NewClient(http.ClientConfig{Address: addr, Timeout: time.Minute, KeepAlive: time.Minute})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	p, err := c.GetPortAvailable("127.0.0.1")
	assert.NoError(t, err)
	assert.NotZero(t, p)
	err = c.StartModule(&config.Module{Config: module.Config{Name: "name"}})
	assert.NoError(t, err)
	err = c.StopModule("name")
	assert.NoError(t, err)
}

func TestAPIUnixUnauthorized(t *testing.T) {
	os.MkdirAll("./var/", 0755)
	defer os.RemoveAll("./var/")
	addr := "unix://./var/test.sock"
	s, err := api.NewServer(&mockEngine{pass: false}, http.ServerConfig{Address: addr, Timeout: time.Minute})
	assert.NoError(t, err)
	defer s.Close()
	err = s.Start()
	assert.NoError(t, err)
	c, err := api.NewClient(http.ClientConfig{Address: addr, Timeout: time.Minute, KeepAlive: time.Minute, Username: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	_, err = c.GetPortAvailable("127.0.0.1")
	assert.EqualError(t, err, "[400] Account (test) unauthorized")
	err = c.StartModule(&config.Module{Config: module.Config{Name: "name"}})
	assert.EqualError(t, err, "[400] Account (test) unauthorized")
	err = c.StopModule("name")
	assert.EqualError(t, err, "[400] Account (test) unauthorized")
}
