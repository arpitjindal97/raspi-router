package main

func IptablesCreate(inter PhysicalInterface) {

	if inter.Mode == "hotspot"{

		ExecuteWait("iptables","-t","nat","-A","POSTROUTING","-o",inter.NatInterface,"-j","MASQUERADE")

		ExecuteWait("iptables", "-A" ,"FORWARD" ,"-i", inter.NatInterface,"-o",
			inter.Name, "-m" ,"state" ,"--state" ,"RELATED,ESTABLISHED" ,"-j" ,"ACCEPT")

		ExecuteWait("iptables","-A","FORWARD","-i",inter.Name,"-o",inter.NatInterface ,"-j","ACCEPT")

	}
}

func IptablesClear(old_inter PhysicalInterface) {

	if old_inter.Mode == "hotspot" {
		ExecuteWait("iptables","-t","nat","-D","POSTROUTING","-o",old_inter.NatInterface,"-j","MASQUERADE")

		ExecuteWait("iptables", "-D" ,"FORWARD" ,"-i", old_inter.NatInterface,"-o",
			old_inter.Name, "-m" ,"state" ,"--state" ,"RELATED,ESTABLISHED" ,"-j" ,"ACCEPT")

		ExecuteWait("iptables","-D","FORWARD","-i",old_inter.Name,"-o",old_inter.NatInterface ,"-j","ACCEPT")

	}
}

func IptablesClearAll(){
	ExecuteWait("iptables","-F")
	ExecuteWait("iptables","-F","-t","nat")

}

func EnableNAT() {
	ExecuteWait("sh" ,"-c" ,"echo 1 > /proc/sys/net/ipv4/ip_forward")
}