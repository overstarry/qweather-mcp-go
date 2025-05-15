package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/overstarry/qweather-mcp-go/api"
)

// RegisterWeatherTools Register weather-related tools
func RegisterWeatherTools(s *server.MCPServer, client *api.Client) {
	// Real-time weather tool
	nowTool := mcp.NewTool("get-weather-now",
		mcp.WithDescription("Real-time weather API provides current weather conditions for cities worldwide. Available data includes: temperature, feels-like temperature, weather conditions, wind direction, wind force level, relative humidity, precipitation, atmospheric pressure, and visibility. Data is updated in real-time, providing the most accurate current weather information."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query weather for"),
		),
	)

	s.AddTool(nowTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city ID
		locationData, err := client.GetLocationByName(cityName)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to query city", err), nil
		}

		if locationData.Code != "200" {
			return mcp.NewToolResultError("Failed to query city, API returned an error"), nil
		}

		if len(locationData.Location) == 0 {
			return mcp.NewToolResultError("No matching city found"), nil
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get real-time weather data
		weatherData, err := client.GetWeatherNow(cityID)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to get real-time weather data", err), nil
		}

		if weatherData.Code != "200" {
			return mcp.NewToolResultError("Failed to get real-time weather data, API returned an error"), nil
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

		return mcp.NewToolResultText(strings.Join(weatherText, "\n")), nil
	})

	// Weather forecast tool
	forecastTool := mcp.NewTool("get-weather-forecast",
		mcp.WithDescription("Weather forecast API provides detailed weather predictions for cities worldwide, supporting forecasts from 3 to 30 days. Available data includes: sunrise/sunset times, moonrise/moonset times, temperature range, weather conditions, wind direction and speed, relative humidity, precipitation, atmospheric pressure, cloud cover, and UV index. Forecasts are updated daily to ensure accuracy."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query weather for"),
		),
		mcp.WithString("days",
			mcp.Required(),
			mcp.Enum("3d", "7d", "10d", "15d", "30d"),
			mcp.Description("Forecast days"),
		),
	)

	s.AddTool(forecastTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		days := mcp.ParseString(request, "days", "3d")

		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city ID
		locationData, err := client.GetLocationByName(cityName)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to query city", err), nil
		}

		if locationData.Code != "200" {
			return mcp.NewToolResultError("Failed to query city, API returned an error"), nil
		}

		if len(locationData.Location) == 0 {
			return mcp.NewToolResultError("No matching city found"), nil
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get weather forecast data
		weatherData, err := client.GetWeatherForecast(cityID, days)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to get weather forecast data", err), nil
		}

		if weatherData.Code != "200" {
			return mcp.NewToolResultError("Failed to get weather forecast data, API returned an error"), nil
		}

		// Format weather forecast information
		forecastText := []string{
			fmt.Sprintf("%s Day Weather Forecast - %s (%s %s):", strings.Replace(days, "d", "", -1), cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
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

		return mcp.NewToolResultText(strings.Join(forecastText, "\n")), nil
	})

	// Minutely precipitation forecast tool
	minutelyTool := mcp.NewTool("get-minutely-precipitation",
		mcp.WithDescription("Minutely precipitation forecast API provides accurate precipitation predictions for the next 2 hours for cities worldwide. Available data includes precipitation type (rain/snow) and amount for each minute. This high-precision forecast is particularly useful for outdoor activity planning and real-time weather monitoring."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query precipitation forecast for"),
		),
	)

	s.AddTool(minutelyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city location coordinates
		locationData, err := client.GetLocationByName(cityName)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to query city", err), nil
		}

		if locationData.Code != "200" {
			return mcp.NewToolResultError("Failed to query city, API returned an error"), nil
		}

		if len(locationData.Location) == 0 {
			return mcp.NewToolResultError("No matching city found"), nil
		}

		// Use the coordinates of the first matching city
		cityInfo := locationData.Location[0]
		location := fmt.Sprintf("%s,%s", cityInfo.Lon, cityInfo.Lat)

		// Get minutely precipitation forecast data
		precipData, err := client.GetMinutelyPrecipitation(location)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to get minutely precipitation forecast data", err), nil
		}

		if precipData.Code != "200" {
			return mcp.NewToolResultError("Failed to get minutely precipitation forecast data, API returned an error"), nil
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

		return mcp.NewToolResultText(strings.Join(precipText, "\n")), nil
	})

	// Hourly weather forecast tool
	hourlyTool := mcp.NewTool("get-hourly-forecast",
		mcp.WithDescription("Hourly weather forecast API provides detailed weather information for the next 24-168 hours for cities worldwide. Available data includes: temperature, weather conditions, wind force, wind speed, wind direction, relative humidity, atmospheric pressure, precipitation probability, dew point temperature, and cloud cover. Forecast data is updated hourly to ensure accuracy."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query weather for"),
		),
		mcp.WithString("hours",
			mcp.Enum("24h", "72h", "168h"),
			mcp.Description("Forecast hours (24h, 72h or 168h)"),
		),
	)

	s.AddTool(hourlyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		hours := mcp.ParseString(request, "hours", "24h")

		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city ID
		locationData, err := client.GetLocationByName(cityName)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to query city", err), nil
		}

		if locationData.Code != "200" {
			return mcp.NewToolResultError("Failed to query city, API returned an error"), nil
		}

		if len(locationData.Location) == 0 {
			return mcp.NewToolResultError("No matching city found"), nil
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get hourly weather forecast data
		hourlyData, err := client.GetHourlyForecast(cityID, hours)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to get hourly weather forecast data", err), nil
		}

		if hourlyData.Code != "200" {
			return mcp.NewToolResultError("Failed to get hourly weather forecast data, API returned an error"), nil
		}

		// Format hourly weather forecast information
		hourlyText := []string{
			fmt.Sprintf("%s Hour Weather Forecast - %s (%s %s):", strings.Replace(hours, "h", "", -1), cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
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

		return mcp.NewToolResultText(strings.Join(hourlyText, "\n")), nil
	})

	// Weather warning tool
	warningTool := mcp.NewTool("get-weather-warning",
		mcp.WithDescription("Weather warning API provides real-time weather warning data issued by official agencies in China and multiple countries/regions worldwide. Data includes warning issuing agency, publication time, warning title, detailed warning information, warning level, warning type, and other relevant information."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query weather warnings for"),
		),
	)

	s.AddTool(warningTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city ID
		locationData, err := client.GetLocationByName(cityName)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to query city", err), nil
		}

		if locationData.Code != "200" {
			return mcp.NewToolResultError("Failed to query city, API returned an error"), nil
		}

		if len(locationData.Location) == 0 {
			return mcp.NewToolResultError("No matching city found"), nil
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get weather warning data
		warningData, err := client.GetWeatherWarning(cityID)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to get weather warning data", err), nil
		}

		if warningData.Code != "200" {
			return mcp.NewToolResultError("Failed to get weather warning data, API returned an error"), nil
		}

		// Check if there are active weather warnings
		if len(warningData.Warning) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("Currently %s (%s %s) has no active weather warnings", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2)), nil
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

		return mcp.NewToolResultText(strings.Join(warningText, "\n")), nil
	})
}
