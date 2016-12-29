package server

import (
	//"time"
	"github.com/gwtony/gapi/userver"
	"github.com/gwtony/gapi/hserver"
	"github.com/gwtony/gapi/tserver"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
	"github.com/gwtony/gapi/log"
)

// Server is A HTTP server
type Server struct {
	haddr   string
	taddr   string
	uaddr   string

	hsch    chan int
	usch    chan int
	tsch    chan int

	hs      *hserver.HttpServer
	us      *userver.UdpServer
	ts      *tserver.TcpServer

	log     log.Log
}

// InitServer inits server
func InitServer(conf *config.Config, log log.Log) (*Server, error) {
	s := &Server{}

	s.log = log

	s.hsch = make(chan int, 1)
	s.usch = make(chan int, 1)
	s.tsch = make(chan int, 1)

	if conf.HttpAddr != "" {
		s.haddr = conf.HttpAddr
		hs, err := hserver.InitHttpServer(conf.HttpAddr, s.log)
		if err != nil {
			s.log.Error("Init http server failed")
				return nil, err
		}
		s.hs = hs
	}

	if conf.UdpAddr != "" {
		s.uaddr = conf.UdpAddr
		us, err := userver.InitUdpServer(conf.UdpAddr, s.log)
		if err != nil {
			s.log.Error("Init udp server failed")
			return nil, err
		}
		s.us = us
	}

	if conf.TcpAddr != "" {
		s.taddr = conf.TcpAddr
		ts, err := tserver.InitTcpServer(conf.TcpAddr, s.log)
		if err != nil {
			s.log.Error("Init tcp server failed")
			return nil, err
		}
		s.ts = ts
	}

	if s.hs == nil && s.us == nil && s.ts == nil {
		s.log.Error("No server inited")
		return nil, errors.InitServerError
	}

	s.log.Debug("Init server done")

	//modules.InitModules(conf, hs, log)

	return s, nil
}

// Run starts server
func (s *Server) Run() error {
	if s.hs != nil {
		go s.hs.Run(s.hsch)
	}
	if s.us != nil {
		go s.us.Run(s.usch)
	}
	if s.ts != nil {
		go s.us.Run(s.tsch)
	}

	//TODO: monitor or something
	select {
		case <-s.hsch:
			s.log.Error("http server run failed")
			break
		case <-s.usch:
			s.log.Error("udp server run failed")
			break
		case <-s.tsch:
			s.log.Error("tcp server run failed")
			break
	}

	return nil
}

//TODO: get more server
func (s *Server) GetHserver() (*hserver.HttpServer) {
	return s.hs
}
