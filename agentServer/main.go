package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(errorStyle.Render("Error "), "Unable to unmarshal JSON data from test request!", "\n", err)
		} else {
			fmt.Println(successStyle.Render("Success "), "Data received from command and control server.")
			fmt.Println("Target: ", data.Target, " | ", "Method: ", data.Method, " | ", "Number of Requests ", data.NumRequests)
			fmt.Println("Starting Test ......")
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
