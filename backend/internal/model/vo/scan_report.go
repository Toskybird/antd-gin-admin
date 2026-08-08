package vo

import "antd-gin-admin-backend/internal/model/entity"

// ScanReportVO is the API view of a Scan Report.
type ScanReportVO struct {
	ReportCode string `json:"report_code"`
	JobCode    string `json:"job_code"`
	Version    int    `json:"version"`
	ObjectKey  string `json:"object_key"`
	Format     string `json:"format"`
	CreatedAt  string `json:"created_at"`
}

func BuildScanReportVO(r *entity.ScanReport) *ScanReportVO {
	if r == nil {
		return nil
	}
	return &ScanReportVO{
		ReportCode: r.ReportCode,
		JobCode:    r.JobCode,
		Version:    r.Version,
		ObjectKey:  r.ObjectKey,
		Format:     r.Format,
		CreatedAt:  r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
