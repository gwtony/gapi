package macedon

import (
	"net"
	"sync"
	"time"
	"git.lianjia.com/lianjia-sysop/napi/log"
)

type PurgeContext struct {
	lock    *sync.RWMutex

	ips     []string
	iplen   int
	port    string
	cmd     string

	timeout time.Duration

	log     log.Log
}

func InitPurgeContext(ips []string, port, cmd string, timeout time.Duration, log log.Log) *PurgeContext {
	pc := &PurgeContext{}

	pc.log = log
	pc.ips = ips
	pc.iplen = len(ips)

	pc.lock = &sync.RWMutex{}

	pc.port = port
	pc.cmd = cmd

	return pc
}

/* Do not return any error */
func (pc *PurgeContext) DoPurge(name string) error {
	pc.log.Debug("Do purge")

	ch := make(chan int, pc.iplen)

	pc.lock.RLock()
	for _, host := range pc.ips {
		pc.log.Debug("Purge ip: %s", host)

		go func(ip string) {
			defer func () { ch <- 1 }()
			conn, err := net.DialTimeout("tcp", ip + ":" + pc.port, pc.timeout)
			if err != nil {
				pc.log.Error("Connect to %s failed", ip + ":" + pc.port)
				return
			}
			defer conn.Close()
			conn.Write([]byte(pc.cmd + " " + name))

		}(host)
	}
	pc.lock.RUnlock()

	for i := 0; i < pc.iplen; i++ {
		<-ch
	}

	pc.log.Debug("All purge done")

	return nil
}
