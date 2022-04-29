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
	// Set uo a general purpose user.
	user := &app.User{
		Username:  "WWW",
		Password:  "*",
		Realname:  "Visitor",
		Title:     "visitor",
		LastLogin: time.Now().Format(time.RFC3339),
		LoginTime: time.Now().Format(time.RFC3339),
	}
	// Set the app values in the struct.
	sys := &app.App{
		Log:  app.Logger,
		User: user,
	}
	// Show in the console app has started.
	sys.Log.Println("Starting service.")

	/* Routers definitions */
	mux := mux.NewRouter()

	/* Middleware */
	mux.Use(routes.MiddleWare)

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
