package main

import (
	"net/http"
	"text/template"
)

// Server Setup Stuff
var serverPort string = ":8080"

// Struct for rendering the report template
type ReportStatistics struct {
	RequestsSent    int
	RequestsBlocked int
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("staticFiles/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", getPage)

	http.ListenAndServe(serverPort, nil)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "GET" {
		http.ServeFile(w, r, "webFiles/index.html")
	} else if method == "POST" {
		//Get the template file ready
		t, _ := template.ParseFiles("webFiles/templates/reportsStatsPage.html")

		//TODO Run the test

		//Fill the data in
		stats := ReportStatistics{100, 50}

		// Run the template
		t.Execute(w, stats)
	} else {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	}
}
