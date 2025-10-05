package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	client := NewClient("https://api.example.com", "test-key")

	if client.BaseURL != "https://api.example.com" {
		t.Errorf("Expected BaseURL to be 'https://api.example.com', got '%s'", client.BaseURL)
	}

	if client.APIKey != "test-key" {
		t.Errorf("Expected APIKey to be 'test-key', got '%s'", client.APIKey)
	}

	if client.LogLevel != LogLevelError {
		t.Errorf("Expected default LogLevel to be LogLevelError, got %v", client.LogLevel)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

// TestSetLogLevel tests log level setting
func TestSetLogLevel(t *testing.T) {
	client := NewClient("https://api.example.com", "test-key")

	client.SetLogLevel(LogLevelDebug)
	if client.LogLevel != LogLevelDebug {
		t.Errorf("Expected LogLevel to be LogLevelDebug, got %v", client.LogLevel)
	}
}

// TestMakeRequest tests basic API request
func TestMakeRequest(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check API key header
		if r.Header.Get("X-QW-Api-Key") != "test-key" {
			t.Error("API key header not set correctly")
		}

		// Check query parameters
		if r.URL.Query().Get("location") != "beijing" {
			t.Error("Query parameter not set correctly")
		}

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"200","data":"test"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	params := map[string]string{
		"location": "beijing",
	}

	body, err := client.MakeRequest("/test", params)
	if err != nil {
		t.Fatalf("MakeRequest failed: %v", err)
	}

	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["code"] != "200" {
		t.Errorf("Expected code '200', got '%s'", response["code"])
	}
}

// TestMakeRequestWithContext tests request with context
func TestMakeRequestWithContext(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"200"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.MakeRequestWithContext(ctx, "/test", map[string]string{})
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// TestMakeRequestHTTPError tests HTTP error handling
func TestMakeRequestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	_, err := client.MakeRequest("/test", map[string]string{})
	if err == nil {
		t.Error("Expected error for non-200 status code")
	}
}

// TestGetLocationByName tests city location lookup
func TestGetLocationByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := LocationResponse{
			Code: APICodeSuccess,
			Location: []Location{
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
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	locationData, err := client.GetLocationByName("Beijing")
	if err != nil {
		t.Fatalf("GetLocationByName failed: %v", err)
	}

	if locationData.Code != APICodeSuccess {
		t.Errorf("Expected code '%s', got '%s'", APICodeSuccess, locationData.Code)
	}

	if len(locationData.Location) == 0 {
		t.Error("Expected at least one location")
	}

	if locationData.Location[0].Name != "Beijing" {
		t.Errorf("Expected city name 'Beijing', got '%s'", locationData.Location[0].Name)
	}
}

// TestGetCityCoordinates tests city coordinates helper function
func TestGetCityCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := LocationResponse{
			Code: APICodeSuccess,
			Location: []Location{
				{
					Name: "Shanghai",
					ID:   "101020100",
					Lat:  "31.2304",
					Lon:  "121.4737",
					Adm1: "Shanghai",
					Adm2: "Shanghai",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	lat, lon, cityInfo, err := client.GetCityCoordinates("Shanghai")
	if err != nil {
		t.Fatalf("GetCityCoordinates failed: %v", err)
	}

	if lat != "31.23" {
		t.Errorf("Expected lat '31.23', got '%s'", lat)
	}

	if lon != "121.47" {
		t.Errorf("Expected lon '121.47', got '%s'", lon)
	}

	if cityInfo.Name != "Shanghai" {
		t.Errorf("Expected city name 'Shanghai', got '%s'", cityInfo.Name)
	}
}

// TestGetWeatherNow tests real-time weather query
func TestGetWeatherNow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := WeatherNowResponse{
			Code:       APICodeSuccess,
			UpdateTime: "2024-01-01T12:00+08:00",
		}
		response.Now.Temp = "20"
		response.Now.FeelsLike = "18"
		response.Now.Text = "Sunny"
		response.Now.WindDir = "North"
		response.Now.WindScale = "3"
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	weatherData, err := client.GetWeatherNow("101010100")
	if err != nil {
		t.Fatalf("GetWeatherNow failed: %v", err)
	}

	if weatherData.Code != APICodeSuccess {
		t.Errorf("Expected code '%s', got '%s'", APICodeSuccess, weatherData.Code)
	}

	if weatherData.Now.Temp != "20" {
		t.Errorf("Expected temp '20', got '%s'", weatherData.Now.Temp)
	}
}

// TestGetWeatherForecast tests weather forecast query
func TestGetWeatherForecast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := WeatherDailyResponse{
			Code:       APICodeSuccess,
			UpdateTime: "2024-01-01T12:00+08:00",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	forecastData, err := client.GetWeatherForecast("101010100", "3d")
	if err != nil {
		t.Fatalf("GetWeatherForecast failed: %v", err)
	}

	if forecastData.Code != APICodeSuccess {
		t.Errorf("Expected code '%s', got '%s'", APICodeSuccess, forecastData.Code)
	}
}

// TestGetAirQuality tests air quality query
func TestGetAirQuality(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AirQualityResponse{
			Code: APICodeSuccess,
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
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	airQualityData, err := client.GetAirQuality("39.90", "116.41")
	if err != nil {
		t.Fatalf("GetAirQuality failed: %v", err)
	}

	if airQualityData.Code != APICodeSuccess {
		t.Errorf("Expected code '%s', got '%s'", APICodeSuccess, airQualityData.Code)
	}

	if len(airQualityData.Indexes) == 0 {
		t.Error("Expected at least one air quality index")
	}
}

// TestGetAirQualityEmptyCode tests handling of empty Code field
func TestGetAirQualityEmptyCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Response with empty code but valid data
		response := `{"indexes":[{"code":"qaqi","name":"QAQI","aqi":50,"aqiDisplay":"50"}],"pollutants":[]}`
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	airQualityData, err := client.GetAirQuality("39.90", "116.41")
	if err != nil {
		t.Fatalf("GetAirQuality failed: %v", err)
	}

	// Should set default code to "unknown"
	if airQualityData.Code != APICodeUnknown {
		t.Errorf("Expected code to be set to '%s', got '%s'", APICodeUnknown, airQualityData.Code)
	}

	if len(airQualityData.Indexes) == 0 {
		t.Error("Expected at least one air quality index")
	}
}

// TestMinFunction tests min helper function
func TestMinFunction(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{5, 10, 5},
		{10, 5, 5},
		{5, 5, 5},
		{0, 1, 0},
		{-1, 0, -1},
	}

	for _, tt := range tests {
		result := min(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}
