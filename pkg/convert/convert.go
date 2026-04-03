package convert

import (
	"log"
	"net/http"
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

// ParseFormFieldInt преобразует поле в целое число, полученное из HTTP запроса.
func ParseFormFieldInt(r *http.Request, fieldName string) int {

	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "0"

		return 0
	}
	value, err := strconv.Atoi(formValue)
	if err != nil {
		log.Printf("Ошибка: [%s] --- поле: [%s] --- значение: [%v]", err.Error(), fieldName, value)

		return 0
	}

	return value
}
