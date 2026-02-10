package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/overstarry/qweather-mcp-go/api"
)

func TestHandleWeatherNow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/weather/now":
			resp := api.WeatherNowResponse{Code: "200", UpdateTime: "2024-01-01T12:00+08:00"}
			resp.Now.Temp = "20"
			resp.Now.FeelsLike = "18"
			resp.Now.Text = "Sunny"
			resp.Now.WindDir = "North"
			resp.Now.WindScale = "3"
			resp.Now.Humidity = "50"
			resp.Now.Precip = "0"
			resp.Now.Pressure = "1013"
			resp.Now.Vis = "10"
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleWeatherNow(client, WeatherNowInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleWeatherNow failed: %v", err)
	}
	if !strings.Contains(out.WeatherInfo, "Current Weather - Beijing") {
		t.Fatalf("WeatherInfo = %q, want to contain %q", out.WeatherInfo, "Current Weather - Beijing")
	}
	if !strings.Contains(out.WeatherInfo, "Temperature: 20°C") {
		t.Fatalf("WeatherInfo = %q, want to contain %q", out.WeatherInfo, "Temperature: 20°C")
	}
}

func TestHandleWeatherForecast_DefaultDays(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/weather/3d":
			resp := api.WeatherDailyResponse{Code: "200", UpdateTime: "2024-01-01T12:00+08:00"}
			resp.Daily = append(resp.Daily, struct {
				FxDate         string `json:"fxDate"`
				Sunrise        string `json:"sunrise"`
				Sunset         string `json:"sunset"`
				Moonrise       string `json:"moonrise"`
				Moonset        string `json:"moonset"`
				MoonPhase      string `json:"moonPhase"`
				MoonPhaseIcon  string `json:"moonPhaseIcon"`
				TempMax        string `json:"tempMax"`
				TempMin        string `json:"tempMin"`
				IconDay        string `json:"iconDay"`
				TextDay        string `json:"textDay"`
				IconNight      string `json:"iconNight"`
				TextNight      string `json:"textNight"`
				Wind360Day     string `json:"wind360Day"`
				WindDirDay     string `json:"windDirDay"`
				WindScaleDay   string `json:"windScaleDay"`
				WindSpeedDay   string `json:"windSpeedDay"`
				Wind360Night   string `json:"wind360Night"`
				WindDirNight   string `json:"windDirNight"`
				WindScaleNight string `json:"windScaleNight"`
				WindSpeedNight string `json:"windSpeedNight"`
				Humidity       string `json:"humidity"`
				Precip         string `json:"precip"`
				Pressure       string `json:"pressure"`
				Vis            string `json:"vis"`
				Cloud          string `json:"cloud"`
				UvIndex        string `json:"uvIndex"`
			}{
				FxDate: "2024-01-02", TempMin: "10", TempMax: "20", TextDay: "Sunny", TextNight: "Clear", Sunrise: "06:00", Sunset: "18:00",
				Precip: "0", Humidity: "50", WindDirDay: "N", WindScaleDay: "3", WindDirNight: "N", WindScaleNight: "2", UvIndex: "2",
			})
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleWeatherForecast(client, WeatherForecastInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleWeatherForecast failed: %v", err)
	}
	if !strings.Contains(out.ForecastInfo, "3 Day Weather Forecast - Beijing") {
		t.Fatalf("ForecastInfo = %q, want to contain %q", out.ForecastInfo, "3 Day Weather Forecast - Beijing")
	}
}

func TestHandleWeatherForecast_InvalidDays(t *testing.T) {
	client := api.NewClient("http://example.com", "test-key")
	_, err := handleWeatherForecast(client, WeatherForecastInput{CityName: "Beijing", Days: "5d"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHandleMinutelyPrecipitation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/minutely/5m":
			resp := api.MinutelyResponse{
				Code:       "200",
				UpdateTime: "2024-01-01T12:00+08:00",
				FxLink:     "https://example.com",
				Summary:    "No precipitation in the next 2 hours",
			}
			resp.Minutely = append(resp.Minutely, struct {
				FxTime string `json:"fxTime"`
				Precip string `json:"precip"`
				Type   string `json:"type"`
			}{
				FxTime: "2024-01-01T12:05+08:00",
				Precip: "0.2",
				Type:   "rain",
			})
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleMinutelyPrecipitation(client, MinutelyPrecipitationInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleMinutelyPrecipitation failed: %v", err)
	}
	if !strings.Contains(out.PrecipitationInfo, "2-Hour Precipitation Forecast:") {
		t.Fatalf("PrecipitationInfo = %q, want to contain %q", out.PrecipitationInfo, "2-Hour Precipitation Forecast:")
	}
	if !strings.Contains(out.PrecipitationInfo, "Data Source: https://example.com") {
		t.Fatalf("PrecipitationInfo = %q, want to contain %q", out.PrecipitationInfo, "Data Source: https://example.com")
	}
}

func TestHandleHourlyForecast_DefaultHours(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/weather/24h":
			resp := api.HourlyResponse{Code: "200", UpdateTime: "2024-01-01T12:00+08:00"}
			resp.Hourly = append(resp.Hourly, struct {
				FxTime    string `json:"fxTime"`
				Temp      string `json:"temp"`
				Icon      string `json:"icon"`
				Text      string `json:"text"`
				Wind360   string `json:"wind360"`
				WindDir   string `json:"windDir"`
				WindScale string `json:"windScale"`
				WindSpeed string `json:"windSpeed"`
				Humidity  string `json:"humidity"`
				Precip    string `json:"precip"`
				Pressure  string `json:"pressure"`
				Cloud     string `json:"cloud"`
				Dew       string `json:"dew"`
			}{
				FxTime: "2024-01-01T13:00+08:00", Temp: "20", Text: "Sunny", WindDir: "North", WindScale: "3", WindSpeed: "10",
				Humidity: "50", Precip: "0", Pressure: "1013", Cloud: "10", Dew: "5",
			})
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleHourlyForecast(client, HourlyForecastInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleHourlyForecast failed: %v", err)
	}
	if !strings.Contains(out.HourlyInfo, "24 Hour Weather Forecast - Beijing") {
		t.Fatalf("HourlyInfo = %q, want to contain %q", out.HourlyInfo, "24 Hour Weather Forecast - Beijing")
	}
	if !strings.Contains(out.HourlyInfo, "Cloud Cover: 10%") {
		t.Fatalf("HourlyInfo = %q, want to contain %q", out.HourlyInfo, "Cloud Cover: 10%")
	}
	if !strings.Contains(out.HourlyInfo, "Dew Point: 5°C") {
		t.Fatalf("HourlyInfo = %q, want to contain %q", out.HourlyInfo, "Dew Point: 5°C")
	}
}

func TestHandleHourlyForecast_InvalidHours(t *testing.T) {
	client := api.NewClient("http://example.com", "test-key")
	_, err := handleHourlyForecast(client, HourlyForecastInput{CityName: "Beijing", Hours: "48h"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHandleWeatherWarning_NoWarnings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/warning/now":
			json.NewEncoder(w).Encode(api.WarningResponse{Code: "200", UpdateTime: "2024-01-01T12:00+08:00"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleWeatherWarning(client, WeatherWarningInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleWeatherWarning failed: %v", err)
	}
	if !strings.Contains(out.WarningInfo, "has no active weather warnings") {
		t.Fatalf("WarningInfo = %q, want to contain %q", out.WarningInfo, "has no active weather warnings")
	}
}

func TestHandleWeatherWarning_WithWarning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/warning/now":
			resp := api.WarningResponse{Code: "200", UpdateTime: "2024-01-01T12:00+08:00"}
			resp.Warning = append(resp.Warning, struct {
				ID            string `json:"id"`
				Sender        string `json:"sender"`
				PubTime       string `json:"pubTime"`
				Title         string `json:"title"`
				StartTime     string `json:"startTime"`
				EndTime       string `json:"endTime"`
				Status        string `json:"status"`
				Severity      string `json:"severity"`
				SeverityColor string `json:"severityColor"`
				Type          string `json:"type"`
				TypeName      string `json:"typeName"`
				Urgency       string `json:"urgency"`
				Certainty     string `json:"certainty"`
				Text          string `json:"text"`
				Related       string `json:"related"`
			}{
				Title: "Test Warning", Sender: "Agency", PubTime: "2024-01-01T12:00+08:00", TypeName: "Storm", Severity: "Severe", SeverityColor: "Red",
				Status: "active", Text: "Stay inside",
			})
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleWeatherWarning(client, WeatherWarningInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleWeatherWarning failed: %v", err)
	}
	if !strings.Contains(out.WarningInfo, "Weather Warnings - Beijing") {
		t.Fatalf("WarningInfo = %q, want to contain %q", out.WarningInfo, "Weather Warnings - Beijing")
	}
	if !strings.Contains(out.WarningInfo, "Warning Title: Test Warning") {
		t.Fatalf("WarningInfo = %q, want to contain %q", out.WarningInfo, "Warning Title: Test Warning")
	}
}

func TestHandleWeatherIndices_Defaults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/v7/indices/1d":
			resp := api.IndicesResponse{Code: "200", UpdateTime: "2024-01-01T12:00+08:00"}
			resp.Daily = append(resp.Daily, struct {
				Date     string `json:"date"`
				Type     string `json:"type"`
				Name     string `json:"name"`
				Level    string `json:"level"`
				Category string `json:"category"`
				Text     string `json:"text"`
			}{
				Date: "2024-01-01", Type: "5", Name: "UV Index", Level: "2", Category: "Low", Text: "No special protection needed.",
			})
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleWeatherIndices(client, WeatherIndicesInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleWeatherIndices failed: %v", err)
	}
	if !strings.Contains(out.IndicesInfo, "1-day Weather Indices - Beijing") {
		t.Fatalf("IndicesInfo = %q, want to contain %q", out.IndicesInfo, "1-day Weather Indices - Beijing")
	}
}

func TestHandleAirQuality_UnknownCodeWithData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/airquality/v1/current/39.90/116.41":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"indexes":[{"code":"qaqi","name":"QAQI","aqi":50,"aqiDisplay":"50"}],"pollutants":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleAirQuality(client, AirQualityInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleAirQuality failed: %v", err)
	}
	if !strings.Contains(out.AirQualityInfo, "Real-time Air Quality - Beijing") {
		t.Fatalf("AirQualityInfo = %q, want to contain %q", out.AirQualityInfo, "Real-time Air Quality - Beijing")
	}
}

func TestHandleAirQualityHourly_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/airquality/v1/hourly/39.90/116.41":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code":"200","hours":[{"forecastTime":"2024-01-01T00:00:00Z","indexes":[{"code":"qaqi","name":"QAQI","aqi":50,"aqiDisplay":"50","level":"1","category":"Good"}],"pollutants":[]} ]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleAirQualityHourly(client, AirQualityHourlyInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleAirQualityHourly failed: %v", err)
	}
	if !strings.Contains(out.HourlyInfo, "24-hour Air Quality Forecast - Beijing") {
		t.Fatalf("HourlyInfo = %q, want to contain %q", out.HourlyInfo, "24-hour Air Quality Forecast - Beijing")
	}
}

func TestHandleAirQualityDaily_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geo/v2/city/lookup":
			json.NewEncoder(w).Encode(api.LocationResponse{
				Code: "200",
				Location: []api.Location{{
					Name: "Beijing", ID: "101010100", Lat: "39.90", Lon: "116.41", Adm1: "Beijing", Adm2: "Beijing",
				}},
			})
		case "/airquality/v1/daily/39.90/116.41":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code":"200","days":[{"forecastStartTime":"2024-01-01T00:00:00Z","forecastEndTime":"2024-01-02T00:00:00Z","indexes":[{"code":"qaqi","name":"QAQI","aqi":50,"aqiDisplay":"50","level":"1","category":"Good"}],"pollutants":[]} ]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-key")
	out, err := handleAirQualityDaily(client, AirQualityDailyInput{CityName: "Beijing"})
	if err != nil {
		t.Fatalf("handleAirQualityDaily failed: %v", err)
	}
	if !strings.Contains(out.DailyInfo, "3-day Air Quality Forecast - Beijing") {
		t.Fatalf("DailyInfo = %q, want to contain %q", out.DailyInfo, "3-day Air Quality Forecast - Beijing")
	}
}
