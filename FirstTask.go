package main

import (
	"io/ioutil"
	"encoding/json"
	"os"
)

type ConfigFile struct {
	PhysicalInterfaces	[]PhysicalInterfaces
	BridgeInterfaces	[]BridgeInterfaces
}

type PhysicalInterfaces struct {
	Name        	string
	IsWifi     		string
	Mode       		string
	BridgeMode		string
	NatInterface	string
	IpModes        	string
	IpAddress      	string
	SubnetMask     	string
	Wpa            	string
	Hostapd        	string
	Dnsmasq        	string
	Info           	BasicInfo
}

type BridgeInterfaces struct {
	Name			string
	IpMode			string
	IpAddress      	string
	SubnetMask     	string
	Info           	BasicInfo
	Slaves			[]string
}

func FirstTask() ConfigFile {

	raw, err := ioutil.ReadFile("config/config.json")

	var file ConfigFile

	if err != nil {

		// File not found
		bb, _ := json.MarshalIndent(file, "", "	")

		// make dir if not exists
		_, err = ioutil.ReadDir("config")
		if err != nil {
			os.Mkdir("config", 0755)
		}

		ioutil.WriteFile("config/config.json", bb, 0666)
		raw = bb
	}

	json.Unmarshal(raw, &file)



	// Get Bridge Interfaces
	out := GetOutput("ip link show type bridge |grep -v link| awk '{print $2}'")

	// Format them
	bridge_names := FormatInterfaceName(out)
	//file.BridgeInterfaces = CorrectBridgeMismatch(file.BridgeInterfaces,interface_names)




	// Get interfaces excluding bridges

	out = GetOutput(
		"ip link | grep -v link | awk '{print $2}'")

	// Format them
	interface_names := FormatInterfaceName(out)

	for _,i := range bridge_names {
		for j := 0;j<len(interface_names);j++ {
			if i == interface_names[j] {
				interface_names = append(interface_names[:j],interface_names[j+1:]...)
			}
		}
	}

	file.PhysicalInterfaces = CorrectInterfaceMismatch(file.PhysicalInterfaces, interface_names)

	if file.BridgeInterfaces == nil {

		file.BridgeInterfaces = []BridgeInterfaces{}
	}



	b, _ := json.MarshalIndent(file, "", "	")

	ioutil.WriteFile("config/config.json", b, 0666)

	// get wpa, hostapd and dnsmasq files
	file = CaptureConfFiles(file)

	return file

}
func CreateDefaultInterface(name string) PhysicalInterfaces {

	var name_default PhysicalInterfaces
	name_default.Name = name

	str := GetOutput("iwconfig " + name)

	if str == "" {
		name_default.IsWifi = "false"
	} else {
		name_default.IsWifi = "true"
	}

	name_default.Mode = "default"
	name_default.BridgeMode = "wpa"
	name_default.NatInterface = ""
	name_default.IpModes = "dhcp"

	return name_default
}
func CaptureConfFiles(file ConfigFile) ConfigFile {

	for i := 0; i < len(file.PhysicalInterfaces); i++ {

		name := file.PhysicalInterfaces[i].Name

		// Dnsmasq
		raw, err := ioutil.ReadFile("config/" + name + "_dnsmasq.conf")
		if err != nil {
			str := "bind-interfaces\n" +
				"server=8.8.8.8\n" +
				"domain-needed\n" +
				"bogus-priv\n" +
				"dhcp-range=192.168.2.2,192.168.2.100,12h"

			raw = []byte(str)

			ioutil.WriteFile("config/"+name+"_dnsmasq.conf", raw, os.FileMode(0644))

		}
		file.PhysicalInterfaces[i].Dnsmasq = string(raw)

		if file.PhysicalInterfaces[i].IsWifi == "false" {
			continue
		}

		// Hostapd
		raw, err = ioutil.ReadFile("config/" + name + "_hostapd.conf")
		if err != nil {
			str := "interface=" + name + "\n" +
				"driver=nl80211\n" +
				"ssid=Raspberry-Hotspot\n" +
				"hw_mode=g\n" +
				"ieee80211n=1\n" +
				"wmm_enabled=1\n" +
				"macaddr_acl=0\n" +
				"ht_capab=[HT40][SHORT-GI-20][DSSS_CCK-40]\n" +
				"channel=6\n" +
				"auth_algs=1\n" +
				"ignore_broadcast_ssid=0\n" +
				"wpa=2\n" +
				"wpa_key_mgmt=WPA-PSK\n" +
				"wpa_passphrase=raspberry\n" +
				"rsn_pairwise=CCMP\n"

			raw = []byte(str)

			ioutil.WriteFile("config/"+name+"_hostapd.conf", raw, os.FileMode(0644))

		}
		file.PhysicalInterfaces[i].Hostapd = string(raw)

		// Wpa Supplicant
		raw, err = ioutil.ReadFile("config/" + name + "_wpa.conf")
		if err != nil {
			str := "ctrl_interface=/run/wpa_supplicant\n" +
				"update_config=1\n" +
				"network={\n" +
				"ssid=\"Your Wifi Name\"\n" +
				"psk=\"Your Password\"\n" +
				"}\n"

			raw = []byte(str)

			ioutil.WriteFile("config/"+name+"_wpa.conf", raw, os.FileMode(0644))

		}
		file.PhysicalInterfaces[i].Wpa = string(raw)

	}

	return file
}

func FormatInterfaceName(out string) []string {

	var interfaceNames []string
	count := 0
	interfaceNames = append(interfaceNames, "")

	for i := 0; i < len(out)-1; i++ {

		s := string(out[i])

		if s == ":" && (out[i+1]) == 10 {
			count++
			if i+2 != len(out) {
				interfaceNames = append(interfaceNames, "")
			}
			i++
			continue
		}

		interfaceNames[count] = interfaceNames[count] + s

	}
	return interfaceNames
}

func CorrectInterfaceMismatch( found []PhysicalInterfaces, actual []string) []PhysicalInterfaces {

	//removing from found
	for i := 0; i < len(found); i++ {

		name := (found)[i].Name
		match := 0
		for _, j := range actual {

			if name == j {
				match = 1
			}
		}
		if match == 0 {
			found = append( found[:i], found[i+1:]...)
			i--
		}
	}

	//adding to found which exists in actual but not in the found
	for i := 0; i < len(actual); i++ {

		match := 0
		for j := 0; j < len(found); j++ {

			if found[j].Name == actual[i] {
				match = 1
			}
		}
		if match == 0 {
			found = append(found[:], CreateDefaultInterface(actual[i]))
		}

	}

	return found
}

func CorrectBridgeMismatch( found []BridgeInterfaces, actual []string) []BridgeInterfaces {

	//adding to found which exists in actual but not in the found
	for i := 0; i < len(actual); i++ {

		if actual[i] == ""{continue}
		match := 0
		for j := 0; j < len(found); j++ {

			if found[j].Name == actual[i] {
				match = 1
			}
		}
		if match == 0 {
			found = append(found[:], GetCurrentInfoBridge(actual[i]))
		}

	}

	return found
}
func GetCurrentInfoBridge(ifname string) BridgeInterfaces {

	var new_ifname BridgeInterfaces
	new_ifname.Name = ifname

	str := GetOutput("ip link show master "+ifname+" | grep -v link | awk '{print $2}'")

	new_ifname.Slaves = FormatInterfaceName(str)

	new_ifname.IpMode = "dhcp"
	new_ifname.IpAddress  = ""
	new_ifname.SubnetMask = ""

	return new_ifname
}