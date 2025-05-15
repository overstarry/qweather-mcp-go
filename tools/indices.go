package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/overstarry/qweather-mcp-go/api"
)

// RegisterIndicesTools Register weather indices related tools
func RegisterIndicesTools(s *server.MCPServer, client *api.Client) {
	// Weather indices tool
	indicesTool := mcp.NewTool("get-weather-indices",
		mcp.WithDescription("Weather life indices forecast API provides various life indices for cities worldwide. Supports 1-day and 3-day forecasts. Available index types:\n\n"+
			"- Type 0: All index types\n"+
			"- Type 1: Sports (indicates suitability for outdoor sports activities)\n"+
			"- Type 2: Car Washing (suggests whether it's suitable to wash cars)\n"+
			"- Type 3: Dressing (provides clothing suggestions based on weather)\n"+
			"- Type 4: Fishing (shows suitability of fishing conditions)\n"+
			"- Type 5: UV (ultraviolet radiation intensity level)\n"+
			"- Type 6: Travel (indicates suitability for travel and sightseeing)\n"+
			"- Type 7: Allergy (allergy and pollen risk level)\n"+
			"- Type 8: Cold (cold risk level)\n"+
			"- Type 9: Comfort (overall comfort level of weather)\n"+
			"- Type 10: Wind (wind conditions and their effects)\n"+
			"- Type 11: Sunglasses (need for wearing sunglasses)\n"+
			"- Type 12: Makeup (weather effects on makeup)\n"+
			"- Type 13: Sunscreen (sunscreen needs)\n"+
			"- Type 14: Traffic (weather effects on traffic conditions)\n"+
			"- Type 15: Sports Watching (suitability for watching outdoor sports)\n"+
			"- Type 16: Air Pollution Diffusion Conditions (air pollution diffusion conditions)\n\n"+
			"Note: Not all cities provide all indices. International cities mainly support types 1, 2, 4, and 5."),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("Name of the city to query weather indices for"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Enum("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"),
			mcp.Description("Type of weather index to retrieve"),
		),
		mcp.WithString("days",
			mcp.Enum("1d", "3d"),
			mcp.Description("Forecast days (1d or 3d)"),
		),
	)

	s.AddTool(indicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cityName := mcp.ParseString(request, "cityName", "")
		indexType := mcp.ParseString(request, "type", "0")
		days := mcp.ParseString(request, "days", "1d")

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

		// Get weather indices data
		indicesData, err := client.GetWeatherIndices(cityID, days, indexType)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to get weather indices data", err), nil
		}

		if indicesData.Code != "200" {
			return mcp.NewToolResultError("Failed to get weather indices data, API returned an error"), nil
		}

		// Format weather indices information
		daysText := "1-day"
		if days == "3d" {
			daysText = "3-day"
		}

		indicesText := []string{
			fmt.Sprintf("%s Weather Indices - %s (%s %s):", daysText, cityInfo.Name, cityInfo.Adm1, cityInfo.Adm2),
			fmt.Sprintf("Last Updated: %s", indicesData.UpdateTime),
			"",
		}

		for _, index := range indicesData.Daily {
			indexInfo := []string{
				fmt.Sprintf("Date: %s", index.Date),
				fmt.Sprintf("Index Type: %s", index.Name),
				fmt.Sprintf("Level: %s", index.Level),
				fmt.Sprintf("Category: %s", index.Category),
				fmt.Sprintf("Recommendation: %s", index.Text),
				"---",
			}
			indicesText = append(indicesText, strings.Join(indexInfo, "\n"))
		}

		return mcp.NewToolResultText(strings.Join(indicesText, "\n")), nil
	})
}
