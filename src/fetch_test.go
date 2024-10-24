package src

import (
	"testing"
)

// TestFetchArtists tests the FetchArtists function.
func TestFetchArtists(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Fetch artists successfully",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchArtists()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchArtists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("FetchArtists() returned no artists")
			}
		})
	}
}

// TestFetchLocations tests the FetchLocations function.
func TestFetchLocations(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "Fetch locations for artist ID 1",
			id:      "1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchLocations(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ArtistId == 0 {
				t.Errorf("FetchLocations() returned empty location for ID %s", tt.id)
			}
		})
	}
}

// TestFetchDates tests the FetchDates function.
func TestFetchDates(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "Fetch dates for artist ID 1",
			id:      "1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchDates(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Id == 0 {
				t.Errorf("FetchDates() returned empty dates for ID %s", tt.id)
			}
		})
	}
}

// TestFetchDatesAndConcerts tests the FetchDatesAndConcerts function.
func TestFetchDatesAndConcerts(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "Fetch dates and concerts for artist ID 1",
			id:      "1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchDatesAndConcerts(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDatesAndConcerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("FetchDatesAndConcerts() returned no concert dates for ID %s", tt.id)
			}
		})
	}
}
