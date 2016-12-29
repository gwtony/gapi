package hserver

import (
	//"fmt"
	"io"
	"time"
	"net/http"
	//"io/ioutil"
	//"encoding/json"
	"github.com/gwtony/gapi/variable"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/errors"
	"github.com/gwtony/gapi/router"
)

// HttpServer http server
type HttpServer struct {
	addr        string
	location    string

	router      *router.Router

	log         log.Log
}

var hserver *HttpServer

// InitHttpServer inits http server
func InitHttpServer(addr string, log log.Log) (*HttpServer, error) {
	hs := &HttpServer{}

	hs.addr = addr
	hs.log  = log

	hs.router = router.InitRouter(log)

	return hs, nil
}

// AddRouter adds http server router
func (hs *HttpServer) AddRouter(url string, h http.Handler) error {
	return hs.router.AddRouter(url, h)
}


// Run runs http server
func (hs *HttpServer) Run(ch chan int) error {
	s := &http.Server{
		Addr:           hs.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        hs.router,
	}

	err := s.ListenAndServe()
	ch<-1
	return err
}

// ReturnError return http error
func ReturnError(r *http.Request, w http.ResponseWriter, msg string, err error, log log.Log) {
	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)

	if err == errors.NoContentError {
		// 204 should not return body(RFC) 
		log.Info("Return Error: (%d) %s to client: %s", http.StatusNoContent, msg, r.RemoteAddr)
		http.Error(w, "", http.StatusNoContent)
		return
	}
	if err == errors.BadRequestError {
		log.Info("Return Error: (%d) %s to client: %s", http.StatusBadRequest, msg, r.RemoteAddr)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if err == errors.ForbiddenError {
		log.Info("Return Error: (%d) %s to client: %s", http.StatusForbidden, msg, r.RemoteAddr)
		http.Error(w, msg, http.StatusForbidden)
		return
	}
	if err == errors.BadGatewayError {
		log.Info("Return Error: (%d) %s to client: %s", http.StatusBadGateway, msg, r.RemoteAddr)
		http.Error(w, msg, http.StatusBadGateway)
		return
	}
	if err == errors.ConflictError {
		log.Info("Return Error: (%d) %s to client: %s", http.StatusConflict, msg, r.RemoteAddr)
		http.Error(w, msg, http.StatusConflict)
		return
	}

	log.Info("Return Error: (%d) %s to client: %s", http.StatusInternalServerError, msg, r.RemoteAddr)
	http.Error(w, msg, http.StatusInternalServerError)
}

// ReturnResponse returns http response
func ReturnResponse(r *http.Request, w http.ResponseWriter, msg string, log log.Log) {
	if msg != "" {
		log.Info("Return Ok: (200) %s to client: %s", msg, r.RemoteAddr)
	} else {
		log.Info("Return Ok: (200) to client: %s", r.RemoteAddr)
	}


	if msg == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, msg)
}

