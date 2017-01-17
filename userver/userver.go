package userver

import (
	//"fmt"
	//"io"
	//"time"
	"net"
	"strings"
	"strconv"
	//"io/ioutil"
	//"encoding/json"
	"github.com/gwtony/gapi/variable"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/errors"
	//"github.com/gwtony/gapi/router"
)

type UdpHandler interface {
	ServUdp([]byte, int)
}
// UdpServer http server
type UdpServer struct {
	ip      net.IP
	port    int
	//nfi     *net.Interface

	handler UdpHandler
	bufSize int

	log     log.Log
}

var userver *UdpServer

// InitUdpServer inits udp server
func InitUdpServer(addr string, log log.Log) (*UdpServer, error) {
	us := &UdpServer{}

	addr_s := strings.Split(addr, ":")
	if len(addr_s) != 2 {
		return nil, errors.InitUdpServerError
	}

	//us.nfi = nil
	//if nfi != "" {
	//	nfi, err := net.InterfaceByName(nfi)
	//	if err != nil {
	//		us.nfi = nfi
	//	}
	//}
	us.ip = net.ParseIP(addr_s[0])
	us.port, _ = strconv.Atoi(addr_s[1])
	us.log  = log
	us.bufSize = variable.UDP_DEFAULT_BUFFER_SIZE

	return us, nil
}

// AddHandler adds udp server handler
func (us *UdpServer) AddHandler(uh UdpHandler) {
	us.handler = uh
}

func (us *UdpServer) SetBuffer(size int) {
	if size > variable.UDP_DEFAULT_BUFFER_SIZE {
		us.bufSize = size
	}
}


// Run runs udp server
func (us *UdpServer) Run(ch chan int) error {
	//TODO: set timeout
	us.log.Debug("udp ip is ", us.ip)
	//if us.nfi != nil {
	//	uc, err := net.ListenMulticastUDP("udp", us.nfi, &net.UDPAddr{IP: us.ip, Port: us.port})
	//} else {
	//	uc, err := net.ListenUDP("udp", &net.UDPAddr{IP: us.ip, Port: us.port})
	//}
	uc, err := net.ListenUDP("udp", &net.UDPAddr{IP: us.ip, Port: us.port})
    if err != nil {
        // handle error
        us.log.Error("Listen udp failed")
		ch<-1
        return err
    }

	buf := make([]byte, us.bufSize)
    for {
        ret, addr, err := uc.ReadFrom(buf)
        if err != nil {
			us.log.Error("Read from %s failed", addr)
            continue
        }
		us.log.Debug("Read %d from address: %s", ret, addr)
        us.handler.ServUdp(buf, ret)
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
//		log.Info("Return Error: (%d) %s to client: %s", http.StatusNoContent, msg, r.RemoteAddr)
//		http.Error(w, "", http.StatusNoContent)
//		return
//	}
//	if err == errors.BadRequestError {
//		log.Info("Return Error: (%d) %s to client: %s", http.StatusBadRequest, msg, r.RemoteAddr)
//		http.Error(w, msg, http.StatusBadRequest)
//		return
//	}
//	if err == errors.ForbiddenError {
//		log.Info("Return Error: (%d) %s to client: %s", http.StatusForbidden, msg, r.RemoteAddr)
//		http.Error(w, msg, http.StatusForbidden)
//		return
//	}
//	if err == errors.BadGatewayError {
//		log.Info("Return Error: (%d) %s to client: %s", http.StatusBadGateway, msg, r.RemoteAddr)
//		http.Error(w, msg, http.StatusBadGateway)
//		return
//	}
//	if err == errors.ConflictError {
//		log.Info("Return Error: (%d) %s to client: %s", http.StatusConflict, msg, r.RemoteAddr)
//		http.Error(w, msg, http.StatusConflict)
//		return
//	}
//
//	log.Info("Return Error: (%d) %s to client: %s", http.StatusInternalServerError, msg, r.RemoteAddr)
//	http.Error(w, msg, http.StatusInternalServerError)
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
//		w.WriteHeader(http.StatusOK)
//		return
//	}
//
//	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)
//	w.WriteHeader(http.StatusOK)
//
//	io.WriteString(w, msg)
//}
//
