package dto

// CreateAssetRequest creates a new Asset.
type CreateAssetRequest struct {
	Name     string `json:"name" binding:"required"`
	RootURL  string `json:"root_url" binding:"required"`
	DeptCode string `json:"dept_code" binding:"required"`
	Remark   string `json:"remark"`
}

// UpdateAssetRequest updates an Asset.
type UpdateAssetRequest struct {
	Name     *string `json:"name"`
	RootURL  *string `json:"root_url"`
	DeptCode *string `json:"dept_code"`
	Remark   *string `json:"remark"`
	Status   *string `json:"status"`
}

// ListAssetRequest filters Asset list.
type ListAssetRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}
