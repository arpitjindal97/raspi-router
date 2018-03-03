package main

import (
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

func StartBridging() {

	for _, i := range File.BridgeInterfaces {

		read_closer := ioutil.NopCloser(bytes.NewBuffer([]byte(i.Name)))

		body := http.Request{Body: read_closer}

		CreateBridge(nil, &body)

		for _, j := range i.Slaves {

			GetOutput("ip link set " + j + " master " + i.Name)
		}

		if i.IpMode == "dhcp" {

			go ExecuteWait("dhcpcd", "-q", "-w", "-t", "0", i.Name)
		} else {

			ExecuteWait("ifconfig", i.Name, i.IpAddress, "netmask", i.SubnetMask)
		}
	}
}

func CreateBridge(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	str := string(body)

	out := GetOutput("ip link ls " + str)

	if strings.Contains(out, "does not exist") == false {
		GetOutput("ip link add name " + str + " type bridge")
	}

	GetOutput("ip link set dev " + str + " up")

	out = GetOutput("ip link ls " + str + " | grep UP")

	go ExecuteWait("dhcpcd", "-q", "-w", "-t", "0", str)

	if w == nil {
		return
	}

	if strings.Contains(out, "does not exists") == true {

		w.Write([]byte("Error encountered while creating " + str))

		return
	}

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)
	var new_bridge BridgeInterfaces
	new_bridge.Name = str
	new_bridge.IpMode = "dhcp"
	new_bridge.Slaves = []string{}
	file.BridgeInterfaces = append(file.BridgeInterfaces[:], new_bridge)
	b, _ := json.MarshalIndent(file, "", "	")
	ioutil.WriteFile("config/config.json", b, 0644)

	File.BridgeInterfaces = file.BridgeInterfaces


	w.Write([]byte(str + " created and up"))

}

func DeleteBridge(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	str := string(body)

	out := GetOutput("ip link ls " + str)

	if strings.Contains(out, "does not exist") == true {
		w.Write([]byte("Error encountered while deleting " + str))
		return
	}

	GetOutput("ip link del " + str)

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	for i := 0; i < len(file.BridgeInterfaces); i++ {

		if file.BridgeInterfaces[i].Name == str {

			file.BridgeInterfaces = append(file.BridgeInterfaces[:i], file.BridgeInterfaces[i+1:]...)
			break
		}
	}

	for i:=0;i<len(file.PhysicalInterfaces);i++ {
		if file.PhysicalInterfaces[i].BridgeMaster == str {
			file.PhysicalInterfaces[i].BridgeMaster = ""
			File.PhysicalInterfaces[i].BridgeMaster = ""
			break
		}
	}

	b, _ := json.MarshalIndent(file, "", "	")

	ioutil.WriteFile("config/config.json", b, 0644)

	File.BridgeInterfaces = file.BridgeInterfaces

	w.Write([]byte(str + " deleted"))
}

func BridgeRemoveSlave(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	str := string(body)

	GetOutput("ip link set " + str + " nomaster")

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	for i := 0; i < len(file.BridgeInterfaces); i++ {

		for j := 0; j < len(file.BridgeInterfaces[i].Slaves); j++ {

			if file.BridgeInterfaces[i].Slaves[j] == str {
				file.BridgeInterfaces[i].Slaves =
					append(file.BridgeInterfaces[i].Slaves[:j], file.BridgeInterfaces[i].Slaves[j+1:]...)
				j--
			}
		}
	}

	for i := 0; i < len(file.PhysicalInterfaces); i++ {
		if file.PhysicalInterfaces[i].Name == str {
			file.PhysicalInterfaces[i].BridgeMaster = ""
			File.PhysicalInterfaces[i].BridgeMaster = ""
			break
		}
	}

	b, _ := json.MarshalIndent(file, "", "	")

	ioutil.WriteFile("config/config.json", b, 0666)

	File.BridgeInterfaces = file.BridgeInterfaces

	w.Write([]byte(str + " removed"))

}

func BridgeAddSlave(w http.ResponseWriter, r *http.Request) {
	type BridgeAdd struct {
		BridgeIfname string
		SlaveIfname  string
	}

	var bridge BridgeAdd
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&bridge)
	if err != nil {
		panic(err)
	}

	GetOutput("ip link set " + bridge.SlaveIfname + " master " + bridge.BridgeIfname)

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	for i := 0; i < len(file.BridgeInterfaces); i++ {
		if file.BridgeInterfaces[i].Name == bridge.BridgeIfname {
			file.BridgeInterfaces[i].Slaves = append(file.BridgeInterfaces[i].Slaves, bridge.SlaveIfname)
			break
		}
	}

	for i := 0; i < len(file.PhysicalInterfaces); i++ {
		if file.PhysicalInterfaces[i].Name == bridge.SlaveIfname {
			file.PhysicalInterfaces[i].BridgeMaster = bridge.BridgeIfname
			File.PhysicalInterfaces[i].BridgeMaster = bridge.BridgeIfname
			break
		}
	}

	b, _ := json.MarshalIndent(file, "", "	")

	ioutil.WriteFile("config/config.json", b, 0666)

	File.BridgeInterfaces = file.BridgeInterfaces

	w.Write([]byte(bridge.SlaveIfname + " added to " + bridge.BridgeIfname))

}
