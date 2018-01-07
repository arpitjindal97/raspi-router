package main

import (
	"os/exec"
	"fmt"
)

func StartTheInterfaces(file ConfigFile) {
	path := "/home/arpit/Desktop/workspace/GoglandProjects/raspbian-router-server/"

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

				fmt.Println(cmd.Process.Pid)


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
