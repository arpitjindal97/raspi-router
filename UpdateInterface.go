package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
)

func UpdateInterface(w http.ResponseWriter, r *http.Request) {

	var rec_interface Interfaces
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&rec_interface)
	if err!=nil{
		panic(err)
	}

	raw,_ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw,&file)

	var i int

	for i=0;i<len(file.NetworkInterfaces);i++ {

		if file.NetworkInterfaces[i].Name  == rec_interface.Name {
			break
		}
	}

	file.NetworkInterfaces[i].Mode = rec_interface.Mode
	file.NetworkInterfaces[i].RouteMode = rec_interface.RouteMode
	file.NetworkInterfaces[i].RouteInterface = rec_interface.RouteInterface

	file.NetworkInterfaces[i].IpModes = rec_interface.IpModes
	file.NetworkInterfaces[i].IpAddress = rec_interface.IpAddress
	file.NetworkInterfaces[i].SubnetMask = rec_interface.SubnetMask
	file.NetworkInterfaces[i].Wpa = ""
	file.NetworkInterfaces[i].Hostapd = ""
	file.NetworkInterfaces[i].Dnsmasq = ""

	b,_ := json.MarshalIndent(file,"","	")

	ioutil.WriteFile("config/config.json",b,0644)

	if file.NetworkInterfaces[i].IsWifi == "true" {

		raw = []byte(rec_interface.Hostapd)
		ioutil.WriteFile("config/"+rec_interface.Name+"_hostapd.conf",raw,os.FileMode(0644))

		raw = []byte(rec_interface.Wpa)
		ioutil.WriteFile("config/"+rec_interface.Name+"_wpa.conf",raw,os.FileMode(0644))
	}
	raw = []byte(rec_interface.Dnsmasq)
	ioutil.WriteFile("config/"+rec_interface.Name+"_dnsmasq.conf",raw,os.FileMode(0644))


	File = FirstTask()
}