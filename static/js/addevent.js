$(document).ready(function() {
    const monthDrop = $("#month")
    const dayDrop = $("#day")
    const yearDrop = $("#year")
    const yearCnt = 2
    const thisYear = new Date().getFullYear()
    const monthNames = [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
    ]

    const timeFromDrop = $("#timefrom")
    const timeToDrop = $("#timeto")
    const timeFrom = ["11am", "12pm"]
    const timeTo = ["11am", "12pm"]

    for( var t = 1; t <= 11; t++) {
        timeFrom.push(t + "pm")
        timeTo.push(t + "pm")
    }

    for( var t = 0; t < timeFrom.length; t++) {
        let timeFromElement = document.createElement("option")
        timeFromElement.value = timeFrom[t]
        timeFromElement.text = timeFrom[t]
        if( t == 3 ){
            timeFromElement.selected = true
        }
        timeFromDrop.append(timeFromElement)

        let timeToElement = document.createElement("option")
        timeToElement.value = timeFrom[t]
        timeToElement.text = timeFrom[t]
        if( t == (timeFrom.length - 2)){
            timeToElement.selected = true
        }
        timeToDrop.append(timeToElement)
    }

    for( var m = 0; m < monthNames.length; m++ ){
        let monthElement = document.createElement("option")
        monthElement.value = m
        monthElement.text = monthNames[m]
        monthDrop.append(monthElement)
    }

    for( var y = 0; y < yearCnt; y++){
        let yearElement = document.createElement("option")
        yearElement.value = thisYear + y
        yearElement.text = thisYear + y
        yearDrop.append(yearElement)
    }

    var d = new Date();
    var month = d.getMonth();
    var year = d.getFullYear();
    var day = d.getDate();
    
    yearDrop.val(year);
    yearDrop.on("change", AdjustDays);
    monthDrop.val(month);
    monthDrop.on("change", AdjustDays);

    AdjustDays();
    dayDrop.val(day)

    function AdjustDays() {
        var year = yearDrop.val();
        var month = parseInt(monthDrop.val()) + 1;
        dayDrop.empty();
    
        //get the last day, so the number of days in that month
        var days = new Date(year, month, 0).getDate();

        for( var d = 1; d <= days; d++ ){
            var dayElement = document.createElement("option")
            dayElement.value = d
            dayElement.text = d
            dayDrop.append(dayElement)
        }
    }

    setInterval(function() {
        var time = new Date();
        $("#doy").text(time.toDateString());
        $("#tod").text(time.toLocaleTimeString());
    }, 1000);
})