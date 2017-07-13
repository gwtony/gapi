package usocket

import (
	//"fmt"
	//"io"
	//"time"
	"os"
	"net"
	//"strings"
	//"strconv"
	//"io/ioutil"
	//"encoding/json"
	"github.com/gwtony/gapi/variable"
	"github.com/gwtony/logger"
	//"github.com/gwtony/gapi/errors"
	//"github.com/gwtony/gapi/router"
)

type UsocketHandler interface {
	ServUsocket([]byte, int)
}
// UsocketServer http server
type UsocketServer struct {
	//ip      net.IP
	//port    int
	//nfi     *net.Interface
	socket  string

	handler UsocketHandler
	bufSize int

	log     logger.Log
}

var usocket *UsocketServer

// InitUsocketServer inits usocket server
func InitUsocketServer(addr string, log logger.Log) (*UsocketServer, error) {
	us := &UsocketServer{}

	//addr_s := strings.Split(addr, ":")
	//if len(addr_s) != 2 {
	//	return nil, errors.InitUsocketServerError
	//}
	us.socket = addr

	//us.ip = net.ParseIP(addr_s[0])
	//us.port, _ = strconv.Atoi(addr_s[1])
	us.log  = log
	us.bufSize = variable.USOCK_DEFAULT_BUFFER_SIZE

	return us, nil
}

// AddHandler adds udp server handler
func (us *UsocketServer) AddHandler(uh UsocketHandler) {
	us.handler = uh
}

func (us *UsocketServer) SetBuffer(size int) {
	if size > variable.USOCK_DEFAULT_BUFFER_SIZE {
		us.bufSize = size
	}
}


// Run runs udp server
func (us *UsocketServer) Run(ch chan int) error {
	//TODO: set timeout
	us.log.Debug("usocket file is ", logger.String("socket", us.socket))
	addr, err := net.ResolveUnixAddr("unixgram", us.socket)
	if err != nil {
		us.log.Error("Cannot resolve unix addr", logger.Err(err))
		ch<-1
		return err
	}
	os.Remove(us.socket)
	uc, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		us.log.Error("Cannot listen to unix domain socket", logger.Err(err))
		ch<-1
		return err
	}
	defer func() {
		uc.Close()
		os.Remove(us.socket)
	}()
	//uc, err := net.ListenUDP("udp", &net.UDPAddr{IP: us.ip, Port: us.port})
    if err != nil {
        // handle error
        us.log.Error("Listen udp failed")
		ch<-1
        return err
    }
	os.Chmod(us.socket, 0777)

	buf := make([]byte, us.bufSize)
    for {
		//this will cause big memory leak
		//buf := make([]byte, us.bufSize)
        ret, addr, err := uc.ReadFrom(buf)
        if err != nil {
			us.log.Error("Read from client failed", logger.String("remote", addr.String()))
            continue
        }
		us.log.Debug("Read from address success", logger.Int("size", ret), logger.String("socket", us.socket))

		//Need not goroutine, or buf will be confused
        us.handler.ServUsocket(buf, ret)
    }

	ch<-0
	return nil
}
