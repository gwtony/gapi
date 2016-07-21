package macedon

import (
	"fmt"
	"os"
	"time"
	"strings"
	"git.lianjia.com/lianjia-sysop/napi/config"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type MacedonConfig struct {
	eaddr      []string /* etcd addr */

	api_loc    string   /* macedon api location */
	loc        string   /* macedon location */

	domain     string

	purge_cmd  string
	purge_ips  []string
	purge_port string
	purge_to   time.Duration
}


func (conf *MacedonConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	eaddr_str, err := cf.C.GetString("macedon", "etcd_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [macedon] Read conf: No etcd_addr")
		return err
	}
	if eaddr_str == "" {
		fmt.Fprintln(os.Stderr, "[Error] [macedon] Empty etcd server address")
		return errors.BadConfigError
	}
	eaddr := strings.Split(eaddr_str, ",")
	for i := 0; i < len(eaddr); i++ {
		if eaddr[i] != "" {
			if !strings.Contains(eaddr[i], ":") {
				conf.eaddr = append(conf.eaddr, eaddr[i] + ":" + DEFAULT_ETCD_PORT)
			} else {
				conf.eaddr = append(conf.eaddr, eaddr[i])
			}
		}
	}

	conf.loc, err = cf.C.GetString("macedon", "location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No macedon_location, use default location", DEFAULT_SKYDNS_LOC)
		conf.loc = DEFAULT_SKYDNS_LOC
	}

	conf.api_loc, err = cf.C.GetString("macedon", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No api_location, use default location", MACEDON_LOC)
		conf.api_loc = MACEDON_LOC
	}

	conf.domain, err = cf.C.GetString("macedon", "domain")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: No domain")
		return err
	}

	conf.purge_cmd, err = cf.C.GetString("macedon", "purge_cmd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: use default purge_cmd")
		conf.purge_cmd = DEFAULT_PURGE_CMD
	}

	purge_to, err := cf.C.GetInt64("macedon", "purge_timeout")
	if err != nil || purge_to <= 0 {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: use default purge_timeout: ", DEFAULT_PURGE_TIMEOUT)
		purge_to = DEFAULT_PURGE_TIMEOUT
	}
	conf.purge_to =  time.Duration(purge_to) * time.Second

	conf.purge_port, err = cf.C.GetString("macedon", "purge_port")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: use default purge_port: ", DEFAULT_PURGE_PORT)
		conf.purge_port = DEFAULT_PURGE_PORT
	}

	ips_str, err := cf.C.GetString("macedon", "purge_ips")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: purge_ip invalid")
		return err
	}

	if ips_str == "" {
		fmt.Fprintln(os.Stderr, "[Info] [Macedon] Read conf: purge_ip is empty")
		return err
	}

	ips := strings.Split(ips_str, ",")
	for i := 0; i < len(ips); i++ {
		if ips[i] != "" {
			conf.purge_ips = append(conf.purge_ips, ips[i])
		}
	}

	return nil
}
