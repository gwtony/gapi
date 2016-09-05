# napi
## Description
Http api for nginx module
* sample
* goblin
* macedon

## Installation
* make ('dist' dir will be created)
* Mysql table (table description in sql/*.sql)

## Config Sample

```
[default]
addr: host:port        #napi server address

log: napi.log          #optional, default is "../log/napi.log"
level: error           #debug, info, error

[macedon]
api_location: /macedon #macedon api location
etcd_addr: addr1,addr2 #etcd address, 'ip1:port,ip2:port...'
domain: domain.name    #dns domain name


[goblin]
mysql_addr: host:port  #mysql address
mysql_dbname: name     #db name
mysql_dbuser: user     #db user
mysql_dbpwd: pwd       #db pwd
location: /goblin      #goblin nginx module admin location
api_location: /goblin  #goblin api location
host: hostname         #nginx hostname

[sample]
mysql_addr: host:port  #mysql address
mysql_dbname: name     #db name
mysql_dbuser: user     #db user
mysql_dbpwd: pwd       #db pwd
location: /sample      #sample nginx module admin location
api_location: /sample  #sample api location
host: hostname         #nginx hostname
```

## Usage

* -f config file
* -h help
* -v version

## Schema

* [goblin schema infomation](modules/goblin/doc/SCHEMA.md)
* [macedon schema infomation](modules/macedon/doc/SCHEMA.md)
* [sample schema infomation](modules/sample/doc/SCHEMA.md)

## Dependency

* [log4go](http://code.google.com/p/log4go)
* [goconfig](https://github.com/msbranco/goconfig)
* [golang/x/ssh](http://golang.org/x/crypto/ssh)
* [mysql](https://github.com/go-sql-driver/mysql)
