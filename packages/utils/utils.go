package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	adminToken        = "DDS_ADMIN_TOKEN"
	carouselImagePath = "static/img/carousel/"
)

var (
	adminPwd = os.Getenv("DDS_ADMIN_PASSWORD")
)

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

// This function is currently not used --> remove?
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
			"/static/img/carousel/IMG_9999.png",
			"/static/img/carousel/IMG_3015.jpg",
		},
	}
}
