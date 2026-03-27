package convert

import (
	"log"
	"strconv"
)

// ConvStrToInt конвертировать строку в число.
func ConvStrToInt(str string) (int, error) {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Ошибка: [%s] --- значение: [%v]", err.Error(), value)

		return 0, err
	}

	return value, nil
}
