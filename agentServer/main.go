package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

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
		fmt.Println("Post Request Recieved")
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
