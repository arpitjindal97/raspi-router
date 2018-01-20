package main

func IptablesCreate(inter Interfaces) {

	if inter.Mode == "hotspot"{

		ExecuteWait("iptables","-t","nat","-A","POSTROUTING","-o",inter.RouteInterface,"-j","MASQUERADE")

		ExecuteWait("iptables", "-A" ,"FORWARD" ,"-i", inter.RouteInterface,"-o",
			inter.Name, "-m" ,"state" ,"--state" ,"RELATED,ESTABLISHED" ,"-j" ,"ACCEPT")

		ExecuteWait("iptables","-A","FORWARD","-i",inter.Name,"-o",inter.RouteInterface ,"-j","ACCEPT")

	}
}

func IptablesClear(old_inter Interfaces) {

	if old_inter.Mode == "hotspot" {
		ExecuteWait("iptables","-t","nat","-D","POSTROUTING","-o",old_inter.RouteInterface,"-j","MASQUERADE")

		ExecuteWait("iptables", "-D" ,"FORWARD" ,"-i", old_inter.RouteInterface,"-o",
			old_inter.Name, "-m" ,"state" ,"--state" ,"RELATED,ESTABLISHED" ,"-j" ,"ACCEPT")

		ExecuteWait("iptables","-D","FORWARD","-i",old_inter.Name,"-o",old_inter.RouteInterface ,"-j","ACCEPT")

	}
}

func IptablesClearAll(){
	ExecuteWait("iptables","-F")
	ExecuteWait("iptables","-F","-t","nat")

}

func EnableNAT() {
	ExecuteWait("sh" ,"-c" ,"echo 1 > /proc/sys/net/ipv4/ip_forward")
}