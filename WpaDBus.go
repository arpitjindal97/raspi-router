package main

import (
	"github.com/godbus/dbus"
	"log"
	"fmt"
)

func DBusCreateInterface(ifname string, driver string, config string, inter Interfaces) {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}

	obj := conn.Object("fi.w1.wpa_supplicant1", //Well known name on the bus (`busctl list` will show these)
		dbus.ObjectPath("/fi/w1/wpa_supplicant1")) //Object path (`busctl tree <well known name>` shows these)

	//Method name is the interface + the method name
	//a{sv} is the same as map[string]interface{}
	//o is the same as dbus.ObjectPath

	var intfPath dbus.ObjectPath

	err = obj.Call("fi.w1.wpa_supplicant1.CreateInterface", 0, map[string]interface{}{
		"Ifname":     ifname,
		"Driver":     driver,
		"ConfigFile": config,
	}).Store(&intfPath)

	//err = obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, "wlp0s29u1u2").Store(&intfPath)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(intfPath)

	//{"Ifname": __import__('gi.repository.GLib', globals(), locals(), ['Variant']).Variant("s","wlp0s29u1u2")}

	//wait for state completed
	DbusDhcpcdRoutine(inter)
}

func DBusRemoveInterface(ifname string) {

	fmt.Println("removing interface " + ifname)
	DbusStopDhcp(ifname)

	conn, err := dbus.SystemBus()
	if err != nil {

		panic(err)
	}
	obj := conn.Object("fi.w1.wpa_supplicant1", //Well known name on the bus (`busctl list` will show these)
		dbus.ObjectPath("/fi/w1/wpa_supplicant1")) //Object path (`busctl tree <well known name>` shows these)

	//Method name is the interface + the method name
	//a{sv} is the same as map[string]interface{}
	//o is the same as dbus.ObjectPath

	var intfPath dbus.ObjectPath

	err = obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, ifname).Store(&intfPath)
	if err != nil {
		panic(err)
	}

	log.Println(intfPath)

	obj.Call("fi.w1.wpa_supplicant1.RemoveInterface", 0, intfPath)

	fmt.Println("interface removed " + ifname)

}

func DbusDhcpcdRoutine(inter Interfaces) {


	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}

	obj := conn.Object("fi.w1.wpa_supplicant1", //Well known name on the bus (`busctl list` will show these)
		dbus.ObjectPath("/fi/w1/wpa_supplicant1")) //Object path (`busctl tree <well known name>` shows these)

	var intfPath dbus.ObjectPath

	obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, inter.Name).Store(&intfPath)

	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='"+string(intfPath)+"',interface='fi.w1.wpa_supplicant1.Interface',member='PropertiesChanged'")

	dbus_objects[inter.Name] = make(chan *dbus.Signal, 10)

	conn.Signal(dbus_objects[inter.Name])

	go func() {

		fmt.Println("routine called")

		for v := range dbus_objects[inter.Name] {

			var mm map[string]interface{}

			dbus.Store(v.Body, &mm)

			for key := range mm {

				if key == "Stop" {
					fmt.Println("stop signal received")
					return
				}

				if key != "State" {
					continue
				}

				fmt.Println(key + "\t" + mm[key].(string))

				if mm[key].(string) == "completed" {

					if inter.IpModes == "dhcp" {
						Systemctl("start", "dhcpcd@"+inter.Name)
					} else {
						ExecuteWait("ifconfig", inter.Name, inter.IpAddress, "netmask", inter.SubnetMask)
					}

					fmt.Println("dhcpcd routine done")
					return
				}

			}

		}
	}()
}

func DbusStopDhcp(ifname string) {
	mm := map[string]interface{}{"Stop": "completed"}

	m := []interface{}{mm}

	dd := dbus.Signal{
		Sender: "",
		Path:   "",
		Name:   "string",
		Body:   m,
	}

	op := GetOutput("ps aux | grep dhcpcd |grep -v grep | grep "+ifname)

	fmt.Println("output: "+op)

	if op == "" {

		fmt.Println("dhcpcd not running")
		dbus_objects[ifname] <- &dd

	}

	Systemctl("stop", "dhcpcd@"+ifname)
}
