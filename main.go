/*
	duckydockside.com - Web Server Pages App
	=========================================

	Complete documentation and user guides are available here:
	https://https://github.com/yveshoebeke/duckysdockside/blob/master/README.md

	@author	yves.hoebeke@gmail.com - 1011001.1110110.1100101.1110011

	@version 1.0.0

	@date 2022-05-01

	(c) 2022 - Ducky's Dockside Bar & Grill, LLC - All Rights Reserved.
*/

package main

/* System libraries */
import (
	"net/http"
	"os"
	"time"

	"duckysdockside.com/packages/app"
	"duckysdockside.com/packages/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var sys app.App

/* Middleware */
func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Research below and other potential actions.
		sys.Log.Printf("URL: %s | Method: %s", r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}

/* End Middleware */

/*
       ^ ^
      (o O)
 ___oOO(.)OOo___
 _______________

 ************************************
 *	Execution start point!!!!!!!!!	*
 *	Structure and Methods 			*
 *	- Setup and start of app.		*
 *	- Serve and Listen.				*
 ************************************

*/
func main() {
	// See app package for assignment.
	defer app.LogFileHandle.Close()
	// Set the app values in the struct.
	sys.Log = app.Logger
	// Show in the console app has started.
	sys.Log.Println("Starting service.")

	/* Routers definitions */
	mux := mux.NewRouter()

	/* Middleware */
	mux.Use(MiddleWare)

	/* Allow static content */
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(app.StaticLocation))))

	/* Handlers */
	mux.HandleFunc("/", routes.Home).Methods(http.MethodGet)
	mux.HandleFunc("/admin", routes.Admin).Methods(http.MethodGet, http.MethodPost)
	mux.HandleFunc("/adminmenu", routes.AdminMenu).Methods(http.MethodGet)
	mux.HandleFunc("/addevent", routes.AddEvent).Methods(http.MethodGet, http.MethodPost)
	mux.HandleFunc("/manageevents", routes.ManageEvents).Methods(http.MethodGet, http.MethodPost)

	/* Server setup and start */
	DuckyDocksideServer := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, mux),
		Addr:         app.ServerPort,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	/*
	**************************************
	* Setup and initialization completed *
	*                                    *
	*         Launch the server!         *
	**************************************
	 */
	sys.Log.Fatal(DuckyDocksideServer.ListenAndServe())
}
