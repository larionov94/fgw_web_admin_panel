package entity

import "time"

// Audit аудит для отслеживания изменений данных
type Audit struct {
	CreatedAt time.Time `json:"createdAt"` // CreatedAt - дата создания записи.
	CreatedBy int       `json:"createdBy"` // CreatedBy - табельный номер сотрудника.
	UpdatedAt time.Time `json:"updatedAt"` // UpdatedAt - дата изменения записи.
	UpdatedBy int       `json:"updatedBy"` // UpdatedBy - табельный номер сотрудника изменивший запись.
}
