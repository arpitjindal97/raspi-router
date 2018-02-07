package main

import (
	"time"
	"os/exec"
	"github.com/godbus/dbus"
	"log"
)


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

	log.Println("Enabling Packet Forwarding")
	EnableNAT()
	log.Println("Clearing all existing rules of iptables")
	IptablesClearAll()

	time.Sleep(time.Second * 2)

	for i := 0; i < len(file.NetworkInterfaces); i++ {

		log.Println("Starting up the interface "+file.NetworkInterfaces[i].Name)

		ExecuteWait("ip", "link", "set", file.NetworkInterfaces[i].Name, "up")

		StartParticularInterface(file.NetworkInterfaces[i])
	}

}

func StartParticularInterface(inter Interfaces) {

	if inter.Name == "lo" {
		log.Println("Ignoring "+inter.Name)
		return
	}

	log.Println("Flushing the existing IP addr and Route of "+inter.Name)

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

		log.Println("WPA Supplicant on "+inter.Name)
		DBusCreateInterface(inter.Name, "nl80211", path+"config/"+inter.Name+"_wpa.conf", inter)

		if inter.IpModes == "dhcp"{
			DbusDhcpcdRoutine(inter)
		} else{

			log.Println("Static IP addr assigned to "+inter.Name)
			ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
		}

	} else {

		log.Println("Hostapd started on "+inter.Name)
		exec.Command("hostapd", path+"config/"+inter.Name+"_hostapd.conf").Start()

		log.Println("Static IP addr assigned to "+inter.Name)
		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)

		//time.Sleep(time.Second*2)

		log.Println("Dnsmasq started on "+inter.Name)
		ExecuteWait("dnsmasq", "--user=root", "--interface="+inter.Name, "-C", path+"config/"+inter.Name+"_dnsmasq.conf")

		log.Println("Configuring IP Tables for "+inter.Name)
		IptablesCreate(inter)

	}
}

func EthStart(inter Interfaces) {

	eth_thread[inter.Name] = "start"

	if inter.Mode == "default" {

		if inter.IpModes == "dhcp" {

			log.Println("Polling for Cable plugin on "+ inter.Name)
			go EthDhcp(inter)

		} else {
			//static Ip address

			log.Println("Static IP addr assigned to "+inter.Name)

			ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
		}
	} else {
		// Hotspot

		//static IP
		log.Println("Static IP addr assigned to "+inter.Name)

		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)

		//dnsmasq
		log.Println("Dnsmasq started on "+inter.Name)
		ExecuteWait("dnsmasq", "--user=root", "--interface="+inter.Name, "-C", path+"config/"+inter.Name+"_dnsmasq.conf")

		//handle routing
		log.Println("Configuring IP Tables for "+inter.Name)
		IptablesCreate(inter)
	}
}
func EthDhcp(inter Interfaces){

	for eth_thread[inter.Name] == "start" {

		carrier := GetOutput("cat /sys/class/net/" + inter.Name + "/carrier")
		if carrier == "1" {
			log.Println("Cable Plugged in on interface "+ inter.Name)
			go ExecuteWait("dhcpcd","-q","-w","-t","0",inter.Name)
			return
		}
		time.Sleep(time.Second * 5)
	}
}
func StopInterface (rec_interface Interfaces) {

	log.Println("Flushing the existing IP addr and route of "+rec_interface.Name)

	ExecuteWait("ip", "addr", "flush", "dev", rec_interface.Name)
	ExecuteWait("ip", "route", "flush", "dev", rec_interface.Name)

	if rec_interface.IsWifi == "true" {

		//if there is any change in wpa, hostapd,dnsmasq then restart
		if rec_interface.Mode == "hotspot" {

			log.Println("Kiling Hostapd and Dnsmasq of "+ rec_interface.Name)
			Kill("hostapd.*" + rec_interface.Name)
			Kill("dnsmasq.*" + rec_interface.Name)

			//clear old rules
			log.Println("Clearing IP table rules of "+rec_interface.Name)
			IptablesClear(rec_interface)

		} else if rec_interface.Mode == "default" {
			DBusRemoveInterface(rec_interface.Name)

		}

	} else {
		Kill("dhcpcd.*" + rec_interface.Name)
		if rec_interface.Mode == "hotspot" {

			log.Println("Kiling Dnsmasq of "+ rec_interface.Name)
			Kill("dnsmasq.*" + rec_interface.Name)

			//clear old rules
			log.Println("Clearing IP table rules of "+rec_interface.Name)
			IptablesClear(rec_interface)

		}
	}
}
