package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
)

func UpdateInterface(w http.ResponseWriter, r *http.Request) {

	// reading response from frontend
	///storing in rec_interface
	var rec_interface PhysicalInterfaces
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&rec_interface)
	if err != nil {
		panic(err)
	}

	//storing the new config into config.json & files
	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	var i int

	for i = 0; i < len(file.PhysicalInterfaces); i++ {

		if file.PhysicalInterfaces[i].Name == rec_interface.Name {
			break
		}
	}

	file.PhysicalInterfaces[i].Mode = rec_interface.Mode
	file.PhysicalInterfaces[i].BridgeMode = rec_interface.BridgeMode
	file.PhysicalInterfaces[i].NatInterface = rec_interface.NatInterface
	rec_interface.IsWifi = file.PhysicalInterfaces[i].IsWifi
	if rec_interface.Mode != "off" {
		file.PhysicalInterfaces[i].IpModes = rec_interface.IpModes
		file.PhysicalInterfaces[i].IpAddress = rec_interface.IpAddress
		file.PhysicalInterfaces[i].SubnetMask = rec_interface.SubnetMask
	}
	file.PhysicalInterfaces[i].Wpa = ""
	file.PhysicalInterfaces[i].Hostapd = ""
	file.PhysicalInterfaces[i].Dnsmasq = ""

	b, _ := json.MarshalIndent(file, "", "	")

	ioutil.WriteFile("config/config.json", b, 0644)

	if file.PhysicalInterfaces[i].IsWifi == "true" {

		raw = []byte(rec_interface.Hostapd)
		ioutil.WriteFile("config/"+rec_interface.Name+"_hostapd.conf", raw, os.FileMode(0644))

		raw = []byte(rec_interface.Wpa)
		ioutil.WriteFile("config/"+rec_interface.Name+"_wpa.conf", raw, os.FileMode(0644))
	}
	raw = []byte(rec_interface.Dnsmasq)
	ioutil.WriteFile("config/"+rec_interface.Name+"_dnsmasq.conf", raw, os.FileMode(0644))



	/*StopInterface(File.PhysicalInterfaces[i])
	time.Sleep(time.Second*2)
	StartParticularInterface(rec_interface)


	File = FirstTask()

	str := "Operation Completed"
	w.Write([]byte(str))*/

}
