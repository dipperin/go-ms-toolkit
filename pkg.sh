#!/usr/bin/env bash

if [[ ! $1 ]];then
    echo "do=tidy  or do=vendor"
    exit 2
fi

export GO111MODULE=on

export GOPROXY=https://goproxy.io

go mod $1
