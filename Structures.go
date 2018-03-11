package main

type ConfigFile struct {
	OSInfo				OSInfo
	PhysicalInterfaces	[]PhysicalInterfaces
	BridgeInterfaces	[]BridgeInterfaces
}

type OSInfo struct{
	DistributionId	string
	Description		string
	Release			string
	Codename		string
	Hostname		string
	KernelRelease	string
	Architecture	string
	ModelName		string
	CPUs			string
	LocalTime		string
	TimeZone		string
	UpTime			string
	UpSince			string
}

type PhysicalInterfaces struct {
	Name        	string
	IsWifi     		string
	Mode       		string
	BridgeMode		string
	BridgeMaster	string
	NatInterface	string
	IpModes        	string
	IpAddress      	string
	SubnetMask     	string
	Wpa            	string
	Hostapd        	string
	Dnsmasq        	string
	Info           	BasicInfo
}

type BridgeInterfaces struct {
	Name			string
	IpMode			string
	IpAddress      	string
	SubnetMask     	string
	Info           	BasicInfo
	Slaves			[]string
}


type BasicInfo struct {
IpAddress        string
BroadcastAddress string
Gateway          string
MacAddress       string
RecvBytes        string
RecvPackts       string
TransBytes       string
TransPackts      string

ConntectedTo string
ApMacAddr    string
BitRate      string
Frequency    string
LinkQuality  string
Channel      string
}


type BridgeSlave struct {
	BridgeIfname string
	SlaveIfname  string
}
