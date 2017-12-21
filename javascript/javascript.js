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