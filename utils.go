package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

func jsonToFile(fileName string, bs []byte, permissions os.FileMode) {
	// https://www.geeksforgeeks.org/how-to-read-and-write-the-files-in-golang/
	// https://stackoverflow.com/questions/24770403/write-struct-to-json-file-using-struct-fields-not-json-keys

	ioutil.WriteFile(fileName, bs, permissions)
}

func printResponse(r *http.Response) {
	fmt.Println("response Status:", r.Status)
	fmt.Println("response Headers:", r.Header)
	fmt.Print("response Body: ")

	io.Copy(os.Stdout, r.Body)
	fmt.Println("")
}
