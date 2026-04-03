package entity

type Sector struct {
	Id         int    `json:"id"`         // Id идентификатор печки.
	NameSector string `json:"nameSector"` // NameSector наименование участка печки.
	VpMlSector string `json:"vpMlSector"` // VpMlSector описание линий на печке.
	AuditRec   Audit  `json:"auditRec"`   // AuditRec аудит.
}
