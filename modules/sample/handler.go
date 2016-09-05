package sample

import (
	"fmt"
	//"bytes"
	//"time"
	"net/http"
	"git.lianjia.com/lianjia-sysop/napi/variable"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type Handler struct {
	add_loc    string
	delete_loc string
	host       string
	log        log.Log
}

func InitHandler(loc, host string, log log.Log) *Handler {
	h := &Handler{}
	h.host = host
	h.log = log
	h.add_loc = loc + SAMPLE_ADD_LOCATION
	h.delete_loc = loc + SAMPLE_DELETE_LOCATION

	return h
}

func (h *Handler) RuleOperate(addr string, args string, op int) error {
	var err error
	var resp *http.Response
	client := &http.Client{}

	switch op {
	case ADD_RULE:
		h.log.Debug("Add rule args is ", args)
		req, err := http.NewRequest("GET", "http://" + addr + h.add_loc + "?" + args, nil)
		if err != nil {
			h.log.Error("New request failed: ", err)
			return errors.InternalServerError
		}
		if h.host != "" {
			req.Host = h.host
			h.log.Debug("Add header %s", h.host)
		}
		resp, err = client.Do(req)
		break
	case DELETE_RULE:
		h.log.Debug("Delete rule args is ", args)
		req, err := http.NewRequest("GET", "http://" + addr + h.delete_loc + "?" + args, nil)
		if err != nil {
			h.log.Error("New request failed: ", err)
			return errors.InternalServerError
		}
		if h.host != "" {
			req.Host = h.host
		}
		resp, err = client.Do(req)
		break
	default: /* Should not reach here */
		h.log.Error("Unknown operate code: ", op)
		return errors.InternalServerError
	}

	if err != nil {
		h.log.Error("Opereate service to nginx failed: ", err)
		return errors.BadGatewayError
	}

	defer resp.Body.Close()

	if resp.StatusCode != variable.HTTP_OK {
		h.log.Error("Opereate http status error: %d", resp.StatusCode)
		return errors.BadGatewayError
	}

	return nil
}

func (h *Handler) RuleAdd(addr, host string, band, expire int, header string) error {
	if addr == "" {
		return nil
	}
	args := fmt.Sprint("host=", host, "&band=", band,
			"&expire=", expire, "&header=", header)

	return h.RuleOperate(addr, args, ADD_RULE)
}

func (h *Handler) RuleDelete(addr, host string, band, expire int, header string) error {
	if addr == "" {
		return nil
	}
	args := fmt.Sprint("host=", host, "&band=", band,
			"&expire=", expire, "&header=", header)

	return h.RuleOperate(addr, args, DELETE_RULE)
}
