package main

import (
	"net/http"
	"os/exec"
	"encoding/json"
	"github.com/godbus/dbus"
	"os"
	"fmt"
)


// Packages needed wireless_tools, iw, net-tools

var File ConfigFile
var dbus_objects map[string] chan *dbus.Signal

var eth_thread map[string] string


var path = "/home/arpit/Desktop/workspace/angular/mdl/"

func main() {

	SetPath()

	File = FirstTask()

	http.HandleFunc("/api/PhysicalInterfaceStart",HandlePhysicalInterStart)
	http.HandleFunc("/api/PhysicalInterfaceStop",HandlePhysicalInterStop)
	http.HandleFunc("/api/PhysicalInterfaceSave",HandlePhysicalInterSave)
	http.HandleFunc("/api/PhysicalInterfaces",PhysicalInterface)
	http.HandleFunc("/api/DeviceInfo",DeviceInfo)
	http.HandleFunc("/api/UpdateInterface",UpdateInterface)
	http.HandleFunc("/api/CreateBridge",CreateBridge)
	http.HandleFunc("/api/DeleteBridge",DeleteBridge)
	http.HandleFunc("/api/BridgeRemoveSlave",BridgeRemoveSlave)
	http.HandleFunc("/api/BridgeAddSlave",BridgeAddSlave)
	http.HandleFunc("/api",Index)

	dbus_objects = make(map[string] chan *dbus.Signal)
	eth_thread = make(map[string] string)

	StartTheInterfaces()

	//StartBridging()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.ListenAndServe(":5000",nil)

}


func Index(w http.ResponseWriter, r *http.Request) {

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
	}

	for i:=0;i<len(File.BridgeInterfaces);i++ {
		File.BridgeInterfaces[i].Info = GetCommonInterfaceInfo(File.BridgeInterfaces[i].Name)
	}


	b,_ := json.MarshalIndent(File,"","	")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(b)
}


func GetOutput(command string) string {
	out, _ := exec.Command("sh", "-c",command).Output()
	str:=string(out)
	if str == ""{
		return ""
	}
	return string(str[:(len(str)-1)])
}

func Kill(pattern string){


	bb := exec.Command("sh","-c","ps aux | grep -E '"+pattern+"' | grep -v grep | awk '{print $2}'")

	b,_ := bb.Output()

	pid := string(b)
	if pid == ""{
		return
	}
	pid = pid[:len(pid)-1]

	bb = exec.Command("kill",pid)
	bb.Start()
	bb.Wait()
}
func PKill(wpa string){

	c1 := exec.Command("pkill",wpa)
	c1.Start()
	c1.Wait()

}

func ExecuteWait(name string, arg ...string){

	cmd := exec.Command(name,arg...)
	cmd.Start()
	cmd.Wait()
}

func Systemctl(action string,service_name string) {
	cmd := exec.Command("systemctl",action,service_name)
	cmd.Start()
	go cmd.Wait()
}

func SetPath() {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	path=pwd+"/"

}