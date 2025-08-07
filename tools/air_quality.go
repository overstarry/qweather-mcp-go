package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/overstarry/qweather-mcp-go/api"
	"github.com/overstarry/qweather-mcp-go/utils"
)

// RegisterAirQualityTools Register air quality related tools
func RegisterAirQualityTools(s *server.MCPServer, client *api.Client) {
	// Real-time air quality tool
	airQualityTool := mcp.NewTool("get-air-quality",
		mcp.WithDescription("Real-time air quality API provides air quality data for specific locations with 1x1 kilometer precision. Includes AQI based on different national/regional local standards, AQI level, color, main pollutants, QWeather universal AQI, pollutant concentrations, sub-indices, health recommendations, and related monitoring station information."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query air quality for"),
		),
	)

	s.AddTool(airQualityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city coordinates
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
		// Format coordinates, keep up to 2 decimal places
		lat := fmt.Sprintf("%.2f", parseFloat(cityInfo.Lat))
		lon := fmt.Sprintf("%.2f", parseFloat(cityInfo.Lon))

		// Get air quality data
		airQualityData, err := client.GetAirQuality(lat, lon)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get air quality data: %v (Coordinates: lat=%s, lon=%s)", err, lat, lon)
			return mcp.NewToolResultError(errorMsg), nil
		}

		if airQualityData.Code != "200" && airQualityData.Code != "unknown" {
			errorMsg := fmt.Sprintf("Failed to get air quality data: API returned error code %s (Coordinates: lat=%s, lon=%s)",
				airQualityData.Code, lat, lon)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Even if Code is "unknown", we need to confirm that the indexes array is not empty
		if airQualityData.Code == "unknown" && len(airQualityData.Indexes) > 0 {
			fmt.Printf("API returned status code 'unknown', but indexes array is not empty. Data is valid, continuing to process request\n")
			fmt.Printf("Data structure analysis: API response contains %d air quality indexes\n", len(airQualityData.Indexes))
		}

		if len(airQualityData.Indexes) == 0 {
			errorMsg := fmt.Sprintf("Failed to get air quality data: API returned success but no air quality indexes found (Coordinates: lat=%s, lon=%s, Code=%s)",
				lat, lon, airQualityData.Code)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Format air quality information
		airQualityText := []string{
			fmt.Sprintf("Real-time Air Quality - %s (%s %s):", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			"",
			"Air Quality Index:",
		}

		for _, index := range airQualityData.Indexes {
			indexInfo := []string{
				fmt.Sprintf("%s: %s", index.Name, index.AqiDisplay),
			}

			if index.Level != "" {
				indexInfo = append(indexInfo, fmt.Sprintf("Level: %s", index.Level))
			}

			if index.Category != "" {
				// Use English directly
				indexInfo = append(indexInfo, fmt.Sprintf("Category: %s", index.Category))
			}

			if index.PrimaryPollutant != nil {
				// Use English directly
				indexInfo = append(indexInfo, fmt.Sprintf("Main Pollutant: %s", index.PrimaryPollutant.Name))
			}

			if index.Health != nil {
				healthInfo := []string{
					"Health Effects:",
					// Use English directly
					fmt.Sprintf("- %s", index.Health.Effect),
					"Health Recommendations:",
					// Use English directly
					fmt.Sprintf("- General Population: %s", index.Health.Advice.GeneralPopulation),
					fmt.Sprintf("- Sensitive Population: %s", index.Health.Advice.SensitivePopulation),
				}
				indexInfo = append(indexInfo, strings.Join(healthInfo, "\n"))
			}

			indexInfo = append(indexInfo, "---")
			airQualityText = append(airQualityText, strings.Join(indexInfo, "\n"))
		}

		airQualityText = append(airQualityText, "", "Pollutant Concentrations:")
		for _, pollutant := range airQualityData.Pollutants {
			// Use English directly
			airQualityText = append(airQualityText, fmt.Sprintf("%s: %.1f%s", pollutant.Name, pollutant.Concentration.Value, pollutant.Concentration.Unit))
		}

		if len(airQualityData.Stations) > 0 {
			airQualityText = append(airQualityText, "", "Related Monitoring Stations:")
			for _, station := range airQualityData.Stations {
				airQualityText = append(airQualityText, fmt.Sprintf("- %s", station.Name))
			}
		}

		return mcp.NewToolResultText(strings.Join(airQualityText, "\n")), nil
	})

	// Hourly air quality forecast tool
	airQualityHourlyTool := mcp.NewTool("get-air-quality-hourly",
		mcp.WithDescription("Hourly air quality forecast API provides air quality data for the next 24 hours, including AQI, pollutant concentrations, sub-indices, and health recommendations. Data includes various air quality standards (such as QAQI, GB-DEFRA, etc.) and specific concentrations of pollutants like PM2.5, PM10, NO2, O3, SO2, etc."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query air quality forecast for"),
		),
	)

	s.AddTool(airQualityHourlyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city coordinates
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
		// Format coordinates, keep up to 2 decimal places
		lat := fmt.Sprintf("%.2f", parseFloat(cityInfo.Lat))
		lon := fmt.Sprintf("%.2f", parseFloat(cityInfo.Lon))

		// Get hourly air quality forecast data
		airQualityData, err := client.GetAirQualityHourly(lat, lon)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get hourly air quality forecast data: %v (Coordinates: lat=%s, lon=%s)", err, lat, lon)
			return mcp.NewToolResultError(errorMsg), nil
		}

		if airQualityData.Code != "200" && airQualityData.Code != "unknown" {
			errorMsg := fmt.Sprintf("Failed to get hourly air quality forecast data: API returned error code %s (Coordinates: lat=%s, lon=%s)",
				airQualityData.Code, lat, lon)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Even if Code is "unknown", we need to confirm that the hours array is not empty
		if airQualityData.Code == "unknown" && len(airQualityData.Hours) > 0 {
			fmt.Printf("API returned status code 'unknown', but hours array is not empty. Data is valid, continuing to process request\n")
			fmt.Printf("Data structure analysis: API response contains %d hours of forecast data\n", len(airQualityData.Hours))
		}

		if len(airQualityData.Hours) == 0 {
			errorMsg := fmt.Sprintf("Failed to get hourly air quality forecast data: API returned success but no forecast hours found (Coordinates: lat=%s, lon=%s, Code=%s)",
				lat, lon, airQualityData.Code)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Format hourly air quality forecast information
		hourlyText := []string{
			fmt.Sprintf("24-hour Air Quality Forecast - %s (%s %s):", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			"",
		}

		for _, hour := range airQualityData.Hours {
			// Parse time and format to local time
			t, err := time.Parse(time.RFC3339, hour.ForecastTime)
			if err != nil {
				t = time.Time{} // Use zero value time
			}
			timeStr := t.Format("2006-01-02 15:04") + " UTC"

			var indexInfos []string
			for _, index := range hour.Indexes {
				var healthInfo string
				if index.Health != nil {
					healthInfo = fmt.Sprintf("Health Effects: %s\nHealth Recommendations:\n  General Population: %s\n  Sensitive Population: %s",
						index.Health.Effect,
						index.Health.Advice.GeneralPopulation,
						index.Health.Advice.SensitivePopulation)
				}

				primaryPollutant := ""
				if index.PrimaryPollutant != nil {
					// Use English directly
					primaryPollutant = fmt.Sprintf("  Main Pollutant: %s", index.PrimaryPollutant.Name)
				}

				indexInfo := []string{
					"Air Quality Index:",
					fmt.Sprintf("  %s: %s", index.Name, index.AqiDisplay),
					fmt.Sprintf("  Level: %s", index.Level),
					fmt.Sprintf("  Category: %s", index.Category),
					primaryPollutant,
					healthInfo,
				}
				indexInfos = append(indexInfos, utils.JoinStrings(indexInfo, "\n"))
			}

			var pollutantInfos []string
			if len(hour.Pollutants) > 0 {
				pollutantInfos = append(pollutantInfos, "Pollutant Concentrations:")
				for _, pollutant := range hour.Pollutants {
					// Use English directly
					pollutantInfos = append(pollutantInfos, fmt.Sprintf("  %s: %.1f%s",
						pollutant.Name,
						pollutant.Concentration.Value,
						pollutant.Concentration.Unit))
				}
			} else {
				pollutantInfos = append(pollutantInfos, "No pollutant data")
			}

			hourInfo := []string{
				fmt.Sprintf("Forecast Time: %s", timeStr),
				strings.Join(indexInfos, "\n"),
				strings.Join(pollutantInfos, "\n"),
				"---",
			}
			hourlyText = append(hourlyText, strings.Join(hourInfo, "\n\n"))
		}

		return mcp.NewToolResultText(strings.Join(hourlyText, "\n")), nil
	})

	// Daily air quality forecast tool
	airQualityDailyTool := mcp.NewTool("get-air-quality-daily",
		mcp.WithDescription("Daily air quality forecast API provides air quality predictions for the next 3 days, including AQI values, pollutant concentrations, and health recommendations. Data includes various air quality standards and specific concentrations of pollutants such as PM2.5, PM10, NO2, O3, SO2, etc."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query air quality forecast for"),
		),
	)

	s.AddTool(airQualityDailyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		if cityName == "" {
			return mcp.NewToolResultError("City name cannot be empty"), nil
		}

		// Query city coordinates
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
		// Format coordinates, keep up to 2 decimal places
		lat := fmt.Sprintf("%.2f", parseFloat(cityInfo.Lat))
		lon := fmt.Sprintf("%.2f", parseFloat(cityInfo.Lon))

		// Get daily air quality forecast data
		airQualityData, err := client.GetAirQualityDaily(lat, lon)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get daily air quality forecast data: %v (Coordinates: lat=%s, lon=%s)", err, lat, lon)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Handle the Code field returned by the API. If it's "unknown" (default value set by GetAirQualityDaily) and the days array is not empty, consider it a successful response
		if airQualityData.Code != "200" && airQualityData.Code != "unknown" {
			errorMsg := fmt.Sprintf("Failed to get daily air quality forecast data: API returned error code %s (Coordinates: lat=%s, lon=%s)",
				airQualityData.Code, lat, lon)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Even if Code is "unknown", we need to confirm that the days array is not empty
		if airQualityData.Code == "unknown" && len(airQualityData.Days) > 0 {
			fmt.Printf("API returned status code 'unknown', but days array is not empty. Data is valid, continuing to process request\n")
			fmt.Printf("Data structure analysis: API response contains %d days of forecast data\n", len(airQualityData.Days))
		}

		if len(airQualityData.Days) == 0 {
			errorMsg := fmt.Sprintf("Failed to get daily air quality forecast data: API returned success but no forecast days found (Coordinates: lat=%s, lon=%s, Code=%s)",
				lat, lon, airQualityData.Code)
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Format daily air quality forecast information
		dailyText := []string{
			fmt.Sprintf("3-day Air Quality Forecast - %s (%s %s):", cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			"",
		}

		for _, day := range airQualityData.Days {
			// Parse time and format
			startTime, err := time.Parse(time.RFC3339, day.ForecastStartTime)
			if err != nil {
				startTime = time.Time{} // Use zero value time
			}
			startTimeStr := startTime.Format("2006-01-02 15:04") + " UTC"

			endTime, err := time.Parse(time.RFC3339, day.ForecastEndTime)
			if err != nil {
				endTime = time.Time{} // Use zero value time
			}
			endTimeStr := endTime.Format("2006-01-02 15:04") + " UTC"

			var indexInfos []string
			for _, index := range day.Indexes {
				var healthInfo string
				if index.Health != nil {
					healthInfo = fmt.Sprintf("Health Effects: %s\nHealth Recommendations:\n  General Population: %s\n  Sensitive Population: %s",
						index.Health.Effect,
						index.Health.Advice.GeneralPopulation,
						index.Health.Advice.SensitivePopulation)
				}

				primaryPollutant := ""
				if index.PrimaryPollutant != nil {
					// Use English directly
					primaryPollutant = fmt.Sprintf("  Main Pollutant: %s", index.PrimaryPollutant.Name)
				}

				indexInfo := []string{
					"Air Quality Index:",
					fmt.Sprintf("  %s: %s", index.Name, index.AqiDisplay),
					fmt.Sprintf("  Level: %s", index.Level),
					fmt.Sprintf("  Category: %s", index.Category),
					primaryPollutant,
					healthInfo,
				}
				indexInfos = append(indexInfos, utils.JoinStrings(indexInfo, "\n"))
			}

			var pollutantInfos []string
			if len(day.Pollutants) > 0 {
				pollutantInfos = append(pollutantInfos, "Pollutant Concentrations:")
				for _, pollutant := range day.Pollutants {
					// Use English directly
					pollutantInfos = append(pollutantInfos, fmt.Sprintf("  %s: %.1f%s",
						pollutant.Name,
						pollutant.Concentration.Value,
						pollutant.Concentration.Unit))
				}
			} else {
				pollutantInfos = append(pollutantInfos, "No pollutant data")
			}

			dayInfo := []string{
				fmt.Sprintf("Forecast Period: %s to %s", startTimeStr, endTimeStr),
				strings.Join(indexInfos, "\n"),
				strings.Join(pollutantInfos, "\n"),
				"---",
			}
			dailyText = append(dailyText, strings.Join(dayInfo, "\n\n"))
		}

		return mcp.NewToolResultText(strings.Join(dailyText, "\n")), nil
	})
}

// parseFloat Parses string to float, returns 0 on error
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
