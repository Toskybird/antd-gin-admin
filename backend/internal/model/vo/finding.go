package vo

import "antd-gin-admin-backend/internal/model/entity"

// FindingVO is the read-only API view of a Finding.
type FindingVO struct {
	FindingCode     string `json:"finding_code"`
	JobCode         string `json:"job_code"`
	RuleCode        string `json:"rule_code"`
	RuleDisplayName string `json:"rule_display_name"`
	Severity        string `json:"severity"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Evidence        string `json:"evidence"`
	Location        string `json:"location"`
	CreatedAt       string `json:"created_at"`
}

func BuildFindingVO(f *entity.Finding, ruleDisplayName string) *FindingVO {
	if f == nil {
		return nil
	}
	return &FindingVO{
		FindingCode:     f.FindingCode,
		JobCode:         f.JobCode,
		RuleCode:        f.RuleCode,
		RuleDisplayName: ruleDisplayName,
		Severity:        f.Severity,
		Title:           f.Title,
		Description:     f.Description,
		Evidence:        f.Evidence,
		Location:        f.Location,
		CreatedAt:       f.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
