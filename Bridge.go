package main

import (
	"strings"
	"io/ioutil"
	"encoding/json"
)

func BridgeInterStart(inter BridgeInterfaces) string {

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

func BridgeInterCreate(inter BridgeInterfaces) string {

	out := GetOutput("ip link ls " + inter.Name)

	if strings.Contains(out, "does not exist") == false {
		GetOutput("ip link add name " + inter.Name + " type bridge")
	}

	GetOutput("ip link set dev " + inter.Name + " up")

	out = GetOutput("ip link ls " + inter.Name + " | grep UP")

	if strings.Contains(out, "does not exists") == true {

		return "Error encountered while creating " + inter.Name

	}

	return inter.Name + " created and up"

}

func BridgeInterDelete(inter BridgeInterfaces) string {

	out := GetOutput("ip link ls " + inter.Name)

	if strings.Contains(out, "does not exist") == true {
		return "Bridge Interface " + inter.Name + " doesn't exists"

	}

	GetOutput("ip link del " + inter.Name)

	out = GetOutput("ip link ls " + inter.Name)
	if strings.Contains(out, "does not exist") == true {
		return inter.Name + " successfully deleted"
	} else {
		return "Error deleting " + inter.Name
	}

	return inter.Name + " deleted"
}

func BridgeInterRemoveSlave(slave_inter string) string {

	GetOutput("ip link set " + slave_inter + " nomaster")

	return slave_inter + " removed"

}

func BridgeInterAddSlave(bridge_inter string, slave_inter string) string {

	GetOutput("ip link set " + slave_inter + " master " + bridge_inter)

	return slave_inter + " added to " + bridge_inter

}
func BridgeInterSave(inter BridgeInterfaces, action string) string {

	raw, _ := ioutil.ReadFile("config/config.json")
	var file ConfigFile
	json.Unmarshal(raw, &file)

	if action == "add" {

		var new_bridge BridgeInterfaces
		new_bridge.Name = inter.Name
		new_bridge.IpMode = "dhcp"
		new_bridge.Slaves = []string{}
		file.BridgeInterfaces = append(file.BridgeInterfaces[:], new_bridge)

	} else if action == "delete" {

		for i := 0; i < len(file.BridgeInterfaces); i++ {

			if file.BridgeInterfaces[i].Name == inter.Name {

				file.BridgeInterfaces = append(file.BridgeInterfaces[:i], file.BridgeInterfaces[i+1:]...)
				break
			}
		}

		for i := 0; i < len(file.PhysicalInterfaces); i++ {
			if file.PhysicalInterfaces[i].BridgeMaster == inter.Name {
				file.PhysicalInterfaces[i].BridgeMaster = ""
				File.PhysicalInterfaces[i].BridgeMaster = ""
				break
			}
		}

	} else if action == "update" {
		for i := 0; i < len(file.BridgeInterfaces); i++ {

			if file.BridgeInterfaces[i].Name == inter.Name {

				file.BridgeInterfaces[i] = inter
				break
			}
		}

	} else
	{
	}

	b, _ := json.MarshalIndent(file, "", "	")
	ioutil.WriteFile("config/config.json", b, 0644)

	File.BridgeInterfaces = file.BridgeInterfaces

	return "Configuration saved"
}

func BridgeInterStop(inter BridgeInterfaces) string {

	if inter.IpMode == "dhcp" {
		if len(inter.Slaves) == 0 {
			return "dhcpcd not active"
		}

		Kill("dhcpcd.*" + inter.Name)

		return "dhcpcd killed"
	} else {

		ExecuteWait("ip", "addr", "flush", "dev", inter.Name)
		return "Flushed static Ip addr"
	}

}
