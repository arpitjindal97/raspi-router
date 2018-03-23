
//css modules
import './styles.css'



import * as $ from 'jquery';
import 'bootstrap/dist/js/bootstrap.js';
import 'material-design-lite'
import 'bootstrap-notify'

const server_ip: string = "/api";

function LoadHtmlDiv(div_id: any, html_file: any) {
    let con = document.getElementById(div_id)
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

function fill_status_page() {

    LoadHtmlDiv("content_div", "device_info.html")

    $.getJSON(server_ip + '/OSInfo', function (data) {

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

window.onload = function(){

    let status_but = document.getElementById("status");

    $('#status').click(fill_status_page);

    status_but.click()

    document.getElementById("interface-toggle").onclick = function () {
        $("#interface-list").collapse('toggle')
    };
}

let current_inter: PhysicalInterface;


function interface_item_clicked(element: HTMLLIElement) {

    LoadHtmlDiv("content_div", "interface.html");

    $.getJSON(server_ip + '/PhysicalInterfaces', function (data) {


        let i;

        document.getElementById("nat_int").innerHTML = "<option value='' id='nat_int_'></option>";

        for (let j = 0; j < data.length; j++) {
            if (data[j]["Name"] == element.innerHTML) {
                i = j;
                continue;
            }

            let opt = document.createElement('option');
            opt.value = data[j]["Name"];
            opt.id = 'nat_int_'+opt.value;
            opt.innerHTML = opt.value;
            $("#nat_int").append(opt);
        }


        current_inter = data[i];

        document.getElementById("interface_name").innerHTML = data[i]["Name"];
        document.getElementById("ip_addr_info").innerHTML = data[i]["Info"]["IpAddress"];
        document.getElementById("broad_addr").innerHTML = data[i]["Info"]["BroadcastAddress"];
        document.getElementById("gate_addr").innerHTML = data[i]["Info"]["Gateway"];
        document.getElementById("mac_addr").innerHTML = data[i]["Info"]["MacAddress"];
        document.getElementById("rec_bytes").innerHTML = data[i]["Info"]["RecvBytes"];
        document.getElementById("rec_packs").innerHTML = data[i]["Info"]["RecvPackts"];
        document.getElementById("trans_bytes").innerHTML = data[i]["Info"]["TransBytes"];
        document.getElementById("trans_packs").innerHTML = data[i]["Info"]["TransPackts"];

        document.getElementById("bridge_mode_" + data[i]["BridgeMode"]).setAttribute("checked", "")

        let element1 = document.getElementById("nat_int_" + data[i]["NatInterface"])

        if (element1 == null)
            document.getElementById("nat_int_").setAttribute("selected", "")
        else
            element1.setAttribute("selected", "")

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
            document.getElementById("wifi_config_div").removeAttribute("style");

            document.getElementById("ip_mode_div").removeAttribute("style");

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
            document.getElementById("dnsmasq_div").removeAttribute("style");
            document.getElementById("hostapd_div").removeAttribute("style");

            document.getElementById("ip_mode_div").removeAttribute("style");

            document.getElementById("bridge_mode_wpa").setAttribute("disabled", "");
            document.getElementById("bridge_mode_hostapd").setAttribute("disabled", "");
            document.getElementById("nat_int").removeAttribute("disabled");
            document.getElementById("mode_default").removeAttribute("checked");
            document.getElementById("mode_off").removeAttribute("checked");
            document.getElementById("mode_bridge").removeAttribute("checked");
            this.setAttribute("checked", "");

        };

        document.getElementById("mode_bridge").onclick = function () {

            document.getElementById("wifi_config_div").removeAttribute("style");
            document.getElementById("ip_dns_outer").setAttribute("style", "display:none");
            document.getElementById("hostapd_div").removeAttribute("style");


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
            document.getElementById("ip_dns_outer").setAttribute("style", "display:none");
            document.getElementById("hostapd_div").setAttribute("style", "display:none");


            document.getElementById("bridge_mode_wpa").setAttribute("disabled", "");
            document.getElementById("bridge_mode_hostapd").setAttribute("disabled", "");
            document.getElementById("nat_int").setAttribute("disabled", "");

            document.getElementById("mode_hotspot").removeAttribute("checked");
            document.getElementById("mode_default").removeAttribute("checked");
            document.getElementById("mode_bridge").removeAttribute("checked");
            this.setAttribute("checked", "");

        };

        document.getElementById("ip_mode_dhcp").onclick = function () {

            document.getElementById("ip_addr").setAttribute("disabled", "");
            document.getElementById("subnet").setAttribute("disabled", "");

            this.setAttribute("checked", "");
            document.getElementById("ip_mode_static").removeAttribute("checked");
        }
        document.getElementById("ip_mode_static").onclick = function () {

            document.getElementById("ip_addr").removeAttribute("disabled");
            document.getElementById("subnet").removeAttribute("disabled");

            document.getElementById("ip_mode_dhcp").removeAttribute("checked");
            this.setAttribute("checked", "");
        }

        document.getElementById("mode_" + data[i]["Mode"]).click();

        document.getElementById("ip_mode_" + data[i]["IpModes"] ).click();

        document.getElementById("ip_addr").setAttribute("value", data[i]["IpAddress"]);
        document.getElementById("subnet").setAttribute("value", data[i]["SubnetMask"]);

        document.getElementById("interface_save_button").onclick = function (ev) {
            saveButtonClicked()
        };

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


$.getJSON(server_ip + '/PhysicalInterfaces', function (data) {

    for (let i = 0; i < data.length; i++) {

        let opt = document.createElement('li');
        opt.value = data[i]["Name"];
        opt.className = "dropdown-item"
        opt.id = "interface-item"

        opt.onclick = function (ev: MouseEvent) {
            interface_item_clicked(opt)
        }

        opt.innerHTML = data[i]["Name"];

        document.getElementById("interface-list").appendChild(opt)
    }

});


export interface PhysicalInterface {
    Name: string,
    IsWifi: string,
    Mode: string,
    BridgeMode: string,
    BridgeMaster: string,
    NatInterface: string,
    IpModes: string,
    IpAddress: string,
    SubnetMask: string,
    Wpa: string,
    Hostapd: string,
    Dnsmasq: string,
    Info: BasicInfo
}


export interface BasicInfo {
    IpAddress: string,
    BroadcastAddress: string,
    Gateway: string,
    MacAddress: string,
    RecvBytes: string,
    RecvPackts: string,
    TransBytes: string,
    TransPackts: string,

    ConntectedTo: string,
    ApMacAddr: string,
    BitRate: string,
    Frequency: string,
    LinkQuality: string,
    Channel: string,
}

 function saveButtonClicked() {


     let objecttosend: PhysicalInterface = {
         Name: document.getElementById("interface_name").innerHTML,
         IsWifi: current_inter.IsWifi,
         Mode: <string>$('[name="mode"]:checked').val(),
         BridgeMode: <string>$('[name="bridge"]:checked').val(),
         BridgeMaster: $("#bridge_master").html(),
         NatInterface: <string>$('#nat_int option:selected').val(),
         IpModes: <string>$('[name="ip_mode"]:checked').val(),
         IpAddress: <string>$("#ip_addr").val(),
         SubnetMask: <string>$("#subnet").val(),
         Wpa: <string>$("#wpa_config_area").val(),
         Hostapd: <string>$("#hostapd_config").val(),
         Dnsmasq: <string>$("#dnsmasq_config").val(),
         Info: null
     };


    if (objecttosend.BridgeMaster != "") {

        console.log("Can't make changes to this interface, first remove it from associated bridge");
        $.notify({
            message: 'Can\'t make changes to '+objecttosend.Name+', first remove it from associated bridge',
            icon: 'glyphicon glyphicon-danger-sign',

        },{
            type: 'danger'
        });
        return
    }


    //Stop the interface, then Save the config, Start the interface
     $.notify({
         message: 'Stopping Interface '+objecttosend.Name,
         icon: 'glyphicon glyphicon-info-sign',

     },{
         type: 'info'
     });

     $.post(server_ip+"/PhysicalInterfaceStop", JSON.stringify(current_inter), function (data,status){

         console.log(data['Message']);

         $.notify({
             message: 'Saving Configuration of '+current_inter.Name,
             icon: 'glyphicon glyphicon-info-sign',

         },{
             type: 'info'
         });


         $.post(server_ip+"/PhysicalInterfaceSave", JSON.stringify(objecttosend), function (data,status){

             console.log(data['Message']);

             $.notify({
                 message: 'Starting Interface '+objecttosend.Name,
                 icon: 'glyphicon glyphicon-info-sign',

             },{
                 type: 'info'
             });

             $.post(server_ip+"/PhysicalInterfaceStart", JSON.stringify(objecttosend), function (data,status){

                 console.log(data['Message']);
                 $.notify({
                     message: 'Successfully started interface '+objecttosend.Name,
                     icon: 'glyphicon glyphicon-success-sign',

                 },{
                     type: 'success'
                 });
             });
         });

     });



}