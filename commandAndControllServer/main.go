package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// lipgloss styles
var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#77DD77"))
var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6961"))

// Server Setup Stuff
var serverAddress string

// Struct for rendering the report template
type ReportStatistics struct {
	RequestsSent    int
	RequestsBlocked int
}

// Struct for the test data
type TestData struct {
	Target      string //What IP / Domain Name to target for the test
	Method      string //What method to send
	NumRequests int    //How many requests to send
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
	method := r.Method
	testIsFinished := false

	//Data for rendering the template
	droppedConnections := 0
	wentThroughConnections := 0

	if method == "GET" {
		if testIsFinished {
			fmt.Println("Test is done austin's lazy ass should render the template")
		} else {
			//fmt.Println("Not done yet good things come to those who wait")
			//TODO Do some user output stuff?
		}
	} else if method == "POST" {
		fmt.Println(droppedConnections, wentThroughConnections)
	}
}

// Utility Functions

/*
* Function:   runTest()
*  Purpose:   Sends the request to the agent servers to start the tests
 */
func sendTestRequest(data TestData) bool {

	//Print a little debugging stuff
	fmt.Println(successStyle.Render("Test Data Received: "), "Target: ", data.Target, " | ", "Method: ", data.Method, " | ", "NumRequests: ", data.NumRequests)

	//make the list of servers
	serverVar := getEnvironmentVariable("SERVERS")
	serverList := strings.Split(serverVar, ",")

	//send the requests to the agentServers

	for index := range serverList {
		currServer := "http://127.0.0.1:80/" //serverList[index]
		fmt.Println(index)
		//Make the struct into a JSON object
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
	return true
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
