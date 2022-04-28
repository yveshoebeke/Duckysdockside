package events

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"duckysdockside.com/packages/utils"
)

var eventsDataFile = os.Getenv("DDS_EVENTSDATAFILE")

// Event structure.
type Event struct {
	Id     int    `json:"id"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Name   string `json:"name"`
	Genre  string `json:"genre"`
	Image  string `json:"image"`
	Active bool   `json:"active"`
}

type Events []Event

type DisplayEvent struct {
	DisplayDate  string `json:"displaydate"`
	DisplayTime  string `json:"displaytime"`
	DisplayName  string `json:"displayname"`
	DisplayGenre string `json:"displaygenre"`
}

type DisplayEvents []DisplayEvent

// Read all of the json data (see below for flag functionality)
func ReadEventsDataFile(activeOnly bool) (Events, error) {
	jsonEvents, _ := ioutil.ReadFile(eventsDataFile)
	var data Events
	// Unmarshal json data into struct.
	err := json.Unmarshal(jsonEvents, &data)
	if err != nil {
		return nil, err
	}
	// if activeOnly flag is true, only return active records for today and after.
	if activeOnly {
		yesterday := time.Now().AddDate(0, 0, -1)
		active := Events{}
		for _, d := range data {
			eventDay, err := time.Parse("01-02-2006", d.Date)
			if err != nil {
				return nil, err
			}
			if d.Active && eventDay.After(yesterday) {
				active = append(active, d)
			}
		}
		return active, nil
	} else {
		return data, nil
	}
}

// Write structure out to json file
func WriteEventsDataFile(data *Events) error {
	// Marshal into json format.
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	// Write to file.
	err = ioutil.WriteFile(eventsDataFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Add new event
func AddEvent(r *http.Request) error {
	r.ParseForm()
	month, _ := strconv.Atoi(r.FormValue("Month"))
	month += 1 // This to compensate JS Date Obj that defines January ar month 0...
	day, _ := strconv.Atoi(r.FormValue("Day"))
	event := &Event{
		Id:     int(time.Now().Unix()),
		Date:   fmt.Sprintf("%02d-%02d-%s", month, day, r.FormValue("Year")),
		Time:   fmt.Sprintf("from %s to %s", r.FormValue("TimeFrom"), r.FormValue("TimeTo")),
		Name:   r.FormValue("Name"),
		Genre:  r.FormValue("Genre"),
		Active: true,
	}
	if r.FormValue("Active") != "true" {
		event.Active = false
	}

	err := UpdateEventsData(event)
	if err != nil {
		return err
	}

	return nil
}

// Update event entry from main
func UpdateEvent(r *http.Request) error {
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("Id"))
	if err != nil {
		return err
	}
	event := &Event{
		Id:     id,
		Date:   r.FormValue("Date"),
		Time:   r.FormValue("Time"),
		Name:   r.FormValue("Name"),
		Genre:  r.FormValue("Genre"),
		Active: true,
	}
	if r.FormValue("Active") != "true" {
		event.Active = false
	}

	err = UpdateEventsDataRecord(event)
	if err != nil {
		return err
	}

	return nil
}

// Update and replace file content with added event.
func UpdateEventsData(event *Event) error {
	events, err := ReadEventsDataFile(false)
	if err != nil {
		return err
	}
	// Add the new event to events slice and save to file.
	events = append(events, *event)
	err = WriteEventsDataFile(&events)
	if err != nil {
		return err
	}

	return nil
}

// UOdate individual event "record" (slice)
func UpdateEventsDataRecord(event *Event) error {
	events, err := ReadEventsDataFile(false)
	if err != nil {
		return err
	}
	// Loop through Events slice until index matches,
	// update each field, break the loop and save it to file.
	for i := 0; i < len(events); i++ {
		if events[i].Id == event.Id {
			events[i].Date = event.Date
			events[i].Time = event.Time
			events[i].Name = event.Name
			events[i].Genre = event.Genre
			events[i].Active = event.Active
			break
		}
	}

	err = WriteEventsDataFile(&events)
	if err != nil {
		return err
	}

	return nil
}

// Remove anything that has an event "days" or older.
func purgeOldEvents(days int) error {
	active := Events{}
	events, err := ReadEventsDataFile(false)
	if err != nil {
		return err
	}
	// Prepare the purge date and purge anything after.
	purgeDate := time.Now().AddDate(0, 0, days)
	for _, event := range events {
		eventDay, err := time.Parse("01-02-2006", event.Date)
		if err != nil {
			return err
		}
		if eventDay.After(purgeDate) {
			active = append(active, event)
		}
	}
	// Save it to file.
	err = WriteEventsDataFile(&active)
	if err != nil {
		return err
	}

	return nil
}

// A little bubble sort to arange events by date.
func sortEventsByDate() error {
	events, err := ReadEventsDataFile(false)
	if err != nil {
		return err
	}
	// Loop until sorted flag stays true (all sorted).
	for {
		sorted := true
		prevDateIndex := 0

		for i := 0; i < len(events); i++ {
			// Convert date to yyyymmdd integer for comparison.
			evDate, _ := time.Parse("01-02-2006", events[i].Date)
			index, _ := strconv.Atoi(evDate.Format("20060102"))
			// Out of order -> swap
			if prevDateIndex > index {
				eventBuffer := events[i-1]
				events[i-1] = events[i]
				events[i] = eventBuffer
				sorted = false
			} else {
				prevDateIndex = index
			}
		}

		if sorted {
			break
		}
	}
	// Save it to file.
	err = WriteEventsDataFile(&events)
	if err != nil {
		return err
	}
	return nil
}

// Cumulative functions for admin's management.
func ManageEvents() error {
	// Purge events that are expired (7 days here).
	err := purgeOldEvents(-7)
	if err != nil {
		return err
	}
	// Sort events.
	err = sortEventsByDate()
	if err != nil {
		return err
	}

	return nil
}

// Format event data for home page display.
func FormatEventData(events Events) (DisplayEvents, error) {
	displayevent := DisplayEvent{}
	displayevents := DisplayEvents{}
	for _, event := range events {
		displayevent.DisplayDate = utils.FormatDisplayDate(event.Date)
		displayevent.DisplayTime = event.Time
		displayevent.DisplayName = event.Name
		displayevent.DisplayGenre = event.Genre

		displayevents = append(displayevents, displayevent)
	}

	return displayevents, nil
}
