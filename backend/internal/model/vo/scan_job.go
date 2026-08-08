package vo

import "antd-gin-admin-backend/internal/model/entity"

// ScanJobVO is the API view of a Scan Job.
type ScanJobVO struct {
	JobCode      string  `json:"job_code"`
	AssetCode    string  `json:"asset_code,omitempty"`
	EntryURL     string  `json:"entry_url"`
	DeptCode     string  `json:"dept_code"`
	Policy       string  `json:"policy"`
	MaxDepth     int     `json:"max_depth"`
	MaxPages     int     `json:"max_pages"`
	Status       string  `json:"status"`
	FindingCount int     `json:"finding_count"`
	TempURL      bool    `json:"temp_url"`
	CreatedAt    string  `json:"created_at"`
	FinishedAt   *string `json:"finished_at,omitempty"`
}

func BuildScanJobVO(j *entity.ScanJob) *ScanJobVO {
	if j == nil {
		return nil
	}
	vo := &ScanJobVO{
		JobCode:      j.JobCode,
		AssetCode:    j.AssetCode,
		EntryURL:     j.EntryURL,
		DeptCode:     j.DeptCode,
		Policy:       j.Policy,
		MaxDepth:     j.MaxDepth,
		MaxPages:     j.MaxPages,
		Status:       j.Status,
		FindingCount: j.FindingCount,
		TempURL:      j.TempURL,
		CreatedAt:    j.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if j.FinishedAt != nil {
		s := j.FinishedAt.Format("2006-01-02 15:04:05")
		vo.FinishedAt = &s
	}
	return vo
}
