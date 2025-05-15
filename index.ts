import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const QWEATHER_API_BASE = process.env.QWEATHER_API_BASE
if (!QWEATHER_API_BASE) {
    throw new Error("QWEATHER_API_BASE env is not set")
}
const QWEATHER_API_KEY = process.env.QWEATHER_API_KEY
if (!QWEATHER_API_KEY) {
    throw new Error("QWEATHER_API_KEY env is not set")
}

// Create server instance
const server = new McpServer({
    name: "qweather",
    version: "1.0.0",
    capabilities: {
        resources: {},
        tools: {},
    },
});

async function makeQWeatherRequest<T>(endpoint: string, params: Record<string, string>, pathParams?: string[]): Promise<T | null> {
    let url: URL;
    if (pathParams && pathParams.length > 0) {
        // Replace placeholders in endpoint with actual values
        let endpointWithParams = endpoint;
        pathParams.forEach(param => {
            endpointWithParams = endpointWithParams.replace('{}', param);
        });
        url = new URL(`${QWEATHER_API_BASE}${endpointWithParams}`);
    } else {
        url = new URL(`${QWEATHER_API_BASE}${endpoint}`);
    }

    Object.entries(params).forEach(([key, value]) => {
        url.searchParams.append(key, value);
    });

    try {
        const response = await fetch(url.toString(), {
            headers: {
                'X-QW-Api-Key': QWEATHER_API_KEY as string
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}, URL: ${url.toString()}`);
        }
        return await response.json() as T;
    } catch (error) {
        console.error("Error making QWeather request:", error, "URL:", url.toString());
        return null;
    }
}

interface QWeatherNowResponse {
    code: string;
    updateTime: string;
    now: {
        obsTime: string;
        temp: string;
        feelsLike: string;
        text: string;
        windDir: string;
        windScale: string;
        humidity: string;
        precip: string;
        pressure: string;
        vis: string;
    };
}

interface QWeatherDailyResponse {
    code: string;
    updateTime: string;
    fxLink: string;
    daily: Array<{
        fxDate: string;
        sunrise: string;
        sunset: string;
        moonrise: string;
        moonset: string;
        moonPhase: string;
        moonPhaseIcon: string;
        tempMax: string;
        tempMin: string;
        iconDay: string;
        textDay: string;
        iconNight: string;
        textNight: string;
        wind360Day: string;
        windDirDay: string;
        windScaleDay: string;
        windSpeedDay: string;
        wind360Night: string;
        windDirNight: string;
        windScaleNight: string;
        windSpeedNight: string;
        humidity: string;
        precip: string;
        pressure: string;
        vis: string;
        cloud: string;
        uvIndex: string;
    }>;
}

interface QWeatherLocationResponse {
    code: string;
    location: Array<{
        name: string;
        id: string;
        lat: string;
        lon: string;
        adm2: string;
        adm1: string;
        country: string;
        type: string;
        rank: string;
    }>;
}

interface QWeatherWarningResponse {
    code: string;
    updateTime: string;
    fxLink: string;
    warning: Array<{
        id: string;
        sender: string;
        pubTime: string;
        title: string;
        startTime: string;
        endTime: string;
        status: string;
        severity: string;
        severityColor: string;
        type: string;
        typeName: string;
        urgency: string;
        certainty: string;
        text: string;
        related: string;
    }>;
}

interface QWeatherIndicesResponse {
    code: string;
    updateTime: string;
    fxLink: string;
    daily: Array<{
        date: string;
        type: string;
        name: string;
        level: string;
        category: string;
        text: string;
    }>;
}

interface QWeatherHourlyResponse {
    code: string;
    updateTime: string;
    fxLink: string;
    hourly: Array<{
        fxTime: string;
        temp: string;
        icon: string;
        text: string;
        wind360: string;
        windDir: string;
        windScale: string;
        windSpeed: string;
        humidity: string;
        precip: string;
        pressure: string;
        cloud: string;
        dew: string;
    }>;
}

interface QWeatherMinutelyResponse {
    code: string;
    updateTime: string;
    fxLink: string;
    summary: string;
    minutely: Array<{
        fxTime: string;
        precip: string;
        type: string;
    }>;
}

interface QWeatherAirQualityResponse {
    code: string;
    metadata: {
        tag: string;
    };
    indexes: Array<{
        code: string;
        name: string;
        aqi: number;
        aqiDisplay: string;
        level?: string;
        category?: string;
        color: {
            red: number;
            green: number;
            blue: number;
            alpha: number;
        };
        primaryPollutant?: {
            code: string;
            name: string;
            fullName: string;
        } | null;
        health?: {
            effect: string;
            advice: {
                generalPopulation: string;
                sensitivePopulation: string;
            };
        };
    }>;
    pollutants: Array<{
        code: string;
        name: string;
        fullName: string;
        concentration: {
            value: number;
            unit: string;
        };
        subIndexes?: Array<{
            code: string;
            aqi: number;
            aqiDisplay: string;
        }>;
    }>;
    stations?: Array<{
        id: string;
        name: string;
    }>;
}

interface QWeatherAirQualityHourlyResponse {
    code: string;
    metadata: {
        tag: string;
    };
    hours: Array<{
        forecastTime: string;
        indexes: Array<{
            code: string;
            name: string;
            aqi: number;
            aqiDisplay: string;
            level?: string;
            category?: string;
            color: {
                red: number;
                green: number;
                blue: number;
                alpha: number;
            };
            primaryPollutant?: {
                code: string;
                name: string;
                fullName: string;
            } | null;
            health?: {
                effect: string;
                advice: {
                    generalPopulation: string;
                    sensitivePopulation: string;
                };
            };
        }>;
        pollutants: Array<{
            code: string;
            name: string;
            fullName: string;
            concentration: {
                value: number;
                unit: string;
            };
            subIndexes?: Array<{
                code: string;
                aqi: number;
                aqiDisplay: string;
            }>;
        }>;
    }>;
}

interface QWeatherAirQualityDailyResponse {
    code: string;
    metadata: {
        tag: string;
    };
    days: Array<{
        forecastStartTime: string;
        forecastEndTime: string;
        indexes: Array<{
            code: string;
            name: string;
            aqi: number;
            aqiDisplay: string;
            level?: string;
            category?: string;
            color: {
                red: number;
                green: number;
                blue: number;
                alpha: number;
            };
            primaryPollutant?: {
                code: string;
                name: string;
                fullName: string;
            } | null;
            health?: {
                effect: string;
                advice: {
                    generalPopulation: string;
                    sensitivePopulation: string;
                };
            };
        }>;
        pollutants: Array<{
            code: string;
            name: string;
            fullName: string;
            concentration: {
                value: number;
                unit: string;
            };
            subIndexes?: Array<{
                code: string;
                aqi: number;
                aqiDisplay: string;
            }>;
        }>;
    }>;
}

server.tool(
    "get-weather-now",
    "Real-time weather API provides current weather conditions for global cities. Available data includes: temperature, feels-like temperature, weather conditions, wind direction, wind force scale, relative humidity, precipitation, atmospheric pressure, and visibility. The data is updated in real-time to provide the most accurate current weather information.",
    {
        cityName: z.string().describe("Name of the city to look up weather for"),
    },
    async ({ cityName }) => {
        // First, look up the city to get its ID
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's ID
        const cityId = locationData.location[0].id;
        const cityInfo = locationData.location[0];

        const weatherData = await makeQWeatherRequest<QWeatherNowResponse>("/v7/weather/now", {
            location: cityId,
        });

        if (!weatherData || weatherData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve current weather data",
                    },
                ],
            };
        }

        const now = weatherData.now;
        const weatherText = [
            `Current Weather for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            `Temperature: ${now.temp}°C (Feels like: ${now.feelsLike}°C)`,
            `Condition: ${now.text}`,
            `Wind: ${now.windDir} Scale ${now.windScale}`,
            `Humidity: ${now.humidity}%`,
            `Precipitation: ${now.precip}mm`,
            `Pressure: ${now.pressure}hPa`,
            `Visibility: ${now.vis}km`,
            `Last Updated: ${weatherData.updateTime}`,
        ].join("\n");

        return {
            content: [
                {
                    type: "text",
                    text: weatherText,
                },
            ],
        };
    }
);

server.tool(
    "get-weather-forecast",
    "Weather forecast API provides detailed weather predictions for global cities, supporting forecasts from 3 to 30 days. Available data includes: sunrise/sunset times, moonrise/moonset times, temperature range, weather conditions, wind direction and speed, relative humidity, precipitation, atmospheric pressure, cloud cover, and UV index. The forecast is updated daily to ensure accuracy.",
    {
        cityName: z.string().describe("Name of the city to look up weather for"),
        days: z.enum(["3d", "7d", "10d", "15d", "30d"]).describe("Number of forecast days"),
    },
    async ({ cityName, days }) => {
        // First, look up the city to get its ID
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's ID
        const cityId = locationData.location[0].id;
        const cityInfo = locationData.location[0];

        const weatherData = await makeQWeatherRequest<QWeatherDailyResponse>(`/v7/weather/${days}`, {
            location: cityId,
        });

        if (!weatherData || weatherData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve weather forecast data",
                    },
                ],
            };
        }

        const forecastText = [
            `${days.replace('d', ' Days')} Weather Forecast for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            `Last Updated: ${weatherData.updateTime}`,
            '',
            ...weatherData.daily.map(day => [
                `Date: ${day.fxDate}`,
                `Temperature: ${day.tempMin}°C ~ ${day.tempMax}°C`,
                `Daytime: ${day.textDay}`,
                `Night: ${day.textNight}`,
                `Sunrise: ${day.sunrise || 'N/A'}  Sunset: ${day.sunset || 'N/A'}`,
                `Precipitation: ${day.precip}mm`,
                `Humidity: ${day.humidity}%`,
                `Wind: Day-${day.windDirDay}(Scale ${day.windScaleDay}), Night-${day.windDirNight}(Scale ${day.windScaleNight})`,
                `UV Index: ${day.uvIndex}`,
                '---'
            ].join('\n'))
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: forecastText,
                },
            ],
        };
    }
);

server.tool(
    "get-minutely-precipitation",
    "Minute-level precipitation forecast API provides precise precipitation predictions for the next 2 hours in global cities. Available data includes precipitation type (rain/snow) and precipitation amount for each minute. This high-precision forecast is particularly useful for outdoor activity planning and real-time weather monitoring.",
    {
        cityName: z.string().describe("Name of the city to look up precipitation forecast for"),
    },
    async ({ cityName }) => {
        // First, look up the city to get its location coordinates
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's coordinates
        const cityInfo = locationData.location[0];
        const location = `${cityInfo.lon},${cityInfo.lat}`;

        const precipData = await makeQWeatherRequest<QWeatherMinutelyResponse>("/v7/minutely/5m", {
            location: location,
        });

        if (!precipData || precipData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve precipitation forecast data",
                    },
                ],
            };
        }

        const precipText = [
            `Minutely Precipitation Forecast - ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            `Forecast Description: ${precipData.summary}`,
            `Last Updated: ${precipData.updateTime}`,
            '',
            '2-Hour Precipitation Forecast:',
            ...precipData.minutely.map(minute =>
                `Time: ${minute.fxTime.split('T')[1].split('+')[0]} - ${minute.type === 'rain' ? 'Rain' : 'Snow'}: ${minute.precip}mm`
            ),
            '',
            `Data Source: ${precipData.fxLink}`,
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: precipText,
                },
            ],
        };
    }
);

server.tool(
    "get-hourly-forecast",
    "Hourly weather forecast API provides detailed weather information for global cities for the next 24-168 hours. Available data includes: temperature, weather conditions, wind force, wind speed, wind direction, relative humidity, atmospheric pressure, precipitation probability, dew point temperature, and cloud cover. The forecast data is updated hourly to ensure accuracy.",
    {
        cityName: z.string().describe("Name of the city to look up weather for"),
        hours: z.enum(["24h", "72h", "168h"]).default("24h").describe("Number of forecast hours (24h, 72h, or 168h)"),
    },
    async ({ cityName, hours }) => {
        // First, look up the city to get its ID
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's ID
        const cityId = locationData.location[0].id;
        const cityInfo = locationData.location[0];

        const hourlyData = await makeQWeatherRequest<QWeatherHourlyResponse>(`/v7/weather/${hours}`, {
            location: cityId,
        });

        if (!hourlyData || hourlyData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve hourly forecast data",
                    },
                ],
            };
        }

        const hourlyText = [
            `${hours.replace('h', '-Hour')} Weather Forecast - ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            `Last Updated: ${hourlyData.updateTime}`,
            '',
            ...hourlyData.hourly.map(hour => [
                `Time: ${hour.fxTime.split('T')[1].split('+')[0]}`,
                `Temperature: ${hour.temp}°C`,
                `Weather: ${hour.text}`,
                `Wind: ${hour.windDir} (Scale ${hour.windScale}, ${hour.windSpeed}km/h)`,
                `Humidity: ${hour.humidity}%`,
                `Precipitation: ${hour.precip}mm`,
                `Pressure: ${hour.pressure}hPa`,
                hour.cloud ? `Cloud Cover: ${hour.cloud}%` : null,
                hour.dew ? `Dew Point: ${hour.dew}°C` : null,
                '---'
            ].filter(Boolean).join('\n'))
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: hourlyText,
                },
            ],
        };
    }
);

server.tool(
    "get-weather-warning",
    "Weather Warning API provides real-time weather warning data issued by official authorities in China and multiple countries/regions worldwide. The data includes warning issuer, publication time, warning title, detailed warning information, warning level, warning type, and other relevant information.",
    {
        cityName: z.string().describe("Name of the city to look up weather warnings for"),
    },
    async ({ cityName }) => {
        // First, look up the city to get its ID
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's ID
        const cityId = locationData.location[0].id;
        const cityInfo = locationData.location[0];

        const warningData = await makeQWeatherRequest<QWeatherWarningResponse>("/v7/warning/now", {
            location: cityId,
        });

        if (!warningData || warningData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve weather warning data",
                    },
                ],
            };
        }

        if (!warningData.warning || warningData.warning.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: `No active weather warnings for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2})`,
                    },
                ],
            };
        }

        const warningText = [
            `Weather Warnings for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            `Last Updated: ${warningData.updateTime}`,
            '',
            ...warningData.warning.map(warning => [
                `Warning Title: ${warning.title}`,
                `Issued By: ${warning.sender}`,
                `Publication Time: ${warning.pubTime}`,
                `Warning Type: ${warning.typeName}`,
                `Severity Level: ${warning.severity} (${warning.severityColor})`,
                `Valid Period: ${warning.startTime || 'Not specified'} to ${warning.endTime || 'Not specified'}`,
                `Status: ${warning.status}`,
                `Details: ${warning.text}`,
                '---'
            ].join('\n'))
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: warningText,
                },
            ],
        };
    }
);

server.tool(
    "get-weather-indices",
    "Weather indices forecast API provides various life indices for cities worldwide. Supports both 1-day and 3-day forecasts. Available indices types:\n\n" +
    "- Type 0: All indices types\n" +
    "- Type 1: Sport (Indicates suitability for outdoor sports activities)\n" +
    "- Type 2: Car Wash (Suggests whether it's suitable for washing cars)\n" +
    "- Type 3: Dressing (Provides clothing recommendations based on weather)\n" +
    "- Type 4: Fishing (Shows how favorable conditions are for fishing)\n" +
    "- Type 5: UV (Ultraviolet radiation intensity level)\n" +
    "- Type 6: Travel (Indicates suitability for traveling and sightseeing)\n" +
    "- Type 7: Allergy (Risk level for allergies and pollen)\n" +
    "- Type 8: Cold Risk (Risk level for catching a cold)\n" +
    "- Type 9: Comfort (Overall comfort level of the weather)\n" +
    "- Type 10: Wind (Wind conditions and their effects)\n" +
    "- Type 11: Sunglasses (Need for wearing sunglasses)\n" +
    "- Type 12: Makeup (How weather affects makeup wear)\n" +
    "- Type 13: Sunscreen (Need for sun protection)\n" +
    "- Type 14: Traffic (Weather impact on traffic conditions)\n" +
    "- Type 15: Sports Spectating (Suitability for watching outdoor sports)\n" +
    "- Type 16: Air Quality (Air pollution diffusion conditions)\n\n" +
    "Note: Not all indices are available for every city. International cities mainly support types 1, 2, 4, and 5.",
    {
        cityName: z.string().describe("Name of the city to look up weather indices for"),
        type: z.enum([
            "0", // All Types
            "1", // Sport
            "2", // Car Wash
            "3", // Dressing
            "4", // Fishing
            "5", // UV
            "6", // Travel
            "7", // Allergy
            "8", // Cold Risk
            "9", // Comfort
            "10", // Wind
            "11", // Sunglasses
            "12", // Makeup
            "13", // Sunscreen
            "14", // Traffic
            "15", // Sports Spectating
            "16", // Air Quality
        ]).describe("Type of weather index to retrieve"),
        days: z.enum(["1d", "3d"]).default("1d").describe("Number of forecast days (1d or 3d)"),
    },
    async ({ cityName, type, days }) => {
        // First, look up the city to get its ID
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's ID
        const cityId = locationData.location[0].id;
        const cityInfo = locationData.location[0];

        const indicesData = await makeQWeatherRequest<QWeatherIndicesResponse>(`/v7/indices/${days}`, {
            location: cityId,
            type: type,
        });

        if (!indicesData || indicesData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve weather indices data",
                    },
                ],
            };
        }

        const indicesText = [
            `${days === "1d" ? "1-Day" : "3-Day"} Weather Indices for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            `Last Updated: ${indicesData.updateTime}`,
            '',
            ...indicesData.daily.map(index => [
                `Date: ${index.date}`,
                `Index Type: ${index.name}`,
                `Level: ${index.level}`,
                `Category: ${index.category}`,
                `Suggestion: ${index.text}`,
                '---'
            ].join('\n'))
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: indicesText,
                },
            ],
        };
    }
);

server.tool(
    "get-air-quality",
    "Real-time Air Quality API provides air quality data for specific locations with 1x1 kilometer precision. Includes AQI based on local standards of different countries/regions, AQI levels, colors, primary pollutants, QWeather universal AQI, pollutant concentrations, sub-indices, health advice, and related monitoring station information.",
    {
        cityName: z.string().describe("Name of the city to look up air quality for"),
    },
    async ({ cityName }) => {
        // First, look up the city to get its coordinates
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's coordinates
        const cityInfo = locationData.location[0];
        // Format coordinates to have at most 2 decimal places
        const lat = Number(cityInfo.lat).toFixed(2);
        const lon = Number(cityInfo.lon).toFixed(2);

        // Update API endpoint to use path parameters (latitude first, then longitude)
        const airQualityEndpoint = `/airquality/v1/current/${lat}/${lon}`;
        const airQualityData = await makeQWeatherRequest<QWeatherAirQualityResponse>(
            airQualityEndpoint,
            {}, // No query parameters needed
            [lat, lon] // Path parameters in correct order: latitude, longitude
        );

        if (!airQualityData || !airQualityData.indexes || airQualityData.indexes.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve air quality data",
                    },
                ],
            };
        }

        // Format output text
        const airQualityText = [
            `Real-time Air Quality for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            '',
            'Air Quality Indices:',
            ...airQualityData.indexes.map(index => [
                `${index.name}: ${index.aqiDisplay}`,
                index.level ? `Level: ${index.level}` : null,
                index.category ? `Category: ${index.category}` : null,
                index.primaryPollutant ? `Primary Pollutant: ${index.primaryPollutant.fullName}` : null,
                index.health ? [
                    'Health Effects:',
                    `- ${index.health.effect}`,
                    'Health Advice:',
                    `- General Population: ${index.health.advice.generalPopulation}`,
                    `- Sensitive Population: ${index.health.advice.sensitivePopulation}`,
                ].join('\n') : null,
                '---'
            ].filter(Boolean).join('\n')),
            '',
            'Pollutant Concentrations:',
            ...airQualityData.pollutants.map(pollutant =>
                `${pollutant.fullName}: ${pollutant.concentration.value}${pollutant.concentration.unit}`
            ),
            '',
            airQualityData.stations ? [
                'Related Monitoring Stations:',
                ...airQualityData.stations.map(station => `- ${station.name}`),
            ].join('\n') : null,
        ].filter(Boolean).join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: airQualityText,
                },
            ],
        };
    }
);

server.tool(
    "get-air-quality-hourly",
    "Hourly Air Quality Forecast API provides air quality data for the next 24 hours, including AQI, pollutant concentrations, sub-indices, and health advice. The data includes various air quality standards (such as QAQI, GB-DEFRA, etc.) and specific concentrations of pollutants like PM2.5, PM10, NO2, O3, SO2.",
    {
        cityName: z.string().describe("Name of the city to look up air quality forecast for"),
    },
    async ({ cityName }) => {
        // First, look up the city to get its coordinates
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's coordinates
        const cityInfo = locationData.location[0];
        // Format coordinates to have at most 2 decimal places
        const lat = Number(cityInfo.lat).toFixed(2);
        const lon = Number(cityInfo.lon).toFixed(2);

        // Use path parameters to call the hourly air quality forecast API (latitude first, then longitude)
        const airQualityHourlyEndpoint = `/airquality/v1/hourly/${lat}/${lon}`;
        const airQualityData = await makeQWeatherRequest<QWeatherAirQualityHourlyResponse>(
            airQualityHourlyEndpoint,
            {}, // No query parameters needed
            [lat, lon] // Path parameters in order: latitude, longitude
        );

        if (!airQualityData || !airQualityData.hours || airQualityData.hours.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve air quality forecast data",
                    },
                ],
            };
        }

        // Format output text
        const hourlyText = [
            `24-Hour Air Quality Forecast for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            '',
            ...airQualityData.hours.map(hour => {
                const time = new Date(hour.forecastTime).toLocaleString('en-US', {
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false,
                    timeZone: 'UTC'
                }) + ' UTC';

                const indexInfo = hour.indexes.map(index => {
                    const healthInfo = index.health ? [
                        `Health Effects: ${translateHealthEffect(index.health.effect)}`,
                        `Health Advice:`,
                        `  General Population: ${translateAdvice(index.health.advice.generalPopulation)}`,
                        `  Sensitive Population: ${translateAdvice(index.health.advice.sensitivePopulation)}`
                    ].join('\n') : '';

                    return [
                        `Air Quality Indices:`,
                        `  ${index.name}: ${index.aqiDisplay}`,
                        `  Level: ${index.level || 'N/A'}`,
                        `  Category: ${translateCategory(index.category || 'Unknown')}`,
                        index.primaryPollutant ? `  Primary Pollutant: ${translatePollutant(index.primaryPollutant.fullName)}` : '',
                        healthInfo
                    ].filter(Boolean).join('\n');
                }).join('\n');

                const pollutantInfo = hour.pollutants.length > 0 ? [
                    'Pollutant Concentrations:',
                    ...hour.pollutants.map(pollutant =>
                        `  ${translatePollutant(pollutant.fullName)}: ${pollutant.concentration.value}${pollutant.concentration.unit}`
                    )
                ].join('\n') : 'No pollutant data available';

                return [
                    `Forecast Time: ${time}`,
                    indexInfo,
                    pollutantInfo,
                    '---'
                ].join('\n\n');
            }).join('\n')
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: hourlyText,
                },
            ],
        };
    }
);

server.tool(
    "get-air-quality-daily",
    "Daily Air Quality Forecast API provides air quality predictions for the next 3 days, including AQI values, pollutant concentrations, and health recommendations. The data includes various air quality standards and specific concentrations of pollutants like PM2.5, PM10, NO2, O3, SO2.",
    {
        cityName: z.string().describe("Name of the city to look up air quality forecast for"),
    },
    async ({ cityName }) => {
        // First, look up the city to get its coordinates
        const locationData = await makeQWeatherRequest<QWeatherLocationResponse>("/geo/v2/city/lookup", {
            location: cityName,
        });

        if (!locationData || locationData.code !== "200") {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to find the specified city",
                    },
                ],
            };
        }

        if (!locationData.location || locationData.location.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "No matching city found",
                    },
                ],
            };
        }

        // Use the first matching city's coordinates
        const cityInfo = locationData.location[0];
        // Format coordinates to have at most 2 decimal places
        const lat = Number(cityInfo.lat).toFixed(2);
        const lon = Number(cityInfo.lon).toFixed(2);

        // Use path parameters to call the daily air quality forecast API
        const airQualityDailyEndpoint = `/airquality/v1/daily/${lat}/${lon}`;
        const airQualityData = await makeQWeatherRequest<QWeatherAirQualityDailyResponse>(
            airQualityDailyEndpoint,
            {}, // No query parameters needed
            [lat, lon] // Path parameters in order: latitude, longitude
        );

        if (!airQualityData || !airQualityData.days || airQualityData.days.length === 0) {
            return {
                content: [
                    {
                        type: "text",
                        text: "Failed to retrieve air quality forecast data",
                    },
                ],
            };
        }

        // Format output text
        const dailyText = [
            `3-Day Air Quality Forecast for ${cityInfo.name} (${cityInfo.adm1} ${cityInfo.adm2}):`,
            '',
            ...airQualityData.days.map(day => {
                const startTime = new Date(day.forecastStartTime).toLocaleString('en-US', {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false,
                    timeZone: 'UTC'
                }) + ' UTC';

                const endTime = new Date(day.forecastEndTime).toLocaleString('en-US', {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false,
                    timeZone: 'UTC'
                }) + ' UTC';

                const indexInfo = day.indexes.map(index => {
                    const healthInfo = index.health ? [
                        `Health Effects: ${index.health.effect}`,
                        `Health Advice:`,
                        `  General Population: ${index.health.advice.generalPopulation}`,
                        `  Sensitive Population: ${index.health.advice.sensitivePopulation}`
                    ].join('\n') : '';

                    return [
                        `Air Quality Index:`,
                        `  ${index.name}: ${index.aqiDisplay}`,
                        `  Level: ${index.level || 'N/A'}`,
                        `  Category: ${index.category || 'Unknown'}`,
                        index.primaryPollutant ? `  Primary Pollutant: ${index.primaryPollutant.fullName}` : '',
                        healthInfo
                    ].filter(Boolean).join('\n');
                }).join('\n\n');

                const pollutantInfo = day.pollutants.length > 0 ? [
                    'Pollutant Concentrations:',
                    ...day.pollutants.map(pollutant =>
                        `  ${pollutant.fullName}: ${pollutant.concentration.value}${pollutant.concentration.unit}`
                    )
                ].join('\n') : 'No pollutant data available';

                return [
                    `Forecast Period: ${startTime} to ${endTime}`,
                    indexInfo,
                    pollutantInfo,
                    '---'
                ].join('\n\n');
            }).join('\n')
        ].join('\n');

        return {
            content: [
                {
                    type: "text",
                    text: dailyText,
                },
            ],
        };
    }
);

// Helper functions for translation
function translateCategory(category: string): string {
    const categoryMap: Record<string, string> = {
        '优': 'Excellent',
        '良': 'Good',
        '轻度污染': 'Light Pollution',
        '中度污染': 'Moderate Pollution',
        '重度污染': 'Heavy Pollution',
        '严重污染': 'Severe Pollution'
    };
    return categoryMap[category] || category;
}

function translatePollutant(pollutant: string): string {
    const pollutantMap: Record<string, string> = {
        '颗粒物（粒径小于等于2.5μm）': 'PM2.5',
        '颗粒物（粒径小于等于10μm）': 'PM10',
        '二氧化氮': 'NO2',
        '臭氧': 'O3',
        '二氧化硫': 'SO2',
        '一氧化碳': 'CO'
    };
    return pollutantMap[pollutant] || pollutant;
}

function translateHealthEffect(effect: string): string {
    const effectMap: Record<string, string> = {
        '空气质量令人满意，基本无空气污染。': 'Air quality is satisfactory with minimal air pollution.',
        '空气质量可接受，但某些污染物可能对极少数异常敏感人群健康有较弱影响。': 'Air quality is acceptable, but some pollutants may have a slight impact on the health of extremely sensitive individuals.',
        '易感人群症状有轻度加剧，健康人群出现刺激症状。': 'Sensitive individuals may experience mild symptom aggravation, while healthy individuals may experience irritation symptoms.',
        '进一步加剧易感人群症状，可能对健康人群心脏、呼吸系统有影响。': 'Further aggravation of symptoms in sensitive individuals, possible effects on the cardiovascular and respiratory systems of healthy individuals.',
        '心脏病和肺病患者症状显著加剧，运动耐受力降低，健康人群普遍出现症状。': 'Significant aggravation of symptoms in patients with heart and lung conditions, reduced exercise tolerance, and general symptoms in healthy individuals.',
        '健康人群运动耐受力降低，有明显强烈症状，提前出现某些疾病。': 'Reduced exercise tolerance in healthy individuals, obvious and severe symptoms, early onset of certain diseases.'
    };
    return effectMap[effect] || effect;
}

function translateAdvice(advice: string): string {
    const adviceMap: Record<string, string> = {
        '各类人群可正常活动。': 'All groups can maintain normal activities.',
        '一般人群可正常活动。': 'General population can maintain normal activities.',
        '极少数异常敏感人群应减少户外活动。': 'Extremely sensitive individuals should reduce outdoor activities.',
        '儿童、老年人及心脏病、呼吸系统疾病患者应减少长时间、高强度的户外锻炼。': 'Children, elderly, and individuals with heart or respiratory conditions should reduce prolonged, high-intensity outdoor exercise.',
        '儿童、老年人及心脏病、呼吸系统疾病患者应停留在室内，停止户外运动，一般人群减少户外运动。': 'Children, elderly, and individuals with heart or respiratory conditions should stay indoors and avoid outdoor activities. General population should reduce outdoor activities.',
        '儿童、老年人和病人应停留在室内，避免体力消耗，一般人群避免户外活动。': 'Children, elderly, and patients should stay indoors and avoid physical exertion. General population should avoid outdoor activities.'
    };
    return adviceMap[advice] || advice;
}

async function main() {
    const transport = new StdioServerTransport();
    await server.connect(transport);
    console.error("Weather MCP Server running on stdio");
}

main().catch((error) => {
    console.error("Fatal error in main():", error);
    process.exit(1);
});