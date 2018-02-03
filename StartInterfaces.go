package main

import (
	"time"
	"fmt"
	"os/exec"
	"github.com/godbus/dbus"
)

var path = "/home/arpit/Desktop/workspace/angular/mdl/"

func StartTheInterfaces(file ConfigFile) {
	//path="/home/pi/Desktop/"

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


	ExecuteWait("ip", "addr", "flush", "dev", inter.Name)
	ExecuteWait("ip", "route", "flush", "dev", inter.Name)

	if inter.Mode == "off" {
		return
	}
	if inter.IsWifi == "false" {
		EthStart(inter)
		return
	}

	//Wifi Interface
	dbus_objects[inter.Name] = make(chan *dbus.Signal, 10)


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
			go ExecuteWait("dhcpcd","-q","-w","-t","0",inter.Name)
			return
		}
		time.Sleep(time.Second * 5)
	}
}
func StopInterface (rec_interface Interfaces) {


	ExecuteWait("ip", "addr", "flush", "dev", rec_interface.Name)
	ExecuteWait("ip", "route", "flush", "dev", rec_interface.Name)

	if rec_interface.IsWifi == "true" {

		//if there is any change in wpa, hostapd,dnsmasq then restart
		if rec_interface.Mode == "hotspot" {

			Kill("hostapd.*" + rec_interface.Name)
			Kill("dnsmasq.*" + rec_interface.Name)

			//clear old rules
			IptablesClear(rec_interface)

		} else if rec_interface.Mode == "default" {
			DBusRemoveInterface(rec_interface.Name)

		}

	} else {
		Kill("dhcpcd.*" + rec_interface.Name)
		if rec_interface.Mode == "hotspot" {

			Kill("dnsmasq.*" + rec_interface.Name)

			//clear old rules
			IptablesClear(rec_interface)

		}
	}
}
