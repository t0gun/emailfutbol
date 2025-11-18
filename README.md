> A small CLI tool that emails only the football fixtures I care about

This tool fetches football fixtures from [API-Football](https://www.api-football.com/) and  filters by league or team IDs.it uses basic smtp setup to email results.i use it with cron locally on my desktop mac. you can have a look if it suits your needs


you can test run using
```bash
go run main.go -config=config.toml
