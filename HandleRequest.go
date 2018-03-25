package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gobuffalo/packr"
)


func HandlePhysicalInterReconfigure(w http.ResponseWriter, r *http.Request) {

	var inter PhysicalInterface

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	for i:=0;i<len(File.PhysicalInterfaces);i++ {

		if File.PhysicalInterfaces[i].Name == inter.Name {

			_ = PhysicalInterStop(File.PhysicalInterfaces[i])
			break
		}
	}

	response := PhysicalInterSave(inter)

	response = PhysicalInterStart(inter)

	SetHeader(&w)
	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}



func Handle_BridgeInterDelete(w http.ResponseWriter,r *http.Request) {

	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterDelete(inter)

	SetHeader(&w)
	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)

}


func Handle_BridgeInterCreate(w http.ResponseWriter,r *http.Request) {
	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterCreate(inter)

	SetHeader(&w)
	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
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

	SetHeader(&w)
	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}

func Handle_BridgeInterStart(w http.ResponseWriter,r *http.Request) {
	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterStart(inter)

	SetHeader(&w)
	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}


func Handle_BridgeInterStop(w http.ResponseWriter,r *http.Request) {
	var inter BridgeInterfaces

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterStop(inter)

	SetHeader(&w)

	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}
func Handle_BridgeInterRemoveSlave(w http.ResponseWriter,r *http.Request) {

	var inter_name string

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter_name)
	if err != nil {
		panic(err)
	}

	response := BridgeInterRemoveSlave(inter_name)

	SetHeader(&w)
	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}
func Handle_BridgeInterAddSlave(w http.ResponseWriter,r *http.Request) {
	var inter BridgeSlave

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&inter)
	if err != nil {
		panic(err)
	}

	response := BridgeInterAddSlave(inter.BridgeIfname,inter.SlaveIfname)

	SetHeader(&w)

	b, _ := json.MarshalIndent(MakeJSON(response), "", "	")
	w.Write(b)
}

func Handle_PhysicalInterfaceName(w http.ResponseWriter,r *http.Request) {

	SetHeader(&w)

	vars := mux.Vars(r)["inter_name"]

	for i,item := range File.PhysicalInterfaces {

		if item.Name == vars{

			File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
			b, _ := json.MarshalIndent(item, "", "	")

			w.Write(b)
			return
		}
	}


}

func Handle_PhysicalInterface(w http.ResponseWriter, r *http.Request) {

	for i := 0; i < len(File.PhysicalInterfaces); i++ {

		File.PhysicalInterfaces[i].Info = GetPhysicalInterfaceInfo(File.PhysicalInterfaces[i])
	}

	b, _ := json.MarshalIndent(File.PhysicalInterfaces, "", "	")


	SetHeader(&w)

	w.Write(b)
}
func Handle_StaticFiles(w http.ResponseWriter, r *http.Request) {

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