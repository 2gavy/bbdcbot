package http

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	bbdcbot "github.com/SKAshwin/bbdcbot"
)

//BBDCService an implementation of bbdcbot.SlotService which
//uses HTTP requests to the BBDC website to book bbdc
//slots
type BBDCService struct {
	aspxanon  *http.Cookie
	sessionID *http.Cookie
	client    *http.Client
	AccountID string
	Logger    *log.Logger
}

//StartService logs into the BBDC website with the given credentials
//and returns a BBDCService object which can then interact with the user's
//BBDC Slots.
//Returns error if log in or cookie fetching could not be achieved
func StartService(nric string, pwd string, accID string) (*BBDCService, error) {
	bs := &BBDCService{
		AccountID: accID,
		Logger:    log.New(os.Stderr, "", log.LstdFlags),
		client:    &http.Client{},
	}

	bs.Logger.Println("Fetching Cookies")
	aspxanon, sessionID, err := bs.fetchCookies()
	if err != nil {
		return &BBDCService{}, errors.New("Error fetching cookies: " + err.Error())
	}

	bs.aspxanon = aspxanon
	bs.sessionID = sessionID
	err = bs.logIn(nric, pwd)
	if err != nil {
		return &BBDCService{}, errors.New("Error logging in: " + err.Error())
	}

	return bs, nil
}

//Book books the given slot by making a HTTP request to the BBDC website
func (bs *BBDCService) Book(slot bbdcbot.Slot) error {
	form := url.Values{}
	form.Add("accId", bs.AccountID)
	form.Add("slot", slot.ID)
	req, err := bs.newFormRequest("POST", "http://www.bbdc.sg/bbdc/b-3c-pLessonBookingDetails.asp", form)
	if err != nil {
		return errors.New("Error creating booking request: " + err.Error())
	}
	_, err = bs.client.Do(req)
	if err != nil {
		return errors.New("Error making booking request: " + err.Error())
	}
	return nil
}

//AvailableSlots returns all available BBDC slots by making a HTTP request to the website
func (bs *BBDCService) AvailableSlots() ([]bbdcbot.Slot, error) {
	//fetching the booking page
	bs.Logger.Println("Fetching booking page")
	req, err := bs.newFormRequest("POST", "http://www.bbdc.sg/bbdc/b-3c-pLessonBooking1.asp",
		bs.bookingForm())
	if err != nil {
		return nil, errors.New("Error creating get slot request: " + err.Error())
	}

	resp, err := bs.client.Do(req)
	if err != nil {
		return nil, errors.New("Error making get slot request: " + err.Error())
	}
	body, _ := ioutil.ReadAll(resp.Body)

	return bs.parseBookingSlotReqBody(string(body))
}

//TODO
func (bs *BBDCService) parseBookingSlotReqBody(body string) ([]bbdcbot.Slot, error) {
	//The data is hidden away in the following function call in the HTML page
	//fetched:
	//doTooltipV(event,0, "03/05/2019 (Fri)","3","11:30","13:10","BBDC");
	//substrs := strings.Split(body, "doTooltipV(")[1:]
	//slots := make([]bbdcbot.Slot, len(substrs))
	//for _, substr := range substrs {
	//bookingData := strings.Split(substr, ",")[0:6]
	//rawDay := bookingData[2] //"03/05/2019 (Fri)"
	//layout := "02/01/2006"
	//day, err := time.Parse(layout, strings.Split(strings.Split(rawDay, "\"")[1], " ")[0])
	//if err != nil {
	//	return nil, errors.New("Error parsing day: " + err.Error())
	//}
	//start := strings.Split(bookingData[4], "\"")[1] //11:30
	//end := strings.Split(bookingData[5], "\"")[1]   //13:10

	//sessionNum := bookingData[3]
	//}
	return nil, nil
}

func (bs *BBDCService) logIn(nric string, pwd string) error {
	bs.Logger.Println("Logging in")
	loginForm := url.Values{}
	loginForm.Add("txtNRIC", nric)
	loginForm.Add("txtPassword", pwd)
	loginForm.Add("btnLogin", " ")
	req, err := bs.newFormRequest("POST", "http://www.bbdc.sg/bbdc/bbdc_web/header2.asp", loginForm)
	if err != nil {
		return errors.New("Error creating request: " + err.Error())
	}

	_, err = bs.client.Do(req)
	if err != nil {
		return errors.New("Error making request: " + err.Error())
	}

	return nil
}

func (bs *BBDCService) newFormRequest(method string, addr string, form url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", "http://www.bbdc.sg/bbdc/bbdc_web/header2.asp",
		strings.NewReader(form.Encode()))
	if err != nil {
		return req, errors.New("Error creating request: " + err.Error())
	}
	req.AddCookie(bs.aspxanon)
	req.AddCookie(bs.sessionID)
	req.AddCookie(&http.Cookie{Name: "language", Value: "en-US"})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (bs *BBDCService) fetchCookies() (*http.Cookie, *http.Cookie, error) {
	resp, err := http.Get("http://www.bbdc.sg/bbweb/default.aspx")
	if err != nil {
		return nil, nil, errors.New("Error fetching aspxanon cookie: " + err.Error())
	}
	aspxanon := resp.Cookies()[0]
	resp, err = http.Get("http://www.bbdc.sg/bbdc/bbdc_web/newheader.asp")
	if err != nil {
		return nil, nil, errors.New("Error fetching sessionID cookie: " + err.Error())
	}
	sessionID := resp.Cookies()[0]
	return aspxanon, sessionID, nil
}

func (bs *BBDCService) bookingForm() url.Values {
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
