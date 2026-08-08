package vo

import "antd-gin-admin-backend/internal/model/entity"

// AssetVO is the API view of an Asset.
type AssetVO struct {
	AssetCode string `json:"asset_code"`
	Name      string `json:"name"`
	RootURL   string `json:"root_url"`
	DeptCode  string `json:"dept_code"`
	Remark    string `json:"remark"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

func BuildAssetVO(a *entity.Asset) *AssetVO {
	if a == nil {
		return nil
	}
	return &AssetVO{
		AssetCode: a.AssetCode,
		Name:      a.Name,
		RootURL:   a.RootURL,
		DeptCode:  a.DeptCode,
		Remark:    a.Remark,
		Status:    a.Status,
		UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
