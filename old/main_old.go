package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var serverUrl string = "http://140.113.89.60:5000"

const useModel = "gpt-3.5-turbo"

func GetRequest(path string) (res []byte, err error) {
	totalPath := strings.Join([]string{serverUrl, path}, "/")

	if path == "" {
		totalPath = serverUrl
	}

	resp, err := http.Get(totalPath)
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return res, err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Read and print the response body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return res, err
	}
	res = body
	return
}

func PostRequest(path string, contentType string, sendBody map[string]string) (res []byte, err error) {

	totalPathUrl := strings.Join([]string{serverUrl, path}, "/")

	if path == "" {
		totalPathUrl = serverUrl
	}

	bytesRes, _ := json.Marshal(sendBody)

	request, err := http.NewRequest("POST", totalPathUrl, bytes.NewBuffer(bytesRes))

	if err != nil {
		fmt.Println("Error creating the request:", err)
		return res, err
	}

	request.Header.Set("Content-Type", contentType)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return res, err
	}

	defer resp.Body.Close()

	// fmt.Println("Response Status:", resp.Status)

	// Read and print the response body
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(resp.Body)
	// fmt.Println("Response Body:", buf.String())
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return res, err
	}
	res = body

	return res, err
}

func main_1() {
	// Replace with the URL you want to send the GET request to

	// Send a GET request
	// res := GetRequest("")
	// fmt.Println(string(res))

	scanner := bufio.NewReader(os.Stdin)
	fmt.Print("please input the server ip and port like [127.0.0.1:5000]:")
	inputIP, _ := scanner.ReadString('\n')

	inputIP = strings.TrimSpace(inputIP)

	serverUrl = fmt.Sprintf("http://%s", inputIP)

	for {
		fmt.Print(">")
		b2 := bufio.NewReader(os.Stdin)
		inputStr, _ := b2.ReadString('\n')

		jsonInput := map[string]string{
			"model":   useModel,
			"message": inputStr,
		}

		res, err := PostRequest("chat",
			"application/json",
			jsonInput,
		)

		if err != nil {
			fmt.Println("Error reading the response:", err)
			return
		}

		var resultJsonMap map[string]string

		err = json.Unmarshal(res, &resultJsonMap)

		if err != nil {
			fmt.Println("Error reading the response:", err)
			return
		}

		fmt.Println(resultJsonMap["message"])
	}

}
