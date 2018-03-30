package main

import (
	"net/http"
	"os/exec"
	"encoding/json"
	"github.com/godbus/dbus"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"time"
	"os"
)

// Packages needed wireless_tools, iw, net-tools

var File ConfigFile
var dbus_objects map[string]chan *dbus.Signal

var eth_thread map[string]string

var mylog *log.Logger

var logpath = "/var/log/raspi-router/log.txt"

func main() {


	err := os.Remove(logpath)

	if err != nil {

		err = os.Mkdir("/var/log/raspi-router",0755)

	}
	file, _ := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	mylog = log.New(file, "", log.LstdFlags|log.Lshortfile)

	log.Println("Check logs at: "+logpath)
	File = FirstTask()

	dbus_objects = make(map[string]chan *dbus.Signal)
	eth_thread = make(map[string]string)

	go StartTheInterfaces()

	muxHttp := mux.NewRouter()

	muxHttp.HandleFunc("/api", Index).Methods("GET")
	muxHttp.HandleFunc("/api/OSInfo", GetDeviceInfo).Methods("GET")
	muxHttp.HandleFunc("/api/PhysicalInterfaces", GetAllPhysicalInterfaces).Methods("GET")

	muxHttp.HandleFunc("/api/PhysicalInterfaces/{inter_name}", GetPhysicalInterface).Methods("GET")
	muxHttp.HandleFunc("/api/PhysicalInterfaces/{inter_name}", PutPhysicalInterface).Methods("PUT")

	muxHttp.HandleFunc("/api/BridgeInterDelete", Handle_BridgeInterDelete)
	muxHttp.HandleFunc("/api/BridgeInterCreate", Handle_BridgeInterCreate)
	muxHttp.HandleFunc("/api/BridgeInterSave", Handle_BridgeInterSave)
	muxHttp.HandleFunc("/api/BridgeInterStart", Handle_BridgeInterStart)
	muxHttp.HandleFunc("/api/BridgeInterStop", Handle_BridgeInterStop)
	muxHttp.HandleFunc("/api/BridgeInterRemoveSlave", Handle_BridgeInterRemoveSlave)
	muxHttp.HandleFunc("/api/BridgeInterAddSlave", Handle_BridgeInterAddSlave)

	muxHttp.PathPrefix("/").Handler(http.StripPrefix("/", http.HandlerFunc(GetStaticFiles))).Methods("GET")

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "PUT"},
		AllowCredentials: true,
	})

	handler := c.Handler(muxHttp)

	srv := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:500",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	mylog.Fatal(srv.ListenAndServe())

}

func Index(w http.ResponseWriter, r *http.Request) {

	go func() { File.OSInfo = DeviceInfo() }()

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
	}

	for i := 0; i < len(File.BridgeInterfaces); i++ {
		File.BridgeInterfaces[i].Info = GetCommonInterfaceInfo(File.BridgeInterfaces[i].Name)
	}

	b, _ := json.MarshalIndent(File, "", "	")

	w.Write(b)
}

func GetOutput(command string) string {
	out, _ := exec.Command("sh", "-c", command).Output()
	str := string(out)
	if str == "" {
		return ""
	}
	return string(str[:(len(str) - 1)])
}

func Kill(pattern string) {

	bb := exec.Command("sh", "-c", "ps aux | grep -E '"+pattern+"' | grep -v grep | awk '{print $2}'")

	b, _ := bb.Output()

	pid := string(b)
	if pid == "" {
		return
	}
	pid = pid[:len(pid)-1]

	bb = exec.Command("kill", pid)
	bb.Start()
	bb.Wait()
}
func PKill(wpa string) {

	c1 := exec.Command("pkill", wpa)
	c1.Start()
	c1.Wait()

}

func ExecuteWait(name string, arg ...string) {

	cmd := exec.Command(name, arg...)
	cmd.Start()
	cmd.Wait()
}

func Systemctl(action string, service_name string) {
	cmd := exec.Command("systemctl", action, service_name)
	cmd.Start()
	go cmd.Wait()
}

func GetPath() string {

	return "/etc/raspi-router/"
}

func MakeJSON(str string) JSONResponse {
	var json JSONResponse
	json.Message = str
	return json
}
