package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gobuffalo/packr"
	"log"
)

func GetAllPhysicalInterfaces(w http.ResponseWriter, r *http.Request) {

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
	}

	b, _ := json.MarshalIndent(File.PhysicalInterfaces, "", "	")

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
}
func GetPhysicalInterface(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type","application/json")
	vars := mux.Vars(r)["inter_name"]

	for i, item := range File.PhysicalInterfaces {

		if item.Name == vars {

			File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
			b, _ := json.MarshalIndent(item, "", "	")

			w.Write(b)
			return
		}
	}

}

func PutPhysicalInterface(w http.ResponseWriter, r *http.Request) {

	var inter PhysicalInterface

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		if File.PhysicalInterfaces[i].Name == inter.Name {

			_ = PhysicalInterStop(File.PhysicalInterfaces[i])
			break
		}
	}

	response := PhysicalInterSave(inter)

	response = PhysicalInterStart(inter)

	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
}


func GetAllBridgeInterfaces(w http.ResponseWriter, r *http.Request) {

	for i := 0; i < len(File.BridgeInterfaces); i++ {
		File.BridgeInterfaces[i].Info = GetCommonInterfaceInfo(File.BridgeInterfaces[i].Name)
	}

	b, _ := json.MarshalIndent(File.BridgeInterfaces, "", "	")

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
}

func GetBridgeInterfaces(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["inter_name"]

	var b []byte
	b = nil
	for  _,item := range File.BridgeInterfaces {
		if  item.Name == name {
			b,_ = json.MarshalIndent(item, "","	")
			break
		}
	}

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
}

func PutBridgeInterfaces(w http.ResponseWriter, r *http.Request) {
	var resp BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&resp)
	if err != nil {
		panic(err)
	}

	var b []byte
	b = nil
	log.Println("Creating Bridge ")
	response := BridgeInterfaceCreate(resp)
	_ = BridgeInterfaceStart(resp)

	b, _ = json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}

func PatchBridgeInterfaces(w http.ResponseWriter, r *http.Request) {
	var resp BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&resp)
	if err != nil {
		panic(err)
	}

	for _,item := range File.BridgeInterfaces {
		if item.Name == resp.Name {
			_ = BridgeInterfaceStop(item)
			break
		}
	}

	var b []byte
	b = nil

	response := BridgeInterfaceUpdate(resp)
	response = BridgeInterfaceStart(resp)

	b, _ = json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}
func DeleteBridgeInterfaces(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["inter_name"]

	var response string

	for _,item := range File.BridgeInterfaces {
		if item.Name == name {
			_ = BridgeInterfaceStop(item)
			response = BridgeInterfaceDelete(item)
			break
		}
	}


	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
}


func GetStaticFiles(w http.ResponseWriter, r *http.Request) {

	box := packr.NewBox("./dist")

	html, err := box.MustBytes(r.URL.Path)

	if err != nil || r.URL.Path == "" {
		w.Write(box.Bytes("index.html"))
		return
	}

	if r.URL.Path[len(r.URL.Path)-3:] == "css" {
		w.Header().Set("Content-Type", "text/css")
	}

	w.Write(html)
}
func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {

	File.OSInfo = DeviceInfo()

	b, _ := json.MarshalIndent(File.OSInfo, "", "	")

	w.Header().Set("Content-Type","application/json")
	w.Write(b)

}
