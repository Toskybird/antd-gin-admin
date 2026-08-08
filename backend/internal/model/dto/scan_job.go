package dto

// CreateScanJobRequest creates a Scan Job from asset or temporary URL.
type CreateScanJobRequest struct {
	Source    string `json:"source" binding:"required"` // asset | temp_url
	AssetCode string `json:"asset_code"`
	EntryURL  string `json:"entry_url"`
	Policy    string `json:"policy" binding:"required"` // quick|standard|deep
	MaxDepth  int    `json:"max_depth"`
	MaxPages  int    `json:"max_pages"`
}

// ListScanJobRequest filters scan job list.
type ListScanJobRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Status   string `form:"status"`
	Policy   string `form:"policy"`
	Keyword  string `form:"keyword"`
}
