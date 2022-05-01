package app

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// Get pertinent env. values.
	StaticLocation = os.Getenv("DDS_STATIC_LOCATION")
	Logfile        = os.Getenv("DDS_LOGFILE")
	ServerPort     = os.Getenv("DDS_SERVER_PORT")
	// Logging.
	Logger        *log.Logger
	LogFileHandle *os.File
)

// App.
type App struct {
	Log  *log.Logger
	User *User
}

// User.
type User struct {
	Username  string
	Password  string
	Realname  string
	Title     string
	LastLogin string
	LoginTime string
}

// Setup.
func init() {
	var err error
	// Logging.
	Logger = log.New()
	Logger.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	Logger.SetLevel(log.InfoLevel)

	// log file set up.
	LogFileHandle, err = os.OpenFile(Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening Logfile: %s -> %v", Logfile, err)
	}

	mw := io.MultiWriter(os.Stdout, LogFileHandle)
	Logger.SetOutput(mw)
}
