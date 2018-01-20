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

	StartTheInterfaces(File)


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
	//cmd.Wait()
}