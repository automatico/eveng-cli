package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type eveServer struct {
	Server   string
	Username string
	Password string
}

func main() {

	config := config()
	cookie := auth(config)
	getStatus(config, cookie)

}

func config() eveServer {
	var eve eveServer

	content, err := ioutil.ReadFile(".eve.json")

	if err != nil {
		log.Fatalf("error while reading %v", err)
	}

	json.Unmarshal([]byte(content), &eve)

	return eve
}

func auth(server eveServer) http.Cookie {

	// Allow insecure TLS certificate
	// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// https://girishjoshi.io/post/implementing-http-basic-authentication-in-golang/
	// https://ednsquare.com/story/how-to-make-http-requests-in-golang------5VIjL3
	// http://networkbit.ch/golang-http-client/
	data, _ := json.Marshal(map[string]string{
		"username": server.Username,
		"password": server.Password,
	})

	req, err := http.NewRequest("POST", fmt.Sprintf(`https://%s/api/auth/login`, server.Server), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// fmt.Println("response Cookie:", resp.Cookies())

	// unetlab_session=a1e6e0b1-3435-4d69-aae6-8b93aa4de746; Path=/api/
	cookieText := getSubstring(resp.Header.Get("Set-Cookie"), "=", ";")
	fmt.Println("cookieText: ", cookieText)

	cookie := http.Cookie{
		Name:  "unetlab_session",
		Value: cookieText,
	}
	return cookie
}

func getStatus(server eveServer, cookie http.Cookie) {
	data, _ := json.Marshal(map[string]string{})

	req, err := http.NewRequest("GET", fmt.Sprintf(`https://%s/api/status`, server.Server), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&cookie)

	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", string(body))
}

func getSubstring(str string, start string, end string) string {
	// https://www.dotnetperls.com/between-before-after-go
	// Get substring between two strings.
	posFirst := strings.Index(str, start)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(str, end)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(start)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return str[posFirstAdjusted:posLast]
}
