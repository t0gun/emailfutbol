package apifutbol

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIClient_GetFixtures(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	mockResponse := APIResponse{
		Response: []*FixturesResponse{
			{
				Fixture: Fixture{Id: 123, Timezone: "Europe/London", Date: "2023-10-10"},
				League:  League{Name: "Premier League"},
				Teams: Teams{
					Home: Home{Name: "Team A"},
					Away: Away{Name: "Team B"},
				},
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/fixtures", r.URL.Path)
		assert.Equal("today", r.URL.Query().Get("date"))
		assert.Equal("Europe/London", r.URL.Query().Get("timezone"))
		assert.Equal("testkey", r.Header.Get("x-rapidapi-key"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := NewAPIClient(ts.URL, "testkey", "today", "Europe/London", ts.Client())
	fixtures, err := client.GetFixtures()
	require.NoError(err)
	assert.Len(fixtures, 1)
	assert.Equal(123, fixtures[0].Fixture.Id)
	assert.Equal("Premier League", fixtures[0].League.Name)
	assert.Equal("Team A", fixtures[0].Teams.Home.Name)
}

func TestAPIClient_GetFixtures_BadUrl(t *testing.T) {
	client := &http.Client{}
	baseUrl := "http://example.com/foo%zz"
	api := NewAPIClient(baseUrl, "testkey", "today", "Europe/London", client)
	_, err := api.GetFixtures()
	require.Error(t, err)

}

func FuzzAPIClient_GetFixtures(f *testing.F) {
	f.Add("http://example.com")
	f.Add("http://example.com/%zz")
	f.Add("://broken-url")

	client := &http.Client{}

	f.Fuzz(func(t *testing.T, url string) {
		api := NewAPIClient(url, "key", "today", "Europe/London", client)
		_, _ = api.GetFixtures()
	})
}
