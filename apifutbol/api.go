package apifutbol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type FixturesResponse struct {
}

func GetFixtures(timezone, apikey, date string) (*http.Response, error) {
	client := &http.Client{}
	params := url.Values{}
	params.Add("timezone", timezone)
	params.Add("from", date)
	params.Add("to", date)
	finalUrl := fmt.Sprintf("https://v3.football.api-sports.io/fixtures?%s", params.Encode())

	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")
	req.Header.Add("x-rapidapi-key", apikey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DecodeResponse(resp *http.Response) (*FixturesResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var fixturesresp *FixturesResponse
	if err = json.Unmarshal(body, &fixturesresp); err != nil {
		return nil, err
	}
	return fixturesresp, nil
}
