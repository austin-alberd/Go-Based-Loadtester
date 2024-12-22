package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// lipgloss styles
var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#77DD77"))
var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6961"))

// Server Setup Stuff
var serverPort string = ":8080"

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
	// Serve static files
	fs := http.FileServer(http.Dir("staticFiles/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", getPage)
	http.HandleFunc("/testData", testData)

	http.ListenAndServe(serverPort, nil)
}

// HTTP Routes
func getPage(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "GET" {
		http.ServeFile(w, r, "webFiles/index.html")
	} else if method == "POST" {
		//Get the template file ready
		t, _ := template.ParseFiles("webFiles/templates/reportsStatsPage.html")

		//Get the data from the form and put it into the struct
		numRequestsSTR := r.FormValue("numRequests")
		numRequests, _ := strconv.Atoi(numRequestsSTR)
		testVals := TestData{r.FormValue("target"), r.FormValue("method"), numRequests}

		runTest(testVals)

		stats := ReportStatistics{100, 50}

		// Run the template
		t.Execute(w, stats)
	} else {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	}
}

func testData(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "GET" {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	} else if method == "POST" {
		//TODO Implement data handling
	}
}

// Utility Functions
func runTest(data TestData) bool {
	fmt.Println(successStyle.Render("Test Data Received: "), "Target: ", data.Target, " | ", "Method: ", data.Method, " | ", "NumRequests: ", data.NumRequests)

	//Load the servers from the environmentVariable
	err := godotenv.Load()
	if err != nil {
		//If there is an error with loading the .env file check if the SERVERS environment variable exists
		_, exists := os.LookupEnv("SERVERS")
		if !exists {
			fmt.Println(errorStyle.Render("Error: "), "Could not find SERVERS environment variable.")
			return false
		}
	} else {
		fmt.Println(successStyle.Render("Success: "), "Servers environment variable loaded")
	}
	//make the list of servers
	serverList := strings.Split(os.Getenv("SERVERS"), ",")
	
	//send the requests to the agentServers

	for index := range serverList{
		fmt.Println(successStyle.Render("Success: "), "Sent start request to agent server: ", serverList[index])
	}
	return true
}
