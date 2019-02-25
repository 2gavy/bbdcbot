package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentalVariables()
	client := &http.Client{}
	//fetching cookies
	aspxanon, sessionID := fetchCookies()

	//logging in
	loginForm := url.Values{}
	loginForm.Add("txtNRIC", os.Getenv("NRIC"))
	loginForm.Add("txtPassword", os.Getenv("PASSWORD"))
	loginForm.Add("btnLogin", " ")
	req, err := http.NewRequest("POST", "http://www.bbdc.sg/bbdc/bbdc_web/header2.asp",
		strings.NewReader(loginForm.Encode()))
	errCheck(err, "Error creating log in request")
	req.AddCookie(aspxanon)
	req.AddCookie(sessionID)
	req.AddCookie(&http.Cookie{Name: "language", Value: "en-US"})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
	errCheck(err, "Error logging in")

	//fetching the booking page
	req, err = http.NewRequest("POST", "http://www.bbdc.sg/bbdc/b-3c-pLessonBooking1.asp",
		strings.NewReader(bookingForm().Encode()))
	req.AddCookie(aspxanon)
	req.AddCookie(sessionID)
	req.AddCookie(&http.Cookie{Name: "language", Value: "en-US"})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	errCheck(err, "Error creating get bookings request")
	resp, err := client.Do(req)
	errCheck(err, "Error fetching booking slots")
	body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("booking.txt", body, 0644)

	//parse booking page to get booking dates
	//The data is hidden away in the following function call in the HTML page
	//fetched:
	//doTooltipV(event,0, "03/05/2019 (Fri)","3","11:30","13:10","BBDC");
	substrs := strings.Split(string(body), "doTooltipV(")[1:]
	for _, substr := range substrs {
		bookingData := strings.Split(substr, ",")[0:6]
		day := bookingData[2]
		sessionNum := bookingData[3]
		if strings.Contains(day, "Sat") || strings.Contains(day, "Sun") || sessionNum == "\"7\"" || sessionNum == "\"8\"" {
			alert("Slot available on " + day + " from " + bookingData[4] + " to " + bookingData[5])
		}
	}

}

func alert(msg string) {
	log.Println(msg)
}

func loadEnvironmentalVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading environmental variables: ")
		log.Fatal(err.Error())
	}
}

func fetchCookies() (*http.Cookie, *http.Cookie) {
	resp, err := http.Get("http://www.bbdc.sg/bbweb/default.aspx")
	errCheck(err, "Error fetching cookies")
	aspxanon := resp.Cookies()[0]
	resp, err = http.Get("http://www.bbdc.sg/bbdc/bbdc_web/newheader.asp")
	errCheck(err, "Error fetching cookies (sessionID)")
	sessionID := resp.Cookies()[0]
	return aspxanon, sessionID
}

func bookingForm() url.Values {
	bookingForm := url.Values{}
	bookingForm.Add("accId", os.Getenv("ACCOUNT_ID"))
	bookingForm.Add("Month", "Feb/2019")
	bookingForm.Add("Month", "Mar/2019")
	bookingForm.Add("Month", "Apr/2019")
	bookingForm.Add("Month", "May/2019")
	bookingForm.Add("Month", "Jun/2019")
	bookingForm.Add("Session", "1")
	bookingForm.Add("Session", "2")
	bookingForm.Add("Session", "3")
	bookingForm.Add("Session", "4")
	bookingForm.Add("Session", "5")
	bookingForm.Add("Session", "6")
	bookingForm.Add("Session", "7")
	bookingForm.Add("Session", "8")
	bookingForm.Add("allSes", "on")
	bookingForm.Add("Day", "2")
	bookingForm.Add("Day", "3")
	bookingForm.Add("Day", "4")
	bookingForm.Add("Day", "5")
	bookingForm.Add("Day", "6")
	bookingForm.Add("Day", "7")
	bookingForm.Add("Day", "1")
	bookingForm.Add("allDay", "")
	bookingForm.Add("defPLVenue", "1")
	bookingForm.Add("optVenue", "1")

	return bookingForm
}

func errCheck(err error, msg string) {
	if err != nil {
		log.Fatal(msg + ": " + err.Error())
	}
}
