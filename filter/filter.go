package filter

import (
	"fmt"
	"github.com/t0gun/emailfutbal/apifutbol"
	"github.com/t0gun/emailfutbal/config"
	"log"
	"time"
)

type MatchCollector struct {
	selection []string
	seen      map[string]bool
}

func NewMatchCollector() *MatchCollector {
	return &MatchCollector{
		selection: []string{},
		seen:      make(map[string]bool),
	}
}

func (m *MatchCollector) addMatches(fixture *apifutbol.FixturesResponse) {
	layout := time.RFC3339
	parsedTime, err := time.Parse(layout, fixture.Fixture.Date)
	if err != nil {
		log.Printf("error parsing time: %v", err)
		return
	}

	formattedTime := parsedTime.Format("Mon 02 Jan, 3:04 PM")
	match := fmt.Sprintf("%s vs %s - %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name, formattedTime)

	if !m.seen[match] {
		m.selection = append(m.selection, match)
		m.seen[match] = true
	}
}

func (m *MatchCollector) Results() []string {
	return m.selection
}

func SelectFixtureByTeams(cfg *config.Config, fixtures []*apifutbol.FixturesResponse, mc *MatchCollector) {
	teams := make(map[int]struct{})
	for _, team := range cfg.Api.Teams {
		teams[team] = struct{}{}
	}
	for _, fixture := range fixtures {
		if _, ok := teams[fixture.Teams.Home.Id]; ok {
			mc.addMatches(fixture)
		}
		if _, ok := teams[fixture.Teams.Away.Id]; ok {
			mc.addMatches(fixture)
		}

	}
}

func SelectFixtureByLeagues(cfg *config.Config, fixtures []*apifutbol.FixturesResponse, mc *MatchCollector) {
	leagues := make(map[int]struct{})
	for _, leagueID := range cfg.Api.Leagues {
		leagues[leagueID] = struct{}{}
	}
	for _, fixture := range fixtures {
		if _, ok := leagues[fixture.League.Id]; ok {
			mc.addMatches(fixture)
		}

	}
}
