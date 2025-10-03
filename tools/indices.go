package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/overstarry/qweather-mcp-go/api"
)

// WeatherIndicesInput input parameters for get-weather-indices tool
type WeatherIndicesInput struct {
	CityName string `json:"cityName" jsonschema:"Name of the city to query weather life indices for (e.g. UV index, comfort level)"`
	Type     string `json:"type,omitempty" jsonschema:"Index type: 0=all, 1=sports, 2=car wash, 3=clothing, 4=fishing, 5=UV, 6=travel, 7=allergy, 8=cold, 9=comfort, 10=wind, 11=sunglasses, 12=makeup, 13=sunscreen, 14=traffic, 15=sports watching, 16=air pollution diffusion. Default is 0 (all)."`
	Days     string `json:"days,omitempty" jsonschema:"Forecast duration: 1d (today) or 3d (3 days). Defaults to 1d if not specified."`
}

// WeatherIndicesOutput output structure for get-weather-indices tool
type WeatherIndicesOutput struct {
	IndicesInfo string `json:"indicesInfo" jsonschema:"Formatted weather life indices including UV, comfort, clothing suggestions, etc."`
}

// RegisterIndicesTools Register weather indices related tools
func RegisterIndicesTools(s *mcp.Server, client *api.Client) {
	// Weather indices tool
	mcp.AddTool(s, &mcp.Tool{
		Name: "get-weather-indices",
		Description: "Weather life indices forecast API provides various life indices for cities worldwide. Supports 1-day and 3-day forecasts. Available index types:\n\n" +
			"- Type 0: All index types\n" +
			"- Type 1: Sports (indicates suitability for outdoor sports activities)\n" +
			"- Type 2: Car Washing (suggests whether it's suitable to wash cars)\n" +
			"- Type 3: Dressing (provides clothing suggestions based on weather)\n" +
			"- Type 4: Fishing (shows suitability of fishing conditions)\n" +
			"- Type 5: UV (ultraviolet radiation intensity level)\n" +
			"- Type 6: Travel (indicates suitability for travel and sightseeing)\n" +
			"- Type 7: Allergy (allergy and pollen risk level)\n" +
			"- Type 8: Cold (cold risk level)\n" +
			"- Type 9: Comfort (overall comfort level of weather)\n" +
			"- Type 10: Wind (wind conditions and their effects)\n" +
			"- Type 11: Sunglasses (need for wearing sunglasses)\n" +
			"- Type 12: Makeup (weather effects on makeup)\n" +
			"- Type 13: Sunscreen (sunscreen needs)\n" +
			"- Type 14: Traffic (weather effects on traffic conditions)\n" +
			"- Type 15: Sports Watching (suitability for watching outdoor sports)\n" +
			"- Type 16: Air Pollution Diffusion Conditions (air pollution diffusion conditions)\n\n" +
			"Note: Not all cities provide all indices. International cities mainly support types 1, 2, 4, and 5.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input WeatherIndicesInput) (*mcp.CallToolResult, WeatherIndicesOutput, error) {
		if input.CityName == "" {
			return nil, WeatherIndicesOutput{}, fmt.Errorf("city name cannot be empty")
		}

		if input.Type == "" {
			input.Type = "0"
		}

		if input.Days == "" {
			input.Days = "1d"
		}

		// Query city ID
		locationData, err := client.GetLocationByName(input.CityName)
		if err != nil {
			return nil, WeatherIndicesOutput{}, fmt.Errorf("failed to query city: %w", err)
		}

		if locationData.Code != "200" {
			return nil, WeatherIndicesOutput{}, fmt.Errorf("failed to query city, API returned an error")
		}

		if len(locationData.Location) == 0 {
			return nil, WeatherIndicesOutput{}, fmt.Errorf("no matching city found")
		}

		// Use the ID of the first matching city
		cityID := locationData.Location[0].ID
		cityInfo := locationData.Location[0]

		// Get weather indices data
		indicesData, err := client.GetWeatherIndices(cityID, input.Days, input.Type)
		if err != nil {
			return nil, WeatherIndicesOutput{}, fmt.Errorf("failed to get weather indices data: %w", err)
		}

		if indicesData.Code != "200" {
			return nil, WeatherIndicesOutput{}, fmt.Errorf("failed to get weather indices data, API returned an error")
		}

		// Format weather indices information
		daysText := "1-day"
		if input.Days == "3d" {
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

		return nil, WeatherIndicesOutput{IndicesInfo: strings.Join(indicesText, "\n")}, nil
	})
}
