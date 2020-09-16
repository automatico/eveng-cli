/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type eveServer struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func printResponse(r *http.Response) {
	fmt.Println("response Status:", r.Status)
	fmt.Println("response Headers:", r.Header)
	fmt.Print("response Body: ")

	io.Copy(os.Stdout, r.Body)
	fmt.Println("")
}

func serverConfig() eveServer {
	var eve eveServer

	content, err := ioutil.ReadFile(".eve.json")

	if err != nil {
		log.Fatalf("error while reading %v", err)
	}

	json.Unmarshal([]byte(content), &eve)

	return eve
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

// labCmd represents the lab command
var labCmd = &cobra.Command{
	Use:   "lab",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("lab called")
		server := serverConfig()
		cookie := server.auth()

		status := server.getStatus(cookie)
		printResponse(status)
	},
}

func init() {
	rootCmd.AddCommand(labCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// labCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// labCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
