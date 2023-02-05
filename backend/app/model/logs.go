package model

type Logs struct {
	AccessedAt string `json:"accessed_at,omitempty"`
	Sequence   string `json:"sequence,omitempty"`
	AccessID   string `json:"access_id,omitempty"`
	AccessName string `json:"access_name,omitempty"`
	IpAddress  string `json:"ip_address,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	Request    string `json:"request,omitempty"`
	Response   string `json:"response,omitempty"`
}

type StatisticsDB struct {
	Request         string `json:"request,omitempty"`
	Count           string `json:"count,omitempty"`
	UniqueUserAgent string `json:"unique_user_agent,omitempty"`
	Response200     string `json:"response_200,omitempty"`
	Response404     string `json:"response_404,omitempty"`
	Response500     string `json:"response_500,omitempty"`
	ResponseETC     string `json:"response_etc,omitempty"`
}

type StatisticsData struct {
	Count           int `json:"count"`
	UniqueUserAgent int `json:"unique_user_agent"`
	Response200     int `json:"response_200"`
	Response404     int `json:"response_404"`
	Response500     int `json:"response_500"`
	ResponseETC     int `json:"response_etc"`
}

type StatisticsGroup struct {
	StatisticsData map[string]StatisticsData `json:"request"`
}
