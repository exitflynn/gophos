package main

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	url = "http://172.16.68.6:8090/login.xml"
)

func login(user, passwd string) string {
	reqBody := strings.NewReader(fmt.Sprintf("mode=191&username=%s&password=%s", user, passwd))
	resp, err := http.Post(url, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if len(resp.Header["Content-Length"]) > 0 && resp.Header["Content-Length"][0] >= "10" && resp.Body != nil {
			buf := make([]byte, 100)
			_, err := resp.Body.Read(buf)
			if err != nil {
				return fmt.Sprintf("Error reading response body: %v", err)
			}
			if buf[90] == 'Y' {
				return fmt.Sprintf("[+] SUCCESS User = %s, Pass = %s", user, passwd)
			}
		}
	}
	if err != nil {
		return fmt.Sprintf("[-] %s Failure :( (Error reading response: %v)", user, err)
	}
	return fmt.Sprintf("[-] %s Failure :( ", user)
}

func logout(user string) string {
	reqBody := strings.NewReader(fmt.Sprintf("mode=193&username=%s", user))
	resp, err := http.Post(url, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()

	// Process logout response if needed

	return ""
}

func main() {
	user := "enroll_id"
	passwd := "password"

	loginResult := login(user, passwd)
	fmt.Println(loginResult)

	// // Call logout function if needed
	// logoutResult := logout(user)
	// fmt.Println(logoutResult)
}
