package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	utils "eveng-cli/src/utils"
)

// EveServer ...
type EveServer struct {
	Server      string `json:"server"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	VerifyHTTPS bool   `json:"verify_https"`
}

// ServerConfig ...
func ServerConfig() EveServer {
	var eve EveServer

	content, err := ioutil.ReadFile(".eve.json")

	if err != nil {
		log.Fatalf("error while reading %v", err)
	}

	json.Unmarshal([]byte(content), &eve)

	return eve
}

// EveServerToJSON ..
// Convert an EveServer struct to a JSON byte slice
func EveServerToJSON(es EveServer) []byte {

	jsonData, err := json.Marshal(es)

	if err != nil {
		log.Println("failed marshaling data: ", err)
	}

	return jsonData

}

// Auth ...
func (es EveServer) Auth() http.Cookie {

	// Allow insecure TLS certificate
	// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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

	client := &http.Client{Timeout: time.Second * 10}
	// client := utils.SetupHTTPClient(false)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	// unetlab_session=a1e6e0b1-3435-4d69-aae6-8b93aa4de746; Path=/api/
	cookieText := utils.GetSubstring(resp.Header.Get("Set-Cookie"), "=", ";")

	cookie := http.Cookie{
		Name:  "unetlab_session",
		Value: cookieText,
	}
	return cookie
}

// GetStatus ...
func (es EveServer) GetStatus(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/status`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

// GetFolders ...
func (es EveServer) GetFolders(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/folders/`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

// GetRoles ...
func (es EveServer) GetRoles(cookie http.Cookie) *http.Response {
	url := fmt.Sprintf(`https://%s/api/list/roles`, es.Server)
	resp := getRequest(url, cookie)
	return resp
}

// GetUsers ...
func (es EveServer) GetUsers(cookie http.Cookie) *http.Response {
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
	// client := utils.SetupHTTPClient(false)
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

// HTTPReturnCodes ...
func HTTPReturnCodes(resp *http.Response) (int, error) {
	statusCode := resp.StatusCode

	// List of HTTP status codes taken from:
	// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes
	switch statusCode {
	// 2XX Status Codes
	// OK
	case 200:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Created
	case 201:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Accepted
	case 202:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Non-Authoritative Information (since HTTP/1.1)
	case 203:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// No Content
	case 204:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Reset Content
	case 205:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Partial Content (RFC 7233)
	case 206:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Multi-Status (WebDAV; RFC 4918)
	case 207:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Already Reported (WebDAV; RFC 5842)
	case 208:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// IM Used (RFC 3229)
	case 226:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil

	// 3XX Status Codes
	// Multiple Choices
	case 300:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Moved Permanently
	case 301:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Found (Previously "Moved temporarily")
	case 302:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// See Other (since HTTP/1.1)
	case 303:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Not Modified (RFC 7232)
	case 304:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Use Proxy (since HTTP/1.1)
	case 305:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Switch Proxy
	case 306:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Temporary Redirect (since HTTP/1.1)
	case 307:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil
	// Permanent Redirect (RFC 7538)
	case 308:
		fmt.Println("Status: ", statusCode)
		return statusCode, nil

	// 4XX Status Codes
	// Bad Request
	case 400:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Unauthorized (RFC 7235)
	case 401:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Payment Required
	case 402:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Forbidden
	case 403:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Not Found
	case 404:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Method Not Allowed
	case 405:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Not Acceptable
	case 406:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Proxy Authentication Required (RFC 7235)
	case 407:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Request Timeout
	case 408:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Conflict
	case 409:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Gone
	case 410:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Length Required
	case 411:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Precondition Failed (RFC 7232)
	case 412:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Payload Too Large (RFC 7231)
	case 413:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// URI Too Long (RFC 7231)
	case 414:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Unsupported Media Type (RFC 7231)
	case 415:
		fmt.Println("Status: ", statusCode)
	// Range Not Satisfiable (RFC 7233)
	case 416:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Expectation Failed
	case 417:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// I'm a teapot (RFC 2324, RFC 7168)
	case 418:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Misdirected Request (RFC 7540)
	case 421:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Unprocessable Entity (WebDAV; RFC 4918)
	case 422:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Locked (WebDAV; RFC 4918)
	case 423:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Failed Dependency (WebDAV; RFC 4918)
	case 424:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Too Early (RFC 8470)
	case 425:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Upgrade Required
	case 426:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Precondition Required (RFC 6585)
	case 428:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Too Many Requests (RFC 6585)
	case 429:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Request Header Fields Too Large (RFC 6585)
	case 431:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Unavailable For Legal Reasons (RFC 7725)
	case 451:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)

	// 5XX HTTP status codes
	// Internal Server Error
	case 500:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Not Implemented
	case 501:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Bad Gateway
	case 502:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Service Unavailable
	case 503:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Gateway Timeout
	case 504:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// HTTP Version Not Supported
	case 505:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Variant Also Negotiates (RFC 2295)
	case 506:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Insufficient Storage (WebDAV; RFC 4918)
	case 507:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Loop Detected (WebDAV; RFC 5842)
	case 508:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Not Extended (RFC 2774)
	case 510:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)
	// Network Authentication Required (RFC 6585)
	case 511:
		fmt.Println("Status: ", statusCode)
		return statusCode, errors.New(resp.Status)

	default:
		fmt.Println("Unknown Status Code :", statusCode)
		return statusCode, errors.New(resp.Status)
	}
	return statusCode, errors.New(resp.Status)
}
