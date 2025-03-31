package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tempfile := filepath.Join(t.TempDir(), "config.toml")
	data := `
[smtp]
server = "smtp.icloud.com"
port = 587
email = "apprentice_py@icloud.com"
password = "12345"


[api]
apikey = "12345"
leagues = [11, 22, 33]
teams = [11, 22, 33]
timezone = "Europe/London"
`
	if err := os.WriteFile(tempfile, []byte(data), 0644); err != nil {
		t.Fatalf("cannot write to file: %v", err)
	}
	want := &Config{
		Api: Apifutbol{
			Apikey:   "12345",
			Teams:    []int{11, 22, 33},
			Leagues:  []int{11, 22, 33},
			Timezone: "Europe/London",
		},
		Smtp: Smtp{
			Server:   "smtp.icloud.com",
			Port:     587,
			Email:    "apprentice_py@icloud.com",
			Password: "12345",
		},
	}
	got, err := LoadConfig(tempfile)
	if err != nil {
		t.Fatalf("Cannot load config: %v", err)
	}

	if !reflect.DeepEqual(*got, *want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
