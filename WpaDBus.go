package main

import (
	"github.com/godbus/dbus"
	"fmt"
)

func DBusCreateInterface(ifname string, driver string, config string, inter PhysicalInterface) {

	conn,err := dbus.SystemBusPrivate()
	if err != nil {
		panic(err)
	}
	err = conn.Auth(nil)
	err = conn.Hello()
	//fmt.Println("reached here")

	obj := conn.Object("fi.w1.wpa_supplicant1", //Well known name on the bus (`busctl list` will show these)
		dbus.ObjectPath("/fi/w1/wpa_supplicant1")) //Object path (`busctl tree <well known name>` shows these)


	var intfPath dbus.ObjectPath

	err = obj.Call("fi.w1.wpa_supplicant1.CreateInterface", 0, map[string]interface{}{
		"Ifname":     ifname,
		"Driver":     driver,
		"ConfigFile": config,
	}).Store(&intfPath)


	if err != nil {
		fmt.Println(config)
		mylog.Println("Looks like there is some problem with wpa_supplicant, Is it installed?")
		panic(err)
	}

	mylog.Println(ifname+" Object path : "+string(intfPath))

	conn.Close()
}

func DBusRemoveInterface(ifname string) {

	mylog.Println("Removing "+ifname+" from WPA Supplicant")

	//stop running dhcpcd
	DbusStopDhcp(ifname)

	conn,_ := dbus.SystemBusPrivate()
	conn.Auth(nil)
	conn.Hello()
	defer conn.Close()

	obj := conn.Object("fi.w1.wpa_supplicant1",
		dbus.ObjectPath("/fi/w1/wpa_supplicant1"))

	var intfPath dbus.ObjectPath

	err := obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, ifname).Store(&intfPath)
	if err != nil {
		mylog.Println("No interface found with this name")
		return
	}

	mylog.Println(ifname+" Object path : "+string(intfPath))

	obj.Call("fi.w1.wpa_supplicant1.RemoveInterface", 0, intfPath)

	mylog.Println(ifname +" Removed")

}

func DbusDhcpcdRoutine(inter PhysicalInterface) {

	conn,_ := dbus.SystemBusPrivate()
	conn.Auth(nil)
	conn.Hello()

	obj := conn.Object("fi.w1.wpa_supplicant1",
		dbus.ObjectPath("/fi/w1/wpa_supplicant1"))

	var intfPath dbus.ObjectPath

	obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, inter.Name).Store(&intfPath)

	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='"+string(intfPath)+"',interface='fi.w1.wpa_supplicant1.Interface',member='PropertiesChanged'")

	//dbus_objects[inter.Name] = make(chan *dbus.Signal, 10)

	conn.Signal(dbus_objects[inter.Name])


	go func() {

		if DbusFetchProperty(inter) == "completed" {

			go ExecuteWait("dhcpcd","-q","-w",inter.Name)

			return
		}

		mylog.Println("Listening on "+inter.Name+" to get connected to a Network")

		outer :for v := range dbus_objects[inter.Name] {

			var mm map[string]interface{}

			dbus.Store(v.Body, &mm)

			for key := range mm {

				if key == "Stop" {
					mylog.Println("Stoping the listener on "+inter.Name)
					break outer
				}
				//fmt.Print(key + "\t")

				//fmt.Println(mm[key])
				if key != "State" {
					continue
				}

				mylog.Println(inter.Name+"\t"+key + "\t" + mm[key].(string))

				if mm[key].(string) == "completed" {

					mylog.Println("Started dhcpcd on "+inter.Name)
					go ExecuteWait("dhcpcd","-q","-w",inter.Name)

					break outer
				}

			}

		}
	}()
	//conn.RemoveSignal(dbus_objects[inter.Name])
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

	//fmt.Println("output: "+op)

	if op == "" {

		//mylog.Println("Sending the signal to ")
		dbus_objects[ifname] <- &dd

	}

	Kill("dhcpcd.*"+ifname)
}

func DbusFetchProperty(inter PhysicalInterface) string {

	conn, err := dbus.SystemBusPrivate()
	conn.Auth(nil)
	conn.Hello()
	defer conn.Close()
	if err != nil {
		mylog.Fatal(err)
	}

	obj := conn.Object("fi.w1.wpa_supplicant1",
		dbus.ObjectPath("/fi/w1/wpa_supplicant1"))

	var intfPath dbus.ObjectPath

	obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, inter.Name).Store(&intfPath)


	obj = conn.Object("fi.w1.wpa_supplicant1", dbus.ObjectPath(string(intfPath)))


	variant,_ := obj.GetProperty("fi.w1.wpa_supplicant1.Interface.State")

	some := variant.String()[1:len(variant.String())-1]
	//fmt.Println(some)
	return some

}