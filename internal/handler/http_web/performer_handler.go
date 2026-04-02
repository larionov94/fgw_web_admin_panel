package http_web

import (
	"fgw_web_admin_panel/internal/api/middleware"
	"fgw_web_admin_panel/internal/handler/page"
	"fgw_web_admin_panel/internal/service"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"net/http"
)

const (
	tmplPerformerHTML = "performers.html"

	titlePerformersPage = "Список сотрудников"
	pagePerformerPanel  = "performers"
)

type PerformerHandler struct {
	performerService *service.PerformerService
	logg             *logg.Logger
	authMiddleware   *middleware.AuthMiddleware
}

func NewPerformerHandler(performerService *service.PerformerService, logg *logg.Logger, authMiddleware *middleware.AuthMiddleware) *PerformerHandler {
	return &PerformerHandler{performerService, logg, authMiddleware}
}

func (p *PerformerHandler) ServeHTTPRouter(mux *http.ServeMux) {
	mux.HandleFunc("/admin/performers", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRoleForAForms([]int{0}, p.ShowPerformersPage)))
}

func (p *PerformerHandler) ShowPerformersPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5006,
			StatusCode: http.StatusMethodNotAllowed,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	performers, err := p.performerService.AllPerformers(r.Context())
	if err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5007,
			StatusCode: http.StatusInternalServerError,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	performerData, err := p.authMiddleware.GetPerformerData(r, p.performerService)
	if err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5008,
			StatusCode: http.StatusInternalServerError,
			Method:     r.Method,
			Path:       r.URL.Path,
		})
	}

	data := page.NewDataPage(&page.Page{
		Title:          titlePerformersPage,
		CurrentPage:    pagePerformerPanel,
		InfoPerformer:  performerData,
		PerformersList: performers,
	})

	page.RenderPages(w, r, tmplStartPageHTML, data, tmplPerformerHTML)
}
