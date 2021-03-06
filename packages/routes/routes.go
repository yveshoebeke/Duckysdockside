package routes

import (
	"html/template"
	"net/http"
	"sync"

	"duckysdockside.com/packages/app"
	"duckysdockside.com/packages/events"
	"duckysdockside.com/packages/menus"
	"duckysdockside.com/packages/utils"

	log "github.com/sirupsen/logrus"
)

var (
	// Templating
	tmpl *template.Template
)

/* Routers */
/* Home page */
func Home(w http.ResponseWriter, r *http.Request) {
	// Predefined here for go-routines
	var (
		err            error
		eventData      events.Events
		eventsData     events.DisplayEvents
		foodMenu       menus.FoodMenu
		carouselImages [][]string
		wxData         utils.WxData
	)

	// Home page template data structure.
	type HomePageData struct {
		EventData      events.DisplayEvents
		CarouselImages [][]string
		FoodMenu       menus.FoodMenu
		Weather        utils.WxData
	}

	// Group for concurrency -> weather, events, food and images data fetching threads.
	wg := new(sync.WaitGroup)
	// 1. Get current local weather data.
	wg.Add(1)
	go func() {
		defer wg.Done()
		wxData, err = utils.CurrentWeather()
		if err != nil {
			log.Error("Weather:", err)
		}
	}()
	// 2. Fetch & format event schedule.
	wg.Add(1)
	go func() {
		defer wg.Done()

		eventData, err = events.ReadEventsDataFile(true)
		if err != nil {
			log.Error(err)
		}

		// Format the event data.
		eventsData, err = events.FormatEventData(eventData)
		if err != nil {
			log.Error(err)
		}
	}()
	// 3. Fetching menu data.
	wg.Add(1)
	go func() {
		defer wg.Done()

		foodMenu, err = menus.ReadFoodMenuDataFile()
		if err != nil {
			log.Error(err)
		}
	}()
	// 4. Fetch image carousel file names.
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Set images below logo here. @todo: make dynamic.
		carouselImages = utils.GetDefaultImages()
	}()
	// Wait for it all to finish.
	wg.Wait()

	// Assign data to the struct the home template.
	homePageData := HomePageData{
		EventData:      eventsData,
		CarouselImages: carouselImages,
		FoodMenu:       foodMenu,
		Weather:        wxData,
	}

	// Serve it to the user.
	tmpl = template.Must(template.ParseFiles(app.TemplateLocation + "home.go.html"))
	tmpl.Execute(w, homePageData)
}

/* Administration Functionality */
// Admin credential check.
func Admin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, app.HtmlLocation+"admin.html")
	case http.MethodPost:
		err := utils.CheckAdminPassword(r)
		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/", http.StatusFound)
		}

		utils.SetTokenHeader(r)
		http.Redirect(w, r, "/adminmenu", http.StatusFound)
	}
}

// Admin option selections.
func AdminMenu(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, app.HtmlLocation+"adminmenu.html")
}

// Add an event.
func AddEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, app.HtmlLocation+"addevent.html")
	case http.MethodPost:
		// Add the posted event data.
		err := events.AddEvent(r)
		if err != nil {
			log.Error(err)
		}
		// Return to page.
		http.ServeFile(w, r, app.HtmlLocation+"addevent.html")
	}
}

// Maintain existing events.
func ManageEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := events.ManageEvents()
		if err != nil {
			log.Println("Manage Events error:", err)
		}
		// Get the newest data from file -> template.
		eventData, err := events.ReadEventsDataFile(false)
		if err != nil {
			log.Error(err)
		}
		// Serve page to user.
		tmpl = template.Must(template.ParseFiles(app.TemplateLocation + "manageevents.go.html"))
		tmpl.Execute(w, eventData)
	case http.MethodPost:
		// Update the incoming posted event data
		err := events.UpdateEvent(r)
		if err != nil {
			log.Error(err)
		}
		// Get the newest data from file -> template.
		eventData, err := events.ReadEventsDataFile(false)
		if err != nil {
			log.Error(err)
		}
		// Return to page.
		tmpl = template.Must(template.ParseFiles(app.TemplateLocation + "manageevents.go.html"))
		tmpl.Execute(w, eventData)
	}
}

/* End Administration routines */
