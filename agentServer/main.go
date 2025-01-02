package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// Struct for the test data
type TestData struct {
	Target      string `json:"target"`      //What IP / Domain Name to target for the test
	Method      string `json:"method"`      //What method to send
	NumRequests int    `json:"numRequests"` //How many requests to send
}

type TestDataReturn struct {
	Accepted int32 `json:"acceptedConnections"` //Connections that went through
	Dropped  int32 `json:"droppedConnections"`  //Connections that did not go through

}

// lipgloss styles
var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#77DD77"))
var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6961"))

// Server Setup Stuff
var serverAddress string
var commandAndControlAddress string

func main() {
	//Load The environment variable
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	//Load the server address
	serverAddress = getEnvironmentVariable("webServerAddressAgent")

	//Load the command and control server address
	commandAndControlAddress = getEnvironmentVariable("commandAndControlServerAddress")

	//HTTP Routes
	http.HandleFunc("/", processTestRequest)

	//Startup the server
	http.ListenAndServe(serverAddress, nil)
}

// HTTP Routes
func processTestRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data TestData // Struct to hold requested data

		// Read the JSON data from the request body and store it in the struct
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(errorStyle.Render("Error "), "Unable to unmarshal JSON data from test request!", "\n", err)
		} else {
			fmt.Println(successStyle.Render("Success "), "Data received from command and control server.")
			fmt.Println("Target: ", data.Target, " | ", "Method: ", data.Method, " | ", "Number of Requests ", data.NumRequests)

			testResults := runTest(data)
			fmt.Println(successStyle.Render("Success "), "Test Completed Successfully")
			fmt.Println("Accepted Connections: ", testResults[0], " | ", "Dropped Connections: ", testResults[1])

			testResultsToSend := TestDataReturn{testResults[0], testResults[1]}

			//Make the struct into a JSON object
			jsonData, _ := json.Marshal(testResultsToSend)

			//Construct the request
			req, _ := http.NewRequest("POST", commandAndControlAddress, bytes.NewBuffer((jsonData)))
			req.Header.Set("Content-Type", "application/json")

			//Send it out
			client := &http.Client{}
			_, err := client.Do(req)
			if err != nil {
				fmt.Println(errorStyle.Render("Error "), "Could not send request")
			}

			fmt.Println(successStyle.Render("Success"), "Sent data to server", commandAndControlAddress)
		}

	} else {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	}
}

//Utility Functions
/*
* Function:   getEnvironmentVariable
*  Purpose:   Get an environment variable
 */
func getEnvironmentVariable(varName string) string {
	_, exists := os.LookupEnv(varName)

	if exists {
		fmt.Println(successStyle.Render("Success "), varName, " has been successfully retrieved")
		return os.Getenv(varName)
	} else {
		fmt.Println(errorStyle.Render("Error ", varName, " could not be retrieved"))
		return ":("
	}
}

func runTest(data TestData) [2]int32 {
	fmt.Println("Starting Test ......")
	var testResults [2]int32
	acceptedConnections := 0
	droppedConnections := 0

	for i := 0; i < data.NumRequests; i++ {
		resp, err := http.Get(data.Target)
		if err != nil {
			fmt.Println(errorStyle.Render("Error "), "Could not send request")
		}
		if resp.StatusCode != http.StatusOK {
			droppedConnections += 1
		} else {
			acceptedConnections += 1
		}
	}

	testResults[0] = int32(acceptedConnections)
	testResults[1] = int32(droppedConnections)

	return testResults
}
