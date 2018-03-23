package main

import (
	"net/http"
	"os/exec"
	"encoding/json"
	"github.com/godbus/dbus"
	"os"
	"fmt"
	"github.com/gorilla/mux"
	"time"
	"log"
)


// Packages needed wireless_tools, iw, net-tools

var File ConfigFile
var dbus_objects map[string] chan *dbus.Signal

var eth_thread map[string] string


var path = "/home/arpit/Desktop/workspace/angular/mdl/"

func main() {

	SetPath()

	File = FirstTask()
	muxHttp := mux.NewRouter()

	muxHttp.HandleFunc("/api/PhysicalInterfaceReconfigure",HandlePhysicalInterReconfigure)

	muxHttp.HandleFunc("/api/PhysicalInterfaces",Handle_PhysicalInterface)
	muxHttp.HandleFunc("/api/PhysicalInterfaces/{inter_name}",Handle_PhysicalInterfaceName)

	muxHttp.HandleFunc("/api/OSInfo",Handle_DeviceInfo)

	muxHttp.HandleFunc("/api/BridgeInterDelete",Handle_BridgeInterDelete)
	muxHttp.HandleFunc("/api/BridgeInterCreate",Handle_BridgeInterCreate)
	muxHttp.HandleFunc("/api/BridgeInterSave",Handle_BridgeInterSave)
	muxHttp.HandleFunc("/api/BridgeInterStart",Handle_BridgeInterStart)
	muxHttp.HandleFunc("/api/BridgeInterStop",Handle_BridgeInterStop)
	muxHttp.HandleFunc("/api/BridgeInterRemoveSlave",Handle_BridgeInterRemoveSlave)
	muxHttp.HandleFunc("/api/BridgeInterAddSlave",Handle_BridgeInterAddSlave)

	muxHttp.HandleFunc("/api",Index)

	dbus_objects = make(map[string] chan *dbus.Signal)
	eth_thread = make(map[string] string)

	StartTheInterfaces()

	//StartBridging()


	muxHttp.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	srv := &http.Server{
		Handler:      muxHttp,
		Addr:         "0.0.0.0:5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}


func Index(w http.ResponseWriter, r *http.Request) {

	go func(){File.OSInfo = GetDeviceInfo()}()

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
	}

	for i:=0;i<len(File.BridgeInterfaces);i++ {
		File.BridgeInterfaces[i].Info = GetCommonInterfaceInfo(File.BridgeInterfaces[i].Name)
	}


	b,_ := json.MarshalIndent(File,"","	")


	SetHeader(&w)

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

func SetHeader(w *http.ResponseWriter){

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func MakeJSON(str string) JSONResponse {
	var json JSONResponse
	json.Message = str
	return json
}