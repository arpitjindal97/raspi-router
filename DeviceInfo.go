package main

import (
	"net/http"

	"encoding/json"
)

type OSInfo struct{
	DistributionId	string
	Description		string
	Release			string
	Codename		string
	Hostname		string
	KernelRelease	string
	Architecture	string
	ModelName		string
	CPUs			string
	LocalTime		string
	TimeZone		string
	UpTime			string
	UpSince			string
}

func DeviceInfo(w http.ResponseWriter, r *http.Request) {

	var info OSInfo

	info.DistributionId = GetOutput("lsb_release -i | awk '{print $3}'")
	info.Description = GetOutput("lsb_release -d | cut -f1 --complement'")
	info.Release = GetOutput(" lsb_release -r | cut -f2")
	info.Codename = GetOutput(" lsb_release -c | cut -f2")

	info.Hostname = GetOutput("hostname")
	info.KernelRelease = GetOutput("uname -r")
	info.Architecture = GetOutput("uname -m")
	info.ModelName = GetOutput("lscpu | grep 'Model name' | awk '{$1=\"\";$2=\"\";print}'")
	info.ModelName = info.ModelName[2:]

	info.CPUs = GetOutput("lscpu | grep CPU\\(s\\): | grep -v node | awk '{print $2}'")

	info.LocalTime = GetOutput("date")
	info.TimeZone = GetOutput("timedatectl | grep zone | awk '{$1=\"\";$2=\"\";print}'")
	info.TimeZone = info.TimeZone[2:]

	info.UpTime = GetOutput("uptime -p")
	info.UpSince = GetOutput("uptime -s")

	b,_ := json.MarshalIndent(info,"","	")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(b)

}
