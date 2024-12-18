package src

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "strings"

    model "tracker/models"
)

var Data model.Data

func FetchArtists() ([]model.Artist, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var artists []model.Artist
    if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
        return nil, err
    }
    return artists, nil
}

func FetchLocations(id string) (model.Location, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
    if err != nil {
        fmt.Println("Error reading the response body:", err)
        return model.Location{}, err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading the response body:", err)
        return model.Location{}, err
    }
    // Unmarshal the JSON data into Go structs
    var data model.AllLocations
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return model.Location{}, err
    }

    Data.Locations = data.Location
    var locations model.Location

    for _, Artistid := range data.Location {
        idNum, _ := strconv.Atoi(id)
        if Artistid.ArtistId == idNum {
            locations = Artistid
            break
        }
    }
    return locations, nil
}

func FetchDates(id string) (model.Date, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
    if err != nil {
        fmt.Println("Error reading the response body:", err)
        return model.Date{}, err
    }
    defer resp.Body.Close()

    var data model.RootDates
    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading the response body:", err)
        return model.Date{}, err
    }
    // Unmarshal the JSON data into Go structs
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return model.Date{}, err
    }

    var dates model.Date

    Data.Dates = data.Tdates

    for _, Artistid := range data.Tdates {
        idNum := strconv.Itoa(Artistid.Id)
        if idNum == id {
            dates = Artistid
        }
    }

    for i, date := range dates.Dates {
        if date[0] == '*' {
            dates.Dates[i] = date[1:]
        }
    }

    return dates, nil
}

func FetchDatesAndConcerts(id string) (model.DatesLocations, error) {
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
    if err != nil {
        fmt.Println("Error reading the response body:", err)
        return nil, err
    }
    defer resp.Body.Close()

    var data model.RootsRelation

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading the response body:", err)
        return nil, err
    }

    // Unmarshal the JSON data into Go structs
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return nil, err
    }

    var datesLocations model.DatesLocations

    for _, Artistid := range data.Relation {
        idNum := strconv.Itoa(Artistid.Id)
        if idNum == id {
            datesLocations = Artistid.Places
        }
    }

    return datesLocations, nil
}

func FetchLocationMap(name string) (string, error) {
    apiKey := os.Getenv("HEREAPI_KEY")

    if strings.Contains(name, "uk") {
        name = name[:len(name)-2] // Slice off the last 2 characters
        name += "_england-united-kingdom"
    }
    
    // Autocomplete request to get place ID
    autoCompleteURL := fmt.Sprintf("https://autocomplete.search.hereapi.com/v1/autocomplete?q=%s&limit=1&apiKey=%s", name, apiKey)
    resp, err := http.Get(autoCompleteURL)
    if err != nil {
        fmt.Println("Error making the autocomplete request:", err)
        return "", err
    }
    defer resp.Body.Close()

    var autoResp model.AutocompleteResponse
    body, _ := io.ReadAll(resp.Body)
    err = json.Unmarshal(body, &autoResp)
    if err != nil || len(autoResp.Items) == 0 {
        fmt.Println("Error parsing autocomplete response:", err)
        return "", err
    }
    placeID := autoResp.Items[0].ID

    

    // Geocoding request to get latitude and longitude
    geocodeURL := fmt.Sprintf("https://lookup.search.hereapi.com/v1/lookup?id=%s&apiKey=%s", placeID, apiKey)
    resp, err = http.Get(geocodeURL)
    if err != nil {
        fmt.Println("Error making the geocode request:", err)
        return "", err
    }
    defer resp.Body.Close()

    var geoResp model.GeocodeResponse
    body, _ = io.ReadAll(resp.Body)
    err = json.Unmarshal(body, &geoResp)
    if err != nil {
        fmt.Println("Error parsing geocode response:", err)
        return "", err
    }

    return fmt.Sprintf("%f %f", geoResp.Position.Lat, geoResp.Position.Lng), nil
}

