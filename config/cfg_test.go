package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	validConfig = `
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
	invalidConfig = "invalid config"

	want = &Config{
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
)

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		data        string
		want        *Config
		expectError bool
	}{
		"valid config": {data: validConfig, want: want, expectError: false},
		"No file":      {data: "", want: want, expectError: true},
		"Invalid toml": {data: invalidConfig, want: want, expectError: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tempfile := filepath.Join(t.TempDir(), "config.toml")

			if tc.data != "" {
				if err := os.WriteFile(tempfile, []byte(tc.data), 0644); err != nil {
					t.Fatalf("cannot write to file: %v", err)
				}
			}

			got, err := LoadConfig(tempfile)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect error, got: %v", err)
				}

				if !reflect.DeepEqual(*got, *tc.want) {
					t.Fatalf("got %v want %v", got, tc.want)
				}

			}

		})
	}
}
