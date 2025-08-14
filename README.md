# emailfutbol

A small CLI tool that emails only the football fixtures I care about

## What It Does

- Fetches fixtures from [API-Football](https://www.api-football.com/)
- Filters by league or team IDs
- Emails the result using a basic SMTP setup

---

## Why I Built It

I don’t check sports apps or news often, but I check my email every morning.  
So I built this to send only the matches I actually care about, for the next day.

The tool is designed for **local cron-based usage**, and built around my own workflow — but the code is modular enough to extend.

---

## Key Design Notes

- **Personal Use Case**  
  Currently built to send **only the next day’s fixtures** for a defined list of teams/leagues.

- **Local Cron Setup**  
  Designed to run as a local job (via `cron` or `systemd` timer).  
  Not intended to be a hosted service.

- **SMTP-Only (No APIs)**  
  Uses basic SMTP with app passwords.  

- **Modular Structure**
    - API client, filtering, and mail logic are in separate packages
    - Easy to swap components or add CLI flags
    - Config-driven via `.toml`

---



## Example Usage

```bash
go run main.go -config=config.toml
