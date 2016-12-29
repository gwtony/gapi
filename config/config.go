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
	HttpAddr   string  /* http server bind address */
	UdpAddr    string  /* udp server bind address */
	TcpAddr    string  /* tcp server bind address */

	Location   string  /* handler location */

	Log        string  /* log file */
	Level      string  /* log level */

	C          *goconf.ConfigFile /* goconfig struct */
}

// ReadConf reads conf from file
func (conf *Config) ReadConf(file string) error {
	if file == "" {
		file = filepath.Join(variable.DEFAULT_CONFIG_PATH, variable.DEFAULT_CONFIG_FILE)
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
		fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No HttpAddr")
		conf.HttpAddr = ""
	}
	conf.TcpAddr, err = conf.C.GetString("default", "tcp_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No TcpAddr")
		conf.UdpAddr = ""
	}
	conf.UdpAddr, err = conf.C.GetString("default", "udp_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No UdpAddr")
		conf.UdpAddr = ""
	}

	conf.Log, err = conf.C.GetString("default", "log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Default] Log not found, use default log file")
		conf.Log = ""
	}
	conf.Level, err = conf.C.GetString("default", "level")
	if err != nil {
		conf.Level = "error"
		fmt.Fprintln(os.Stderr, "[Info] [Default] Level not found, use default log level error")
	}

	return nil
}

