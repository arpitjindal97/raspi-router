package main

import (
	"time"
	"os/exec"
	"github.com/godbus/dbus"
	"encoding/json"
	"io/ioutil"
	"os"
)

func StartTheInterfaces() {

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

	mylog.Println("Enabling Packet Forwarding")
	EnableNAT()
	mylog.Println("Clearing all existing rules of iptables")
	IptablesClearAll()

	time.Sleep(time.Second * 2)

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		mylog.Println("Starting up the interface " + File.PhysicalInterfaces[i].Name)

		ExecuteWait("ip", "link", "set", File.PhysicalInterfaces[i].Name, "up")
		ExecuteWait("ip", "link", "set", File.PhysicalInterfaces[i].Name, "nomaster")

		PhysicalInterStart(File.PhysicalInterfaces[i])
	}

	for i := 0; i < len(File.BridgeInterfaces); i++ {

		BridgeInterfaceCreate(File.BridgeInterfaces[i])

		BridgeInterfaceStart(File.BridgeInterfaces[i])

	}

}

func PhysicalInterStart(inter PhysicalInterface) string {

	if inter.Name == "lo" {
		mylog.Println("Ignoring " + inter.Name)
		return ""
	}

	mylog.Println("Flushing the existing IP addr and Route of " + inter.Name)

	ExecuteWait("ip", "addr", "flush", "dev", inter.Name)
	ExecuteWait("ip", "route", "flush", "dev", inter.Name)

	if inter.IsWifi == "false" {
		return PhysicalInterStartEth(inter)
	} else {
		return PhysicalInterStartWlan(inter)
	}

}

func PhysicalInterStartEth(inter PhysicalInterface) string {

	eth_thread[inter.Name] = "start"

	if inter.Mode == "default" {

		if inter.IpModes == "dhcp" {

			mylog.Println("Polling for Cable plugin on " + inter.Name)
			go PhysicalInterDhcpEth(inter)

		} else {
			//static Ip address

			mylog.Println("Static IP addr assigned to " + inter.Name)

			ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
			ExecuteWait("route", "add", "default", "gw", inter.Gateway, inter.Name)
		}
	} else if inter.Mode == "hotspot" {
		// Hotspot

		//static IP
		mylog.Println("Static IP addr assigned to " + inter.Name)

		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)

		//dnsmasq
		mylog.Println("Dnsmasq started on " + inter.Name)
		ExecuteWait("dnsmasq", "--user=root", "--interface="+inter.Name, "-C", GetPath()+inter.Name+"_dnsmasq.conf")

		//handle routing
		mylog.Println("Configuring IP Tables for " + inter.Name)
		IptablesCreate(inter)
	} else if inter.Mode == "bridge" {

		// Nothing to do
	} else {
		ExecuteWait("ip", "link", "set", inter.Name, "down")
	}

	return inter.Name + " started"
}
func PhysicalInterDhcpEth(inter PhysicalInterface) {

	for eth_thread[inter.Name] == "start" {

		carrier := GetOutput("cat /sys/class/net/" + inter.Name + "/carrier")
		if carrier == "1" {
			mylog.Println("Cable Plugged in on interface " + inter.Name)
			go ExecuteWait("dhcpcd", "-q", "-w", "-t", "0", inter.Name)
			return
		}
		time.Sleep(time.Second * 5)
	}
}

func PhysicalInterStartWlan(inter PhysicalInterface) string {
	//Wifi Interface
	dbus_objects[inter.Name] = make(chan *dbus.Signal, 10)

	if inter.Mode == "default" {

		mylog.Println("WPA Supplicant on " + inter.Name)
		DBusCreateInterface(inter.Name, "nl80211", GetPath()+inter.Name+"_wpa.conf", inter)

		if inter.IpModes == "dhcp" {
			DbusDhcpcdRoutine(inter)
		} else {

			mylog.Println("Static IP addr assigned to " + inter.Name)
			ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
			ExecuteWait("route", "add", "default", "gw", inter.Gateway, inter.Name)
		}

	} else if inter.Mode == "hotspot" {

		mylog.Println("Hostapd started on " + inter.Name)
		exec.Command("hostapd", GetPath()+inter.Name+"_hostapd.conf").Start()

		mylog.Println("Static IP addr assigned to " + inter.Name)
		ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)

		//time.Sleep(time.Second*2)

		mylog.Println("Dnsmasq started on " + inter.Name)
		ExecuteWait("dnsmasq", "--user=root", "--interface="+inter.Name, "-C", GetPath()+inter.Name+"_dnsmasq.conf")

		mylog.Println("Configuring IP Tables for " + inter.Name)
		IptablesCreate(inter)

	} else if inter.Mode == "bridge" {

		if inter.BridgeMode == "wpa" {

			mylog.Println("WPA Supplicant on " + inter.Name)
			DBusCreateInterface(inter.Name, "nl80211", GetPath()+inter.Name+"_wpa.conf", inter)

		} else if inter.BridgeMode == "hostapd" {

			mylog.Println("Hostapd started on " + inter.Name)
			exec.Command("hostapd", GetPath()+inter.Name+"_hostapd.conf").Start()
		}

	} else {

		ExecuteWait("ip", "link", "set", inter.Name, "down")

	}
	return inter.Name + " started"
}
func PhysicalInterStop(inter PhysicalInterface) string {

	if inter.IsWifi == "true" {

		if inter.Mode == "hotspot" {

			mylog.Println("Killing Hostapd and Dnsmasq of " + inter.Name)
			Kill("hostapd.*" + inter.Name)
			Kill("dnsmasq.*" + inter.Name)

			//clear old rules
			mylog.Println("Clearing IP table rules of " + inter.Name)
			IptablesClear(inter)

		} else if inter.Mode == "default" {
			DBusRemoveInterface(inter.Name)

		} else if inter.Mode == "bridge" {
			if inter.BridgeMode == "wpa" {

				DBusRemoveInterface(inter.Name)

			} else if inter.BridgeMode == "hostapd" {

				mylog.Println("Kiling Hostapd of " + inter.Name)
				Kill("hostapd.*" + inter.Name)
			}
		}

	} else {
		Kill("dhcpcd.*" + inter.Name)
		if inter.Mode == "hotspot" {

			mylog.Println("Killing Dnsmasq of " + inter.Name)
			Kill("dnsmasq.*" + inter.Name)

			//clear old rules
			mylog.Println("Clearing IP table rules of " + inter.Name)
			IptablesClear(inter)

		}
	}
	return inter.Name + " stopped"
}

func PhysicalInterSave(inter PhysicalInterface) string {

	var orig *PhysicalInterface

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		if File.PhysicalInterfaces[i].Name == inter.Name {
			orig = &File.PhysicalInterfaces[i]
			break
		}
	}

	*orig = inter
	(*orig).Hostapd = ""
	(*orig).Wpa = ""
	(*orig).Dnsmasq = ""
	(*orig).Info = BasicInfo{}

	File.OSInfo = OSInfo{}

	b, _ := json.MarshalIndent(File, "", "	")

	ioutil.WriteFile(GetPath()+"config.json", b, 0644)

	(*orig).Hostapd = inter.Hostapd
	(*orig).Wpa = inter.Wpa
	(*orig).Dnsmasq = inter.Dnsmasq

	if (*orig).IsWifi == "true" {

		raw := []byte(inter.Hostapd)
		ioutil.WriteFile(GetPath()+inter.Name+"_hostapd.conf", raw, os.FileMode(0644))

		raw = []byte(inter.Wpa)
		ioutil.WriteFile(GetPath()+inter.Name+"_wpa.conf", raw, os.FileMode(0644))
	}
	raw := []byte(inter.Dnsmasq)
	ioutil.WriteFile(GetPath()+inter.Name+"_dnsmasq.conf", raw, os.FileMode(0644))

	return "Configuration Saved"
}
