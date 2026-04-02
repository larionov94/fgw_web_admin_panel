package convert

import "time"

// FormatDateTime форматирует дату в формате ДД.ММ.ГГГГ ЧЧ:ММ.
func FormatDateTime(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.Format("02.01.2006 15:04:05")
}

// FormatDate форматирует дату в формате ДД.ММ.ГГГГ.
func FormatDate(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.Format("02.01.2006")
}
