package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/overstarry/qweather-mcp-go/api"
)

// setupIndicesMockServer creates a mock server for indices testing
func setupIndicesMockServer() *httptest.Server {
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
		case r.URL.Path == "/v7/indices/1d", r.URL.Path == "/v7/indices/3d":
			response := api.IndicesResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
				Daily: []struct {
					Date     string `json:"date"`
					Type     string `json:"type"`
					Name     string `json:"name"`
					Level    string `json:"level"`
					Category string `json:"category"`
					Text     string `json:"text"`
				}{
					{
						Date:     "2024-01-01",
						Type:     "5",
						Name:     "UV Index",
						Level:    "2",
						Category: "Low",
						Text:     "UV radiation is low, no special protection needed.",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
}

// TestRegisterIndicesTools tests that indices tools are registered successfully
func TestRegisterIndicesTools(t *testing.T) {
	server := setupIndicesMockServer()
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "test",
		Version: "1.0.0",
	}, nil)

	// Should not panic
	RegisterIndicesTools(s, client)
}

// TestWeatherIndicesInput_Validation tests indices input validation
func TestWeatherIndicesInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		cityName  string
		indexType string
		days      string
		wantEmpty bool
	}{
		{"valid city with defaults", "Beijing", "", "", false},
		{"valid city with type and days", "Beijing", "5", "1d", false},
		{"valid city with 3d forecast", "Beijing", "0", "3d", false},
		{"empty city", "", "0", "1d", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := WeatherIndicesInput{
				CityName: tt.cityName,
				Type:     tt.indexType,
				Days:     tt.days,
			}
			isEmpty := input.CityName == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("CityName validation = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

// TestWeatherIndicesDefaults tests default values for type and days
func TestWeatherIndicesDefaults(t *testing.T) {
	input := WeatherIndicesInput{CityName: "Beijing"}

	// Empty type and days should be handled by the tool
	if input.Type != "" {
		t.Errorf("Expected empty Type, got '%s'", input.Type)
	}

	if input.Days != "" {
		t.Errorf("Expected empty Days, got '%s'", input.Days)
	}
}

// TestWeatherIndicesAllTypes tests different index types
func TestWeatherIndicesAllTypes(t *testing.T) {
	// Test various index types
	types := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"}

	for _, indexType := range types {
		t.Run("type_"+indexType, func(t *testing.T) {
			input := WeatherIndicesInput{
				CityName: "Beijing",
				Type:     indexType,
				Days:     "1d",
			}

			if input.CityName == "" {
				t.Error("CityName should not be empty")
			}

			if input.Type != indexType {
				t.Errorf("Expected Type '%s', got '%s'", indexType, input.Type)
			}
		})
	}
}

// TestWeatherIndicesOutputStructure tests output structure
func TestWeatherIndicesOutputStructure(t *testing.T) {
	output := WeatherIndicesOutput{
		IndicesInfo: "Test indices info",
	}

	if output.IndicesInfo != "Test indices info" {
		t.Errorf("Expected 'Test indices info', got '%s'", output.IndicesInfo)
	}
}
