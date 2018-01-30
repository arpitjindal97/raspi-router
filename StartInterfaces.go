package main

import (
	"time"
	"fmt"
	"os/exec"
	"github.com/godbus/dbus"
)

var path = "/home/arpit/Desktop/workspace/angular/mdl/"

func StartTheInterfaces(file ConfigFile) {

	Systemctl("stop", "wpa_supplicant")
	//Systemctl("disable","wpa_supplicant")

	Systemctl("stop", "dhcpcd")
	//Systemctl("disable","dhcpcd")

	Systemctl("stop", "hostapd")
	//Systemctl("disable","hostapd")

	Systemctl("stop", "dnsmasq")
	//Systemctl("disable","dnsmasq")

	PKill("wpa_supplicant")
	PKill("dhcpcd")
	PKill("hostapd")
	PKill("dnsmasq")

	EnableNAT()
	IptablesClearAll()

	time.Sleep(time.Second * 2)

	for i := 0; i < len(file.NetworkInterfaces); i++ {

		ExecuteWait("ip", "link", "set", file.NetworkInterfaces[i].Name, "up")

		StartParticularInterface(file.NetworkInterfaces[i])
	}

}

func StartParticularInterface(inter Interfaces) {

	if inter.Name == "lo" {
		return
	}

	if inter.IsWifi == "false" {
		EthStart(inter)
		return
	}

	//Wifi Interface
	dbus_objects[inter.Name] = make(chan *dbus.Signal, 10)

	ExecuteWait("ip", "addr", "flush", "dev", inter.Name)
	ExecuteWait("ip", "route", "flush", "dev", inter.Name)

	if inter.Mode == "default" {

		DBusCreateInterface(inter.Name, "nl80211", path+"config/"+inter.Name+"_wpa.conf", inter)

		if inter.IpModes == "dhcp"{
			DbusDhcpcdRoutine(inter)
		} else{

			ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
		}

	} else {

		exec.Command("hostapd", path+"config/"+inter.Name+"_hostapd.conf").Start()

		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)

		//time.Sleep(time.Second*2)
		ExecuteWait("dnsmasq", "--user=root", "--interface="+inter.Name, "-C", path+"config/"+inter.Name+"_dnsmasq.conf")
		fmt.Println("starting hostapd on " + inter.Name)

		IptablesCreate(inter)

	}
}

func EthStart(inter Interfaces) {

	ExecuteWait("ip", "addr", "flush", "dev", inter.Name)
	ExecuteWait("ip", "route", "flush", "dev", inter.Name)

	eth_thread[inter.Name] = "start"

	if inter.Mode == "default" {

		if inter.IpModes == "dhcp" {

			go EthDhcp(inter)

		} else {
			//static Ip address
			ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
		}
	} else {
		// Hotspot

		//static IP
		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)

		//dnsmasq
		ExecuteWait("dnsmasq", "--user=root", "--interface="+inter.Name, "-C", path+"config/"+inter.Name+"_dnsmasq.conf")

		//handle routing
		IptablesCreate(inter)
	}
}
func EthDhcp(inter Interfaces){

	for eth_thread[inter.Name] == "start" {

		carrier := GetOutput("cat /sys/class/net/" + inter.Name + "/carrier")
		if carrier == "1" {
			Systemctl("start", "dhcpcd@"+inter.Name)
			return
		}
		time.Sleep(time.Second * 5)
	}
}