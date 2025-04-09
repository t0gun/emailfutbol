package main

import (
	"fmt"
	"github.com/ogundiyantobiloba/emailfutbal/apifutbol"
	"github.com/ogundiyantobiloba/emailfutbal/config"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Error Loading config:%s", err)
	}

	date, err := config.GetTomorrowsDate(cfg.Api.Timezone)
	if err != nil {
		log.Fatalf("Error getting date:%s", err)
	}

	client := &http.Client{}
	api := apifutbol.NewAPIClient(cfg.Api.Apiurl, cfg.Api.Apikey, date, cfg.Api.Timezone, client)
	fixtures, err := api.GetFixtures()
	selection := selectFixtureByTeams(cfg, fixtures)
	fmt.Println(selection)

}

func selectFixtureByTeams(cfg *config.Config, fixtures []*apifutbol.FixturesResponse) []string {
	var selection []string
	seen := make(map[string]bool)

	for _, team := range cfg.Api.Teams {
		for _, f := range fixtures {
			fixture := f.Teams

			if team == fixture.Home.Id || team == fixture.Away.Id {
				key := fmt.Sprintf("%s-vs-%s-%s", fixture.Home.Id, fixture.Away.Id, f.Fixture.Timestamp)

				if !seen[key] {
					match := fmt.Sprintf("%s vs %s - %s", fixture.Home.Name, fixture.Away.Name, f.Fixture.Date)
					selection = append(selection, match)
					seen[key] = true
				}
			}
		}
	}

	return selection
}

func selectFixtureByLeagues(cfg *config.Config, fixtures []*apifutbol.FixturesResponse) []string {

	return []string{""}
}
