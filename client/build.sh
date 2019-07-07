#!/usr/bin/env bash

export GOPATH=$(pwd)
SERVICE_URL=$1
if [[ -z "${SERVICE_URL}" ]]
then
	echo "Please set SERVICE_URL"
	exit -1
fi

mkdir ./build | true
rm -rf ./build/*

# cd src/github.com/usb-stick-client && dep ensure && cd -
go build -ldflags "-X github.com/usb-stick-client/model.ServiceUrl=$SERVICE_URL" -o ./build/usb github.com/usb-stick-client