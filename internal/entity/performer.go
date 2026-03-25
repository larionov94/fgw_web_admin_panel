package entity

import "time"

// Performer структура сотрудника.
type Performer struct {
	Id            int           `json:"id"`            // Id - идентификатор.
	SectorId      int           `json:"sectorId"`      // SectorId - идентификатор участка печки.
	FIO           string        `json:"fio"`           // FIO - фамилия, имя, отчество.
	TabNum        int           `json:"tabNum"`        // TabNum - табельный номер.
	Barcode       string        `json:"barcode"`       // Barcode - код пропуска. (указан на карточке пропуска).
	AccessBarcode *string       `json:"accessBarcode"` // AccessBarcode - код пропуска для доступа к сканеру.
	Passwd        string        `json:"passwd"`        // Passwd - пароль.
	IssuedAt      *time.Time    `json:"issuedAt"`      // IssuedAt - дата выдачи допуска к сканеру.
	Archive       bool          `json:"archive"`       // Archive - архив.
	PerformerRole PerformerRole `json:"performerRole"`
	AuditRec      Audit         `json:"auditRec"`    // AuditRec - аудит.
	AuthSuccess   bool          `json:"authSuccess"` // AuthSuccess авторизация.
}

// PerformerRole роли сотрудника.
type PerformerRole struct {
	RoleIdAForms   int    `json:"roleIdAForms"`   // RoleIdAForms - идентификатор роли для доступа к AForms.
	RoleIdAFGW     int    `json:"roleIdAFGW"`     // RoleIdAFGW - идентификатор роли для доступа к AFGW. (ТЛК).
	RoleNameAForms string `json:"roleNameAForms"` // RoleNameAForms наименование роли для доступа к AForms.
	RoleNameAFGW   string `json:"roleNameAFGW"`   // RoleNameAFGW наименование роли для доступа к AFGW.
}

// PerformerAuth данные авторизации сотрудника.
type PerformerAuth struct {
	Success   bool      `json:"success"`   // Success авторизация true/false.
	Performer Performer `json:"performer"` // Performer данные сотрудника только при Success=true.
	Message   string    `json:"message"`   // Message текстовое сообщение о результате.
}
