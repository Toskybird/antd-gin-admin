package dto

// SaveEmailAlertRequest is the PUT body for the singleton 邮件告警约定.
type SaveEmailAlertRequest struct {
	Enabled       bool     `json:"enabled"`
	Host          string   `json:"host"`
	Port          int      `json:"port"`
	Encryption    string   `json:"encryption"`
	Username      string   `json:"username"`
	Password      *string  `json:"password"`
	SenderName    string   `json:"sender_name"`
	SenderAddress string   `json:"sender_address"`
	Recipients    []string `json:"recipients"`
}
