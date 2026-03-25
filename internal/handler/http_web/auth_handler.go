package http_web

import (
	"fgw_web_admin_panel/internal/service"
	"fgw_web_admin_panel/pkg/logg"
	"html/template"
	"net/http"
)

type AuthHandler struct {
	performerService service.PerformerService
	logg             *logg.Logger
}

func NewAuthHandler(performerService service.PerformerService, logg *logg.Logger) *AuthHandler {
	return &AuthHandler{performerService, logg}
}

func (a *AuthHandler) ServeHTTPRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", a.StartPage)
}

func (a *AuthHandler) StartPage(w http.ResponseWriter, r *http.Request) {
	parseHTML, err := template.ParseFiles("web/templates/html/auth.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	parseHTML.Execute(w, a)
}
