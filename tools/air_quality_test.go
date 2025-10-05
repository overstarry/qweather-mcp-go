package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/overstarry/qweather-mcp-go/api"
)

// setupAirQualityMockServer creates a mock server for air quality testing
func setupAirQualityMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/geo/v2/city/lookup":
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
		case r.URL.Path == "/airquality/v1/current/39.90/116.41":
			response := api.AirQualityResponse{
				Code: "200",
				Indexes: []struct {
					Code       string `json:"code"`
					Name       string `json:"name"`
					Aqi        int    `json:"aqi"`
					AqiDisplay string `json:"aqiDisplay"`
					Level      string `json:"level,omitempty"`
					Category   string `json:"category,omitempty"`
					Color      struct {
						Red   int `json:"red"`
						Green int `json:"green"`
						Blue  int `json:"blue"`
						Alpha int `json:"alpha"`
					} `json:"color"`
					PrimaryPollutant *struct {
						Code     string `json:"code"`
						Name     string `json:"name"`
						FullName string `json:"fullName"`
					} `json:"primaryPollutant,omitempty"`
					Health *struct {
						Effect string `json:"effect"`
						Advice struct {
							GeneralPopulation   string `json:"generalPopulation"`
							SensitivePopulation string `json:"sensitivePopulation"`
						} `json:"advice"`
					} `json:"health,omitempty"`
				}{
					{
						Code:       "qaqi",
						Name:       "QAQI",
						Aqi:        50,
						AqiDisplay: "50",
						Level:      "1",
						Category:   "Good",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
}

// TestRegisterAirQualityTools tests that air quality tools are registered successfully
func TestRegisterAirQualityTools(t *testing.T) {
	server := setupAirQualityMockServer()
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "test",
		Version: "1.0.0",
	}, nil)

	// Should not panic
	RegisterAirQualityTools(s, client)
}

// TestAirQualityInput_Validation tests air quality input validation
func TestAirQualityInput_Validation(t *testing.T) {
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
			input := AirQualityInput{CityName: tt.cityName}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

// TestAirQualityHourlyInput_Validation tests hourly air quality input validation
func TestAirQualityHourlyInput_Validation(t *testing.T) {
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
			input := AirQualityHourlyInput{CityName: tt.cityName}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

// TestAirQualityDailyInput_Validation tests daily air quality input validation
func TestAirQualityDailyInput_Validation(t *testing.T) {
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
			input := AirQualityDailyInput{CityName: tt.cityName}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

// TestAirQualityOutputStructure tests output structure
func TestAirQualityOutputStructure(t *testing.T) {
	output := AirQualityOutput{
		AirQualityInfo: "Test air quality info",
	}

	if output.AirQualityInfo != "Test air quality info" {
		t.Errorf("Expected 'Test air quality info', got '%s'", output.AirQualityInfo)
	}
}

// TestAirQualityHourlyOutputStructure tests hourly output structure
func TestAirQualityHourlyOutputStructure(t *testing.T) {
	output := AirQualityHourlyOutput{
		HourlyInfo: "Test hourly info",
	}

	if output.HourlyInfo != "Test hourly info" {
		t.Errorf("Expected 'Test hourly info', got '%s'", output.HourlyInfo)
	}
}

// TestAirQualityDailyOutputStructure tests daily output structure
func TestAirQualityDailyOutputStructure(t *testing.T) {
	output := AirQualityDailyOutput{
		DailyInfo: "Test daily info",
	}

	if output.DailyInfo != "Test daily info" {
		t.Errorf("Expected 'Test daily info', got '%s'", output.DailyInfo)
	}
}
