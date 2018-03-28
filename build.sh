#!/usr/bin/env bash


build() {
    packr
    rm -rf bin
    mkdir bin
    GOOS=linux GOARCH=arm GOARM=7 go build -o bin/raspi-router-armv7
    GOOS=linux GOARCH=amd64 go build -o bin/raspi-router-amd64
    GOOS=linux GOARCH=386 go build -o bin/raspi-router-i386
    rm main-packr.go
}

update() {
    rm -rf dist
    wget https://github.com/arpitjindal97/raspi-router-frontend/releases/download/v1.0.0/raspi-router-frontend-v1.0.0.tar.gz
    tar -xvf raspi-router-frontend-*.tar.gz
    rm raspi-router-frontend-*.tar.gz
}
run() {
	./bin/raspi-router-amd64
}

if [ "$1" == "update" ]; then
    update
elif [ "$1" == "build" ]; then
    build
elif [ "$1" == "run" ]; then
	run
else
    update
    build
fi
