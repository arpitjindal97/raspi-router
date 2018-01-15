package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"fmt"
)

func UpdateInterface(w http.ResponseWriter, r *http.Request) {

	// reading response from frontend
	///storing in rec_interface
	var rec_interface Interfaces
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&rec_interface)
	if err!=nil{
		panic(err)
	}


	//storing the new config into config.json & files
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
	rec_interface.IsWifi = file.NetworkInterfaces[i].IsWifi
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


	//if there is any change in wpa, hostapd,dnsmasq then restart

	if rec_interface.Mode != File.NetworkInterfaces[i].Mode {

		if rec_interface.Mode == "default" {

			// kill("hostapd")
			// kill("dnsmasq")
			StartParticularInterface(rec_interface)

		} else {
			kill("wpa_supplicant")
			kill("dhcpcd")
			time.Sleep(time.Second*2)
		}

		StartParticularInterface(rec_interface)

	} else {
		if rec_interface.Wpa != File.NetworkInterfaces[i].Wpa {

			kill("wpa_supplicant")
			kill("dhcpcd")
			time.Sleep(time.Second*2)
			
			StartParticularInterface(rec_interface)

		} else if rec_interface.Hostapd != File.NetworkInterfaces[i].Hostapd ||
			rec_interface.Dnsmasq != File.NetworkInterfaces[i].Dnsmasq {

			// kill("hostapd")
			// kill("dnsmasq")
			StartParticularInterface(rec_interface)

		}
	}


	File = FirstTask()

}
