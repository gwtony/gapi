#!/bin/bash
if [ $# -gt 0 -a "$1" == "debug" ]; then
	echo "Debug version"
	$GOROOT/bin/go build -o dist/bin/napi -gcflags '-N -l' example/napi_main.go
else 
	echo "Release version"
	$GOROOT/bin/go build -o dist/bin/napi -ldflags '-s -w' example/napi_main.go
fi
if [ $? -eq 0 ]; then
	cat example/napi.conf.tmpl modules/*/config.conf.tmpl > dist/conf/napi.conf
	cp scripts/start_napi.sh dist/bin/
	echo "Build done: check dist dir"
else
	echo "Build failed"
fi 
