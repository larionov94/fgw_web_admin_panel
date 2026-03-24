package entity

// Role структура роли.
type Role struct {
	Id          int    `json:"id"`          // Id идентификатор роли.
	NameRole    string `json:"nameRole"`    // NameRole наименование роли.
	Description string `json:"description"` // Description описание роли.
	AuditRec    Audit  `json:"auditRec"`    // AuditRec аудит.
}
