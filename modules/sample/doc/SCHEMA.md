# Sample
## Description
Sample api supports add header to mark sample requests, it work with ngx_sample_router module
* Feature:
	- add sample rule
	- delete sample rule
	- read sample rule
	- add ngx server
	- delete ngx server
	- read ngx server

## API

### Add sample rule

```
 Request
POST /sample/rule/add HTTP/1.1
Content-Type: application/json
{
    “type”: “sample”,                                 //string, MUST be "sample"
    “product”: “test”,                                //string, product name
    “match”: {                               
        “band”: [1, 2],                               //[]int, which band to be marked, requests default splits to 10 bands
        “host”: [“a.test.com”, “b.test.com”],         //[]string, hostnames to execute sample
        “expire”: 1467602662                          //int, expire timestamp
    },
    “action”: {
        “type”: “insert_header”,                      //string, MUST be “insert_header”
        “value”: “test:test_a”                        //string, header to add
    }
}
Response
HTTP/1.1 200 OK
{
    “ruleid”: “75613e00-b431-404a-848b-afa72d3ed8f0”  //delete key
}
```

### Delete sample rule


```
Request
POST /sample/rule/delete HTTP/1.1
Content-Type: application/json
{
    “ruleid”: “75613e00-b431-404a-848b-afa72d3ed8f0”  //delete key
}
Response
HTTP/1.1 200 OK
```

* If ruleid not exists, return 204

### Read sample rule
Read sample rule by hostname, support wildcard

```
Request
POST /sample/rule/read HTTP/1.1
Content-Type: application/json
{
    “host”: “a.test.com”,                             //string, hostname
    “type”: “equal”                                   //string, "equal"（exact search） or "like"（fuzzy search）
}

Response
HTTP/1.1 200 OK
[
    {
        “type”: “sample”,
        “product”: “test”,
        “match”: {
            “band”: [1, 2],
            “host”: [“a.test.com”, “b.test.com”],
            “expire”: 1467602662
        },
        “action”: {
            “type”: “insert_header”,
            “value”: “test:test_a”
        }
    }
]

```

* If rule not found, return 204.
* Expired rules will not return

### Add nginx server


```
Request
POST /sample/server/add HTTP/1.1
Content-Type: application/json
{
    "addr": "192.168.0.1:80",   //string, server address in "ip:port"
    "product": "test"           //string, product name
}

Response
HTTP/1.1 200 OK
```

### Delete nginx server


```
Request
POST /sample/server/delete HTTP/1.1
Content-Type: application/json
{
    "addr": "192.168.0.1:80",   //string, server address in "ip:port"
    "product": "test"           //string, product name
}
Response
HTTP/1.1 200 OK
```

### Read nginx server

```
Request
POST /sample/server/read HTTP/1.1
Content-Type: application/json
{
    "addr": "192.168.0.1:80",   //string, server address in "ip:port"
}

Response
HTTP/1.1 200 OK
{
    "addr": "192.168.0.1:80",   //string, server address in "ip:port"
    "product": "test"           //string, product name
}
```

### Clear expire rules


```
Request
POST /sample/rule/clear HTTP/1.1
Content-Type: application/json
{
    "expire": 1467968060,       //int, expire timestamp, rules expired before this timestamp will be clean, MUST be a past timestamp
    "token": "xxxxxx"           //string, authorization token
}
Response
HTTP/1.1 200 OK
```


## Error code
200: success
204: rules not found
400: bad request
404: location not found
500: api server error
502: upstream server error

