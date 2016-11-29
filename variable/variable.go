package variable
import (
	"time"
)

const (
	VERSION                 = "0.1 alpha"

	DEFAULT_LOG_PATH        = "../log/macedon.log"

	DEFAULT_CONTENT_HEADER  = "application/json;charset=utf-8"

	DEFAULT_CONFIG_PATH     = "../conf"
	DEFAULT_CONFIG_FILE     = "macedon.conf"

	HTTP_OK                 = 200
	HTTP_CREATED            = 201
	HTTP_CONFLICT           = 409
	HTTP_NOT_FOUND          = 404

	DEFAULT_QUIT_WAIT_TIME = time.Millisecond * 200
)
