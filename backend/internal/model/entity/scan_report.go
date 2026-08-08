package entity

import (
	"time"

	"gorm.io/gorm"
)

// ScanReport is one generated HTML report version for a Scan Job.
type ScanReport struct {
	ID         int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ReportCode string         `json:"report_code" gorm:"column:report_code;size:32;uniqueIndex"`
	JobCode    string         `json:"job_code" gorm:"column:job_code;size:32;index"`
	Version    int            `json:"version" gorm:"column:version;not null"`
	ObjectKey  string         `json:"object_key" gorm:"column:object_key;size:512"`
	Format     string         `json:"format" gorm:"size:16;not null;default:html"`
	CreatedBy  string         `json:"created_by" gorm:"column:created_by;size:32"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (ScanReport) TableName() string { return "scan_reports" }
