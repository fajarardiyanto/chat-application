package response

type OperationalHoursResponse struct {
	IsOperationalHours bool            `json:"is_operational_hours"`
	OooMessage         string          `json:"ooo_message"`
	OperationalTime    OperationalTime `json:"operational_time"`
}

type OperationalTime struct {
	OpenHour    int32 `json:"open_hour"`
	OpenMinute  int32 `json:"open_minute"`
	CloseHour   int32 `json:"close_hour"`
	CloseMinute int32 `json:"close_minute"`
	OpenAllDay  bool  `json:"open_all_day"`
	CloseAllDay bool  `json:"close_all_day"`
}
