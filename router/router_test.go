package router

import (
	//"net/http"
	//"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/test"
	"testing"

)

func TestInitRouterOk(t *testing.T) {
//	w, _ := test_generate_rr("GET", "/test", nil)
	log := test.TestInitlog()
	r := InitRouter(log)
	if r == nil {
		t.Fatal("init router failed")
	}
	t.Log("init router done")
}

func TestAddRouterOk(t *testing.T) {
	log := test.TestInitlog()
	r := InitRouter(log)
	r.AddRouter("/test", &test.Thandler{})
	t.Log("add router done")
}
func TestServeHTTPOk(t *testing.T) {
	w, req := test.TestGenerateRR("GET", "/test", nil)
	log := test.TestInitlog()
	r := InitRouter(log)
	r.AddRouter("/test", &test.Thandler{})
	r.ServeHTTP(w, req)
	t.Log("serve http done")
}
