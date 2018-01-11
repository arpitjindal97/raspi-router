package main

import (
	"os/exec"
	"time"
)

func StartTheInterfaces(file ConfigFile) {
	path := "/home/arpit/Desktop/workspace/angular/mdl/"

	exec.Command("systemctl","stop","wpa_supplicant").Start()
	exec.Command("systemctl","disable","wpa_supplicant").Start()

	exec.Command("systemctl","stop","dhcpcd").Start()
	exec.Command("systemctl","disable","dhcpcd").Start()

	exec.Command("systemctl","stop","hostapd").Start()
	exec.Command("systemctl","disable","hostapd").Start()

	exec.Command("systemctl","stop","dnsmasq").Start()
	exec.Command("systemctl","disable","dnsmasq").Start()

	kill("wpa_supplicant")
	kill("dhcpcd")
	kill("hostapd")
	kill("dnsmasq")

	time.Sleep(time.Second*2)

	for i := 0; i < len(file.NetworkInterfaces); i++ {

		inter := file.NetworkInterfaces[i]

		if inter.Name == "lo" {
			continue
		}

		exec.Command("ip","link","set",inter.Name,"up").Start()

		if inter.Mode == "default" {

			if inter.IsWifi == "true" {

				cmd := exec.Command("wpa_supplicant","-B","-i",inter.Name,"-c",path+"config/"+inter.Name+"_wpa.conf")
				cmd.Start()

			}

		} else {
			if inter.IsWifi == "true" {

				exec.Command("sh", "-c", "hostapd config/"+inter.Name+"_hostapd.conf").Start()
			}

			exec.Command("sh", "-c", "dnsmasq").Start()

		}



		if inter.IpModes == "dhcp" {

			exec.Command( "dhcpcd","-t","0",inter.Name).Start()
		} else {

			exec.Command("sh", "-c", "assign static ip addr").Start()
		}


	}

}



func kill(wpa string){

	c1 := exec.Command("pkill",wpa)
	c1.Start()
	c1.Wait()


}