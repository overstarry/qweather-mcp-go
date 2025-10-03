package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/overstarry/qweather-mcp-go/api"
)

// WeatherNowInput input parameters for get-weather-now tool
type WeatherNowInput struct {
	CityName string `json:"cityName" jsonschema:"required" jsonschema_description:"Name of the city to query weather for"`
}

// WeatherNowOutput output structure for get-weather-now tool
type WeatherNowOutput struct {
	WeatherInfo string `json:"weatherInfo" jsonschema_description:"Formatted current weather information"`
}

// WeatherForecastInput input parameters for get-weather-forecast tool
type WeatherForecastInput struct {
	CityName string `json:"cityName" jsonschema:"required" jsonschema_description:"Name of the city to query weather for"`
	Days     string `json:"days" jsonschema:"required" jsonschema:"enum=3d,enum=7d,enum=10d,enum=15d,enum=30d" jsonschema_description:"Forecast days"`
}

// WeatherForecastOutput output structure for get-weather-forecast tool
type WeatherForecastOutput struct {
	ForecastInfo string `json:"forecastInfo" jsonschema_description:"Formatted weather forecast information"`
}

// MinutelyPrecipitationInput input parameters for get-minutely-precipitation tool
type MinutelyPrecipitationInput struct {
	CityName string `json:"cityName" jsonschema:"required" jsonschema_description:"Name of the city to query precipitation forecast for"`
}

// MinutelyPrecipitationOutput output structure for get-minutely-precipitation tool
type MinutelyPrecipitationOutput struct {
	PrecipitationInfo string `json:"precipitationInfo" jsonschema_description:"Formatted minutely precipitation forecast"`
}

// HourlyForecastInput input parameters for get-hourly-forecast tool
type HourlyForecastInput struct {
	CityName string `json:"cityName" jsonschema:"required" jsonschema_description:"Name of the city to query weather for"`
	Hours    string `json:"hours,omitempty" jsonschema_description:"Forecast hours (24h, 72h, or 168h). Defaults to 24h if not specified."`
}

// HourlyForecastOutput output structure for get-hourly-forecast tool
type HourlyForecastOutput struct {
	HourlyInfo string `json:"hourlyInfo" jsonschema_description:"Formatted hourly weather forecast"`
}

// WeatherWarningInput input parameters for get-weather-warning tool
type WeatherWarningInput struct {
	CityName string `json:"cityName" jsonschema:"required" jsonschema_description:"Name of the city to query weather warnings for"`
}

// WeatherWarningOutput output structure for get-weather-warning tool
type WeatherWarningOutput struct {
	WarningInfo string `json:"warningInfo" jsonschema_description:"Formatted weather warning information"`
}

// RegisterWeatherTools Register weather-related tools
func RegisterWeatherTools(s *mcp.Server, client *api.Client) {
	// Real-time weather tool
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get-weather-now",
		Description: "Real-time weather API provides current weather conditions for cities worldwide. Available data includes: temperature, feels-like temperature, weather conditions, wind direction, wind force level, relative humidity, precipitation, atmospheric pressure, and visibility. Data is updated in real-time, providing the most accurate current weather information.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input WeatherNowInput) (*mcp.CallToolResult, WeatherNowOutput, error) {
		if input.CityName == "" {
			return nil, WeatherNowOutput{}, fmt.Errorf("city name cannot be empty")
		}

		// Query city ID
		locationData, err := client.GetLocationByName(input.CityName)
		if err != nil {
			return nil, WeatherNowOutput{}, fmt.Errorf("failed to query city: %w", err)
		}

		if locationData.Code != "200" {
			return nil, WeatherNowOutput{}, fmt.Errorf("failed to query city, API returned an error")
		}

		if len(locationData.Location) == 0 {
			return nil, WeatherNowOutput{}, fmt.Errorf("no matching city found")
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get real-time weather data
		weatherData, err := client.GetWeatherNow(cityID)
		if err != nil {
			return nil, WeatherNowOutput{}, fmt.Errorf("failed to get real-time weather data: %w", err)
		}

		if weatherData.Code != "200" {
			return nil, WeatherNowOutput{}, fmt.Errorf("failed to get real-time weather data, API returned an error")
		}

		// Format weather information
		now := weatherData.Now
		weatherText := []string{
			fmt.Sprintf("Current Weather - %s (%s %s):", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			fmt.Sprintf("Temperature: %s°C (Feels like: %s°C)", now.Temp, now.FeelsLike),
			fmt.Sprintf("Weather Condition: %s", now.Text),
			fmt.Sprintf("Wind Direction: %s Wind Force: %s", now.WindDir, now.WindScale),
			fmt.Sprintf("Humidity: %s%%", now.Humidity),
			fmt.Sprintf("Precipitation: %smm", now.Precip),
			fmt.Sprintf("Pressure: %shPa", now.Pressure),
			fmt.Sprintf("Visibility: %skm", now.Vis),
			fmt.Sprintf("Last Updated: %s", weatherData.UpdateTime),
		}

		return nil, WeatherNowOutput{WeatherInfo: strings.Join(weatherText, "\n")}, nil
	})

	// Weather forecast tool
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get-weather-forecast",
		Description: "Weather forecast API provides detailed weather predictions for cities worldwide, supporting forecasts from 3 to 30 days. Available data includes: sunrise/sunset times, moonrise/moonset times, temperature range, weather conditions, wind direction and speed, relative humidity, precipitation, atmospheric pressure, cloud cover, and UV index. Forecasts are updated daily to ensure accuracy.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input WeatherForecastInput) (*mcp.CallToolResult, WeatherForecastOutput, error) {
		if input.CityName == "" {
			return nil, WeatherForecastOutput{}, fmt.Errorf("city name cannot be empty")
		}

		if input.Days == "" {
			input.Days = "3d"
		}

		// Query city ID
		locationData, err := client.GetLocationByName(input.CityName)
		if err != nil {
			return nil, WeatherForecastOutput{}, fmt.Errorf("failed to query city: %w", err)
		}

		if locationData.Code != "200" {
			return nil, WeatherForecastOutput{}, fmt.Errorf("failed to query city, API returned an error")
		}

		if len(locationData.Location) == 0 {
			return nil, WeatherForecastOutput{}, fmt.Errorf("no matching city found")
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get weather forecast data
		weatherData, err := client.GetWeatherForecast(cityID, input.Days)
		if err != nil {
			return nil, WeatherForecastOutput{}, fmt.Errorf("failed to get weather forecast data: %w", err)
		}

		if weatherData.Code != "200" {
			return nil, WeatherForecastOutput{}, fmt.Errorf("failed to get weather forecast data, API returned an error")
		}

		// Format weather forecast information
		forecastText := []string{
			fmt.Sprintf("%s Day Weather Forecast - %s (%s %s):", strings.Replace(input.Days, "d", "", -1), cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			fmt.Sprintf("Last Updated: %s", weatherData.UpdateTime),
			"",
		}

		for _, day := range weatherData.Daily {
			dayForecast := []string{
				fmt.Sprintf("Date: %s", day.FxDate),
				fmt.Sprintf("Temperature: %s°C ~ %s°C", day.TempMin, day.TempMax),
				fmt.Sprintf("Day: %s", day.TextDay),
				fmt.Sprintf("Night: %s", day.TextNight),
				fmt.Sprintf("Sunrise: %s  Sunset: %s", day.Sunrise, day.Sunset),
				fmt.Sprintf("Precipitation: %smm", day.Precip),
				fmt.Sprintf("Humidity: %s%%", day.Humidity),
				fmt.Sprintf("Wind: Day-%s(Force %s), Night-%s(Force %s)", day.WindDirDay, day.WindScaleDay, day.WindDirNight, day.WindScaleNight),
				fmt.Sprintf("UV Index: %s", day.UvIndex),
				"---",
			}
			forecastText = append(forecastText, strings.Join(dayForecast, "\n"))
		}

		return nil, WeatherForecastOutput{ForecastInfo: strings.Join(forecastText, "\n")}, nil
	})

	// Minutely precipitation forecast tool
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get-minutely-precipitation",
		Description: "Minutely precipitation forecast API provides accurate precipitation predictions for the next 2 hours for cities worldwide. Available data includes precipitation type (rain/snow) and amount for each minute. This high-precision forecast is particularly useful for outdoor activity planning and real-time weather monitoring.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MinutelyPrecipitationInput) (*mcp.CallToolResult, MinutelyPrecipitationOutput, error) {
		if input.CityName == "" {
			return nil, MinutelyPrecipitationOutput{}, fmt.Errorf("city name cannot be empty")
		}

		// Query city location coordinates
		locationData, err := client.GetLocationByName(input.CityName)
		if err != nil {
			return nil, MinutelyPrecipitationOutput{}, fmt.Errorf("failed to query city: %w", err)
		}

		if locationData.Code != "200" {
			return nil, MinutelyPrecipitationOutput{}, fmt.Errorf("failed to query city, API returned an error")
		}

		if len(locationData.Location) == 0 {
			return nil, MinutelyPrecipitationOutput{}, fmt.Errorf("no matching city found")
		}

		// Use the coordinates of the first matching city
		cityInfo := locationData.Location[0]
		location := fmt.Sprintf("%s,%s", cityInfo.Lon, cityInfo.Lat)

		// Get minutely precipitation forecast data
		precipData, err := client.GetMinutelyPrecipitation(location)
		if err != nil {
			return nil, MinutelyPrecipitationOutput{}, fmt.Errorf("failed to get minutely precipitation forecast data: %w", err)
		}

		if precipData.Code != "200" {
			return nil, MinutelyPrecipitationOutput{}, fmt.Errorf("failed to get minutely precipitation forecast data, API returned an error")
		}

		// Format precipitation forecast information
		precipText := []string{
			fmt.Sprintf("Minutely Precipitation Forecast - %s (%s %s):", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			fmt.Sprintf("Forecast Description: %s", precipData.Summary),
			fmt.Sprintf("Last Updated: %s", precipData.UpdateTime),
			"",
			"2-Hour Precipitation Forecast:",
		}

		for _, minute := range precipData.Minutely {
			timeStr := strings.Split(strings.Split(minute.FxTime, "T")[1], "+")[0]
			precipType := "Rain"
			if minute.Type == "snow" {
				precipType = "Snow"
			}
			precipText = append(precipText, fmt.Sprintf("Time: %s - %s: %smm", timeStr, precipType, minute.Precip))
		}

		precipText = append(precipText, "", fmt.Sprintf("Data Source: %s", precipData.FxLink))

		return nil, MinutelyPrecipitationOutput{PrecipitationInfo: strings.Join(precipText, "\n")}, nil
	})

	// Hourly weather forecast tool
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get-hourly-forecast",
		Description: "Hourly weather forecast API provides detailed weather information for the next 24-168 hours for cities worldwide. Available data includes: temperature, weather conditions, wind force, wind speed, wind direction, relative humidity, atmospheric pressure, precipitation probability, dew point temperature, and cloud cover. Forecast data is updated hourly to ensure accuracy.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input HourlyForecastInput) (*mcp.CallToolResult, HourlyForecastOutput, error) {
		if input.CityName == "" {
			return nil, HourlyForecastOutput{}, fmt.Errorf("city name cannot be empty")
		}

		if input.Hours == "" {
			input.Hours = "24h"
		}

		// Query city ID
		locationData, err := client.GetLocationByName(input.CityName)
		if err != nil {
			return nil, HourlyForecastOutput{}, fmt.Errorf("failed to query city: %w", err)
		}

		if locationData.Code != "200" {
			return nil, HourlyForecastOutput{}, fmt.Errorf("failed to query city, API returned an error")
		}

		if len(locationData.Location) == 0 {
			return nil, HourlyForecastOutput{}, fmt.Errorf("no matching city found")
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get hourly weather forecast data
		hourlyData, err := client.GetHourlyForecast(cityID, input.Hours)
		if err != nil {
			return nil, HourlyForecastOutput{}, fmt.Errorf("failed to get hourly weather forecast data: %w", err)
		}

		if hourlyData.Code != "200" {
			return nil, HourlyForecastOutput{}, fmt.Errorf("failed to get hourly weather forecast data, API returned an error")
		}

		// Format hourly weather forecast information
		hourlyText := []string{
			fmt.Sprintf("%s Hour Weather Forecast - %s (%s %s):", strings.Replace(input.Hours, "h", "", -1), cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			fmt.Sprintf("Last Updated: %s", hourlyData.UpdateTime),
			"",
		}

		for _, hour := range hourlyData.Hourly {
			timeStr := strings.Split(strings.Split(hour.FxTime, "T")[1], "+")[0]
			hourForecast := []string{
				fmt.Sprintf("Time: %s", timeStr),
				fmt.Sprintf("Temperature: %s°C", hour.Temp),
				fmt.Sprintf("Weather: %s", hour.Text),
				fmt.Sprintf("Wind Direction: %s (Force %s, %skm/h)", hour.WindDir, hour.WindScale, hour.WindSpeed),
				fmt.Sprintf("Humidity: %s%%", hour.Humidity),
				fmt.Sprintf("Precipitation: %smm", hour.Precip),
				fmt.Sprintf("Pressure: %shPa", hour.Pressure),
			}

			if hour.Cloud != "" {
				hourForecast = append(hourForecast, fmt.Sprintf("Cloud Cover: %s%%", hour.Cloud))
			}

			if hour.Dew != "" {
				hourForecast = append(hourForecast, fmt.Sprintf("Dew Point: %s°C", hour.Dew))
			}

			hourForecast = append(hourForecast, "---")
			hourlyText = append(hourlyText, strings.Join(hourForecast, "\n"))
		}

		return nil, HourlyForecastOutput{HourlyInfo: strings.Join(hourlyText, "\n")}, nil
	})

	// Weather warning tool
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get-weather-warning",
		Description: "Weather warning API provides real-time weather warning data issued by official agencies in China and multiple countries/regions worldwide. Data includes warning issuing agency, publication time, warning title, detailed warning information, warning level, warning type, and other relevant information.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input WeatherWarningInput) (*mcp.CallToolResult, WeatherWarningOutput, error) {
		if input.CityName == "" {
			return nil, WeatherWarningOutput{}, fmt.Errorf("city name cannot be empty")
		}

		// Query city ID
		locationData, err := client.GetLocationByName(input.CityName)
		if err != nil {
			return nil, WeatherWarningOutput{}, fmt.Errorf("failed to query city: %w", err)
		}

		if locationData.Code != "200" {
			return nil, WeatherWarningOutput{}, fmt.Errorf("failed to query city, API returned an error")
		}

		if len(locationData.Location) == 0 {
			return nil, WeatherWarningOutput{}, fmt.Errorf("no matching city found")
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get weather warning data
		warningData, err := client.GetWeatherWarning(cityID)
		if err != nil {
			return nil, WeatherWarningOutput{}, fmt.Errorf("failed to get weather warning data: %w", err)
		}

		if warningData.Code != "200" {
			return nil, WeatherWarningOutput{}, fmt.Errorf("failed to get weather warning data, API returned an error")
		}

		// Check if there are active weather warnings
		if len(warningData.Warning) == 0 {
			warningInfo := fmt.Sprintf("Currently %s (%s %s) has no active weather warnings", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2)
			return nil, WeatherWarningOutput{WarningInfo: warningInfo}, nil
		}

		// Format weather warning information
		warningText := []string{
			fmt.Sprintf("Weather Warnings - %s (%s %s):", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			fmt.Sprintf("Last Updated: %s", warningData.UpdateTime),
			"",
		}

		for _, warning := range warningData.Warning {
			startTime := warning.StartTime
			if startTime == "" {
				startTime = "Not specified"
			}
			endTime := warning.EndTime
			if endTime == "" {
				endTime = "Not specified"
			}

			warningInfo := []string{
				fmt.Sprintf("Warning Title: %s", warning.Title),
				fmt.Sprintf("Issuing Agency: %s", warning.Sender),
				fmt.Sprintf("Publication Time: %s", warning.PubTime),
				fmt.Sprintf("Warning Type: %s", warning.TypeName),
				fmt.Sprintf("Severity: %s (%s)", warning.Severity, warning.SeverityColor),
				fmt.Sprintf("Valid Period: %s to %s", startTime, endTime),
				fmt.Sprintf("Status: %s", warning.Status),
				fmt.Sprintf("Details: %s", warning.Text),
				"---",
			}
			warningText = append(warningText, strings.Join(warningInfo, "\n"))
		}

		return nil, WeatherWarningOutput{WarningInfo: strings.Join(warningText, "\n")}, nil
	})
}
