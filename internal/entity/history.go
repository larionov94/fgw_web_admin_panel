package entity

// HistoryPerformer история данных сотрудника.
type HistoryPerformer struct {
	PerformerId int    `json:"performerId"`
	Hostname    string `json:"hostname"`
	IpAddress   string `json:"ipAddress"`
	TraceId     string `json:"traceId"`
	FIO         string `json:"fio"`
	RoleName    string `json:"roleName"`
	CreatedAt   string `json:"createdAt"`
	CreatedBy   string `json:"createdBy"`
}
