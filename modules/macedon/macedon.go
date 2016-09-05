package macedon

import (
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/config"
)

// InitContext inits macedon context
func InitContext(conf *config.Config, hs *hserver.HttpServer, log log.Log) error {
	cf := &MacedonConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("Macedon parse config failed")
		return err
	}

	h := InitHandler(cf.eaddr, cf.loc, log)

	apiLoc := cf.apiLoc
	domain := cf.domain

	pc := InitPurgeContext(h, cf.purgeCmd, cf.purgeTo, log)

	hs.AddRouter(apiLoc + MACEDON_ADD_LOC, &AddHandler{h: h, domain: domain, log: log})
	hs.AddRouter(apiLoc + MACEDON_DELETE_LOC, &DeleteHandler{h: h, domain: domain, pc: pc, log: log})
	hs.AddRouter(apiLoc + MACEDON_READ_LOC, &ReadHandler{h: h, domain: domain, log: log})
	hs.AddRouter(apiLoc + MACEDON_ADD_SERVER_LOC, &AddServerHandler{h: h, pc: pc, log: log})
	hs.AddRouter(apiLoc + MACEDON_DELETE_SERVER_LOC, &DeleteServerHandler{h: h, pc: pc, log: log})
	hs.AddRouter(apiLoc + MACEDON_READ_SERVER_LOC, &ReadServerHandler{h: h, log: log})

	return nil
}
