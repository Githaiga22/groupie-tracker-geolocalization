package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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
		InternalServerHandler(w)
		log.Println("Template 2 parsing error: ", err)
		return

	}
	err = tmpl.Execute(w, locations)
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
		log.Println("here here")
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
	tmpl,err := template.ParseFiles("templates/error.html")
	if err != nil{
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
		InternalServerHandler(w)
	}
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
