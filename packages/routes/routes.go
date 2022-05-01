package routes

import (
	"html/template"
	"net/http"
	"strings"
	"sync"

	"duckysdockside.com/packages/app"
	"duckysdockside.com/packages/events"
	"duckysdockside.com/packages/menus"
	"duckysdockside.com/packages/utils"

	log "github.com/sirupsen/logrus"
)

var (
	// Templating
	tmpl    = template.Must(template.New("").Funcs(funcMap).ParseGlob("html/*"))
	funcMap = template.FuncMap{
		"hasHTTP": func(myUrl string) string {
			if strings.Contains(myUrl, "://") {
				return myUrl
			}

			return "https://" + myUrl
		},
		"userStatus": func(myStatus int) string {
			return ""
		},
	}
)

// Setup
func init() {
}

/* Routers */
/* Home page */
func Home(w http.ResponseWriter, r *http.Request) {
	// Predefined here for go-routines
	var (
		err          error
		eventData    events.Events
		eventsData   events.DisplayEvents
		foodMenu     menus.FoodMenu
		randomImages []string
	)

	// Home page template data structure.
	type HomePageData struct {
		EventData      events.DisplayEvents
		CarouselImages []string
		FoodMenu       menus.FoodMenu
	}

	// Group for concurrency -> events, food and images data fetching threads.
	wg := new(sync.WaitGroup)
	// 1. Fetch & format event schedule.
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
	// 2. Fetching menu data.
	wg.Add(1)
	go func() {
		defer wg.Done()

		foodMenu, err = menus.ReadFoodMenuDataFile()
		if err != nil {
			log.Error(err)
		}
	}()
	// 3. Fetch image carousel file names.
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Set images below logo here.
		randomImages = utils.GetDefaultImages()
		// Todo: Replace above^^^^ with -> randomImages := utils.GetRandomCarouselImages(3).
	}()
	// Wait for it all to finish.
	wg.Wait()

	// Assign data to the struct the home template.
	homePageData := HomePageData{
		EventData:      eventsData,
		CarouselImages: randomImages,
		FoodMenu:       foodMenu,
	}
	// Serve it to the user.
	tmpl.ExecuteTemplate(w, "home.go.html", homePageData)
}

/* Administration Functionality */
// Admin credential check.
func Admin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, app.StaticLocation+"/html/admin.html")
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
	http.ServeFile(w, r, app.StaticLocation+"/html/adminmenu.html")
}

// Add an event.
func AddEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, app.StaticLocation+"/html/addevent.html")
	case http.MethodPost:
		// Add the posted event data.
		err := events.AddEvent(r)
		if err != nil {
			log.Error(err)
		}
		// Return to page.
		http.ServeFile(w, r, app.StaticLocation+"/html/addevent.html")
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
		// Get the newest data from file -> template
		eventData, err := events.ReadEventsDataFile(false)
		if err != nil {
			log.Error(err)
		}
		tmpl.ExecuteTemplate(w, "manageevents.go.html", eventData)
	case http.MethodPost:
		// Update the incoming posted event data
		err := events.UpdateEvent(r)
		if err != nil {
			log.Error(err)
		}
		// Get the newest data from file -> template
		eventData, err := events.ReadEventsDataFile(false)
		if err != nil {
			log.Error(err)
		}
		// Return to page.
		tmpl.ExecuteTemplate(w, "manageevents.go.html", eventData)
	}
}

/* End Administration routines */

/* Middleware */
func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Research below and other potential actions.
		// app.Log.Printf("User: %s | URL: %s | Method: %s", app.User.Username, r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}

/* End Middleware */
