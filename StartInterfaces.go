package main

import (
	"os/exec"
	"time"
	"fmt"
)

func StartTheInterfaces(file ConfigFile) {

	ExecuteWait("systemctl","stop","wpa_supplicant")
	ExecuteWait("systemctl","disable","wpa_supplicant")

	ExecuteWait("systemctl","stop","dhcpcd")
	ExecuteWait("systemctl","disable","dhcpcd")

	ExecuteWait("systemctl","stop","hostapd")
	ExecuteWait("systemctl","disable","hostapd")

	ExecuteWait("systemctl","stop","dnsmasq")
	ExecuteWait("systemctl","disable","dnsmasq")

	kill("wpa_supplicant")
	kill("dhcpcd")
	kill("hostapd")
	kill("dnsmasq")

	time.Sleep(time.Second*2)

	for i := 0; i < len(file.NetworkInterfaces); i++ {

		ExecuteWait("ip","link","set",file.NetworkInterfaces[i].Name,"up")

		StartParticularInterface(file.NetworkInterfaces[i])
	}

}

func StartParticularInterface(inter Interfaces) {


	path := "/home/arpit/Desktop/workspace/angular/mdl/"

	if inter.Name == "lo" {
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

		exec.Command( "dhcpcd","-t","0",inter.Name).Start()

	} else {

		exec.Command("sh", "-c", "assign static ip addr").Start()
	}

}

func kill(wpa string){

	c1 := exec.Command("pkill",wpa)
	c1.Start()
	c1.Wait()

}

func ExecuteWait(name string, arg ...string){

	cmd := exec.Command(name,arg...)
	cmd.Start()
	cmd.Wait()
}