package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func ChatPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

type RequestBody struct {
	Prompt string `json:"prompt"`
}

func SendMessage(w http.ResponseWriter, r *http.Request, state, message string) {
	// TODO: send the error message
	response := map[string]string{
		"status":  "fail",
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
	log.Fatal(http.ListenAndServe(":8085", nil))
}
