package http

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
