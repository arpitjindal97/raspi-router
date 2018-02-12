package main

import "strings"

func StartBridging() {

	for _,i := range File.BridgeInterfaces {

		out := GetOutput("ip link ls "+i.Name)

		if strings.Contains(out,"does not exist") == true {
			GetOutput("ip link add name "+i.Name+" type bridge")
		}

		GetOutput("ip link set dev "+i.Name+" up")

		for _,j := range i.Slaves {

			GetOutput("ip link set "+j+" master "+i.Name)
		}

		if i.IpMode == "dhcp" {

			go ExecuteWait("dhcpcd", "-q", "-w", "-t", "0", i.Name)
		} else {

			ExecuteWait("ifconfig", i.Name, i.IpAddress, "netmask", i.SubnetMask)
		}
	}
}