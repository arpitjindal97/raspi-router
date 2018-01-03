$('#interface-toggle').click(function(){

    $("#interface-list").collapse('toggle');

});


$('#status').click(function() {
    fill_status_page();
});

function fill_status_page(){

    LoadHtmlDiv("content_div","device_info.html")

    $.getJSON(server_ip+'/device_info', function(data) {

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


var server_ip="http://192.168.1.42:8084";
$(document).ready(function(){


    fill_status_page();

    $.getJSON(server_ip+'/interfaces', function(data) {

        document.getElementById("interface-list").innerHTML=

            "<ul class=\"flex-column nav dropdown\" >\n" ;


        for(var i=0;i<data.length;i++){
            document.getElementById("interface-list").innerHTML+=

                "<li class=\"dropdown-item\" id=\"interface-item\" onclick='interface_item_clicked(this)'>"+data[i]["Name"]+"</li>\n";
        }

        document.getElementById("interface-list").innerHTML+= "</ul>";


    });

});




function interface_item_clicked(element) {

    LoadHtmlDiv("content_div","interface.html");

    $.getJSON(server_ip+'/interfaces', function(data) {

        var i;

        document.getElementById("route_int").innerHTML="<option value='' id='route_int_'></option>";

        for(var j=0;j<data.length;j++)
        {
            if(data[j]["Name"] == element.innerHTML)
            {
                i=j;
                continue;
            }
            document.getElementById("route_int").innerHTML+=
                "<option value=\""  +data[j]["Name"]+
                "\" id='route_int_"+data[j]["Name"]+"'>"  +data[j]["Name"]+  "</option>";
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

        document.getElementById("route_mode_"+data[i]["RouteMode"]).setAttribute("checked","checked")

        document.getElementById("route_int_"+data[i]["RouteInterface"]).setAttribute("selected","selected")

        document.getElementById("conn_to").innerHTML = data[i]["Info"]["ConntectedTo"];
        document.getElementById("ap_mac_addr").innerHTML = data[i]["Info"]["ApMacAddr"];
        document.getElementById("bit_rate").innerHTML = data[i]["Info"]["BitRate"];
        document.getElementById("frequency").innerHTML = data[i]["Info"]["Frequency"];
        document.getElementById("link_quality").innerHTML = data[i]["Info"]["LinkQuality"];
        document.getElementById("channel").innerHTML = data[i]["Info"]["Channel"];

        $("#wpa_config_area").val(data[i]["Wpa"]);
        $('#hostapd_config').val(data[i]["Hostapd"]);
        $('#dnsmasq_config').val(data[i]["Dnsmasq"]);


        document.getElementById("mode_default").onclick = function(){

            document.getElementById("dnsmasq_div").setAttribute("style","display:none");
            document.getElementById("hostapd_div").setAttribute("style","display:none");
            document.getElementById("ip_mode_hotspot_div").setAttribute("style","display:none");
            document.getElementById("wifi_config_div").removeAttribute("style");
            document.getElementById("ip_mode_default_div").removeAttribute("style");
            document.getElementById("route_mode_nat").setAttribute("disabled","");
            document.getElementById("route_mode_bridge").setAttribute("disabled","");
            document.getElementById("route_int").setAttribute("disabled","");

        };

        document.getElementById("mode_hotspot").onclick = function(){

            document.getElementById("wifi_config_div").setAttribute("style","display:none");
            document.getElementById("ip_mode_default_div").setAttribute("style","display:none");
            document.getElementById("dnsmasq_div").removeAttribute("style");
            document.getElementById("hostapd_div").removeAttribute("style");
            document.getElementById("ip_mode_hotspot_div").removeAttribute("style");
            document.getElementById("route_mode_nat").removeAttribute("disabled");
            document.getElementById("route_mode_bridge").removeAttribute("disabled");
            document.getElementById("route_int").removeAttribute("disabled");

        };

        document.getElementById("mode_"+data[i]["Mode"]).click()

        document.getElementById("ip_mode_dhcp_default").onclick = function(){

            document.getElementById("ip_addr_static_default").setAttribute("disabled","");
            document.getElementById("subnet_static_default").setAttribute("disabled","");
        }
        document.getElementById("ip_mode_static_default").onclick = function(){

            document.getElementById("ip_addr_static_default").removeAttribute("disabled");
            document.getElementById("subnet_static_default").removeAttribute("disabled");
        }


        document.getElementById("ip_mode_"+data[i]["IpModes"]+"_"+data[i]["Mode"]).click();
        document.getElementById("ip_addr_static_"+data[i]["Mode"]).setAttribute("value",data[i]["IpAddress"]);
        document.getElementById("subnet_static_"+data[i]["Mode"]).setAttribute("value",data[i]["SubnetMask"]);

    });
}


function LoadHtmlDiv(div_id,html_file)
{
    var con = document.getElementById(div_id)
        ,   xhr = new XMLHttpRequest();

    xhr.onreadystatechange = function (e) {
        if (xhr.readyState == 4 && xhr.status == 200) {
            con.innerHTML = xhr.responseText;
        }
    }

    xhr.open("GET", html_file, true);
    xhr.setRequestHeader('Content-type', 'text/html');
    xhr.send();
}