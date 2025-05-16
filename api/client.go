package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client QWeather API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	LogLevel   LogLevel
}

// NewClient Create a new API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		LogLevel:   LogLevelError, // Default to error logging only
	}
}

// SetLogLevel sets the logging level
func (c *Client) SetLogLevel(level LogLevel) {
	c.LogLevel = level
}

// MakeRequest Send API request
func (c *Client) MakeRequest(endpoint string, params map[string]string, pathParams ...string) ([]byte, error) {
	var urlStr string

	// Handle path parameters
	if len(pathParams) > 0 {
		// Replace placeholders in the path
		endpointWithParams := endpoint
		for _, param := range pathParams {
			endpointWithParams = strings.Replace(endpointWithParams, "{}", param, 1)
		}
		urlStr = fmt.Sprintf("%s%s", c.BaseURL, endpointWithParams)
	} else {
		urlStr = fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	}

	// Create URL and add query parameters
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key to request header
	req.Header.Add("X-QW-Api-Key", c.APIKey)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed, status code: %d, URL: %s", resp.StatusCode, u.String())
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Output response information based on log level
	if c.LogLevel >= LogLevelDebug {
		// Output full response at debug level
		bodyPreview := string(body)
		if len(bodyPreview) > 1000 {
			bodyPreview = bodyPreview[:1000] + "... (truncated)"
		}
		fmt.Printf("API Response [%s]: Status=%d, Body=%s\n", endpoint, resp.StatusCode, bodyPreview)
	} else if c.LogLevel >= LogLevelInfo {
		// Output only status code and endpoint at info level
		fmt.Printf("API Response [%s]: Status=%d\n", endpoint, resp.StatusCode)
	}

	return body, nil
}

// GetLocationByName Get location information by city name
func (c *Client) GetLocationByName(cityName string) (*LocationResponse, error) {
	params := map[string]string{
		"location": cityName,
	}

	data, err := c.MakeRequest("/geo/v2/city/lookup", params)
	if err != nil {
		return nil, err
	}

	var response LocationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse location data: %w", err)
	}

	return &response, nil
}

// GetWeatherNow Get real-time weather
func (c *Client) GetWeatherNow(locationID string) (*WeatherNowResponse, error) {
	params := map[string]string{
		"location": locationID,
	}

	data, err := c.MakeRequest("/v7/weather/now", params)
	if err != nil {
		return nil, err
	}

	var response WeatherNowResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse real-time weather data: %w", err)
	}

	return &response, nil
}

// GetWeatherForecast Get weather forecast
func (c *Client) GetWeatherForecast(locationID, days string) (*WeatherDailyResponse, error) {
	params := map[string]string{
		"location": locationID,
	}

	data, err := c.MakeRequest(fmt.Sprintf("/v7/weather/%s", days), params)
	if err != nil {
		return nil, err
	}

	var response WeatherDailyResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse weather forecast data: %w", err)
	}

	return &response, nil
}

// GetMinutelyPrecipitation Get minutely precipitation forecast
func (c *Client) GetMinutelyPrecipitation(location string) (*MinutelyResponse, error) {
	params := map[string]string{
		"location": location,
	}

	data, err := c.MakeRequest("/v7/minutely/5m", params)
	if err != nil {
		return nil, err
	}

	var response MinutelyResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse minutely precipitation data: %w", err)
	}

	return &response, nil
}

// GetHourlyForecast Get hourly weather forecast
func (c *Client) GetHourlyForecast(locationID, hours string) (*HourlyResponse, error) {
	params := map[string]string{
		"location": locationID,
	}

	data, err := c.MakeRequest(fmt.Sprintf("/v7/weather/%s", hours), params)
	if err != nil {
		return nil, err
	}

	var response HourlyResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse hourly weather data: %w", err)
	}

	return &response, nil
}

// GetWeatherWarning Get weather warnings
func (c *Client) GetWeatherWarning(locationID string) (*WarningResponse, error) {
	params := map[string]string{
		"location": locationID,
	}

	data, err := c.MakeRequest("/v7/warning/now", params)
	if err != nil {
		return nil, err
	}

	var response WarningResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse weather warning data: %w", err)
	}

	return &response, nil
}

// GetWeatherIndices Get weather life indices
func (c *Client) GetWeatherIndices(locationID, days, indexType string) (*IndicesResponse, error) {
	params := map[string]string{
		"location": locationID,
		"type":     indexType,
	}

	data, err := c.MakeRequest(fmt.Sprintf("/v7/indices/%s", days), params)
	if err != nil {
		return nil, err
	}

	var response IndicesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse life indices data: %w", err)
	}

	return &response, nil
}

// GetAirQuality Get real-time air quality
func (c *Client) GetAirQuality(lat, lon string) (*AirQualityResponse, error) {
	endpoint := fmt.Sprintf("/airquality/v1/current/%s/%s", lat, lon)

	data, err := c.MakeRequest(endpoint, map[string]string{}, lat, lon)
	if err != nil {
		return nil, err
	}

	var response AirQualityResponse
	if err := json.Unmarshal(data, &response); err != nil {
		if c.LogLevel >= LogLevelError {
			rawData := string(data)
			fmt.Printf("Error parsing air quality data: %v, raw data: %s\n", err, rawData[:min(len(rawData), 500)])
		}
		return nil, fmt.Errorf("failed to parse air quality data: %w", err)
	}

	// Check if Code field is empty and log warning
	if response.Code == "" {
		if c.LogLevel >= LogLevelInfo {
			fmt.Printf("WARNING: Empty Code field in AirQualityResponse\n")
		}
		// Set a default error code for empty codes to prevent confusion in error messages
		response.Code = "unknown"
	}

	return &response, nil
}

// GetAirQualityHourly Get hourly air quality forecast
func (c *Client) GetAirQualityHourly(lat, lon string) (*AirQualityHourlyResponse, error) {
	endpoint := fmt.Sprintf("/airquality/v1/hourly/%s/%s", lat, lon)

	data, err := c.MakeRequest(endpoint, map[string]string{}, lat, lon)
	if err != nil {
		return nil, err
	}

	var response AirQualityHourlyResponse
	if err := json.Unmarshal(data, &response); err != nil {
		if c.LogLevel >= LogLevelError {
			rawData := string(data)
			fmt.Printf("Error parsing hourly air quality data: %v, raw data: %s\n", err, rawData[:min(len(rawData), 500)])
		}
		return nil, fmt.Errorf("failed to parse hourly air quality data: %w", err)
	}

	// Check if Code field is empty and log warning
	if response.Code == "" {
		if c.LogLevel >= LogLevelInfo {
			fmt.Printf("WARNING: Empty Code field in AirQualityHourlyResponse\n")
		}
		// Set a default error code for empty codes to prevent confusion in error messages
		response.Code = "unknown"
	}

	return &response, nil
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetAirQualityDaily Get daily air quality forecast
func (c *Client) GetAirQualityDaily(lat, lon string) (*AirQualityDailyResponse, error) {
	endpoint := fmt.Sprintf("/airquality/v1/daily/%s/%s", lat, lon)

	data, err := c.MakeRequest(endpoint, map[string]string{}, lat, lon)
	if err != nil {
		return nil, err
	}

	var response AirQualityDailyResponse
	if err := json.Unmarshal(data, &response); err != nil {
		if c.LogLevel >= LogLevelError {
			rawData := string(data)
			fmt.Printf("Error parsing daily air quality data: %v, raw data: %s\n", err, rawData[:min(len(rawData), 500)])
		}
		return nil, fmt.Errorf("failed to parse daily air quality data: %w", err)
	}

	// Check if Code field is empty and log warning
	if response.Code == "" {
		if c.LogLevel >= LogLevelInfo {
			fmt.Printf("WARNING: Empty Code field in AirQualityDailyResponse\n")
		}
		// Set a default error code for empty codes to prevent confusion in error messages
		response.Code = "unknown"
	}

	return &response, nil
}
