$('#interface-toggle').click(function () {

    $("#interface-list").collapse('toggle');

});


$('#status').click(function () {
    fill_status_page();
});

function fill_status_page() {

    LoadHtmlDiv("content_div", "device_info.html")

    $.getJSON(server_ip + '/device_info', function (data) {

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
    LoadHtmlDiv("content_div", "bridge.html")
}


var server_ip = "/api";
$(document).ready(function () {


    fill_status_page();

    $.getJSON(server_ip + '/interfaces', function (data) {

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

    $.getJSON(server_ip + '/interfaces', function (data) {

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

    console.log(selectedMode)
    var json_obj={
        "Name":name, "Mode":selectedMode,"BridgeMode":selectedBridgeMode,"NatInterface":nat_int,
        "IpModes":ip_mode,"IpAddress":ip_addr,"SubnetMask":subnet_addr,"Wpa":wpa_config,"Hostapd":hostapd_config,"Dnsmasq":dnsmasq_config,
        "IsWifi":"","Info":null
    };

    $.post(server_ip+"/update_interface", JSON.stringify(json_obj));
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