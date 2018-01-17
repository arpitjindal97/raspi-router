package main

import (
	"time"
	"fmt"
)

func StartTheInterfaces(file ConfigFile) {

	Systemctl("stop","wpa_supplicant")
	Systemctl("disable","wpa_supplicant")

	Systemctl("stop","dhcpcd")
	Systemctl("disable","dhcpcd")

	Systemctl("stop","hostapd")
	Systemctl("disable","hostapd")

	Systemctl("stop","dnsmasq")
	Systemctl("disable","dnsmasq")

	PKill("wpa_supplicant")
	PKill("dhcpcd")
	PKill("hostapd")
	PKill("dnsmasq")

	time.Sleep(time.Second*2)

	for i := 0; i < len(file.NetworkInterfaces); i++ {

		ExecuteWait("ip","link","set",file.NetworkInterfaces[i].Name,"up")

		StartParticularInterface(file.NetworkInterfaces[i])
	}

}

func StartParticularInterface(inter Interfaces) {


	path := "/home/arpit/Desktop/workspace/angular/mdl/"

	if inter.Name == "lo" || inter.Name == "enp7s0"{
		return
	}

	ExecuteWait("ip","addr","flush","dev",inter.Name)
	ExecuteWait("ip","route","flush","dev",inter.Name)

	if inter.Mode == "default" {

		if inter.IsWifi == "true" {

			ExecuteWait("wpa_supplicant","-B","-i",inter.Name,"-c",path+"config/"+inter.Name+"_wpa.conf")

		}

	} else {
		if inter.IsWifi == "true" {

			//exec.Command("sh", "-c", "hostapd config/"+inter.Name+"_hostapd.conf").Start()
			fmt.Println("starting hostapd on "+inter.Name)
		}

		//exec.Command("sh", "-c", "dnsmasq").Start()

	}


	if inter.IpModes == "dhcp" {

		//exec.Command( "dhcpcd","-t","0",inter.Name).Start()

		Systemctl("start","dhcpcd@"+inter.Name)

	} else {

		//exec.Command("sh", "-c", "assign static ip addr").Start()
	}

}
