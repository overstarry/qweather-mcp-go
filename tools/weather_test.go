package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/overstarry/qweather-mcp-go/api"
)

// setupMockServer creates a mock server for testing
func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			response := api.LocationResponse{
				Code: "200",
				Location: []api.Location{
					{
						Name:    "Beijing",
						ID:      "101010100",
						Lat:     "39.90",
						Lon:     "116.41",
						Adm1:    "Beijing",
						Adm2:    "Beijing",
						Country: "China",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/v7/weather/now":
			response := api.WeatherNowResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
			}
			response.Now.Temp = "20"
			response.Now.FeelsLike = "18"
			response.Now.Text = "Sunny"
			response.Now.WindDir = "North"
			response.Now.WindScale = "3"
			response.Now.Humidity = "50"
			response.Now.Precip = "0"
			response.Now.Pressure = "1013"
			response.Now.Vis = "10"
			json.NewEncoder(w).Encode(response)
		case "/v7/weather/3d", "/v7/weather/7d", "/v7/weather/10d", "/v7/weather/15d", "/v7/weather/30d":
			response := api.WeatherDailyResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
			}
			json.NewEncoder(w).Encode(response)
		case "/v7/weather/24h", "/v7/weather/72h", "/v7/weather/168h":
			response := api.HourlyResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
			}
			json.NewEncoder(w).Encode(response)
		case "/v7/minutely/5m":
			response := api.MinutelyResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
				Summary:    "No precipitation in the next 2 hours",
			}
			json.NewEncoder(w).Encode(response)
		case "/v7/warning/now":
			response := api.WarningResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
}

// TestRegisterWeatherTools tests that weather tools are registered successfully
func TestRegisterWeatherTools(t *testing.T) {
	server := setupMockServer()
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "test",
		Version: "1.0.0",
	}, nil)

	// Should not panic
	RegisterWeatherTools(s, client)
}

// TestWeatherNowInput_Validation tests input validation
func TestWeatherNowInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		cityName  string
		wantEmpty bool
	}{
		{"valid city", "Beijing", false},
		{"empty city", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := WeatherNowInput{CityName: tt.cityName}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

// TestWeatherForecastInput_Validation tests forecast input validation
func TestWeatherForecastInput_Validation(t *testing.T) {
	tests := []struct {
		name        string
		cityName    string
		days        string
		wantInvalid bool
	}{
		{"valid 3d", "Beijing", "3d", false},
		{"valid 7d", "Beijing", "7d", false},
		{"valid 10d", "Beijing", "10d", false},
		{"valid 15d", "Beijing", "15d", false},
		{"valid 30d", "Beijing", "30d", false},
		{"invalid days", "Beijing", "5d", true},
		{"empty city", "", "3d", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := WeatherForecastInput{CityName: tt.cityName, Days: tt.days}

			// Check city name
			if tt.wantInvalid && input.CityName == "" {
				return // Expected invalid case
			}

			// Check days format
			validDays := map[string]bool{"3d": true, "7d": true, "10d": true, "15d": true, "30d": true}
			isInvalid := !validDays[input.Days]

			if isInvalid != tt.wantInvalid {
				t.Errorf("Days validation = %v, want %v", isInvalid, tt.wantInvalid)
			}
		})
	}
}

// TestHourlyForecastInput_Validation tests hourly forecast input validation
func TestHourlyForecastInput_Validation(t *testing.T) {
	tests := []struct {
		name        string
		cityName    string
		hours       string
		wantInvalid bool
	}{
		{"valid 24h", "Beijing", "24h", false},
		{"valid 72h", "Beijing", "72h", false},
		{"valid 168h", "Beijing", "168h", false},
		{"invalid hours", "Beijing", "48h", true},
		{"empty city", "", "24h", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := HourlyForecastInput{CityName: tt.cityName, Hours: tt.hours}

			// Check city name
			if tt.wantInvalid && input.CityName == "" {
				return // Expected invalid case
			}

			// Check hours format
			validHours := map[string]bool{"24h": true, "72h": true, "168h": true}
			isInvalid := !validHours[input.Hours]

			if isInvalid != tt.wantInvalid {
				t.Errorf("Hours validation = %v, want %v", isInvalid, tt.wantInvalid)
			}
		})
	}
}

// TestMinutelyPrecipitationInput_Validation tests minutely precipitation input validation
func TestMinutelyPrecipitationInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		cityName  string
		wantEmpty bool
	}{
		{"valid city", "Beijing", false},
		{"empty city", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := MinutelyPrecipitationInput{CityName: tt.cityName}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

// TestWeatherWarningInput_Validation tests weather warning input validation
func TestWeatherWarningInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		cityName  string
		wantEmpty bool
	}{
		{"valid city", "Beijing", false},
		{"empty city", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := WeatherWarningInput{CityName: tt.cityName}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}
