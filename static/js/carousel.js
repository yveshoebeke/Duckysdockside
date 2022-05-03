$(document).ready(function() {
    $("#carousel1").show()
    $("#carousel2").hide()
    $("#carousel3").hide()
    var carousel = 1;
    setInterval(function() {
    id = "#carousel" + carousel
        $(id).fadeTo("slow", 0.05,function(){
            $(id).hide()
            carousel += 1
            if (carousel > 3) {
                carousel = 1
            }
            id = "#carousel" + carousel
            $(id).fadeTo("slow", 1)
        })
   }, 20000);
})