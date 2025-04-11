package main

import (
	"fmt"
	"github.com/ogundiyantobiloba/emailfutbal/apifutbol"
	"github.com/ogundiyantobiloba/emailfutbal/config"
	"github.com/ogundiyantobiloba/emailfutbal/mailer"
	"log"
	"net/http"
	"strings"
	"time"
)

type MatchCollector struct {
	selection []string
	seen      map[string]bool
}

func newMatchCollector() *MatchCollector {
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

func selectFixtureByTeams(cfg *config.Config, fixtures []*apifutbol.FixturesResponse, mc *MatchCollector) {
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

func selectFixtureByLeagues(cfg *config.Config, fixtures []*apifutbol.FixturesResponse, mc *MatchCollector) {
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

func sendMail(mc *MatchCollector, cfg *config.Config) error {

	m := mailer.NewMailer(cfg.Smtp.Server, cfg.Smtp.Port, cfg.Smtp.Email, cfg.Smtp.Password)
	m.Recipient = cfg.Smtp.Email // For now personal need
	m.Subject = "Tomorrow's Fixtures"
	body := strings.Join(mc.Results(), "\n\n")
	m.Body = body
	if err := m.SendEmail(); err != nil {
		return err
	}
	return nil
}

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
	mc := newMatchCollector()
	fixtures, err := api.GetFixtures()
	selectFixtureByTeams(cfg, fixtures, mc)
	selectFixtureByLeagues(cfg, fixtures, mc)

	if len(mc.Results()) == 0 {
		log.Println("No matches found for tomorrow. No email sent.")
		return
	}

	if err = sendMail(mc, cfg); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	} else {
		fmt.Println("Email sent successfully")
	}
}
