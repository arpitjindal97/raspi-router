package main

import (
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"os"
)

type ConfigFile struct {
	NetworkInterfaces []Interfaces
}
type Interfaces struct{
	Name 			string
	IsWifi 			string
	Mode 			string
	RouteMode 		string
	RouteInterface 	string
	IpModes 		string
	IpAddress 		string
	SubnetMask 		string
	Wpa 			string
	Hostapd 		string
	Dnsmasq 		string
	Info			BasicInfo
}

func FirstTask() ConfigFile{

	out, _ := exec.Command("sh","-c","ip link | grep -v link | cut -f 2 -d ' '").Output()

	interface_names := GetInterfaceName(out)

	raw,err := ioutil.ReadFile("config/config.json")

	var file ConfigFile

	if err !=nil{
		bb,_ := json.MarshalIndent(file,"","	")

		_,err = ioutil.ReadDir("config")
		if err !=nil {
			os.Mkdir("config",0755)
		}

		ioutil.WriteFile("config/config.json",bb,0644)
		raw = bb
	}

	json.Unmarshal(raw,&file)

	for i:=0; i< len(file.NetworkInterfaces);i++ {

		name := file.NetworkInterfaces[i].Name
		match:=0
		for _,j := range interface_names{

			if name == j {
				match=1
			}
		}
		if match == 0{
			file.NetworkInterfaces = append(file.NetworkInterfaces[:i],file.NetworkInterfaces[i+1:]...)
			i--
		}
	}

	for i:=0;i< len(interface_names);i++{

		match:=0
		for j:=0;j<len(file.NetworkInterfaces);j++{

			if file.NetworkInterfaces[j].Name ==  interface_names[i] {
				match=1
			}
		}
		if match == 0{
			file.NetworkInterfaces = append(file.NetworkInterfaces[:],CreateDefaultInterface(interface_names[i]))
		}

	}

	b,_ := json.MarshalIndent(file,"","	")

	ioutil.WriteFile("config/config.json",b,0644)

	//b,_ = json.MarshalIndent(file,"","	")
	file = CaptureConfFiles(file)

	return file

}
func CreateDefaultInterface(name string) Interfaces {

	var name_default Interfaces
	name_default.Name = name

	str := GetOutput("iwconfig "+name)

	if str == ""{
		name_default.IsWifi = "false"
	} else {
		name_default.IsWifi = "true"
	}

	name_default.Mode = "default"
	name_default.RouteMode = "nat"
	name_default.RouteInterface = ""
	name_default.IpModes = "dhcp"

	return name_default
}
func CaptureConfFiles(file ConfigFile) ConfigFile{

	for i:=0;i<len(file.NetworkInterfaces);i++ {

		name:= file.NetworkInterfaces[i].Name

		// Dnsmasq
		raw,err := ioutil.ReadFile("config/"+name+"_dnsmasq.conf")
		if err !=nil {
			str:="interface="+name+"\n"+
			"bind-interfaces\n"+
			"server=8.8.8.8\n"+
			"domain-needed\n"+
			"bogus-priv\n"+
			"dhcp-range=192.168.2.2,192.168.2.100,12h"

			raw = []byte(str)

			ioutil.WriteFile("config/"+name+"_dnsmasq.conf",raw,os.FileMode(0644))

		}
		file.NetworkInterfaces[i].Dnsmasq = string(raw)

		if file.NetworkInterfaces[i].IsWifi=="false" {
			continue
		}

		// Hostapd
		raw,err = ioutil.ReadFile("config/"+name+"_hostapd.conf")
		if err !=nil {
			str:="interface="+name+"\n" +
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

			ioutil.WriteFile("config/"+name+"_hostapd.conf",raw,os.FileMode(0644))

		}
		file.NetworkInterfaces[i].Hostapd = string(raw)


		// Wpa Supplicant
		raw,err = ioutil.ReadFile("config/"+name+"_wpa.conf")
		if err !=nil {
			str:=	"ctrl_interface=/run/wpa_supplicant\n" +
					"update_config=1\n" +
					"network={\n" +
					"ssid=\"Your Wifi Name\"\n" +
					"psk=\"Your Password\"\n" +
					"}\n"

			raw = []byte(str)

			ioutil.WriteFile("config/"+name+"_wpa.conf",raw,os.FileMode(0644))

		}
		file.NetworkInterfaces[i].Wpa = string(raw)

	}

	return file
}


func GetInterfaceName(out []byte) []string {


	var interfaceNames []string
	count :=0
	interfaceNames = append(interfaceNames,"")

	for i:=0; i<len(string(out)) ; i++ {

		s := string(out[i])

		if s==":" && (out[i+1]) == 10 {
			count++
			if i+2 != len(string(out)) {
				interfaceNames = append(interfaceNames, "")
			}
			i++;
			continue
		}

		interfaceNames[count] = interfaceNames[count] + s

	}
	return interfaceNames
}