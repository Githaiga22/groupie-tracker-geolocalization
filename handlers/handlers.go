package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "tracker/models"
	"tracker/src"
)

var (
	AllArtistInfo      []model.Data
	fetchDatesFunc     = src.FetchDates
	fetchLocationsFunc = src.FetchLocations
)

func DateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dates" {
		notFoundHandler(w)
		return
	}

	if r.Method != http.MethodGet {
		wrongMethodHandler(w)
		return
	}

	id := r.FormValue("id")

	idNum, _ := strconv.Atoi(id)

	if idNum <= 0 || idNum > 52 {
		badRequestHandler(w)
		return
	}

	dates, err := fetchDatesFunc(id)
	if err != nil {
		InternalServerHandler(w)
		log.Println(err)
		return
	}

	// Check if the handler is running in "test mode" to skip template rendering
	if os.Getenv("TEST_MODE") == "true" {
		// If we're in test mode, return a simple mock response instead of rendering a template
		fmt.Fprintln(w, "Mocked template rendering with dates:", dates)
		return
	}
	tmpl, err := template.ParseFiles("templates/dates.html")
	if err != nil {
		InternalServerHandler(w)
		log.Println("Template 2 parsing error: ", err)
		return

	}
	err = tmpl.Execute(w, dates)
	if err != nil {
		log.Println("Template 2 execution error: ", err)
		return
	}
}

func LocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/locations" {
		notFoundHandler(w)
		return
	}

	if r.Method != http.MethodGet {
		wrongMethodHandler(w)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		badRequestHandler(w)
		return
	}

	idNum, _ := strconv.Atoi(id)

	if idNum <= 0 || idNum > 52 {
		badRequestHandler(w)
		return
	}

	locations, err := fetchLocationsFunc(id)
	if err != nil {
		InternalServerHandler(w)
		log.Println(err)
		return
	}

	// Check if the handler is running in "test mode" to skip template rendering
	if os.Getenv("TEST_MODE") == "true" {
		// If we're in test mode, return a simple mock response instead of rendering a template
		fmt.Fprintln(w, "Mocked template rendering with dates:", locations)
		return
	}

	tmpl, err := template.ParseFiles("templates/locations.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template 2 parsing error: ", err)
		return

	}

	type tempLocations struct {
		Name      string
		Locations []string
	}

	// create a map[string]string{locationName:[lat;long]}
	locationMap := []tempLocations{}
	if len(locations.Locations) > 0 {
		for _, locationNames := range locations.Locations {
			latLong, err := src.FetchLocationMap(locationNames)
			if err != nil {
				InternalServerHandler(w)
				log.Println("Error retrieving coordinates: ", err)
				return
			}
			splitCoordinates := strings.Split(latLong, " ")
			// println(locationNames,splitCoordinates[0], "and", splitCoordinates[1])
			newLocation := tempLocations{locationNames, splitCoordinates}

			// println(newLocation.Name, newLocation.Locations[0], newLocation.Locations[1])

			locationMap = append(locationMap, newLocation)
			// locationMap[locationNames] = []string{splitCoordinates[0],splitCoordinates[1]}
		}
	}

	locationMapJSON, err := json.Marshal(locationMap)
	if err != nil {
		log.Println("Error serializing LocationMap:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"LocationMap": string(locationMapJSON),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Template 2 execution error: ", err)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artist" {
		notFoundHandler(w)
		return
	}

	if r.Method != http.MethodGet {
		wrongMethodHandler(w)
		return
	}

	id := r.URL.Query().Get("id")

	datesAndConcerts, err := src.FetchDatesAndConcerts(id)
	if err != nil {
		InternalServerHandler(w)
		log.Println(err)
		return
	}

	idNum, _ := strconv.Atoi(id)
	if idNum <= 0 || idNum > 52 {
		badRequestHandler(w)
		return
	}
	idNum -= 1

	if len(AllArtistInfo) == 0 {
		r.URL.Path = "/"
		r.Method = http.MethodGet
		HomepageHandler(w, r)
		return
	}

	AllArtistInfo[idNum].DateAndLocation = datesAndConcerts

	Data := AllArtistInfo[idNum]

	// fetch artists details
	tmpl, err := template.ParseFiles("templates/artistPage.html")
	if err != nil {
		InternalServerHandler(w)
		log.Println("Template 2 parsing error: ", err)
		return

	}
	err = tmpl.Execute(w, Data)
	if err != nil {
		log.Println("Template 2 execution error: ", err)
		return
	}
}

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w)
		return
	}

	if r.Method != http.MethodGet {
		wrongMethodHandler(w)
		return
	}

	if len(AllArtistInfo) == 0 {

		artists, err := src.FetchArtists()
		if err != nil {
			InternalServerHandler(w)
			log.Println(err)
			return
		}

		for _, artistsInfo := range artists {
			var tempdate model.Data
			tempdate.Name = artistsInfo.Name
			tempdate.Id = artistsInfo.Id
			tempdate.FirstAlbum = artistsInfo.FirstAlbum
			tempdate.CreationDate = artistsInfo.CreationDate
			tempdate.Image = artistsInfo.Image
			tempdate.Members = artistsInfo.Members
			AllArtistInfo = append(AllArtistInfo, tempdate)
		}
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Println("Template 1 parsing error:", err)
			InternalServerHandler(w)
			return
		}

		err = tmpl.Execute(w, AllArtistInfo)
		if err != nil {
			if err != http.ErrHandlerTimeout {
				InternalServerHandler(w)
				log.Println("Template 1 execution error: ", err)
			}
		}
	}
}

func renderErrorPage(w http.ResponseWriter, statusCode int, title, message string) {
	w.WriteHeader(statusCode)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Println("Error page parsing error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}
	data := struct {
		Title   string
		Message string
	}{
		Title:   title,
		Message: message,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error: page execution:", err)
		http.Error(w, "Internal SErver Error", http.StatusInternalServerError)
	}
}

// Response represents the API key response structure
type Response struct {
	APIKey string `json:"apiKey"`
}

// GetApiKey handles the API key retrieval endpoint
func GetApiKey(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("HEREAPI_KEY")
	if apiKey == "" {
		// Set JSON content type even for errors
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "API key not set",
		})
		return
	}

	// Set content type before writing response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{APIKey: apiKey})
}

func notFoundHandler(w http.ResponseWriter) {
	renderErrorPage(w, http.StatusNotFound, "404 Not Found", "The page you are looking for does not exist.")
}

func wrongMethodHandler(w http.ResponseWriter) {
	renderErrorPage(w, http.StatusMethodNotAllowed, " Method Not Allowed", "Try  the home page")
}

func InternalServerHandler(w http.ResponseWriter) {
	renderErrorPage(w, http.StatusInternalServerError, " Internal Server Error", "Completely our mistake.")
}

func badRequestHandler(w http.ResponseWriter) {
	renderErrorPage(w, http.StatusBadRequest, " Bad Request Error", " Try the home page")
}
