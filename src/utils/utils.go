package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// SetupHTTPClient ...
func SetupHTTPClient(verifyHTTPS bool) *http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: verifyHTTPS}
	client := &http.Client{Transport: customTransport}

	return client
}

// GetSubstring takes a string and returns the a
// substring between two string characters.
// https://www.dotnetperls.com/between-before-after-go
func GetSubstring(str string, start string, end string) string {

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

// JSONToFileWriter takes a slice of bytes and writes it to a file.
// Implements the Writer interface ?
// https://www.geeksforgeeks.org/how-to-read-and-write-the-files-in-golang/
// https://stackoverflow.com/questions/24770403/write-struct-to-json-file-using-struct-fields-not-json-keys
func JSONToFileWriter(fileName string, bs []byte, permissions os.FileMode) {

	ioutil.WriteFile(fileName, bs, permissions)
}

// CookieToJSONFile ...
func CookieToJSONFile(filename string, c http.Cookie) {
	data, _ := json.MarshalIndent(c, "", " ")

	JSONToFileWriter(filename, data, 0600)
}

// JSONCookieFileToStruct ...
func JSONCookieFileToStruct(filename string) (http.Cookie, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	cookie := http.Cookie{}

	err = json.Unmarshal([]byte(file), &cookie)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return cookie, err
}

// PrintResponse takes a HTTP response and prints it out the the terminal.
func PrintResponse(r *http.Response) {
	fmt.Println("response Status:", r.Status)
	fmt.Println("response Headers:", r.Header)
	fmt.Print("response Body: ")

	io.Copy(os.Stdout, r.Body)
	fmt.Println("")
}
