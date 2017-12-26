$('#interface-toggle').click(function(){

    var str = document.getElementById("interface-toggle").getAttribute("for");


    document.getElementById(str.substr(1,str.length)).innerHTML=

    "<ul class=\"flex-column nav dropdown\" >\n" +
        "\n" +
        "                    <li class=\"dropdown-item\">wlan0</li>\n" +
        "                    <li class=\"dropdown-item\">eth0</li>\n" +
        "                    <li class=\"dropdown-item\">lo</li>\n" +
        "\n" +
        "                </ul>";



    $(str).collapse('toggle');

});

$(document).ready(function(){

    var content="interface=wlan0\n" +
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

        $('#dnsmasq_config').val(content);

});