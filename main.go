package main

import (
	"net/http"
	"os/exec"
	"encoding/json"
)


// Packages needed wireless_tools, iw, net-tools

var File ConfigFile

func main() {


	File = FirstTask()

	http.HandleFunc("/api/interfaces",NetworkInterface)
	http.HandleFunc("/api/device_info",DeviceInfo)
	http.HandleFunc("/api/update_interface",UpdateInterface)
	http.HandleFunc("/api",Index)

	//StartTheInterfaces(File)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080",nil)

}


func Index(w http.ResponseWriter, r *http.Request) {
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
