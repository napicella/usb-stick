#!/usr/bin/env bash
set -x

export GOPATH=$(pwd)
SERVICE_URL=$1
if [[ -z "${SERVICE_URL}" ]]
then
	echo "Please set SERVICE_URL"
	exit -1
fi

mkdir -p ./build/linux | true
mkdir -p ./build/darwin | true
mkdir -p ./build/windows | true

rm -rf ./build/*

cd src/github.com/usb-stick-client && dep ensure && cd -

env GOOS=linux GOARCH=amd64 \
go build -ldflags "-X github.com/usb-stick-client/model.ServiceUrl=$SERVICE_URL" -o ./build/linux/usb github.com/usb-stick-client

env GOOS=darwin GOARCH=amd64 \
go build -ldflags "-X github.com/usb-stick-client/model.ServiceUrl=$SERVICE_URL" -o ./build/darwin/usb github.com/usb-stick-client

env GOOS=windows GOARCH=amd64 \
go build -ldflags "-X github.com/usb-stick-client/model.ServiceUrl=$SERVICE_URL" -o ./build/windows/usb.exe github.com/usb-stick-client