package config

import (
	"os"
	"fmt"
	//"time"
	"path/filepath"
	goconf "github.com/msbranco/goconfig"
	"github.com/gwtony/gapi/errors"
	"github.com/gwtony/gapi/variable"
)

// Config of server
type Config struct {
	HttpAddr    string  /* http server bind address */
	UdpAddr     string  /* udp server bind address */
	TcpAddr     string  /* tcp server bind address */
	UsocketAddr string  /* usocket server bind address */

	//UdpNFI     string  /* udp server multicast receive interface */

	Location    string  /* handler location */

	Log         string  /* log file */
	Level       string  /* log level */
	RotateLine  int     /* log rotate line */

	File        string  /* config file */
	C           *goconf.ConfigFile /* goconfig struct */
}

func (conf *Config) SetConf(file string) {
	conf.File = filepath.Join(variable.DEFAULT_CONFIG_PATH, file)
}
// ReadConf reads conf from file
func (conf *Config) ReadConf(file string) error {
	if file == "" {
		if conf.File == "" {
			file = filepath.Join(variable.DEFAULT_CONFIG_PATH, variable.DEFAULT_CONFIG_FILE)
		} else {
			file = conf.File
		}
	}

	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf file %s failed", file)
		return err
	}
	conf.C = c
	return nil
}

// ParseConf parses config
func (conf *Config) ParseConf() error {
	var err error

	if conf.C == nil {
		fmt.Fprintln(os.Stderr, "[Error] Must read config first")
		return errors.BadConfigError
	}

	conf.HttpAddr, err = conf.C.GetString("default", "http_addr")
	if err != nil {
		//fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No HttpAddr")
		conf.HttpAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on http addr:", conf.HttpAddr)
	}

	conf.TcpAddr, err = conf.C.GetString("default", "tcp_addr")
	if err != nil {
		//fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No TcpAddr")
		conf.UdpAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on tcp addr:", conf.TcpAddr)
	}
	conf.UdpAddr, err = conf.C.GetString("default", "udp_addr")
	if err != nil {
		//fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No UdpAddr")
		conf.UdpAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on udp addr:", conf.UdpAddr)
	}

	conf.UsocketAddr, err = conf.C.GetString("default", "usocket_addr")
	if err != nil {
		conf.UsocketAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on usocket addr:", conf.UsocketAddr)
	}
	//conf.UdpAddr, err = conf.C.GetString("default", "udp_interface")
	//if err != nil {
	//	conf.UdpNIF = ""
	//} else {
	//	fmt.Fprintln(os.Stderr, "[Info] [Default] use udp network interface:", conf.UdpNIF)
	//}

	conf.Log, err = conf.C.GetString("default", "log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Default] log not found, use default log file")
		conf.Log = ""
	}
	conf.Level, err = conf.C.GetString("default", "level")
	if err != nil {
		conf.Level = "error"
		fmt.Fprintln(os.Stderr, "[Info] [Default] level not found, use default log level error")
	}
	rline, err := conf.C.GetInt64("default", "rotate_line")
	if err != nil {
		rline = variable.DEFAULT_ROTATE_LINE
		fmt.Fprintln(os.Stderr, "[Info] [Default] rotate_line not found, use default", rline)
	}
	conf.RotateLine = int(rline)

	return nil
}

