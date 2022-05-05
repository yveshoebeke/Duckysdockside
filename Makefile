SRCDIR = /Users/yves/Projects/Duckysdockside/

build:
	go build -o $(SRCDIR)main main.go

run:
	export DDS_LOGFILE="/Users/yves/Projects/Duckysdockside/logs/duckysdockside.log"
	export DDS_EVENTSDATAFILE="/Users/yves/Projects/Duckysdockside/data/events.json"
	export DDS_FOODMENUDATAFILE="/Users/yves/Projects/Duckysdockside/data/foodmenu.json"
	export DDS_STATIC_LOCATION="static/"
	export DDS_HTML_LOCATION="html/"
	export DDS_TEMPLATE_LOCATION="templates/"
	export DDS_SERVER_PORT=":80"
	export DDS_WEATHERAPI="71c3f677cbb242889f4173533220505"
	export DDS_ADMIN_PASSWORD="DocksideSecrets101!"

	./$(SCRDIR)main

all: build run
