package main

import (
	"net/http"
	"encoding/json"
)

type BasicInfo struct {
	IpAddress        string
	BroadcastAddress string
	Gateway          string
	MacAddress       string
	RecvBytes        string
	RecvPackts       string
	TransBytes       string
	TransPackts      string

	ConntectedTo string
	ApMacAddr    string
	BitRate      string
	Frequency    string
	LinkQuality  string
	Channel      string
}

func NetworkInterface(w http.ResponseWriter, r *http.Request) {

	for i := 0; i < len(File.NetworkInterfaces); i++ {

		File.NetworkInterfaces[i].Info = GetBasicInfo(File.NetworkInterfaces[i])
	}

	b, _ := json.MarshalIndent(File.NetworkInterfaces, "", "	")



	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(b)
}

func GetBasicInfo(netInt Interfaces) BasicInfo {

	var info BasicInfo

	info.IpAddress = GetOutput("ip addr show " + netInt.Name + " | grep -v inet6 | awk '/inet/{print $2}'")

	info.Gateway = GetOutput("route -n | grep " + netInt.Name + " | grep UG | awk '{print $2}'")

	info.BroadcastAddress = GetOutput("ip addr show " + netInt.Name + " | grep inet | awk '/brd/ {print $4}'")

	info.MacAddress = GetOutput("ip addr show " + netInt.Name + " | awk '/ether/{print $2}'")

	info.RecvPackts = GetOutput("ip -s link show " + netInt.Name + " | tail -n 3 | head -n 1 | awk '{print $2}'")

	info.RecvBytes = GetOutput("ip -s link show " + netInt.Name + " | tail -n 3 | head -n 1 | awk '{print $1}'")

	info.TransPackts = GetOutput("ip -s link show " + netInt.Name + " | tail -n 1 | awk '{print $2}'")

	info.TransBytes = GetOutput("ip -s link show " + netInt.Name + " | tail -n 1 | awk '{print $1}'")


	if netInt.IsWifi == "true" {

		info.ApMacAddr = GetOutput("iw dev " + netInt.Name + " link | awk '/Connected to/{print $3}'")

		info.ConntectedTo = GetOutput("iw dev " + netInt.Name + " link | awk '/SSID/{print $2}'")

		info.Frequency = GetOutput("iw dev " + netInt.Name + " link | awk '/freq/{print $2}'")

		info.BitRate = GetOutput("iw dev " + netInt.Name + " link | awk '/bitrate/{print $3}'")

		info.LinkQuality = GetOutput("iwconfig " + netInt.Name + " | grep Quality | awk '{print $2}' | cut -d '=' -f 2")

		info.Channel = GetOutput("iw dev " + netInt.Name + " info | grep channel | awk '{print $2}'")
	}
	return info
}
