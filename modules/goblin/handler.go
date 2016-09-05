package goblin

import (
	"fmt"
	"bytes"
	"strconv"
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
	h.log = log
	h.add_loc = loc + GOBLIN_ADD_LOCATION
	h.delete_loc = loc + GOBLIN_DELETE_LOCATION
	h.host = host

	return h
}

func (h *Handler) RuleOperate(addr string, args *bytes.Buffer, op int) error {
	var err error
	var resp *http.Response

	client := &http.Client{}
	switch op {
	case ADD_RULE:
		h.log.Debug("Add rule args is ", args)
		req, err := http.NewRequest("POST", "http://" + addr + h.add_loc, args)
		if err != nil {
			h.log.Error("New request failed: ", err)
			return errors.InternalServerError
		}
		req.Header.Add("Content-Type", variable.FORM_CONTENT_HEADER)
		req.Header.Add("Content-Length", strconv.Itoa(args.Len()))

		if h.host != "" {
			req.Host = h.host
			h.log.Debug("Add header %s", h.host)
		}
		resp, err = client.Do(req)
	case DELETE_RULE:
		h.log.Debug("Delete rule args is ", args)
		req, err := http.NewRequest("POST", "http://" + addr + h.delete_loc, args)
		if err != nil {
			h.log.Error("New request failed: ", err)
			return errors.InternalServerError
		}
		req.Header.Add("Content-Type", variable.FORM_CONTENT_HEADER)
		req.Header.Add("Content-Length", strconv.Itoa(args.Len()))

		if h.host != "" {
			req.Host = h.host
			h.log.Debug("Add header %s", h.host)
		}
		resp, err = client.Do(req)
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

func (h *Handler) RuleAdd(addr, ip, uuid string, expire int, action string) error {
	if ip == "" {
		ip = EMPTY_IP
	}
	args := fmt.Sprint("startip=", ip, "&endip=", ip)
	if uuid != "" && len(uuid) > 1 {
		args = fmt.Sprint(args, "&uuid=", uuid)
	}
	args = fmt.Sprint(args, "&expire=", expire, "&punish=", action, "&punish_arg=0\r\n")

	/* post data */
	data := bytes.NewBufferString(args)

	return h.RuleOperate(addr, data, ADD_RULE)
}

func (h *Handler) RuleDelete(addr, ip, uuid, action string) error {
	if ip == "" {
		ip = EMPTY_IP
	}
	args := fmt.Sprint("startip=", ip, "&endip=", ip)
	if uuid != "" && len(uuid) > 1 {
		args = fmt.Sprint(args, "&uuid=", uuid)
	}

	args = fmt.Sprint(args, "&punish=", action, "\r\n")

	data := bytes.NewBufferString(args)
	return h.RuleOperate(addr, data, DELETE_RULE)
}
