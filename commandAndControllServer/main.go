package main

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/lipgloss"
)

// Lipgloss styles
var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6961"))
var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#77DD77"))

// Server Setup Stuff
var serverPort string = ":8080"

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("staticFiles/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", getPage)

	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		fmt.Println(errorStyle.Render("Error") + " Could Not Start Server")
	} else {
		fmt.Println(successStyle.Render("Success") + " Server has started and is running on localhost:" + serverPort)
	}
}

func getPage(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "GET" {
		http.ServeFile(w, r, "webFiles/index.html")
	} else if method == "POST" {
		//Run the Test
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>Chicken Nuggets</h1>")
	} else {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	}
}

func launchTest(writer http.ResponseWriter, request *http.Request) bool {

}
