package macedon

import (
	"net"
	"fmt"
	"time"
	"math/rand"
	"strings"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

// AddHandler Add record handler
type AddHandler struct {
	h      *Handler
	domain string
	log    log.Log
}

// DeleteHandler Delete record handler
type DeleteHandler struct {
	h      *Handler
	domain string
	pc     *PurgeContext
	log    log.Log
}

// ReadHandler Read record handler
type ReadHandler struct {
	h      *Handler
	domain string
	log    log.Log
}

// ReadServerHandler Read server record handler
type ReadServerHandler struct {
	h   *Handler
	log log.Log
}

// AddServerHandler Add server record handler
type AddServerHandler struct {
	h   *Handler
	pc  *PurgeContext
	log log.Log
}

// DeleteServerHandler Delete server record handler
type DeleteServerHandler struct {
	h   *Handler
	pc  *PurgeContext
	log log.Log
}

// ServeHTTP router interface
func (handler *AddHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var isArpa bool
	var arr []string
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Add record request: ", data)

	/* Check input */
	if data.Name == "" || data.Address == "" {
		handler.log.Error("Name or address invalid")
		http.Error(w, "Name or address invalid", http.StatusBadRequest)
		return
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		handler.log.Error("Name invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = strings.Split(data.Name, ".")
		total := len(arr)
		for i := 0; i < total / 2; i++ {
			tmp := arr[i]
			arr[i] = arr[total - 1 - i]
			arr[total - 1 - i] = tmp
		}
	}
	rec := strings.Join(arr, "/")

	if !isArpa {
		rec = rec + "/" + fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))
	}
	handler.log.Debug("Rec is %s", rec)

	_, err = handler.h.Add(rec, data.Address, data.Ttl, isArpa, false)

	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, nil, handler.log)
}

// ServeHTTP router interface
func (handler *DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var isArpa bool
	var arr []string

	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Delete record request: %s", data)

	/* Check input */
	if data.Name == "" || data.Address == "" {
		handler.log.Error("Name or address invalid")
		http.Error(w, "Name or address invalid", http.StatusBadRequest)
		return
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		handler.log.Error("Name invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = strings.Split(data.Name, ".")
		total := len(arr)
		for i := 0; i < total / 2; i++ {
			tmp := arr[i]
			arr[i] = arr[total - 1 - i]
			arr[total - 1 - i] = tmp
		}
	}
	rec := strings.Join(arr, "/")

	handler.log.Debug("Rec is %s", rec)

	resp, err := handler.h.Read(rec, isArpa, false)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	var found []string
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
			json.Unmarshal([]byte(v.Value), &respRec)
			if strings.Compare(respRec.Host, data.Address) == 0 {
				found = append(found, strings.TrimPrefix(v.Key, DEFAULT_TRIM_KEY))
			}
		}
	} else {
		if resp.Node.Value != "" {
			respRec := RecValue{}
			json.Unmarshal([]byte(resp.Node.Value), &respRec)
			if strings.Compare(respRec.Host, data.Address) == 0 {
				found = append(found, strings.TrimPrefix(resp.Node.Key, DEFAULT_TRIM_ARPA_KEY))
			}
		}
	}

	if len(found) == 0 {
		hserver.ReturnError(w, errors.NoContentError, handler.log)
		return
	}

	for _, v := range found {
		_, err = handler.h.Delete(v, isArpa, false)
		if err != nil {
			hserver.ReturnError(w, err, handler.log)
			return
		}
	}

	go handler.pc.DoPurge(data.Name)

	hserver.ReturnResponse(w, nil, handler.log)
}

// ServeHTTP router interface
func (handler *ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var isArpa bool
	var arr []string

	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Read record request: ", data)

	/* Check input */
	if data.Name == "" {
		handler.log.Error("Name invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		handler.log.Error("Name invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = strings.Split(data.Name, ".")
		total := len(arr)
		for i := 0; i < total / 2; i++ {
			tmp := arr[i]
			arr[i] = arr[total - 1 - i]
			arr[total - 1 - i] = tmp
		}
	}
	rec := strings.Join(arr, "/")

	handler.log.Debug("Rec is %s", rec)

	resp, err := handler.h.Read(rec, isArpa, false)

	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	eresp := &MacedonResponse{}
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
		json.Unmarshal([]byte(v.Value), &respRec)
			addr := MacedonRequest{}
			addr.Name = data.Name
			addr.Address = respRec.Host
			addr.Ttl = respRec.Ttl
			eresp.Result = append(eresp.Result, addr)
		}
	} else {
		if resp.Node.Value != "" {
			respRec := RecValue{}
			json.Unmarshal([]byte(resp.Node.Value), &respRec)
			addr := MacedonRequest{}
			addr.Name = data.Name
			addr.Address = respRec.Host
			addr.Ttl = respRec.Ttl
			eresp.Result = append(eresp.Result, addr)
		}
	}

	if len(eresp.Result) == 0 {
		hserver.ReturnError(w, errors.NoContentError, handler.log)
		return
	}

	hserver.ReturnResponse(w, eresp, handler.log)
}

// ServeHTTP router interface
func (handler *AddServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}

	handler.log.Info("Add server request from %s ", req.RemoteAddr, data)

	if data.Address == "" || net.ParseIP(data.Address) == nil {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Address invalid", http.StatusBadRequest)
		return
	}
	if !strings.Contains(data.Address, ":") {
		data.Address = data.Address + ":" + DEFAULT_PURGE_PORT
	}

	rec := fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))

	handler.log.Debug("Rec is %s", rec)

	_, err = handler.h.Add(rec, data.Address, 0, false, true)

	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}
	handler.pc.AddServer(data.Address)

	hserver.ReturnResponse(w, nil, handler.log)
}

// ServeHTTP router interface
func (handler *ReadServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//if req.Method != "POST" {
	//	handler.log.Error("Method invalid: %s", req.Method)
	//	http.Error(w, "Method invalid", http.StatusBadRequest)
	//	return
	//}

	//result, err:= ioutil.ReadAll(req.Body)
	//if err != nil {
	//	handler.log.Error("Read from request body failed: %s", err)
	//	http.Error(w, "Parse from body failed", http.StatusBadRequest)
	//	return
	//}
	//req.Body.Close()

	//data := &ServerRequest{}
	//err = json.Unmarshal(result, &data)
	//if err != nil {
	//	handler.log.Error("Parse from request body failed:%s ", err)
	//	http.Error(w, "Parse from body failed", http.StatusBadRequest)
	//	return
	//}
	resp, err := handler.h.Read("", false, true)

	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	eresp := &ServerResponse{}
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
		json.Unmarshal([]byte(v.Value), &respRec)
			addr := ServerRequest{}
			addr.Address = respRec.Host
			eresp.Result = append(eresp.Result, addr)
		}
	}

	if len(eresp.Result) == 0 {
		hserver.ReturnError(w, errors.NoContentError, handler.log)
		return
	}

	hserver.ReturnResponse(w, eresp, handler.log)
}

// ServeHTTP router interface
func (handler *DeleteServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}

	handler.log.Info("Delete server request from %s ", req.RemoteAddr, data)

	if data.Address == "" || net.ParseIP(data.Address) == nil {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Address invalid", http.StatusBadRequest)
		return
	}

	if !strings.Contains(data.Address, ":") {
		data.Address = data.Address + ":" + DEFAULT_PURGE_PORT
	}

	resp, err := handler.h.Read("", false, true)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	var found []string
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
			json.Unmarshal([]byte(v.Value), &respRec)
			if strings.Compare(respRec.Host, data.Address) == 0 {
				handler.log.Error("delete server key is %s", v.Key)
				found = append(found, strings.TrimPrefix(v.Key, DEFAULT_TRIM_SERVER_KEY))
			}
		}
	}

	if len(found) == 0 {
		hserver.ReturnError(w, errors.NoContentError, handler.log)
		return
	}

	for _, v := range found {
		_, err = handler.h.Delete(v, false, true)
		if err != nil {
			hserver.ReturnError(w, err, handler.log)
			return
		}
	}

	handler.pc.DeleteServer(data.Address)

	hserver.ReturnResponse(w, nil, handler.log)
}
