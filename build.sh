#!/bin/bash

set -eu

if [ ! -f "build.sh" ]; then
        echo "$0 must be run from the root of the repository."
	    exit 2
fi
export GO111MODULE=on
export GOPROXY=https://goproxy.io
export GOPROXY=https://athens.azurefd.net
export GOPROXY=https://gocenter.io
export GOPROXY=https://proxy.golang.org
export GOPROXY=https://goproxy.cn
#export GOPROXY=https://gonexus.dev
#export GOPROXY=https://mirrors.aliyun.com/goproxy/

go run build/ci.go install smw
