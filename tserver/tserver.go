package tserver

import (
	//"fmt"
	//"io"
	//"time"
	"net"
	//"strings"
	//"strconv"
	//"io/ioutil"
	//"encoding/json"
	//"github.com/gwtony/gapi/variable"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/errors"
	//"github.com/gwtony/gapi/router"
)

type TcpHandler interface {
	ServTcp(net.Conn)
}

// TcpServer http server
type TcpServer struct {
	//ip      net.IP
	//port    int
	addr  string

	handler TcpHandler
	//bufSize int

	log   log.Log
}

var tserver *TcpServer

// InitTcpServer inits udp server
func InitTcpServer(addr string, log log.Log) (*TcpServer, error) {
	ts := &TcpServer{}

	ts.addr = addr
	//addr_s := strings.Split(addr, ":")
	//if len(addr_s) != 2 {
	//	return nil, errors.InitTcpServerError
	//}
	//ts.ip = net.ParseIP(addr_s[0])
	//ts.port, _ = strconv.Atoi(addr_s[1])
	ts.log  = log
	//ts.bufSize = variable.UDP_DEFAULT_BUFFER_SIZE

	return ts, nil
}

// AddHandler adds tcp server handler
func (ts *TcpServer) AddHandler(th TcpHandler) {
	ts.handler = th
}

//func (ts.*TcpServer) SetBuffer(size int) {
//	if size > variable.UDP_DEFAULT_BUFFER_SIZE {
//		ts.bufSize = size
//	}
//}


// Run runs tcp server
func (ts *TcpServer) Run(ch chan int) error {
	//TODO: set timeout
	if ts.handler == nil {
		ts.log.Error("No tcp handler")
		ch<-1
		return errors.NoHandlerError
	}

	ln, err := net.Listen("tcp", ts.addr)
	if err != nil {
		ts.log.Error("Listen tcp failed", err)
		ch<-1
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			ts.log.Error("Accept tcp failed", err)
		}
		go ts.handler.ServTcp(conn)
	}

	ch<-0
	return nil
}

//// ReturnError return http error
//func ReturnError(r *http.Request, w http.ResponseWriter, msg string, err error, log log.Log) {
//	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)
//
//	if err == errors.NoContentError {
//		// 204 should not return body(RFC) 
//		log.Info("Return Error: (%d) %s to client: %s", http.Statts.oContent, msg, r.RemoteAddr)
//		http.Error(w, "", http.Statts.oContent)
//		return
//	}
//	if err == errors.BadRequestError {
//		log.Info("Return Error: (%d) %s to client: %s", http.Statts.adRequest, msg, r.RemoteAddr)
//		http.Error(w, msg, http.Statts.adRequest)
//		return
//	}
//	if err == errors.ForbiddenError {
//		log.Info("Return Error: (%d) %s to client: %s", http.Statts.orbidden, msg, r.RemoteAddr)
//		http.Error(w, msg, http.Statts.orbidden)
//		return
//	}
//	if err == errors.BadGatewayError {
//		log.Info("Return Error: (%d) %s to client: %s", http.Statts.adGateway, msg, r.RemoteAddr)
//		http.Error(w, msg, http.Statts.adGateway)
//		return
//	}
//	if err == errors.ConflictError {
//		log.Info("Return Error: (%d) %s to client: %s", http.Statts.onflict, msg, r.RemoteAddr)
//		http.Error(w, msg, http.Statts.onflict)
//		return
//	}
//
//	log.Info("Return Error: (%d) %s to client: %s", http.Statts.nternalServerError, msg, r.RemoteAddr)
//	http.Error(w, msg, http.Statts.nternalServerError)
//}
//
//// ReturnResponse returns http response
//func ReturnResponse(r *http.Request, w http.ResponseWriter, msg string, log log.Log) {
//	if msg != "" {
//		log.Info("Return Ok: (200) %s to client: %s", msg, r.RemoteAddr)
//	} else {
//		log.Info("Return Ok: (200) to client: %s", r.RemoteAddr)
//	}
//
//
//	if msg == "" {
//		w.WriteHeader(http.Statts.K)
//		return
//	}
//
//	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)
//	w.WriteHeader(http.Statts.K)
//
//	io.WriteString(w, msg)
//}
//
