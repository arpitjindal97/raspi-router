$('#interface-toggle').click(function(){

    $("#interface-list").collapse('toggle');

});

var server_ip="http://103.43.103.19:8084";
$(document).ready(function(){

    /*var content="interface=wlan0\n" +
        "driver=nl80211\n" +
        "ssid=Raspberry-Hotspot\n" +
        "hw_mode=g\n" +
        "ieee80211n=1\n" +
        "wmm_enabled=1\n" +
        "macaddr_acl=0\n" +
        "ht_capab=[HT40][SHORT-GI-20][DSSS_CCK-40]\n" +
        "channel=6\n" +
        "auth_algs=1\n" +
        "ignore_broadcast_ssid=0\n" +
        "wpa=2\n" +
        "wpa_key_mgmt=WPA-PSK\n" +
        "wpa_passphrase=raspberry\n" +
        "rsn_pairwise=CCMP";

        $('#hostapd_config').val(content);


        content="interface=wlan0\n" +
            "bind-interfaces\n" +
            "server=8.8.8.8\n" +
            "domain-needed\n" +
            "bogus-priv\n" +
            "dhcp-range=192.168.2.2,192.168.2.100,12h";

        $('#dnsmasq_config').val(content);*/

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


    $.getJSON(server_ip+'/interfaces', function(data) {

        document.getElementById("interface-list").innerHTML=

            "<ul class=\"flex-column nav dropdown\" >\n" ;


        for(var i=0;i<data.length;i++){
            document.getElementById("interface-list").innerHTML+=

                "<li class=\"dropdown-item\" id='interface-item'>"+data[i]["Name"]+"</li>\n";
        }

        document.getElementById("interface-list").innerHTML+= "</ul>";


    });

});



$("#interface-item").click(function(){

    console.log(this.innerHTML);
    console.log("not working")
});


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