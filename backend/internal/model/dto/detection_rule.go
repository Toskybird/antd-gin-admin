package dto

// ListDetectionRuleRequest filters detection rule list.
type ListDetectionRuleRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Enabled  *bool  `form:"enabled"`
}

// UpdateDetectionRuleStatusRequest toggles enablement.
type UpdateDetectionRuleStatusRequest struct {
	Enabled bool `json:"enabled"`
}
