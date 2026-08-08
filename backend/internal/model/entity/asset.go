package entity

import (
	"time"

	"gorm.io/gorm"
)

// Asset is a managed Web application entry ledger record.
type Asset struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetCode   string         `json:"asset_code" gorm:"column:asset_code;size:32;uniqueIndex"`
	Name        string         `json:"name" gorm:"size:128"`
	RootURL     string         `json:"root_url" gorm:"column:root_url;size:512"`
	DeptCode    string         `json:"dept_code" gorm:"column:dept_code;size:32;index"`
	Remark      string         `json:"remark" gorm:"size:512"`
	Status      string         `json:"status" gorm:"size:16;index"` // active | disabled
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Asset) TableName() string { return "scan_assets" }
