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
