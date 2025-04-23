package main

import (
	"flag"
	"fmt"
	"github.com/t0gun/emailfutbal/apifutbol"
	"github.com/t0gun/emailfutbal/config"
	"github.com/t0gun/emailfutbal/filter"
	"github.com/t0gun/emailfutbal/mailer"
	"log"
	"net/http"
	"strings"
)

func main() {
	configPath := flag.String("config", "config.toml", "Path to the config file")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error Loading config:%s", err)
	}

	date, err := config.GetTomorrowsDate(cfg.Api.Timezone)
	if err != nil {
		log.Fatalf("Error getting date:%s", err)
	}

	client := &http.Client{}
	api := apifutbol.NewAPIClient(cfg.Api.Apiurl, cfg.Api.Apikey, date, cfg.Api.Timezone, client)
	mc := filter.NewMatchCollector()
	fixtures, err := api.GetFixtures()

	if err != nil {
		log.Fatalf("cant get fixtures: %d", err)
	}

	filter.SelectFixtureByTeams(cfg, fixtures, mc)
	filter.SelectFixtureByLeagues(cfg, fixtures, mc)

	if len(mc.Results()) == 0 {
		log.Println("No matches found for tomorrow. No email sent.")
		return
	}

	if err = sendMail(mc, cfg); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	fmt.Println("Email sent successfully")

}

func sendMail(mc *filter.MatchCollector, cfg *config.Config) error {

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
