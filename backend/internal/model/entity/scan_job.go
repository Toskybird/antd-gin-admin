package entity

import (
	"time"

	"gorm.io/gorm"
)

// ScanJob is one scan execution unit against a single scan target.
type ScanJob struct {
	ID              int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	JobCode         string         `json:"job_code" gorm:"column:job_code;size:32;uniqueIndex"`
	AssetCode       string         `json:"asset_code" gorm:"column:asset_code;size:32;index"`
	EntryURL        string         `json:"entry_url" gorm:"column:entry_url;size:512"`
	DeptCode        string         `json:"dept_code" gorm:"column:dept_code;size:32;index"`
	CreatedBy       string         `json:"created_by" gorm:"column:created_by;size:32;index"`
	Policy          string         `json:"policy" gorm:"size:16"` // quick | standard | deep
	MaxDepth        int            `json:"max_depth" gorm:"column:max_depth"`
	MaxPages        int            `json:"max_pages" gorm:"column:max_pages"`
	Status          string         `json:"status" gorm:"size:16;index"` // queued|running|succeeded|failed|cancelled
	FindingCount    int            `json:"finding_count" gorm:"column:finding_count;not null;default:0"`
	TempURL         bool           `json:"temp_url" gorm:"column:temp_url;not null;default:false"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	FinishedAt      *time.Time     `json:"finished_at" gorm:"column:finished_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (ScanJob) TableName() string { return "scan_jobs" }
