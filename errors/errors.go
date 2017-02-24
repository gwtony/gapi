package errors

import (
	"errors"
	"encoding/json"
)

var (
	// Server internal errors
	// ParseOptionError server parse option error
	ParseOptionError  = errors.New("Parse Option Error")
	// ReadConfigError server read config error
	ReadConfigError      = errors.New("Read Config Error")
	// ParseConfigError server parse config error
	ParseConfigError     = errors.New("Parse Config Error")
	// InitLogError server init log error
	InitLogError         = errors.New("Init Log Error")
	// InitServerError init server error
	InitServerError      = errors.New("Init Server Error")
	// LookupHostError lookup hostname error
	LookupHostError      = errors.New("Lookup Host Error")

	InitHttpServerError  = errors.New("Init Http Server Error")
	InitTcpServerError   = errors.New("Init Tcp Server Error")
	InitUdpServerError   = errors.New("Init Udp Server Error")
	InitUsocketServerError   = errors.New("Init Usocket Server Error")

	NoHandlerError       = errors.New("No Handler Error")

	// Http server errors
	// BadConfigError http bad config error
	BadConfigError       = errors.New("Bad Config")
	// UnauthorizedError http unauthorized error
	UnauthorizedError    = errors.New("Unauthorized")
	// NoContentError http No Content error
	NoContentError       = errors.New("No Content")
	// BadRequestError http bad request error
	BadRequestError      = errors.New("Bad Request")
	// NotAcceptableError http tot acceptable error
	NotAcceptableError   = errors.New("Not Acceptable")
	// ForbiddenError http forbidden error
	ForbiddenError       = errors.New("Forbidden")
	// ConflictError http conflict error
	ConflictError        = errors.New("Conflict")
	// InternalServerError http internal server error
	InternalServerError  = errors.New("Internal Server Error")
	// BadGatewayError http bad gateway error
	BadGatewayError      = errors.New("Bad Gateway")
)

type jerr struct {
	Error string
}

func Jerror(msg string) string {
	var je jerr
	je.Error = msg
	res, _ := json.Marshal(je)
	return string(res)
}

