package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/gorilla/mux"
)

func HandlePhysicalInterStart(w http.ResponseWriter, r *http.Request) {

	var inter PhysicalInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := PhysicalInterStart(inter)

	w.Write([]byte(response))
}


func HandlePhysicalInterStop(w http.ResponseWriter, r *http.Request) {

	var inter PhysicalInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := PhysicalInterStop(inter)

	w.Write([]byte(response))
}

func HandlePhysicalInterSave (w http.ResponseWriter,r *http.Request) {
	var inter PhysicalInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	var orig *PhysicalInterfaces

	for i:=0;i<len(File.PhysicalInterfaces);i++ {

		if File.PhysicalInterfaces[i].Name == inter.Name {
			orig = &File.PhysicalInterfaces[i]
			break
		}
	}

	*orig = inter
	(*orig).Hostapd = ""
	(*orig).Wpa = ""
	(*orig).Dnsmasq = ""
	(*orig).Info = BasicInfo{}

	File.OSInfo = OSInfo{}

	b, _ := json.MarshalIndent(File, "", "	")

	ioutil.WriteFile("config/config.json", b, 0644)

	(*orig).Hostapd = inter.Hostapd
	(*orig).Wpa = inter.Wpa
	(*orig).Dnsmasq = inter.Dnsmasq

	if (*orig).IsWifi == "true" {

		raw := []byte(inter.Hostapd)
		ioutil.WriteFile("config/"+inter.Name+"_hostapd.conf", raw, os.FileMode(0644))

		raw = []byte(inter.Wpa)
		ioutil.WriteFile("config/"+inter.Name+"_wpa.conf", raw, os.FileMode(0644))
	}
	raw := []byte(inter.Dnsmasq)
	ioutil.WriteFile("config/"+inter.Name+"_dnsmasq.conf", raw, os.FileMode(0644))

	w.Write([]byte("Configuration saved"))
}



//Handlers for Bridge Interfaces

type BridgeSlave struct {
	BridgeIfname string
	SlaveIfname  string
}

func Handle_BridgeInterDelete(w http.ResponseWriter,r *http.Request) {

	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterDelete(inter)

	w.Write([]byte(response))

}


func Handle_BridgeInterCreate(w http.ResponseWriter,r *http.Request) {
	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterCreate(inter)

	w.Write([]byte(response))
}

func Handle_BridgeInterSave(w http.ResponseWriter,r *http.Request) {

	type BridgeSaveActions struct{
		BridgeInter BridgeInterfaces
		Action		string
	}
	var resp BridgeSaveActions

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&resp)
	if err != nil {
		panic(err)
	}

	response := BridgeInterSave(resp.BridgeInter,resp.Action)

	w.Write([]byte(response))
}

func Handle_BridgeInterStart(w http.ResponseWriter,r *http.Request) {
	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterStart(inter)

	w.Write([]byte(response))
}


func Handle_BridgeInterStop(w http.ResponseWriter,r *http.Request) {
	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterStop(inter)

	w.Write([]byte(response))
}
func Handle_BridgeInterRemoveSlave(w http.ResponseWriter,r *http.Request) {

	var inter_name string

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter_name)
	if err != nil {
		panic(err)
	}

	response := BridgeInterRemoveSlave(inter_name)

	w.Write([]byte(response))
}
func Handle_BridgeInterAddSlave(w http.ResponseWriter,r *http.Request) {
	var inter BridgeSlave

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterAddSlave(inter.BridgeIfname,inter.SlaveIfname)

	w.Write([]byte(response))
}

func Handle_PhysicalInterfaceName(w http.ResponseWriter,r *http.Request) {
	vars := mux.Vars(r)["inter_name"]

	for _,item := range File.PhysicalInterfaces {

		if item.Name == vars{
			b, _ := json.MarshalIndent(item, "", "	")
			w.Write([]byte(b))
			return
		}
	}


}
