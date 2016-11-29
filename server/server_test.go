package server

import (
	"testing"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/test"

)

func TestInitServerOK(t *testing.T) {
	conf := &config.Config{}
	conf.Addr = "localhost:80"
	log := test.TestInitlog()
	InitServer(conf, log)
	t.Log("init server done")
}
