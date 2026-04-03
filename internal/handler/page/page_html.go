package page

import (
	"fgw_web_admin_panel/internal/api/middleware"
	"fgw_web_admin_panel/internal/entity"
	"fgw_web_admin_panel/pkg/convert"
	"fgw_web_admin_panel/pkg/msg"
	"fmt"
	"html/template"
	"net/http"
)

const (
	pathToTmplDefault = "web/templates/html/"
	prefixAFormsTmpl  = "admin/"
	tmplErrorHTML     = "error.html"
)

// ErrorPage страница с описанием ошибки.
type ErrorPage struct {
	Title      string
	ErrMsg     error
	MsgCode    string
	StatusCode int
	Method     string
	Path       string
}

type Page struct {
	Title          string                    // Title название страницы.
	CurrentPage    string                    // CurrentPage ключ страницы для отображения.
	InfoPerformer  *middleware.PerformerData // InfoPerformer информация о сотруднике.
	PerformersList []*entity.Performer       // PerformersList список сотрудников.
	RolesList      []*entity.Role            // RolesList список ролей.
	SectorsList    []*entity.Sector          // SectorsList список участков печей.
}

type DataPage struct {
	Page *Page // Page данные страницы.
}

func NewDataPage(page *Page) *DataPage {
	return &DataPage{Page: page}
}

// RenderPage отображение одной страницы.
func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	pathToTmpl := fmt.Sprintf("%s%s", pathToTmplDefault, tmpl)

	parseTmpl, err := template.New(tmpl).Funcs(template.FuncMap{}).ParseFiles(pathToTmpl)
	if err != nil {
		renderErrorDirectly(w, ErrorPage{
			Title:      "Ошибка",
			ErrMsg:     err,
			MsgCode:    msg.EH5001,
			StatusCode: http.StatusNotFound,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		renderErrorDirectly(w, ErrorPage{
			Title:      "Ошибка",
			ErrMsg:     err,
			MsgCode:    msg.EH5002,
			StatusCode: http.StatusInternalServerError,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}
}

// RenderPages отображение страниц связанные между собой.
func RenderPages(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, addTmpl ...string) {
	pathToTemplates := []string{fmt.Sprintf("%s%s%s", pathToTmplDefault, prefixAFormsTmpl, tmpl)}

	for _, templates := range addTmpl {
		pathToTemplates = append(pathToTemplates, fmt.Sprintf("%s%s%s", pathToTmplDefault, prefixAFormsTmpl, templates))
	}

	parseTmpl, err := template.New(tmpl).Funcs(template.FuncMap{
		"formatDateTime": convert.FormatDateTime,
		"formatDate":     convert.FormatDate,
	}).ParseFiles(pathToTemplates...)
	if err != nil {
		renderErrorDirectly(w, ErrorPage{
			Title:      "Ошибка",
			ErrMsg:     err,
			MsgCode:    msg.EH5001,
			StatusCode: http.StatusNotFound,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		renderErrorDirectly(w, ErrorPage{
			Title:      "Ошибка",
			ErrMsg:     err,
			MsgCode:    msg.EH5002,
			StatusCode: http.StatusInternalServerError,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}
}

// RenderPageError отображает шаблон с описанием ошибки.
func RenderPageError(w http.ResponseWriter, r *http.Request, errPage ErrorPage) {
	data := struct {
		ErrPage ErrorPage
	}{
		ErrPage: errPage,
	}

	RenderPage(w, r, tmplErrorHTML, data)
}

func renderErrorDirectly(w http.ResponseWriter, errPage ErrorPage) {
	w.WriteHeader(errPage.StatusCode)

	errorHTML := fmt.Sprintf(
		`
<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; }
        .error { color: #d9534f; border: 1px solid #d9534f; padding: 15px; margin: 10px 0; }
    </style>
</head>
<body>
    <h1>Ошибка %v</h1>
    <div class="error">
        <strong>Статус код:</strong> %d<br>
        <strong>Сообщение:</strong> %s<br>
        <strong>Метод:</strong> %s<br>
        <strong>Путь:</strong> %s
    </div>
</body>
</html>`, errPage.Title, errPage.ErrMsg, errPage.StatusCode, errPage.MsgCode, errPage.Method, errPage.Path)

	_, err := fmt.Fprintln(w, errorHTML)
	if err != nil {
		return
	}
}
