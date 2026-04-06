package http_web

import (
	"fgw_web_aforms_panel/internal/api"
	"fgw_web_aforms_panel/internal/api/middleware"
	"fgw_web_aforms_panel/internal/entity"
	"fgw_web_aforms_panel/internal/handler/page"
	"fgw_web_aforms_panel/internal/service"
	"fgw_web_aforms_panel/pkg/convert"
	"fgw_web_aforms_panel/pkg/logg"
	"fgw_web_aforms_panel/pkg/msg"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-uuid"
)

const (
	tmplRedirectHTML = "redirect.html"
	tmplAuthHTML     = "auth.html"

	urlAforms             = "/aforms"
	urlAuth               = "/auth"
	urlLogin              = "/login"
	urlLogoutTempRedirect = "/logout-temp-redirect"
	urlTempRedirect       = "/temp-redirect"
	tmplStartPageHTML     = "index.html"

	exitMsg  = "Выход"
	entryMsg = "Вход"

	titleAformsPanelPage = "Панель форм-комплектов"
	pageAformsPanel      = "dashboard"
)

var UUIDString string

const (
	RedirectDelayFast    = 100  // 0.1 секунда
	RedirectDelayNormal  = 300  // 0.3 секунды
	FallbackDelayDefault = 3000 // 3 секунды
)

type RedirectData struct {
	Title           string
	Message         string
	NoScriptMessage string
	TargetURL       string
	CurrentURL      string
	TempURL         string
	Delay           int
	FallbackDelay   int
	ClearHistory    bool
	AddTempState    bool
}

type AuthHandler struct {
	performerService service.PerformerUseCase
	historyService   service.HistoryUseCase
	logg             *logg.Logger
	authMiddleware   *middleware.AuthMiddleware
}

func NewAuthHandler(performerService service.PerformerUseCase, historyService service.HistoryUseCase, logg *logg.Logger, authMiddleware *middleware.AuthMiddleware) *AuthHandler {
	return &AuthHandler{performerService, historyService, logg, authMiddleware}
}

func (a *AuthHandler) ServeHTTPRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", a.ShowAuthForm)
	mux.HandleFunc("/login", a.LoginPage)
	mux.HandleFunc("/auth", a.AuthPerformerHTML)
	mux.HandleFunc("/logout", a.Logout)
	mux.HandleFunc("/aforms", a.authMiddleware.RequireAuth(a.authMiddleware.RequireRoleForAForms([]int{0}, a.StartPage)))
}

func (a *AuthHandler) AuthPerformerHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

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

	performerIdStr := r.FormValue("performerTabNum")
	performerPass := r.FormValue("performerPassword")

	if performerIdStr == "" || performerPass == "" {
		http.Error(w, "PerformerId or PerformerPass is empty", http.StatusBadRequest)
		page.RenderPageError(w, r, page.ErrorPage{
			ErrMsg:     nil,
			MsgCode:    msg.EH5003,
			StatusCode: http.StatusUnauthorized,
		})

		return
	}

	performerId, err := convert.ConvStrToInt(performerIdStr)
	if err != nil {
		a.logg.LogE(msg.EH5004, err, logg.SkipNofS)

		return
	}

	authResult, err := a.performerService.AuthPerformerWithData(r.Context(), performerId, performerPass)
	if err != nil {
		if authResult != nil && !authResult.Success {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
		} else {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(""), http.StatusFound)
		}

		return
	}

	if authResult.Success {
		err := a.authMiddleware.CreateAuthSecuritySession(w, r, &middleware.PerformerData{
			PerformerFIO:          authResult.Performer.FIO,
			PerformerTabNum:       performerId,
			PerformerRoleAForms:   authResult.Performer.PerformerRole.RoleNameAForms,
			PerformerRoleAFormsId: authResult.Performer.PerformerRole.RoleIdAForms,
		})

		if err != nil {
			a.authMiddleware.SetSecurityHeaders(w)

			page.RenderPageError(w, r, page.ErrorPage{
				ErrMsg:     err,
				MsgCode:    msg.EH5005,
				StatusCode: http.StatusUnauthorized,
			})

			return
		}

		UUIDStr, err := uuid.GenerateUUID()
		if err != nil {
			return
		}

		SetUUIDStr(UUIDStr)

		if err := a.historyService.AddHistoryOfEntryAndExit(r.Context(), &entity.HistoryPerformer{
			PerformerId: authResult.Performer.TabNum,
			Hostname:    a.logg.HostName(),
			IpAddress:   a.logg.IPAddr(),
			TraceId:     UUIDStr,
			FIO:         authResult.Performer.FIO,
			RoleName:    authResult.Performer.PerformerRole.RoleNameAForms,
			EntryExit:   entryMsg,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:   authResult.Performer.TabNum,
		}); err != nil {
			a.logg.LogE("", err, logg.SkipNofS)

			return

		}
		a.sendLoginSuccessPage(w, r)
	} else {
		http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
	}
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := api.Store.Get(r, api.GetSessionName())
	if err != nil {
		//a.sendLogoutPageWithHistoryClear(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	performerID, ok := session.Values[api.GetPerformerKey()].(int)
	if !ok {
		performerID = 0
	}

	fio, _ := session.Values[api.GetPerformerFIOKey()].(string)
	roleName, _ := session.Values[api.GetRoleAFormsNameKey()].(string)

	history := &entity.HistoryPerformer{
		PerformerId: performerID,
		Hostname:    a.logg.HostName(),
		IpAddress:   a.logg.IPAddr(),
		TraceId:     UUIDString,
		FIO:         fio,
		RoleName:    roleName,
		EntryExit:   exitMsg,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:   performerID,
	}

	if err := a.historyService.AddHistoryOfEntryAndExit(r.Context(), history); err != nil {
		a.logg.LogE("", err, logg.SkipNofS)
	}

	if token, ok := session.Values[api.GetTokenKey()].(string); ok {
		if mw, ok := interface{}(a.authMiddleware).(interface{ RemoveSessionToken(token string) }); ok {
			mw.RemoveSessionToken(token)
		}
	}

	for key := range session.Values {
		delete(session.Values, key)
	}

	session.Options.MaxAge = -1
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode

	if err = session.Save(r, w); err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     api.GetSessionName(),
		Value:    "",
		Path:     api.GetPathToDefault(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	a.logg.LogI(msg.ISS203, logg.SkipNofS)

	a.sendLogoutPageWithHistoryClear(w, r)
}

func (a *AuthHandler) ShowAuthForm(w http.ResponseWriter, r *http.Request) {
	session, err := api.Store.Get(r, api.GetSessionName())
	if err == nil {
		if auth, ok := session.Values[api.GetAuthPerformer()].(bool); ok && auth {
			a.safeRedirectBasedOnRole(w, r)

			return
		}
	}

	a.LoginPage(w, r)
}

func (a *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	a.authMiddleware.SetSecurityHeaders(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	errMsg := r.URL.Query().Get("error")
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errMsg,
	}

	page.RenderPage(w, r, tmplAuthHTML, data)
}

// Обновленный sendLoginSuccessPage.
func (a *AuthHandler) sendLoginSuccessPage(w http.ResponseWriter, r *http.Request) {
	target := urlAforms

	data := RedirectData{
		Title:           "Успешный вход",
		Message:         "Вход выполнен успешно. Выполняется безопасное перенаправление...",
		NoScriptMessage: "Включите JavaScript для безопасного перехода.",
		TargetURL:       target,
		CurrentURL:      urlAuth,
		TempURL:         urlLogoutTempRedirect,
		Delay:           RedirectDelayNormal,
		FallbackDelay:   2000,
		ClearHistory:    true,
		AddTempState:    true,
	}

	a.renderRedirectPage(w, r, data)
}

// safeRedirectBasedOnRole с использованием общего шаблона
func (a *AuthHandler) safeRedirectBasedOnRole(w http.ResponseWriter, r *http.Request) {
	target := urlAforms

	data := RedirectData{
		Title:           "Перенаправление",
		Message:         "Вы уже авторизованы. Выполняется безопасное перенаправление...",
		NoScriptMessage: "Включите JavaScript для безопасного перехода.",
		TargetURL:       target,
		CurrentURL:      r.URL.Path,
		TempURL:         urlTempRedirect,
		Delay:           RedirectDelayFast,
		FallbackDelay:   FallbackDelayDefault,
		ClearHistory:    true,
		AddTempState:    false,
	}

	a.renderRedirectPage(w, r, data)
}

// Обновленный sendLogoutPageWithHistoryClear
func (a *AuthHandler) sendLogoutPageWithHistoryClear(w http.ResponseWriter, r *http.Request) {
	data := RedirectData{
		Title:           "Выход из системы",
		Message:         "Вы успешно вышли из системы. Выполняется безопасное перенаправление на страницу входа...",
		NoScriptMessage: "Включите JavaScript для безопасного выхода.",
		TargetURL:       urlLogin,
		CurrentURL:      r.URL.Path,
		TempURL:         urlLogoutTempRedirect,
		Delay:           RedirectDelayNormal,
		FallbackDelay:   FallbackDelayDefault,
		ClearHistory:    true,
		AddTempState:    true,
	}

	a.renderRedirectPage(w, r, data)
}

func (a *AuthHandler) renderRedirectPage(w http.ResponseWriter, r *http.Request, data RedirectData) {
	if data.Title == "" {
		data.Title = "Перенаправление"
	}
	if data.Message == "" {
		data.Message = "Выполняется безопасное перенаправление..."
	}
	if data.NoScriptMessage == "" {
		data.NoScriptMessage = "Включите JavaScript для безопасного перехода."
	}
	if data.CurrentURL == "" {
		data.CurrentURL = r.URL.Path
	}
	if data.Delay == 0 {
		data.Delay = RedirectDelayNormal
	}
	if data.FallbackDelay == 0 {
		data.FallbackDelay = FallbackDelayDefault
	}

	a.authMiddleware.SetSecurityHeaders(w)

	page.RenderPage(w, r, tmplRedirectHTML, data)
}

func (a *AuthHandler) StartPage(w http.ResponseWriter, r *http.Request) {
	performerData, err := a.authMiddleware.GetPerformerData(r, a.performerService)
	if err != nil {
		a.sendLogoutPageWithHistoryClear(w, r)

		return
	}

	data := page.NewDataPage(&page.Page{
		Title:         titleAformsPanelPage,
		CurrentPage:   pageAformsPanel,
		InfoPerformer: performerData,
	})

	page.RenderPages(w, r, tmplStartPageHTML, data)

	return
}

// SetUUIDStr устанавливает текущий trace_id
func SetUUIDStr(uuidStr string) {
	UUIDString = uuidStr
}

// GetUUIDStr возвращает текущий trace_id.
func GetUUIDStr() string {
	return UUIDString
}
