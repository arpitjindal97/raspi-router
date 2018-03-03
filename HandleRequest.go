package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
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