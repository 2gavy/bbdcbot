package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentalVariables()

	//fetching cookies
	resp, err := http.Get("http://www.bbdc.sg/bbweb/default.aspx")
	errCheck(err, "Error fetching cookies")
	aspxanon := resp.Cookies()[0]
	resp, err = http.Get("http://www.bbdc.sg/bbdc/bbdc_web/newheader.asp")
	errCheck(err, "Error fetching cookies (sessionID)")
	sessionID := resp.Cookies()[0]

	//logging in
	loginForm := url.Values{}
	loginForm.Add("txtNRIC", os.Getenv("NRIC"))
	loginForm.Add("txtPassword", os.Getenv("PASSWORD"))
	loginForm.Add("btnLogin", "+")
	req, err := http.NewRequest("POST", "http://www.bbdc.sg/bbdc/bbdc_web/header2.asp",
		strings.NewReader(loginForm.Encode()))
	errCheck(err, "Error creating log in request")
	req.AddCookie(aspxanon)
	req.AddCookie(sessionID)
	req.AddCookie(&http.Cookie{Name: "language", Value: "en-US"})
	client := &http.Client{}
	resp, err = client.Do(req)
	errCheck(err, "Error logging in")
	log.Println(resp.ContentLength)
	log.Println(req.ContentLength)
}

func loadEnvironmentalVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading environmental variables: ")
		log.Fatal(err.Error())
	}
}

func errCheck(err error, msg string) {
	if err != nil {
		log.Fatal(msg + ": " + err.Error())
	}
}
