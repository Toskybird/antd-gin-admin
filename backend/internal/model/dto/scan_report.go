package dto

// CreateScanReportRequest generates a new report version for a job.
type CreateScanReportRequest struct {
	JobCode string `json:"job_code" binding:"required"`
}

// ListScanReportRequest filters scan report list.
type ListScanReportRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	JobCode  string `form:"job_code"`
}
