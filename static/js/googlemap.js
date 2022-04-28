// Initialize and add the map
function initGoogleMap() {
    // The location of Ducky's Dockside
    const duckysdockside = { lat: 28.08851088815534, lng: -81.54721230468368 };
    const centerpoint = { lat: 28.078365554675464,lng: -81.57906601575843 };
    // The map, centered at Ducky's Dockside
    const map = new google.maps.Map(document.getElementById("DDS_Google_map"), {
      zoom: 11,
      center: centerpoint,
    });
    // The marker, positioned at duckysdockside
    const marker = new google.maps.Marker({
      position: duckysdockside,
      map: map,
    });
  }