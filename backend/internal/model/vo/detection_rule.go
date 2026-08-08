package vo

import "antd-gin-admin-backend/internal/model/entity"

// DetectionRuleVO is the API view of a Detection Rule.
type DetectionRuleVO struct {
	RuleCode    string `json:"rule_code"`
	DisplayName string `json:"display_name"`
	Enabled     bool   `json:"enabled"`
	UpdatedAt   string `json:"updated_at"`
}

func BuildDetectionRuleVO(r *entity.DetectionRule) *DetectionRuleVO {
	if r == nil {
		return nil
	}
	return &DetectionRuleVO{
		RuleCode:    r.RuleCode,
		DisplayName: r.DisplayName,
		Enabled:     r.Enabled,
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
