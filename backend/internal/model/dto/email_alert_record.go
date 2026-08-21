package dto

// EmailAlertRecordPageRequest filters 告警记录 pagination.
type EmailAlertRecordPageRequest struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	SendStatus string `form:"send_status"`
	Module     string `form:"module"`
	StartTime  *int64 `form:"start_time"`
	EndTime    *int64 `form:"end_time"`
}
