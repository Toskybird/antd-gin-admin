package entity

import (
	"time"

	"gorm.io/gorm"
)

// Finding is one detection result produced by a Scan Job.
type Finding struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	FindingCode string         `json:"finding_code" gorm:"column:finding_code;size:32;uniqueIndex"`
	JobCode     string         `json:"job_code" gorm:"column:job_code;size:32;index"`
	RuleCode    string         `json:"rule_code" gorm:"column:rule_code;size:128;index"`
	Severity    string         `json:"severity" gorm:"size:16;index"` // critical|high|medium|low|info
	Title       string         `json:"title" gorm:"size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Evidence    string         `json:"evidence" gorm:"type:text"`
	Location    string         `json:"location" gorm:"size:512"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Finding) TableName() string { return "scan_findings" }
