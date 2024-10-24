package models

////////////////////////////////////////////////////////////
// artists api

type Artist struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relation"`
}

// //////////////////////////////////////////////////////////
// locations and dates api
type RootsRelation struct {
	Relation []DatesLocation `json:"index"`
}

type DatesLocation struct {
	Id     int            `json:"id"`
	Places DatesLocations `json:"datesLocations"`
}

type DatesLocations map[string][]string

// //////////////////////////////////////////////////////////
// dates api
type RootDates struct {
	Tdates []Date `json:"index"`
}
type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

// //////////////////////////////////////////////////////////
// locations api
type AllLocations struct {
	Location []Location `json:"index"` // Match index key in the JSON response
}
type Location struct {
	ArtistId  int      `json:"id"`
	Locations []string `json:"locations"`
	Date      string   `json:"dates"`
}

////////////////////////////////////////////////////////////

type Data struct {
	Name            string
	Id              int
	Image           string
	Members         []string
	CreationDate    int
	FirstAlbum      string
	DateAndLocation DatesLocations
	Dates           []Date
	Locations       []Location
}
