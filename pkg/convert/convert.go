package convert

import (
	"log"
	"strconv"
)

// ConvStrToInt конвертировать строку в число.
func ConvStrToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [  ] --- значение: [%v]", err.Error(), value)

		return 0
	}

	return value
}
