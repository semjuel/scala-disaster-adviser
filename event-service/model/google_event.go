package model

type CalendarResponse struct {
	Items []CalendarEvent `json:"items"`
}

type CalendarEvent struct {
	Uuid     int64
	Id       string   `json:"id"`
	Summary  string   `json:"summary"`
	Location string   `json:"location"`
	Start    DateTime `json:"start"`
	End      DateTime `json:"end"`
	User     User
}

type DateTime struct {
	DateTime string `json:"dateTime"`
	Date     string `json:"date"`
}
