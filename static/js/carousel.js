$(document).ready(function() {
    var intervalTime = 10000
    var noOfCarousels = 4
    var maxFade = 0.05
    var carouselId = "#carousel_"

    for (var currentCarousel = noOfCarousels; currentCarousel > 0; currentCarousel--){
        id = carouselId + currentCarousel
        $(id).hide()
    }
    $(carouselId + currentCarousel).show()

    setInterval(function() {
    id = carouselId + currentCarousel
        $(id).fadeTo("slow", maxFade, function(){
            $(id).hide()
            currentCarousel += 1
            if (currentCarousel == noOfCarousels) {
                currentCarousel = 0
            }
            id = carouselId + currentCarousel
            $(id).fadeTo("slow", 1)
        })
   }, intervalTime);
})