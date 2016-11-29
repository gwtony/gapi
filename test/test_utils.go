package test
import (
	"io"
	//"os"
	//"fmt"
	"net/http"
	//"io/ioutil"
	//"path/filepath"
	"net/http/httptest"
	//"github.com/gwtony/gapi/log"
)
// TLog test for log
type TLog struct {
}

// Info log info
func (l *TLog) Info(arg0 interface{}, args ...interface{}) {}
// Debug log debug
func (l *TLog) Debug(arg0 interface{}, args ...interface{}) {}
// Error log error
func (l *TLog) Error(arg0 interface{}, args ...interface{}) {}

// Thandler test http handler
type Thandler struct {
}
func (t *Thandler)ServeHTTP (w http.ResponseWriter, req *http.Request) {
}

// TestNapiConf test napi config
const TestNapiConf string = `
[default]
addr: 172.30.23.39:8888

log: napi.log
level: debug

location: /dns
`

// TestInitlog test init log
func TestInitlog() (* TLog) {
	return &TLog{}
}

// TestGenerateRR generate request and response
func TestGenerateRR(method, uri string, body io.Reader) (*httptest.ResponseRecorder, *http.Request){
	r, _ := http.NewRequest(method, uri, body)
	w := httptest.NewRecorder()
	return w, r
}
