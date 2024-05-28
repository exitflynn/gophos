package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	loginURL  = "http://172.16.68.6:8090/login.xml"
	logoutURL = "http://172.16.68.6:8090/logout.xml"
)

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

func login(username, password string) (string, error) {
	body := fmt.Sprintf("mode=191&username=%s&password=%s", username, password)
	response, err := sendPostRequest(loginURL, body)
	if err != nil {
		return "", err
	}

	if strings.Contains(response, "You are signed in as") {
		return fmt.Sprintf("[+] SUCCESS User = %s, Pass = %s", username, password), nil
	}

	return fmt.Sprintf("[-] %s Failure :( ", username), nil
}

func logout(username string) (string, error) {
	body := fmt.Sprintf("mode=193&username=%s", username)
	_, err := sendPostRequest(logoutURL, body)
	if err != nil {
		return "", err
	}

	return "[+] Successfully logged out", nil
}

func main() {
	username := "enroll_id"
	password := "password"

	loginResult, err := login(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(loginResult)

	// Call logout function if needed
	logoutResult, err := logout(username)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(logoutResult)
}

