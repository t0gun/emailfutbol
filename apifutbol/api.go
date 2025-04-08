package apifutbol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	APIResponse struct {
		Response []*FixturesResponse `json:"response"`
	}

	FixturesResponse struct {
		Fixture Fixture `json:"fixture"`
		League  League  `json:"league"`
		Teams   Teams   `json:"teams"`
		Goals   Goals   `json:"goals"`
		Score   Score   `json:"score"`
	}

	Fixture struct {
		Id        int    `json:"id"`
		Timezone  string `json:"timezone"`
		Date      string `json:"date"`
		Timestamp int    `json:"timestamp"`
		Venue     Venue  `json:"venue"`
		Status    Status `json:"status"`
	}

	Venue struct {
		Name string `json:"name"`
		City string `json:"city"`
	}

	Status struct {
		Long    string   `json:"long"`
		Short   string   `json:"short"`
		Elapsed *float64 `json:"elapsed"`
	}

	League struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Season  int    `json:"season"`
		Round   string `json:"round"`
	}

	Teams struct {
		Home Home `json:"home"`
		Away Away `json:"away"`
	}

	Home struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Logo   string `json:"logo"`
		Winner *bool  `json:"winner"`
	}

	Away struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Logo   string `json:"logo"`
		Winner *bool  `json:"winner"`
	}
	Goals struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	}

	Score struct {
		Halftime  Halftime  `json:"halftime"`
		Fulltime  Fulltime  `json:"fulltime"`
		Extratime Extratime `json:"extratime"`
		Penalty   Penalty   `json:"penalty"`
	}

	Halftime struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	}

	Fulltime struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	}

	Extratime struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	}

	Penalty struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	}
)

func GetFixtures(timezone, apikey, date string) ([]*FixturesResponse, error) {
	client := &http.Client{}
	url := buildFixturesUrl(timezone, date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")
	req.Header.Add("x-rapidapi-key", apikey)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var wrapper APIResponse
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.Response, nil

}

func buildFixturesUrl(timezone, date string) string {
	params := url.Values{}
	params.Add("date", date)
	params.Add("timezone", timezone)
	return fmt.Sprintf("https://v3.football.api-sports.io/fixtures?%s", params.Encode())
}
