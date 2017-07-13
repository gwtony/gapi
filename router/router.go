package router

import (
	"net/http"
	"github.com/gwtony/logger"
)

// Router is HTTP router
type Router struct {
	handlers map[string]http.Handler
	log logger.Log
}

// InitRouter inits router
func InitRouter(log logger.Log) *Router {
	r := &Router{}
	r.handlers = make(map[string]http.Handler)
	r.log = log

	return r
}

// AddRouter adds a url router
func (r *Router) AddRouter(url string, handler http.Handler) error {
	if _, ok := r.handlers[url]; ok {
		r.log.Error("Url has been added", logger.String("url", url))
		//TODO: add some error
		return nil
	}
	r.handlers[url] = handler

	return nil
}

// ServeHTTP routers service
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h, ok := r.handlers[req.URL.Path]; ok {
		h.ServeHTTP(w, req)
	} else {
		//if r.NotFound != nil {
		//	r.NotFound.ServeHTTP(w, req)
		//	return
		//}

		//logger.Info.Printf("%s Not Found", req.URL.Path)
		http.Error(w, "URL Not Found", 404)
	}
}


