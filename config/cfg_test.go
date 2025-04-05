package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
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

	config = &Config{
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
	assert := assert.New(t)
	require := require.New(t)
	tests := map[string]struct {
		data        string
		want        *Config
		expectError bool
	}{
		"valid config": {data: validConfig, want: config, expectError: false},
		"No file":      {data: "", want: config, expectError: true},
		"Invalid toml": {data: invalidConfig, want: config, expectError: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tempfile := filepath.Join(t.TempDir(), "config.toml")

			if tc.data != "" {
				if err := os.WriteFile(tempfile, []byte(tc.data), 0644); err != nil {
					require.Nil(err)
				}
			}

			got, err := LoadConfig(tempfile)

			if tc.expectError {
				require.NotNil(err)
			} else {
				require.Nil(err)
				assert.Equal(*got, *tc.want)
			}

		})
	}
}

func TestGetTomorrowsDate(t *testing.T) {
	t.Run("valid timezone", func(t *testing.T) {
		timezone := "Europe/London"
		loc, _ := time.LoadLocation(timezone)
		want := time.Now().In(loc).AddDate(0, 0, 1).Format("2006-01-02")
		got, _ := GetTomorrowsDate(timezone)
		assert.Equal(t, want, got)
		t.Parallel()
	})

	t.Run("invalid timezone", func(t *testing.T) {
		_, err := GetTomorrowsDate("incorrect/timezone")
		require.NotNil(t, err)
		t.Parallel()
	})
}
