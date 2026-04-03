package http_web

import (
	"errors"
	"fgw_web_admin_panel/internal/api/middleware"
	"fgw_web_admin_panel/internal/entity"
	"fgw_web_admin_panel/internal/handler/page"
	"fgw_web_admin_panel/internal/service"
	"fgw_web_admin_panel/pkg/convert"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"net/http"
)

const (
	tmplPerformerHTML    = "performers.html"
	tmplPerformerUpdHTML = "performer_upd.html"

	titlePerformersPage = "Список сотрудников"
	pagePerformerPanel  = "performers"

	titlePerformerUpdPage = "Редактирование сотрудника"
	pagePerformerUpdPanel = "performer_upd"

	urlPerformersHtml = "/admin/performers"
)

type PerformerHandler struct {
	performerService service.PerformerUseCase
	roleService      service.RoleUseCase
	logg             *logg.Logger
	authMiddleware   *middleware.AuthMiddleware
}

func NewPerformerHandler(performerService service.PerformerUseCase, roleService service.RoleUseCase, logg *logg.Logger, authMiddleware *middleware.AuthMiddleware) *PerformerHandler {
	return &PerformerHandler{performerService, roleService, logg, authMiddleware}
}

func (p *PerformerHandler) ServeHTTPRouter(mux *http.ServeMux) {
	mux.HandleFunc("/admin/performers", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRoleForAForms([]int{0}, p.ShowPerformersPage)))
	mux.HandleFunc("/admin/performers/upd", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRoleForAForms([]int{0}, p.UpdPerformerPage)))
}

// ShowPerformersPage отображает список сотрудников на веб странице.
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

	page.RenderPages(w, r, tmplStartPageHTML, data, tmplPerformerHTML, tmplPerformerUpdHTML)
}

func (p *PerformerHandler) UpdPerformerPage(w http.ResponseWriter, r *http.Request) {
	performerData, err := p.authMiddleware.GetPerformerData(r, p.performerService)
	if err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5008,
			StatusCode: http.StatusInternalServerError,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	switch r.Method {
	case http.MethodGet:
		p.getUpdPerformerPage(w, r, performerData)
	case http.MethodPost:
		p.postUpdPerformerPage(w, r, performerData)
	default:
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5006,
			StatusCode: http.StatusMethodNotAllowed,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}
}

func (p *PerformerHandler) getUpdPerformerPage(w http.ResponseWriter, r *http.Request, performerData *middleware.PerformerData) {
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

	performer, err := p.ensurePerformerExist(r, "id")
	if err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5012,
			StatusCode: http.StatusInternalServerError,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	data := page.NewDataPage(&page.Page{
		Title:          titlePerformerUpdPage,
		CurrentPage:    pagePerformerUpdPanel,
		InfoPerformer:  performerData,
		PerformersList: []*entity.Performer{performer},
	})

	page.RenderPages(w, r, tmplStartPageHTML, data, tmplPerformerHTML, tmplPerformerUpdHTML)
}

func (p *PerformerHandler) postUpdPerformerPage(w http.ResponseWriter, r *http.Request, performerData *middleware.PerformerData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodPost {
		page.RenderPageError(w, r, page.ErrorPage{
			MsgCode:    msg.EH5006,
			StatusCode: http.StatusMethodNotAllowed,
			Method:     r.Method,
			Path:       r.URL.Path,
		})

		return
	}

	if err := r.ParseForm(); err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			ErrMsg:     err,
			MsgCode:    msg.EH5001,
			StatusCode: http.StatusBadRequest,
		})

		return
	}

	idPerformer, err := p.extractIdParamFromValue(r, "id")
	if err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			ErrMsg:     err,
			MsgCode:    msg.EH5009,
			StatusCode: http.StatusBadRequest,
		})

		return
	}

	fieldAccessBarcode := r.FormValue("AccessBarcode")
	fieldIssuedAt := convert.ParseToMSSQLDateTime(r.FormValue("IssuedAt"))

	if fieldAccessBarcode == "" {
		fieldIssuedAt = nil
	}

	performer := &entity.Performer{
		AccessBarcode: &fieldAccessBarcode,
		PerformerRole: entity.PerformerRole{
			RoleIdAForms: convert.ParseFormFieldInt(r, "RoleIdAForms"),
			RoleIdAFGW:   convert.ParseFormFieldInt(r, "RoleIdAFGW"),
		},
		IssuedAt: fieldIssuedAt,
		AuditRec: entity.Audit{
			UpdatedBy: performerData.PerformerTabNum,
		},
	}

	if err = p.performerService.UpdPerformer(r.Context(), idPerformer, performer); err != nil {
		page.RenderPageError(w, r, page.ErrorPage{
			ErrMsg:     err,
			MsgCode:    msg.EH5010,
			StatusCode: http.StatusInternalServerError,
		})

		return
	}

	http.Redirect(w, r, urlPerformersHtml, http.StatusSeeOther)
}

// ensurePerformerExist обеспечивает существование продукции.
func (p *PerformerHandler) ensurePerformerExist(r *http.Request, paramId string) (*entity.Performer, error) {
	idPerformer, err := p.extractIdParamFromValue(r, paramId)
	if err != nil {
		p.logg.LogE(msg.EH5011, err, logg.SkipNofS)

		return nil, err
	}

	performer, err := p.performerService.FindPerformerById(r.Context(), idPerformer)
	if err != nil {
		p.logg.LogE("", err, logg.SkipNofS)

		return nil, err
	}

	if performer == nil {
		p.logg.LogW(msg.WH4001, logg.SkipNofS)

		return nil, errors.New(msg.WH4001)
	}

	return performer, nil
}

// extractIdParamFromValue извлекает параметр идентификатора из формы.
func (p *PerformerHandler) extractIdParamFromValue(r *http.Request, paramId string) (int, error) {
	idParamStr := r.FormValue(paramId)
	if idParamStr == "" {
		p.logg.LogW(msg.WH4000, logg.SkipNofS)

		return 0, nil
	}

	idParam, err := convert.ConvStrToInt(idParamStr)
	if err != nil {
		return 0, err
	}

	return idParam, nil
}
