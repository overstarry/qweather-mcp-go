package api

// LocationResponse City query response
type LocationResponse struct {
	Code     string     `json:"code"`
	Location []Location `json:"location"`
}

// Location City information
type Location struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
	Adm2    string `json:"adm2"`
	Adm1    string `json:"adm1"`
	Country string `json:"country"`
	Type    string `json:"type"`
	Rank    string `json:"rank"`
}

// WeatherNowResponse Real-time weather response
type WeatherNowResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Text      string `json:"text"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Pressure  string `json:"pressure"`
		Vis       string `json:"vis"`
	} `json:"now"`
}

// WeatherDailyResponse Weather forecast response
type WeatherDailyResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Daily      []struct {
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
	} `json:"daily"`
}

// MinutelyResponse Minutely precipitation forecast response
type MinutelyResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Summary    string `json:"summary"`
	Minutely   []struct {
		FxTime string `json:"fxTime"`
		Precip string `json:"precip"`
		Type   string `json:"type"`
	} `json:"minutely"`
}

// HourlyResponse Hourly weather forecast response
type HourlyResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Hourly     []struct {
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
	} `json:"hourly"`
}

// WarningResponse Weather warning response
type WarningResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Warning    []struct {
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
	} `json:"warning"`
}

// IndicesResponse Weather life indices response
type IndicesResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Daily      []struct {
		Date     string `json:"date"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Level    string `json:"level"`
		Category string `json:"category"`
		Text     string `json:"text"`
	} `json:"daily"`
}

// AirQualityResponse Real-time air quality response
type AirQualityResponse struct {
	Code     string `json:"code"`
	Metadata struct {
		Tag string `json:"tag"`
	} `json:"metadata"`
	Indexes []struct {
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
	} `json:"indexes"`
	Pollutants []struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		FullName      string `json:"fullName"`
		Concentration struct {
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		} `json:"concentration"`
		SubIndexes []struct {
			Code       string `json:"code"`
			Aqi        int    `json:"aqi"`
			AqiDisplay string `json:"aqiDisplay"`
		} `json:"subIndexes,omitempty"`
	} `json:"pollutants"`
	Stations []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"stations,omitempty"`
}

// AirQualityHourlyResponse Hourly air quality forecast response
type AirQualityHourlyResponse struct {
	Code     string `json:"code"`
	Metadata struct {
		Tag string `json:"tag"`
	} `json:"metadata"`
	Hours []struct {
		ForecastTime string `json:"forecastTime"`
		Indexes      []struct {
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
		} `json:"indexes"`
		Pollutants []struct {
			Code          string `json:"code"`
			Name          string `json:"name"`
			FullName      string `json:"fullName"`
			Concentration struct {
				Value float64 `json:"value"`
				Unit  string  `json:"unit"`
			} `json:"concentration"`
			SubIndexes []struct {
				Code       string `json:"code"`
				Aqi        int    `json:"aqi"`
				AqiDisplay string `json:"aqiDisplay"`
			} `json:"subIndexes,omitempty"`
		} `json:"pollutants"`
	} `json:"hours"`
}

// AirQualityDailyResponse Daily air quality forecast response
type AirQualityDailyResponse struct {
	Code     string `json:"code"`
	Metadata struct {
		Tag string `json:"tag"`
	} `json:"metadata"`
	Days []struct {
		ForecastStartTime string `json:"forecastStartTime"`
		ForecastEndTime   string `json:"forecastEndTime"`
		Indexes           []struct {
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
		} `json:"indexes"`
		Pollutants []struct {
			Code          string `json:"code"`
			Name          string `json:"name"`
			FullName      string `json:"fullName"`
			Concentration struct {
				Value float64 `json:"value"`
				Unit  string  `json:"unit"`
			} `json:"concentration"`
			SubIndexes []struct {
				Code       string `json:"code"`
				Aqi        int    `json:"aqi"`
				AqiDisplay string `json:"aqiDisplay"`
			} `json:"subIndexes,omitempty"`
		} `json:"pollutants"`
	} `json:"days"`
}
