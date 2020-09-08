package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type eveServer struct {
	ipv4Address string
	username    string
	password    string
}

func main() {
	eve := eveServer{
		ipv4Address: "",
		username:    "",
		password:    "",
	}

	// Allow insecure TLS certificate
	// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get(fmt.Sprintf("http://%s", eve.ipv4Address))

	if err != nil {
		fmt.Println("dah who turned out the lights")
	}

	fmt.Println(resp)
}
