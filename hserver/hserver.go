package hserver

import (
	//"fmt"
	"io"
	"time"
	"net/http"
	//"io/ioutil"
	//"encoding/json"
	"github.com/gwtony/gapi/variable"
	"github.com/gwtony/logger"
	"github.com/gwtony/gapi/errors"
	"github.com/gwtony/gapi/router"
)

// HttpServer http server
type HttpServer struct {
	addr        string
	location    string

	router      *router.Router

	log         logger.Log
}

var hserver *HttpServer

// InitHttpServer inits http server
func InitHttpServer(addr string, log logger.Log) (*HttpServer, error) {
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
func ReturnError(r *http.Request, w http.ResponseWriter, msg string, err error, log logger.Log) {
	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)

	if err == errors.NoContentError {
		// 204 should not return body(RFC) 
		log.Info("Return error to client", logger.Int("status", http.StatusNoContent), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, "", http.StatusNoContent)
		return
	}
	if err == errors.BadRequestError {
		log.Info("Return error to client", logger.Int("status", http.StatusBadRequest), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if err == errors.ForbiddenError {
		log.Info("Return error to client", logger.Int("status", http.StatusForbidden), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, msg, http.StatusForbidden)
		return
	}
	if err == errors.BadGatewayError {
		log.Info("Return error to client", logger.Int("status", http.StatusBadGateway), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, msg, http.StatusBadGateway)
		return
	}
	if err == errors.ConflictError {
		log.Info("Return error to client", logger.Int("status", http.StatusConflict), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, msg, http.StatusConflict)
		return
	}
	if err == errors.UnauthorizedError {
		log.Info("Return error to client", logger.Int("status", http.StatusUnauthorized), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	if err == errors.NotAcceptableError {
		log.Info("Return error to client", logger.Int("status", http.StatusNotAcceptable), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
		http.Error(w, msg, http.StatusNotAcceptable)
		return
	}

	log.Info("Return error to client", logger.Int("status", http.StatusInternalServerError), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
	http.Error(w, msg, http.StatusInternalServerError)
}

// ReturnResponse returns http response
func ReturnResponse(r *http.Request, w http.ResponseWriter, msg string, log logger.Log) {
	if msg != "" {
		log.Info("Return ok: (200) to client", logger.Int("status", 200), logger.String("message", msg), logger.String("remote", r.RemoteAddr))
	} else {
		log.Info("Return ok: (200) to client", logger.Int("status", 200), logger.String("remote", r.RemoteAddr))
	}

	if msg == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", variable.DEFAULT_CONTENT_HEADER)
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, msg)
}

