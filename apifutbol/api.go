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

type APIClient struct {
	Client   *http.Client
	Timezone string
	Apikey   string
	Date     string
	BaseUrl  string
}

func (api *APIClient) GetFixtures() ([]*FixturesResponse, error) {
	fullUrl := api.buildFixturesUrl()

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")
	req.Header.Add("x-rapidapi-key", api.Apikey)
	req.Header.Add("Accept", "application/json")
	resp, err := api.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	defer resp.Body.Close()
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

func (api *APIClient) buildFixturesUrl() string {
	params := url.Values{}
	params.Add("date", api.Date)
	params.Add("timezone", api.Timezone)
	return fmt.Sprintf("%s/fixtures?%s", api.BaseUrl, params.Encode())
}

func NewAPIClient(baseURL, apikey, date, timezone string, client *http.Client) *APIClient {
	return &APIClient{
		BaseUrl:  baseURL,
		Apikey:   apikey,
		Client:   client,
		Timezone: timezone,
		Date:     date,
	}
}
