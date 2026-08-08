package dto

// ListFindingRequest filters Finding list.
type ListFindingRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	JobCode  string `form:"job_code"`
	RuleCode string `form:"rule_code"`
	Severity string `form:"severity"`
	Keyword  string `form:"keyword"`
}
