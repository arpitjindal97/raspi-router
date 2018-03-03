$('#interface-toggle').click(function () {

    $("#interface-list").collapse('toggle');

});


$('#status').click(function () {
    fill_status_page();
});

function fill_status_page() {

    LoadHtmlDiv("content_div", "device_info.html")

    $.getJSON(server_ip + '/DeviceInfo', function (data) {

        document.getElementById("distID").innerHTML = data["DistributionId"];
        document.getElementById("desc").innerHTML = data["Description"];
        document.getElementById("release").innerHTML = data["Release"];
        document.getElementById("codename").innerHTML = data["Codename"];
        document.getElementById("hostname").innerHTML = data["Hostname"];
        document.getElementById("kernel_rel").innerHTML = data["KernelRelease"];
        document.getElementById("arch").innerHTML = data["Architecture"];
        document.getElementById("model_name").innerHTML = data["ModelName"];
        document.getElementById("cores").innerHTML = data["CPUs"];
        document.getElementById("local_time").innerHTML = data["LocalTime"];
        document.getElementById("timezone").innerHTML = data["TimeZone"];
        document.getElementById("up_time").innerHTML = data["UpTime"];
        document.getElementById("up_since").innerHTML = data["UpSince"];

    });
}

$('#bridge').click(function () {
    fill_bridge_page();
});

function fill_bridge_page() {
    LoadHtmlDiv("content_div", "bridge.html");

    $.getJSON(server_ip + '', function (data) {

        var list_bridge_interface = document.getElementById("list_bridge_interface");

        for(var i=0;i< data["BridgeInterfaces"].length ;i++) {

            var opt = document.createElement('option');
            opt.value = data["BridgeInterfaces"][i]["Name"];
            opt.innerHTML = opt.value;
            list_bridge_interface.appendChild(opt)
        }

        document.getElementById("create_bridge").onclick = function () {
            CreateBridge();
        };

        document.getElementById("bridge_del_buttom").onclick = function () {
            DeleteBridge();
        };
        document.getElementById("bridge_independent_add").onclick = function () {
            AddSlaveBridge();
        };
        document.getElementById("bridge_slave_remove").onclick = function () {
            RemoveSlaveBridge();
        };


        list_bridge_interface.addEventListener("change", BridgeSelectionChange);

        BridgeSelectionChange()


    });
}


var server_ip = "/api";
$(document).ready(function () {


    fill_status_page();

    $.getJSON(server_ip + '/PhysicalInterfaces', function (data) {

        document.getElementById("interface-list").innerHTML =

            "<ul class=\"flex-column nav dropdown\" >\n";


        for (var i = 0; i < data.length; i++) {
            document.getElementById("interface-list").innerHTML +=

                "<li class=\"dropdown-item\" id=\"interface-item\" onclick='interface_item_clicked(this)'>" + data[i]["Name"] + "</li>\n";
        }

        document.getElementById("interface-list").innerHTML += "</ul>";


    });

});


function interface_item_clicked(element) {

    LoadHtmlDiv("content_div", "interface.html");

    $.getJSON(server_ip + '/PhysicalInterfaces', function (data) {

        var i;

        document.getElementById("nat_int").innerHTML = "<option value='' id='nat_int_'></option>";

        for (var j = 0; j < data.length; j++) {
            if (data[j]["Name"] == element.innerHTML) {
                i = j;
                continue;
            }
            document.getElementById("nat_int").innerHTML +=
                "<option value=\"" + data[j]["Name"] +
                "\" id='nat_int_" + data[j]["Name"] + "'>" + data[j]["Name"] + "</option>";
        }

        document.getElementById("interface_name").innerHTML = data[i]["Name"];
        document.getElementById("ip_addr").innerHTML = data[i]["Info"]["IpAddress"];
        document.getElementById("broad_addr").innerHTML = data[i]["Info"]["BroadcastAddress"];
        document.getElementById("gate_addr").innerHTML = data[i]["Info"]["Gateway"];
        document.getElementById("mac_addr").innerHTML = data[i]["Info"]["MacAddress"];
        document.getElementById("rec_bytes").innerHTML = data[i]["Info"]["RecvBytes"];
        document.getElementById("rec_packs").innerHTML = data[i]["Info"]["RecvPackts"];
        document.getElementById("trans_bytes").innerHTML = data[i]["Info"]["TransBytes"];
        document.getElementById("trans_packs").innerHTML = data[i]["Info"]["TransPackts"];

        document.getElementById("bridge_mode_" + data[i]["BridgeMode"]).setAttribute("checked", "")

        var element1 = document.getElementById("nat_int_" + data[i]["NatInterface"])

        if (element1 == null )
            document.getElementById("nat_int_").setAttribute("selected", "")
        else
            element1.setAttribute("selected","")

        document.getElementById("conn_to").innerHTML = data[i]["Info"]["ConntectedTo"];
        document.getElementById("ap_mac_addr").innerHTML = data[i]["Info"]["ApMacAddr"];
        document.getElementById("bit_rate").innerHTML = data[i]["Info"]["BitRate"];
        document.getElementById("frequency").innerHTML = data[i]["Info"]["Frequency"];
        document.getElementById("link_quality").innerHTML = data[i]["Info"]["LinkQuality"];
        document.getElementById("channel").innerHTML = data[i]["Info"]["Channel"];

        $("#wpa_config_area").val(data[i]["Wpa"]);
        $('#hostapd_config').val(data[i]["Hostapd"]);
        $('#dnsmasq_config').val(data[i]["Dnsmasq"]);


        document.getElementById("mode_default").onclick = function () {

            document.getElementById("dnsmasq_div").setAttribute("style", "display:none");
            document.getElementById("hostapd_div").setAttribute("style", "display:none");
            document.getElementById("ip_mode_hotspot_div").setAttribute("style", "display:none");
            document.getElementById("wifi_config_div").removeAttribute("style");
            document.getElementById("ip_mode_default_div").removeAttribute("style");
            document.getElementById("bridge_mode_wpa").setAttribute("disabled", "");
            document.getElementById("bridge_mode_hostapd").setAttribute("disabled", "");
            document.getElementById("nat_int").setAttribute("disabled", "");
            document.getElementById("mode_hotspot").removeAttribute("checked");
            document.getElementById("mode_off").removeAttribute("checked");
            document.getElementById("mode_bridge").removeAttribute("checked");
            this.setAttribute("checked", "");

        };

        document.getElementById("mode_hotspot").onclick = function () {

            document.getElementById("wifi_config_div").setAttribute("style", "display:none");
            document.getElementById("ip_mode_default_div").setAttribute("style", "display:none");
            document.getElementById("dnsmasq_div").removeAttribute("style");
            document.getElementById("hostapd_div").removeAttribute("style");
            document.getElementById("ip_mode_hotspot_div").removeAttribute("style");
            document.getElementById("bridge_mode_wpa").setAttribute("disabled","");
            document.getElementById("bridge_mode_hostapd").setAttribute("disabled","");
            document.getElementById("nat_int").removeAttribute("disabled");
            document.getElementById("mode_default").removeAttribute("checked");
            document.getElementById("mode_off").removeAttribute("checked");
            document.getElementById("mode_bridge").removeAttribute("checked");
            this.setAttribute("checked", "");

        };

        document.getElementById("mode_bridge").onclick = function () {

            document.getElementById("wifi_config_div").removeAttribute("style");
            document.getElementById("dnsmasq_div").setAttribute("style", "display:none");
            document.getElementById("hostapd_div").removeAttribute("style");

            document.getElementById("ip_mode_default_div").setAttribute("style", "display:none");
            document.getElementById("ip_mode_hotspot_div").setAttribute("style", "display:none");

            document.getElementById("bridge_mode_wpa").removeAttribute("disabled");
            document.getElementById("bridge_mode_hostapd").removeAttribute("disabled");
            document.getElementById("nat_int").setAttribute("disabled", "");

            document.getElementById("mode_hotspot").removeAttribute("checked");
            document.getElementById("mode_default").removeAttribute("checked");
            document.getElementById("mode_bridge").removeAttribute("checked");
            this.setAttribute("checked", "");

        };


        document.getElementById("mode_off").onclick = function () {

            document.getElementById("wifi_config_div").setAttribute("style", "display:none");
            document.getElementById("dnsmasq_div").setAttribute("style", "display:none");
            document.getElementById("hostapd_div").setAttribute("style", "display:none");

            document.getElementById("ip_mode_default_div").setAttribute("style", "display:none");
            document.getElementById("ip_mode_hotspot_div").setAttribute("style", "display:none");

            document.getElementById("bridge_mode_wpa").setAttribute("disabled", "");
            document.getElementById("bridge_mode_hostapd").setAttribute("disabled", "");
            document.getElementById("nat_int").setAttribute("disabled", "");

            document.getElementById("mode_hotspot").removeAttribute("checked");
            document.getElementById("mode_default").removeAttribute("checked");
            document.getElementById("mode_bridge").removeAttribute("checked");
            this.setAttribute("checked", "");

        };

        document.getElementById("ip_mode_dhcp_default").onclick = function () {

            document.getElementById("ip_addr_static_default").setAttribute("disabled", "");
            document.getElementById("subnet_static_default").setAttribute("disabled", "");
            document.getElementById("ip_addr_static_hotspot").setAttribute("disabled", "");
            document.getElementById("subnet_static_hotspot").setAttribute("disabled", "");
            this.setAttribute("checked","");
            document.getElementById("ip_mode_dhcp_hotspot").setAttribute("checked","");
            document.getElementById("ip_mode_static_default").removeAttribute("checked");
            document.getElementById("ip_mode_static_hotspot").removeAttribute("checked");
        }
        document.getElementById("ip_mode_static_default").onclick = function () {

            document.getElementById("ip_addr_static_default").removeAttribute("disabled");
            document.getElementById("subnet_static_default").removeAttribute("disabled");
            document.getElementById("ip_addr_static_hotspot").removeAttribute("disabled");
            document.getElementById("subnet_static_hotspot").removeAttribute("disabled");
            document.getElementById("ip_mode_dhcp_default").removeAttribute("checked");
            document.getElementById("ip_mode_dhcp_hotspot").removeAttribute("checked");
            this.setAttribute("checked","");
            document.getElementById("ip_mode_static_hotspot").setAttribute("checked","");
        }
        document.getElementById("ip_mode_dhcp_hotspot").onclick = function () {

            document.getElementById("ip_mode_dhcp_default").click();
        }
        document.getElementById("ip_mode_static_hotspot").onclick = function () {

            document.getElementById("ip_mode_static_default").click()
        }

        document.getElementById("mode_" + data[i]["Mode"]).click();

        document.getElementById("ip_mode_" + data[i]["IpModes"] + "_default").click();
        document.getElementById("ip_mode_" + data[i]["IpModes"] + "_hotspot").click();
        document.getElementById("ip_addr_static_default").setAttribute("value", data[i]["IpAddress"]);
        document.getElementById("ip_addr_static_hotspot").setAttribute("value", data[i]["IpAddress"]);
        document.getElementById("subnet_static_default").setAttribute("value", data[i]["SubnetMask"]);
        document.getElementById("subnet_static_hotspot").setAttribute("value", data[i]["SubnetMask"]);

        document.getElementById("interface_save_button").onclick = function (ev) {
            sendData()
        }

        document.getElementById("bridge_mode_wpa").onclick = function (ev) {
            document.getElementById("bridge_mode_hostapd").removeAttribute("checked")
            document.getElementById("bridge_mode_wpa").setAttribute("checked", "")
        }

        document.getElementById("bridge_mode_hostapd").onclick = function (ev) {
            document.getElementById("bridge_mode_wpa").removeAttribute("checked")
            document.getElementById("bridge_mode_hostapd").setAttribute("checked", "")
        }

        document.getElementById("bridge_master").innerHTML = data[i]["BridgeMaster"];

    });
}

function sendData() {

    var name = document.getElementById("interface_name").innerHTML;

    var modes = document.getElementsByName("mode");
    var selectedMode;
    for (var i = 0; i < modes.length; i++) {
        if (modes.item(i).hasAttribute("checked") == true) {
            selectedMode = modes.item(i).getAttribute("value");
        }
    }

    var bridge_modes = document.getElementsByName("bridge");
    var selectedBridgeMode;
    for (var i = 0; i < bridge_modes.length; i++) {
        if (bridge_modes.item(i).hasAttribute("checked") == true) {
            selectedBridgeMode = bridge_modes.item(i).getAttribute("value");
        }
    }
    var nat_int = document.getElementById("nat_int")

    nat_int = nat_int.options[nat_int.selectedIndex].text;

    var wpa_config = $("#wpa_config_area").val()
    var hostapd_config = $("#hostapd_config").val()
    var dnsmasq_config = $("#dnsmasq_config").val()

    var ip_mode ="dhcp";

    var ip_addr = "";
    var subnet_addr = "";

    if (selectedMode != "off" && selectedMode != "bridge") {
        ip_mode = document.getElementsByName("ip_mode_" + selectedMode);

        for (var i = 0; i < ip_mode.length; i++) {
            if (ip_mode.item(i).hasAttribute("checked") == true) {
                ip_mode = ip_mode.item(i).getAttribute("value");
                break;
            }
        }
         ip_addr = $("#ip_addr_static_" + selectedMode).val();
         subnet_addr = $("#subnet_static_" + selectedMode).val();
    }

    var bridge_master = $("#bridge_master").val();

    console.log(selectedMode)
    var json_obj={
        "Name":name, "Mode":selectedMode,"BridgeMode":selectedBridgeMode,"BridgeMaster":bridge_master,"NatInterface":nat_int,
        "IpModes":ip_mode,"IpAddress":ip_addr,"SubnetMask":subnet_addr,"Wpa":wpa_config,"Hostapd":hostapd_config,"Dnsmasq":dnsmasq_config,
        "IsWifi":"","Info":null
    };

    $.post(server_ip+"/UpdateInterface", JSON.stringify(json_obj));
}

function LoadHtmlDiv(div_id, html_file) {
    var con = document.getElementById(div_id)
        , xhr = new XMLHttpRequest();

    xhr.onreadystatechange = function (e) {
        if (xhr.readyState == 4 && xhr.status == 200) {
            con.innerHTML = xhr.responseText;
        }
    }

    xhr.open("GET", html_file, true);
    xhr.setRequestHeader('Content-type', 'text/html');
    xhr.send();
}


function CreateBridge(){

    var bridge_name = $("#create_bridge_name").val();

    $.post(server_ip+"/CreateBridge", bridge_name, function (data,status){

            progress(80);
            console.log(data)
            fill_bridge_page();
            //progress(100);
    });

    progress(40);
}

function progress(percent){
    document.querySelector('#p1').addEventListener('mdl-componentupgraded', function() {
        this.MaterialProgress.setProgress(percent);
    });
}

function DeleteBridge(){

    var bridge_int = document.getElementById("list_bridge_interface")

    var bridge_name = bridge_int.options[bridge_int.selectedIndex].text;
    $.post(server_ip+"/DeleteBridge", bridge_name, function (data,status){

        progress(80);
        console.log(data)
        fill_bridge_page();
        //progress(100);
    });

}

function BridgeSelectionChange() {

    $.getJSON(server_ip + '', function (data) {

        var list_bridge_interface = document.getElementById("list_bridge_interface");

        var element = data["BridgeInterfaces"][list_bridge_interface.selectedIndex];

        document.getElementById("ip_addr").innerHTML = element["Info"]["IpAddress"];
        document.getElementById("broad_addr").innerHTML = element["Info"]["BroadcastAddress"];
        document.getElementById("gate_addr").innerHTML = element["Info"]["Gateway"];
        document.getElementById("mac_addr").innerHTML = element["Info"]["MacAddress"];
        document.getElementById("rec_bytes").innerHTML = element["Info"]["RecvBytes"];
        document.getElementById("rec_packs").innerHTML = element["Info"]["RecvPackts"];
        document.getElementById("trans_bytes").innerHTML = element["Info"]["TransBytes"];
        document.getElementById("trans_packs").innerHTML = element["Info"]["TransPackts"];

        document.getElementById("bridge_"+element["IpMode"]).setAttribute("checked","");

        document.getElementById("ip_addr_static").innerHTML = element["IpAddress"];
        document.getElementById("subnet_static").innerHTML = element["SubnetMask"];

        var slaves_int = document.getElementById("slaves_int");

        for(var i=0;i< element["Slaves"].length ;i++) {

            var opt = document.createElement('option');
            opt.value = element["Slaves"][i];
            opt.innerHTML = opt.value;
            slaves_int.appendChild(opt)
        }

        element = document.getElementById("bridge_independent_int");

        for(var i=0;i< data["PhysicalInterfaces"].length;i++) {

            if (data["PhysicalInterfaces"][i]["Mode"] == "bridge" ) {

                var opt = document.createElement('option');
                opt.value = data["PhysicalInterfaces"][i]["Name"];
                opt.innerHTML = opt.value;
                element.appendChild(opt)
            }

        }

    });
}

function AddSlaveBridge() {

    var bridge_int = document.getElementById("list_bridge_interface");

    var bridge_name = bridge_int.options[bridge_int.selectedIndex].text;


    bridge_int = document.getElementById("bridge_independent_int");

    var slave_name = bridge_int.options[bridge_int.selectedIndex].text;

    var json = {"BridgeIfname":bridge_name,"SlaveIfname":slave_name};

    $.post(server_ip+"/BridgeAddSlave", JSON.stringify(json), function (data,status){

        progress(80);
        console.log(data)
        fill_bridge_page();
        //progress(100);
    });
}

function RemoveSlaveBridge() {

    var bridge_int = document.getElementById("slaves_int");

    var bridge_name = bridge_int.options[bridge_int.selectedIndex].text;


    $.post(server_ip+"/BridgeRemoveSlave", bridge_name, function (data,status){

        progress(80);
        console.log(data)
        fill_bridge_page();
        //progress(100);
    });
}