package server

import (
	"github.com/gwtony/gapi/hserver"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/log"
)

// Server is A HTTP server
type Server struct {
	addr    string

	hs      *hserver.HttpServer
	log     log.Log
}

// InitServer inits server
func InitServer(conf *config.Config, log log.Log) (*Server, error) {
	s := &Server{}

	s.log = log
	s.addr = conf.Addr

	hs, err := hserver.InitHttpServer(conf.Addr, s.log)
	if err != nil {
		s.log.Error("Init http server failed")
		return nil, err
	}
	s.hs = hs

	s.log.Debug("Init http server done")

	//modules.InitModules(conf, hs, log)

	return s, nil
}

// Run starts server
func (s *Server) Run() error {
	err := s.hs.Run()
	if err != nil {
		s.log.Error("Server run failed: ", err)
		return err
	}

	return nil
}

func (s *Server) GetHserver() (*hserver.HttpServer) {
	return s.hs
}
