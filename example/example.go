package main

import (
	"io"
	"fmt"
	"github.com/gwtony/gapi/api"
)

type Handler struct {
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}


func main() {
	err := api.Init()
	if err != nil {
		fmt.Println("Init api failed")
		return
	}
	config := api.GetConfig()
	log := api.GetLog()

	api.AddHandler("/test", &Handler{})
	log.Info("Add handler done")

	api.Run()
}
