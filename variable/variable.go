package variable
import (
	"time"
)

const (
	VERSION                   = "0.3 alpha"

	DEFAULT_LOG_PATH          = "../log/gapi.log"

	DEFAULT_CONTENT_HEADER    = "application/json;charset=utf-8"

	DEFAULT_CONFIG_PATH       = "../conf"
	DEFAULT_CONFIG_FILE       = "gapi.conf"

	HTTP_OK                   = 200
	HTTP_CREATED              = 201
	HTTP_CONFLICT             = 409
	HTTP_NOT_FOUND            = 404

	DEFAULT_QUIT_WAIT_TIME    = time.Millisecond * 200
	UDP_DEFAULT_BUFFER_SIZE   = 65536
	USOCK_DEFAULT_BUFFER_SIZE = 65536
	DEFAULT_ROTATE_SIZE       = 500 //size in MB
	DEFAULT_BACKUP_SIZE       = 7
)
