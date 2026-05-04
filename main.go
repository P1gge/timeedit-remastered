package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Structures

type Schedule struct {
	Records []struct {
		ID   int `json:"id"`
		Type struct {
			ID int `json:"id"`
		} `json:"type"`
	} `json:"records"`
}

type Reservation struct {
	StartDate string   `json:"startdate"`
	StartTime string   `json:"starttime"`
	EndDate   string   `json:"enddate"`
	EndTime   string   `json:"endtime"`
	Columns   []string `json:"columns"`
}

type Response struct {
	Reservations []Reservation `json:"reservations"`
}

type Event struct {
	Start time.Time
	End   time.Time

	CourseCodes []string
	CourseNames []string

	Title    string
	Comment  string
	Activity string

	Rooms []string

	ClassCodes []string
	ClassNames []string
}

// Helper Functions

func get(arr []string, i int) string {
	if i < 0 || i >= len(arr) {
		return ""
	}
	return arr[i]
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// Parser Functions

func toEvent(r Reservation) (Event, error) {
	layout := "2006-01-02 15:04"

	start, err := time.Parse(layout, r.StartDate+" "+r.StartTime)
	if err != nil {
		return Event{}, err
	}

	end, err := time.Parse(layout, r.EndDate+" "+r.EndTime)
	if err != nil {
		return Event{}, err
	}

	col := r.Columns

	return Event{
		Start: start,
		End:   end,

		CourseCodes: splitCSV(get(col, 0)),
		CourseNames: splitCSV(get(col, 1)),

		Title:    get(col, 2),
		Comment:  get(col, 3),
		Activity: get(col, 4),

		Rooms: splitCSV(get(col, 5)),

		ClassCodes: splitCSV(get(col, 6)),
		ClassNames: splitCSV(get(col, 7)),
	}, nil
}

// Fetch function

func FetchEventsFromURL(url string) ([]Event, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response %s", resp.Status)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	reservations := response.Reservations

	events := make([]Event, 0, len(reservations))

	for _, r := range reservations {
		e, err := toEvent(r)
		if err != nil {
			continue
		}
		events = append(events, e)
	}

	return events, nil
}

func fetchObjectIds(search string, amount int) ([]string, error) {
	if amount < 0 || amount > 100 {
		return nil, fmt.Errorf("amount outside of range")
	}

	base := "https://cloud.timeedit.net/chalmers/web/public/objects.json?"

	flags := "max=" + strconv.Itoa(amount)
	flags += "&sid=3"
	flags += "&search_text=" + url.QueryEscape(search)
	flags += "&types=10"

	url := base + flags

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response %s", resp.Status)
	}

	var schedule Schedule
	err = json.NewDecoder(resp.Body).Decode(&schedule)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, record := range schedule.Records {
		combined := fmt.Sprintf("%d.%d", record.ID, record.Type.ID)
		ids = append(ids, combined)
	}

	return ids, nil
}

// Main

func main() {
	//url := "https://cloud.timeedit.net/chalmers/web/public/ri667XQ1091Z58Qv3Z0Yb6Z6y4YQ200nQYe1u2gQZ0.json"

	/*
		events, err := FetchEventsFromURL(url)
		if err != nil {
			panic(err)
		}

			for _, e := range events {
				fmt.Println(e.Title, e.Start)
			}
	*/

	ids, err := fetchObjectIds("eda333", 100)
	if err != nil {
		panic(err)
	}
	for _, id := range ids {
		fmt.Println(id)
	}
}
