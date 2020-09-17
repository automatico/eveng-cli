package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type eveServer struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// func main() {
//
// 	server := serverConfig()
// 	cookie := server.auth()

// List Users
// response := server.getRequest(fmt.Sprintf(`https://%s/api/users/`, server.Server), cookie)

// List Folders
// response := server.getRequest(fmt.Sprintf(`https://%s/api/folders/`, server.Server), cookie)

// List Roles
//response := server.getRequest(fmt.Sprintf(`https://%s/api/list/roles`, server.Server), cookie)

// List templates
// response := server.getRequest(fmt.Sprintf(`https://%s/api/list/templates/`, server.Server), cookie)

// List Lab
// docs are wrong
// response := server.getRequest(fmt.Sprintf(`https://%s/api/labs/ccie-sp/tess.unl`, server.Server), cookie)
// response := server.getRequest(fmt.Sprintf(`https://%s/api/folders/`, server.Server), cookie)

// List nodes in a lab
// response := server.getRequest(fmt.Sprintf(`https://%s/api/labs/ccie-sp/tess.unl/nodes`, server.Server), cookie)

// get a node in a lab
// response := server.getRequest(fmt.Sprintf(`https://%s/api/labs/ccie-sp/tess.unl/nodes/2`, server.Server), cookie)

// Stop all nodes in a lab
// response := server.getRequest(fmt.Sprintf(`https://%s/api/labs/ccie-sp/tess.unl/nodes/1/stop`, server.Server), cookie)

// Start all nodes in a lab
// response := server.getRequest(fmt.Sprintf(`https://%s/api/labs/ccie-sp/tess.unl/nodes/start`, server.Server), cookie)

// status := server.getStatus(cookie)
// printResponse(status)

// folders := server.getFolders(cookie)
// printResponse(folders)

// roles := server.getRoles(cookie)
// printResponse(roles)

// users := server.getUsers(cookie)
// printResponse(users)

// Write JSON to file
// jsonData := eveServerToJSON(server)
// jsonToFile("output.json", jsonData, 0600)
// }

func serverConfig() eveServer {
	var eve eveServer

	content, err := ioutil.ReadFile(".eve.json")

	if err != nil {
		log.Fatalf("error while reading %v", err)
	}

	json.Unmarshal([]byte(content), &eve)

	return eve
}

// Convert an eveServer struct to a JSON byte slice
func eveServerToJSON(es eveServer) []byte {

	jsonData, err := json.Marshal(es)

	if err != nil {
		log.Println("failed marshaling data: ", err)
	}

	return jsonData

}

func (es eveServer) auth() http.Cookie {

	// Allow insecure TLS certificate
	// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// https://girishjoshi.io/post/implementing-http-basic-authentication-in-golang/
	// https://ednsquare.com/story/how-to-make-http-requests-in-golang------5VIjL3
	// http://networkbit.ch/golang-http-client/
	data, _ := json.Marshal(map[string]string{
		"username": es.Username,
		"password": es.Password,
		"html5":    "0", // Set this or you will get HTML links for the Node URL's
	})

	loginURL := fmt.Sprintf(`https://%s/api/auth/login`, es.Server)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(data))
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

	// unetlab_session=a1e6e0b1-3435-4d69-aae6-8b93aa4de746; Path=/api/
	cookieText := getSubstring(resp.Header.Get("Set-Cookie"), "=", ";")

	cookie := http.Cookie{
		Name:  "unetlab_session",
		Value: cookieText,
	}
	return cookie
}

func (es eveServer) getStatus(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/status`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

func (es eveServer) getFolders(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/folders/`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

func (es eveServer) getRoles(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/list/roles`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

func (es eveServer) getUsers(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/users/`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

func getRequest(url string, cookie http.Cookie) *http.Response {
	data, _ := json.Marshal(map[string]string{})

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

	return resp

}

func postRequest() {

}

func putRequest() {

}

func deleteRequest() {

}
