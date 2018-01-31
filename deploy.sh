#!/usr/bin/env bash
GOOS=linux GOARCH=arm GOARM=7 go build *.go

zip this.zip -r static DeviceInfo

scp this.zip pi@192.168.0.3:/home/pi/Desktop/

rm  this.zip

rm DeviceInfo