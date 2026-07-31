package toplogger

import "time"

// Climbs

type pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type gym struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
}

type holdColor struct {
	ID    string `json:"id"`
	Color string `json:"color"`
}

type climb struct {
	ID        string    `json:"id"`
	Grade     int       `json:"grade"`
	ClimbType string    `json:"climbType"`
	HoldColor holdColor `json:"holdColor"`
}

type climbDay struct {
	ID          string `json:"id"`
	StatsAtDate string `json:"statsAtDate"`
	Gym         gym    `json:"gym"`
}

type climbDayPaginated struct {
	Pagination pagination `json:"pagination"`
	Data       []climbDay `json:"data"`
}

type climbLog struct {
	TickType  int   `json:"tickType"`  // Semantics not fully documented; 1 == top (observed)
	TickIndex int   `json:"tickIndex"` // -1 == skip, > 0 == repeat
	Climb     climb `json:"climb"`
}

type climbLogPaginated struct {
	Pagination pagination `json:"pagination"`
	Data       []climbLog `json:"data"`
}

// Tokens

type accessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type refreshToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Token struct {
	Access  accessToken  `json:"access"`
	Refresh refreshToken `json:"refresh"`
}

// Error

type originalError struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

type extension struct {
	Code          string        `json:"code"`
	OriginalError originalError `json:"originalError"`
}

type cError struct {
	Message   string    `json:"message"`
	Extension extension `json:"extensions"`
}
