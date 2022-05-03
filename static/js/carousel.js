$(document).ready(function() {
    var noOfCarousels = $(".carousel_row").length
    var maxFade = 0.05
    var carouselId = "#carousel"

    for (var currentCarousel = noOfCarousels; currentCarousel > 1; currentCarousel--){
        id = carouselId + currentCarousel
        $(id).hide()
    }
    $(carouselId + currentCarousel).show()

    setInterval(function() {
    id = carouselId + currentCarousel
        $(id).fadeTo("slow", maxFade, function(){
            $(id).hide()
            currentCarousel += 1
            if (currentCarousel > noOfCarousels) {
                currentCarousel = 1
            }
            id = carouselId + currentCarousel
            $(id).fadeTo("slow", 1)
        })
   }, 20000);
})