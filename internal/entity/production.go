package entity

type Production struct {
	Id                 int     `json:"id"`                 // Id ид продукции.
	NameOfPacking      string  `json:"nameOfPacking"`      // NameOfPacking наименование с комментариями для упаковщика.
	Designation        string  `json:"designation"`        // Designation обозначение конструкторского наименования.
	NameOfLabel        string  `json:"nameOfLabel"`        // NameOfLabel наименование для печати на этикетке.
	IsFood             bool    `json:"isFood"`             // IsFood пищевая/не пищевая.
	IsDeclared         bool    `json:"isDeclared"`         // IsDeclared декларированная/не декларированная.
	IsSun              bool    `json:"isSun"`              // IsSun беречь от солнца.
	IsUmbrella         bool    `json:"isUmbrella"`         // IsUmbrella беречь от влаги.
	IsParty            bool    `json:"isParty"`            // IsParty партионная/не партионная.
	IsPerfumery        bool    `json:"isPerfume"`          // IsPerfumery парфюмерия/не парфюмерия.
	Color              string  `json:"color"`              // Color цвет продукции.
	GL                 int     `json:"gl"`                 // GL петля Мёбиуса.
	Article            string  `json:"article"`            // Article артикул продукции.
	SapCode            string  `json:"sapCode"`            // SapCode САП код.
	ItemsPerRows       int     `json:"itemsPerRows"`       // ItemsPerRows кол-во в ряду.
	RowsPerPack        int     `json:"rowsPerPack"`        // RowsPerPack кол-во рядов.
	PalletWeight       float64 `json:"palletWeight"`       // PalletWeight вес упаковки.
	HWD                string  `json:"HWD"`                // HWD высота х ширина х глубина.
	Comm               string  `json:"comm"`               // Comm комментарий.
	BatchNumber        int     `json:"batchNumber"`        // BatchNumber номер текущей партии.
	LabelDate          string  `json:"labelDate"`          // LabelDate дата на этикетку выпуска партии.
	BatchRealDate      string  `json:"batchRealDate"`      // BatchRealDate реальная дата создания, выпуск партии.
	BatchNumberingMode int     `json:"batchNumberingMode"` // BatchNumberingMode нумерация партии даты (1-руч. 2-авт. 3-с указанной даты).
	ExpiryMonths       int     `json:"expiryMonths"`       // ExpiryMonths срок годности.
	VP                 int     `json:"vp"`                 // VP номер ванной печи.
	ML                 int     `json:"ml"`                 // ML номер линии печи..
	Archive            bool    `json:"archive"`            // Archive архив/не архив.
	ProductType        string  `json:"productType"`        // ProductType тип продукции "декларированная".
	Status             bool    `json:"status"`             // Status статус продукции.
	AuditRecord        Audit   `json:"auditRecord"`        // AuditRecord аудит.
}
