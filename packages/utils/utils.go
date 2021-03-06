package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	adminToken        = "DDS_ADMIN_TOKEN"
	weatherApiToken   = "DDS_WEATHERAPI"
	carouselImagePath = "static/img/carousel/"
)

var (
	wxURL    = "http://api.weatherapi.com/v1/current.json?key={{token}}&q=Haines%20City&aqi=no"
	astroURL = "http://api.weatherapi.com/v1/astronomy.json?key={{token}}&q=Haines%20City&dt={{today}}"
	adminPwd = os.Getenv("DDS_ADMIN_PASSWORD")
	client   *http.Client
)

// For data reference, see: https://www.weatherapi.com/api-explorer.aspx
type WxData struct {
	Location interface{}   `json:"location"`
	Current  CurrentWxData `json:"current"`
}

type CurrentWxData struct {
	LastUpdateEpoch int       `json:"last_updated_epoch"`
	TempF           float64   `json:"temp_f"`
	Condition       Condition `json:"condition"`
	WindMPH         float64   `json:"wind_mph"`
	WindDirection   string    `json:"wind_dir"`
	Humidity        int       `json:"humidity"`
	Sunset          string    `json:"sunset"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type AstronomyData struct {
	Location  Location `json:"location"`
	Astronomy Astro    `json:"astronomy"`
}

type Location struct {
	Name    string `json:"name"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type Astro struct {
	AstroData AstroData `json:"astro"`
}

type AstroData struct {
	Sunrise      string `json:"sunrise"`
	Sunset       string `json:"sunset"`
	AstroFiller2 interface{}
}

// Check if correct password was provided to gin access to admin menu and admin functionalities.
func CheckAdminPassword(r *http.Request) error {
	// Check if passphrase is correct.
	r.ParseForm()
	if r.FormValue("Password") != adminPwd {
		return fmt.Errorf("error: %s", "bad password")
	}

	// If so, create bearer token.
	token := make([]byte, 256)
	_, err := rand.Read(token)
	if err != nil {
		return err
	}

	// Set it for verification
	err = os.Setenv(adminToken, fmt.Sprintf("%x", token))
	if err != nil {
		return err
	}

	return nil
}

// Admin security (obsolete?).
func SetTokenHeader(r *http.Request) {
	token := os.Getenv(adminToken)
	bearer := fmt.Sprintf("Bearer %x", token)
	r.Header.Add("Authorization", bearer)
}

// Admin security (obsolete?).
func RemoveBearerToken() error {
	_, found := os.LookupEnv(adminToken)
	if found {
		_ = os.Unsetenv(adminToken)
	}

	return nil
}

// Format event date for homepage display.
func FormatDisplayDate(date string) string {
	// Day of week representation
	var daysOfWeek = map[time.Weekday]string{
		time.Sunday:    "Sun",
		time.Monday:    "Mon",
		time.Tuesday:   "Tue",
		time.Wednesday: "Wed",
		time.Thursday:  "Thu",
		time.Friday:    "Fri",
		time.Saturday:  "Sat",
	}
	// Month of year representation
	var monthsOfYear = map[time.Month]string{
		time.January:   "Jan",
		time.February:  "Feb",
		time.March:     "Mar",
		time.April:     "Apr",
		time.May:       "May",
		time.June:      "Jun",
		time.July:      "Jul",
		time.August:    "Aug",
		time.September: "Sep",
		time.October:   "Oct",
		time.November:  "Nov",
		time.December:  "Dec",
	}

	dtObj, _ := time.Parse("01-02-2006", date)

	dow := daysOfWeek[dtObj.Weekday()]
	moy := monthsOfYear[dtObj.Month()]
	day := dtObj.Day()

	return fmt.Sprintf("%s, %s %d", dow, moy, day)
}

// NOTE: Code below needs to be double-checked and paired with images of appropriate size.
// This function is currently not used --> remove?
// Reads static/img/carousel directory
func GetCarouselImages() []fs.FileInfo {
	carouselImageFiles, err := ioutil.ReadDir(carouselImagePath)
	if err != nil {
		fmt.Println(err)
	}
	// Sanitize slice by removing hidden files and sub-directories.
	for {
		done := true
		for i := 0; i < len(carouselImageFiles); i++ {
			if carouselImageFiles[i].Name()[0] == 46 || carouselImageFiles[i].IsDir() {
				carouselImageFiles[i] = carouselImageFiles[len(carouselImageFiles)-1]
				carouselImageFiles = carouselImageFiles[:len(carouselImageFiles)-1]
				done = false
				break
			}
		}

		if done {
			break
		}
	}

	return carouselImageFiles
}

// This function is currently not used (see GetDefaultImages() ) --> remove?
func GetRandomCarouselImages(count int) []string {
	var (
		result []string
		images = GetCarouselImages()
	)

	fmt.Println("==>", len(images))

	for {
		ok := true
		rand.Seed(time.Now().UnixNano())
		newIndex := rand.Intn(len(images) - 1)
		fmt.Println("~~~", newIndex, "~~~")

		for i := 0; i < len(result); i++ {
			if result[i] == images[newIndex].Name() {
				ok = false
				break
			}
		}

		if ok {
			result = append(result, images[newIndex].Name())
		}

		if len(result) == count {
			break
		}
	}
	fmt.Println(result)

	return result
}

// Temp. solution to populate the carousel section of the home page.
func GetDefaultImages() [][]string {
	return [][]string{
		{
			"/static/img/carousel/output-image1647987426793.JPG",
			"/static/img/carousel/output-image1647987375846.JPG",
			"/static/img/carousel/imagejpeg-0.JPG",
		},
		{
			"/static/img/carousel/IMG_3018.jpg",
			"/static/img/carousel/output-image1647987394780.JPG",
			"/static/img/carousel/IMG_2995.jpg",
		},
		{
			"/static/img/carousel/IMG_3016.jpg",
			"/static/img/carousel/IMG_2848.jpg",
			"/static/img/carousel/IMG_2998.jpg",
		},
		{
			"/static/img/carousel/IMG_2849.jpg",
			"/static/img/carousel/IMG_2991.jpg",
			"/static/img/carousel/IMG_3015.jpg",
		},
	}
}

func CurrentWeather() (WxData, error) {
	var wxData WxData
	var wxURL = strings.Replace(wxURL, "{{token}}", os.Getenv(weatherApiToken), 1)

	err := getJson(wxURL, &wxData)
	if err != nil {
		return WxData{}, err
	}
	// rounding the temperature & windspeed... (no decimals needed).
	wxData.Current.TempF = math.Round(wxData.Current.TempF)
	wxData.Current.WindMPH = math.Round(wxData.Current.WindMPH)

	var astronomyData AstronomyData
	var astroURL = strings.Replace(astroURL, "{{token}}", os.Getenv(weatherApiToken), 1)
	astroURL = strings.Replace(astroURL, "{{today}}", time.Now().Format("2006-01-02"), 1)

	err = getJson(astroURL, &astronomyData)
	if err != nil {
		return WxData{}, err
	}
	// Add the sunset time to the weather object.
	wxData.Current.Sunset = astronomyData.Astronomy.AstroData.Sunset

	return wxData, nil
}

func getJson(url string, target interface{}) error {
	client = &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&target)
}
