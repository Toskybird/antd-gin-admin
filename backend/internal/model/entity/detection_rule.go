package entity

import (
	"time"

	"gorm.io/gorm"
)

// DetectionRule is a system-synced check catalog entry.
type DetectionRule struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	RuleCode    string         `json:"rule_code" gorm:"column:rule_code;size:128;uniqueIndex"`
	DisplayName string         `json:"display_name" gorm:"column:display_name;size:255"`
	Enabled     bool           `json:"enabled" gorm:"not null;default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (DetectionRule) TableName() string { return "scan_detection_rules" }
