package convert

import (
	"log"
	"strings"
	"time"
)

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

func ParseToMSSQLDateTime(date string) *time.Time {

	date = strings.TrimSpace(date)

	date = strings.Replace(date, "T", " ", 1)

	var t time.Time
	var err error

	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
	}

	for _, layout := range layouts {
		t, err = time.Parse(layout, date)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Printf("Ошибка: [%s] --- поле: [%s]", err.Error(), date)
	}

	return &t
}
