package toplogger

import "time"

// Climbs

type pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type gym struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	IconPath    string       `json:"iconPath"`
	ClimbGroups []climbGroup `json:"climbGroups"`
}

type climbGroup struct {
	ID    string `json:"id"`
	Color string `json:"color"`
}

type holdColor struct {
	ID    string `json:"id"`
	Color string `json:"color"`
}

type climbGroupClimb struct {
	ID           string `json:"id"`
	ClimbGroupID string `json:"climbGroupId"`
}

type climb struct {
	ID               string            `json:"id"`
	Grade            int               `json:"grade"`
	ClimbType        string            `json:"climbType"`
	HoldColor        holdColor         `json:"holdColor"`
	ClimbGroupClimbs []climbGroupClimb `json:"climbGroupClimbs"`
}

type climbUserDaysBoulder struct {
	ID        string `json:"id"`
	TickType  int    `json:"tickType"` // 1 == top, 2 == flash
	WasRepeat bool   `json:"wasRepeat"`
	Climb     climb  `json:"climb"`
}

type climbDay struct {
	ID                    string                 `json:"id"`
	StatsAtDate           string                 `json:"statsAtDate"`
	Gym                   gym                    `json:"gym"`
	ClimbUserDaysBoulders []climbUserDaysBoulder `json:"climbUserDaysBoulders"`
}

type climbDayPaginated struct {
	Pagination pagination `json:"pagination"`
	Data       []climbDay `json:"data"`
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
