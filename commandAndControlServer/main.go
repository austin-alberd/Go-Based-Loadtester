package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss" // Makes the terminal look pretty
	"github.com/joho/godotenv"          // Environment Variables in Go
)

// lipgloss styles
var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#77DD77"))
var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6961"))

// Global Vars
var serverAddress string        // Server address
var reportData ReportStatistics // Report data for rendering the template

var totalServers int
var testDataReceived int

// Struct for rendering the report template
type ReportStatistics struct {
	RequestsSent    int
	RequestsBlocked int
}

// Struct for the test data
type TestData struct {
	Target      string `json:"target"`      //What IP / Domain Name to target for the test
	Method      string `json:"method"`      //What method to send
	NumRequests int    `json:"numRequests"` //How many requests to send
}

// Struct to handle the data returned by the test
type TestDataReturn struct {
	Accepted int32 `json:"acceptedConnections"` //Connections that went through
	Dropped  int32 `json:"droppedConnections"`  //Connections that did not go through

}

func main() {
	//Load the environment variables.
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	//Get the server address from the environment variables
	serverAddress = getEnvironmentVariable("webServerAddressCaC")

	// Serve static files
	fs := http.FileServer(http.Dir("staticFiles/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	//Serve the pages
	http.HandleFunc("/", indexRoute)
	http.HandleFunc("/testData", dataReturn)

	//Start the server
	http.ListenAndServe(serverAddress, nil)
}

// HTTP Routes
func indexRoute(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "GET" {
		http.ServeFile(w, r, "webFiles/index.html")
	} else if method == "POST" {

		//Get the data from the form and put it into the struct
		numRequestsSTR := r.FormValue("numRequests")
		numRequests, _ := strconv.Atoi(numRequestsSTR)
		testVals := TestData{r.FormValue("target"), r.FormValue("method"), numRequests}

		//Start the test
		sendTestRequest(testVals)

	} else {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	}
}

func dataReturn(w http.ResponseWriter, r *http.Request) {
	var data TestDataReturn //Struct to hold all of the data returned by the remote server

	if r.Method == "GET" {
		if totalServers == testDataReceived && testDataReceived != 0 {
			t, _ := template.ParseFiles("webFiles/templates/reportsStatsPage.html")
			t.Execute(w, reportData)
		}

	} else if r.Method == "POST" {

		//Unmarshal the data
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &data)

		fmt.Println(successStyle.Render("Success"), " Data received from the remote server", data.Accepted, " | ", data.Dropped)
		reportData.RequestsBlocked += int(data.Dropped)
		reportData.RequestsSent += int(data.Accepted)

		testDataReceived += 1
	}
}

// Utility Functions

/*
* Function:   runTest()
*  Purpose:   Sends the request to the agent servers to start the tests
 */
func sendTestRequest(data TestData) {
	//Reset the variables
	reportData = ReportStatistics{0, 0}
	testDataReceived = 0
	totalServers = 0

	//Print a little debugging stuff
	fmt.Println(successStyle.Render("Test Data Received: "), "Target: ", data.Target, " | ", "Method: ", data.Method, " | ", "NumRequests: ", data.NumRequests)

	//make the list of servers
	serverVar := getEnvironmentVariable("SERVERS")
	serverList := strings.Split(serverVar, ",")

	//send the requests to the agentServers

	for index := range serverList {
		currServer := serverList[index]
		totalServers += 1

		//Make the test data into a JSON object to be sent to the agent Server
		jsonData, _ := json.Marshal(data)

		//Construct the request
		req, _ := http.NewRequest("POST", currServer, bytes.NewBuffer((jsonData)))
		req.Header.Set("Content-Type", "application/json")

		//Send it out
		client := &http.Client{}
		_, err := client.Do(req)
		if err != nil {
			fmt.Println(errorStyle.Render("Error "), "Could not send request to ", currServer, "\n", err)
		} else {
			fmt.Println(successStyle.Render("Success: "), "Sent start request to agent server: ", currServer)
		}

	}
}

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
