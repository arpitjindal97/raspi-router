package main

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"log"
)

func BridgeInterfaceStart(inter BridgeInterfaces) string {

	GetOutput("ip link set dev " + inter.Name + " up")

	out := GetOutput("ip link ls " + inter.Name + " | grep UP")

	if strings.Contains(out, "does not exists") == true {

		return "Error encountered while starting " + inter.Name

	}

	for _, item := range inter.Slaves {

		ExecuteWait("ip", "link", "set", item, inter.Name)
	}

	if inter.IpMode == "dhcp" {
		if len(inter.Slaves) == 0 {
			return "dhcpcd not started, no slave attached"
		}

		go ExecuteWait("dhcpcd", "-q", "-w", "-t", "0", inter.Name)
		return "dhcpcd started"
	} else {

		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
		return "Assigned static Ip address"
	}

}

func BridgeInterfaceCreate(inter BridgeInterfaces) string {

	out := GetOutput("ip link ls " + inter.Name)

	log.Println(out + " " + inter.Name)

	if strings.Contains(out, "does not exist") == false {
		return "Error! interface already exists"
	}

	GetOutput("ip link add name " + inter.Name + " type bridge")

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	file.BridgeInterfaces = append(file.BridgeInterfaces,inter)

	b, _ := json.MarshalIndent(file, "", "	")
	ioutil.WriteFile("config/config.json", b, 0644)

	File.BridgeInterfaces = file.BridgeInterfaces

	return inter.Name + " created"

}

func BridgeInterfaceDelete(inter BridgeInterfaces) string {

	out := GetOutput("ip link ls " + inter.Name)

	if strings.Contains(out, "does not exist") == true {
		return "Bridge Interface " + inter.Name + " doesn't exists"

	}

	// update config file and File instance
	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)
	for i := 0; i < len(file.BridgeInterfaces); i++ {

		if file.BridgeInterfaces[i].Name == inter.Name {

			file.BridgeInterfaces = append(file.BridgeInterfaces[:i], file.BridgeInterfaces[i+1:]...)
			break
		}
	}
	b, _ := json.MarshalIndent(file, "", "	")
	ioutil.WriteFile("config/config.json", b, 0644)
	File.BridgeInterfaces = file.BridgeInterfaces

	GetOutput("ip link del " + inter.Name)

	out = GetOutput("ip link ls " + inter.Name)

	if strings.Contains(out, "does not exist") == true {
		return inter.Name + " successfully deleted"
	} else {
		return "Error deleting " + inter.Name
	}

	return inter.Name + " deleted"
}

func BridgeInterfaceUpdate(inter BridgeInterfaces) string {

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	for i := 0; i < len(file.BridgeInterfaces); i++ {

		if file.BridgeInterfaces[i].Name == inter.Name {

			file.BridgeInterfaces[i] = inter
			break
		}
	}

	b, _ := json.MarshalIndent(file, "", "	")
	ioutil.WriteFile("config/config.json", b, 0644)

	File.BridgeInterfaces = file.BridgeInterfaces

	return "Configuration saved"
}

func BridgeInterfaceStop(inter BridgeInterfaces) string {

	for _, item := range inter.Slaves {

		ExecuteWait("ip", "link", "set", item, "nomaster")
	}

	var message string

	if inter.IpMode == "dhcp" {
		if len(inter.Slaves) == 0 {
			message = "dhcpcd not active"
		}

		Kill("dhcpcd.*" + inter.Name)

		message = "dhcpcd killed"
	} else {

		ExecuteWait("ip", "addr", "flush", "dev", inter.Name)
		message = "Flushed static Ip addr"
	}

	GetOutput("ip link set dev " + inter.Name + " down")
	return message
}
