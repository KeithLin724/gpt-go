package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func PostRequest2(path string, contentType string, sendBody map[string]string) (res []byte, err error) {

	totalPathUrl := path

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

type FetchResult struct {
	URL     string
	State   string
	Connect bool

	mx sync.Mutex
}

var globalResultChan = FetchResult{}

func fetchHTTP(url string) {
	resp, err := http.Get(url)
	state := "Success"
	connect := true

	if err != nil {
		state = fmt.Sprintf("Error: %s", err)
		connect = false
	}

	globalResultChan.mx.Lock()

	globalResultChan.URL = url
	globalResultChan.State = state
	globalResultChan.Connect = connect

	globalResultChan.mx.Unlock()

	if connect {
		defer resp.Body.Close()
	}
}

func ChatPage(w http.ResponseWriter, r *http.Request) {
	// url := os.Getenv("SERVER_URL")
	// resultChan := make(chan FetchResult)
	// fmt.Println(url)
	// go func() {
	// 	for {
	// 		// fmt.Println("Hello")
	// 		fetchHTTP(url, resultChan)
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		result := <-resultChan
	// 		fmt.Printf("%s\n", result.State)
	// 		if result.Connect {
	// 			http.ServeFile(w, r, "./index.html")
	// 		} else {
	// 			http.ServeFile(w, r, "./error.html")
	// 		}
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		fmt.Println("hello")
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	globalResultChan.mx.Lock()
	connect := globalResultChan.Connect
	globalResultChan.mx.Unlock()

	if connect {

		http.ServeFile(w, r, "./index.html")
	} else {
		http.ServeFile(w, r, "./error.html")

	}

	fmt.Println("hello")

}

type RequestBody struct {
	Prompt string `json:"prompt"`
}

// for get the server state

func SendMessage(w http.ResponseWriter, r *http.Request, state, message string) {
	// TODO: send the error message
	response := map[string]string{
		"status":  state,
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// for the use send to the lab server
func SendApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		// TODO: send the error message
		SendMessage(w, r, "fail", "Method not allowed")
		return
	}

	// TODO: Get Chat Url
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		SendMessage(w, r, "fail", err.Error())
		return
	}

	chapAPIUrl := os.Getenv("SERVER_API_URL")

	// Decode the JSON request body into the RequestBody struct
	var requestBody RequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)

		// TODO: send the error message
		SendMessage(w, r, "fail", err.Error())
		return
	}

	// TODO: get the request for the ./index.html and display
	fmt.Printf("Received Prompt: %s\n", requestBody.Prompt)

	//TODO: send the request to the sever get gpt request
	jsonInput := map[string]string{
		"prompt": requestBody.Prompt,
	}
	res, err := PostRequest2(
		chapAPIUrl,
		"application/json",
		jsonInput,
	)

	if err != nil {
		fmt.Println("Error reading the response:", err)
		SendMessage(w, r, "fail", err.Error())
		return
	}
	//TODO:decode the request
	var resultJsonMap map[string]string

	err = json.Unmarshal(res, &resultJsonMap)

	if err != nil {
		fmt.Println("Error reading the response:", err)
		SendMessage(w, r, "fail", err.Error())
		return
	}

	resultStringMessage := resultJsonMap["message"]

	// TODO: send the message to the index.html
	SendMessage(w, r, "success", resultStringMessage)
	// response := map[string]string{
	// 	"status":  "success",
	// 	"message": resultStringMessage,
	// }
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)

}

func main() {
	http.HandleFunc("/", ChatPage)
	http.HandleFunc("/send", SendApi)

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	url := os.Getenv("SERVER_URL")
	fmt.Println(url)

	go func() {
		for {
			// invFunc(testCheck)

			fetchHTTP(url)

			globalResultChan.mx.Lock()
			fmt.Println(globalResultChan.State)
			globalResultChan.mx.Unlock()
			time.Sleep(10 * time.Second)

		}

	}()

	log.Fatal(http.ListenAndServe(":8085", nil))
}
