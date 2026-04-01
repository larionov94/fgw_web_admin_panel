package middleware

import (
	"errors"
	"fgw_web_admin_panel/internal/api"
	"fgw_web_admin_panel/internal/service"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

const (
	expiresCache        = 60 * time.Minute // expiresCache время жизни кэша.
	maxAgeClearSession  = -1               // maxAgeClearSession устанавливает время жизни сессии.
	maxAgeSession       = 15 * time.Minute //  maxAgeSession устанавливает максимальное время жизни сессии.
	prefixTmpl          = "web/templates/html/"
	tmplForceLogoutHTML = "force_logout.html"
)

type UserSession struct {
	PerformerFIO   string    // PerformerFIO фио сотрудника.
	RoleAFormsName string    // RoleAFormsName наименование роли.
	Expires        time.Time // Expires срок действия истекает.
}

type PerformerData struct {
	PerformerFIO          string
	PerformerTabNum       int
	PerformerRoleAForms   string
	PerformerRoleAFormsId int
}

type AuthMiddleware struct {
	store             *sessions.CookieStore // store сохраняем сеансы, используя безопасные cookie.
	sessionName       string                // sessionName наименование сеанса.
	performerKey      string                // performerKey ключ сеанса пользователя.
	performerFIO      string                // performerFIO фио сотрудника.
	roleAFormsKey     string                // roleAFormsKey ключ роли сотрудника для AForms.
	roleAFormsNameKey string                // roleAFormsNameKey ключ для имени роли AForms.
	logg              *logg.Logger          // logg логирование.
	userCache         map[int]*UserSession  // userCache кэш пользователя.
	cacheMu           sync.RWMutex          // cacheMu блокировка взаимного исключения чтения и записи.
}

func NewAuthMiddleware(store *sessions.CookieStore, logg *logg.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		store:             store,
		sessionName:       api.GetSessionName(),
		performerKey:      api.GetPerformerKey(),
		performerFIO:      api.GetPerformerFIOKey(),
		roleAFormsKey:     api.GetRoleAFormsKey(),
		roleAFormsNameKey: api.GetRoleAFormsNameKey(),
		logg:              logg,
	}
}

// GetPerformerData получаем данные о сотруднике из сессии, иначе из БД.
func (m *AuthMiddleware) GetPerformerData(r *http.Request, performerService service.PerformerUseCase) (*PerformerData, error) {
	// 1. Получаем данные из сессии.
	performerTabNum, ok1 := m.GetPerformerTabNum(r)
	performerRoleAFormsId, ok2 := m.GetPerformerRoleAFormsId(r)
	if !ok1 || !ok2 {
		m.logg.LogW(msg.WSS400, logg.SkipNofS)

		return nil, errors.New("Пользователь не авторизован ")
	}

	var performerFIO string
	var performerRoleAFormsName string

	// 2. Проверяем кэш.
	m.cacheMu.RLock()
	if cached, exists := m.userCache[performerTabNum]; exists && time.Now().Before(cached.Expires) {
		performerFIO = cached.PerformerFIO
		performerRoleAFormsName = cached.RoleAFormsName
		m.cacheMu.RUnlock()

		return &PerformerData{
			PerformerFIO:          performerFIO,
			PerformerTabNum:       performerTabNum,
			PerformerRoleAForms:   performerRoleAFormsName,
			PerformerRoleAFormsId: performerRoleAFormsId,
		}, nil
	}
	m.cacheMu.RUnlock()

	// 3. Загрузить данные.
	ctx := r.Context()
	performer, err := performerService.FindPerformerByTabNum(ctx, performerTabNum)
	if err != nil {
		m.logg.LogE(msg.ESS501, err, logg.SkipNofS)

		return nil, err
	}

	if performer == nil {
		m.logg.LogW(msg.WSS401, logg.SkipNofS)

		return nil, err
	}

	performerFIO = performer.FIO
	performerRoleAFormsName = performer.PerformerRole.RoleNameAForms

	// 4. Сохраняем в кэш.
	m.cacheMu.Lock()
	if m.userCache == nil {
		m.userCache = make(map[int]*UserSession)
	}
	m.userCache[performerTabNum] = &UserSession{
		PerformerFIO:   performerFIO,
		RoleAFormsName: performerRoleAFormsName,
		Expires:        time.Now().Add(expiresCache),
	}
	m.cacheMu.Unlock()

	return &PerformerData{
		PerformerFIO:          performerFIO,
		PerformerTabNum:       performerTabNum,
		PerformerRoleAForms:   performerRoleAFormsName,
		PerformerRoleAFormsId: performerRoleAFormsId,
	}, nil
}

// clearSessionCookie очищает куки сессии.
func (m *AuthMiddleware) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     api.GetSessionName(),
		Value:    "",
		Path:     api.GetPathToDefault(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// CreateAuthSecuritySession создает безопасную сессию.
func (m *AuthMiddleware) CreateAuthSecuritySession(w http.ResponseWriter, r *http.Request, performerData *PerformerData) error {
	m.logg.LogIf(logg.SkipNofS, "%s: %s", msg.ISS201, performerData.PerformerFIO)

	session, err := m.store.Get(r, api.GetSessionName())
	if err != nil {
		m.logg.LogWf(logg.SkipNofS, "%s: %v", msg.WSS403, err)

		m.clearSessionCookie(w)
		session, _ = m.store.New(r, api.GetSessionName())
	}

	if session == nil {
		m.logg.LogEf(logg.SkipNofS, err, "%s", msg.ESS509)

		return err
	}

	token := api.GenerateSessionToken()

	session.Values[api.GetPerformerKey()] = performerData.PerformerTabNum
	session.Values[api.GetPerformerFIOKey()] = performerData.PerformerFIO
	session.Values[api.GetRoleAFormsKey()] = performerData.PerformerRoleAFormsId
	session.Values[api.GetRoleAFormsNameKey()] = performerData.PerformerRoleAForms
	session.Values[api.GetAuthPerformer()] = true
	session.Values[api.GetTokenKey()] = token
	session.Values[api.GetCreateAtKey()] = time.Now().Unix()
	session.Values[api.GetLastActivityKey()] = time.Now().Unix()

	session.Options = &sessions.Options{
		Path:     api.GetPathToDefault(),
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	m.SetSecurityHeaders(w)

	m.logg.LogIf(logg.SkipNofS, "%s: %s", msg.ISS202, performerData.PerformerFIO)

	return session.Save(r, w)
}

// RequireAuth - основной middleware для проверки аутентификации.
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.SetSecurityHeaders(w)

		session, err := m.getSecureSession(r)
		if err != nil {
			m.forceLogoutAndRedirect(w, r, msg.ESS502)

			return
		}

		if session == nil {
			m.forceLogoutAndRedirect(w, r, msg.ESS508)

			return
		}

		if auth, ok := session.Values[api.GetAuthPerformer()].(bool); !ok || !auth {
			m.forceLogoutAndRedirect(w, r, msg.ESS507)

			return
		}

		if m.isSessionExpired(session) {
			m.forceLogoutAndRedirect(w, r, msg.ESS506)

			return
		}

		m.updateSessionActivity(session, w, r)

		if r.Header.Get("Accept") == "text/html" {
			m.addHistoryManagementScript(w)
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRoleForAForms - middleware для проверки роли AForms
func (m *AuthMiddleware) RequireRoleForAForms(requiredRoles []int, next http.HandlerFunc) http.HandlerFunc {
	allowedRoles := make(map[int]bool, len(requiredRoles))
	for _, role := range requiredRoles {
		allowedRoles[role] = true
	}

	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, api.GetSessionName())
		if err != nil {
			m.logg.LogE(msg.ESS502, err, logg.SkipNofS)
			http.Redirect(w, r, api.GetPathToDefault(), http.StatusFound)

			return
		}

		performerRoleAFormsId, ok := session.Values[api.GetRoleAFormsKey()].(int)
		if !ok {
			m.forceLogoutAndRedirect(w, r, msg.ESS505)

			return
		}

		if !allowedRoles[performerRoleAFormsId] {
			m.logg.LogHttpEf(http.StatusForbidden, r.Method, r.URL.Path, nil, logg.SkipNofS, "%s", msg.ESS510)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetPerformerTabNum получает табельный номер из сессии.
func (m *AuthMiddleware) GetPerformerTabNum(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, api.GetSessionName())
	if err != nil {
		return 0, false
	}

	performerTabNum, ok := session.Values[api.GetPerformerKey()].(int)

	return performerTabNum, ok
}

// GetPerformerRoleAFormsId получает роль из сессии для AForms.
func (m *AuthMiddleware) GetPerformerRoleAFormsId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, api.GetSessionName())
	if err != nil {
		return 0, false
	}

	performerRole, ok := session.Values[api.GetRoleAFormsKey()].(int)

	return performerRole, ok
}

// getSecureSession - безопасное получение сессии с валидацией.
func (m *AuthMiddleware) getSecureSession(r *http.Request) (*sessions.Session, error) {
	session, err := m.store.Get(r, api.GetSessionName())
	if err != nil {
		return nil, err
	}

	if session.IsNew {
		return nil, nil
	}

	return session, nil
}

// isSessionExpired - проверяет истечение срока действия сессии.
func (m *AuthMiddleware) isSessionExpired(session *sessions.Session) bool {
	if createdAt, ok := session.Values[api.GetCreateAtKey()].(int64); ok {
		createTime := time.Unix(createdAt, 0)

		maxAge := maxAgeSession
		if customMaxAge, ok := session.Values[api.GetMaxAgeKey()].(int); ok {
			maxAge = time.Duration(customMaxAge) * time.Second
		}

		return time.Since(createTime) > maxAge
	}

	return true
}

// updateSessionActivity - обновление времени активности.
func (m *AuthMiddleware) updateSessionActivity(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	session.Values[api.GetLastActivityKey()] = time.Now().Unix()

	// Устанавливаем куку с коротким временем жизни для браузера.
	if cookie, err := r.Cookie(api.GetCookieActivityCheckKey()); err != nil || cookie.Value != api.GetCookieActivityKey() {
		http.SetCookie(w, &http.Cookie{
			Name:     api.GetCookieActivityCheckKey(),
			Value:    api.GetCookieActivityKey(),
			Path:     api.GetPathToDefault(),
			MaxAge:   -1,
			HttpOnly: false,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	if err := session.Save(r, w); err != nil {
		m.logg.LogE(msg.ESS504, err, logg.SkipNofS)

		return
	}
}

// SetSecurityHeaders - установка заголовков безопасности.
func (m *AuthMiddleware) SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

// forceLogoutAndRedirect - принудительный выход и редирект с очисткой истории.
func (m *AuthMiddleware) forceLogoutAndRedirect(w http.ResponseWriter, r *http.Request, reason string) {
	m.logg.LogWf(logg.SkipNofS, "%s: %s ", msg.WSS402, reason)

	if session, err := m.store.Get(r, api.GetSessionName()); err == nil {
		session.Options.MaxAge = maxAgeClearSession
		for key := range session.Values {
			delete(session.Values, key)
		}

		if err = session.Save(r, w); err != nil {
			return
		}
	}

	m.SetSecurityHeaders(w)

	tmpl, err := template.ParseFiles(prefixTmpl + tmplForceLogoutHTML)
	if err != nil {
		http.Error(w, msg.EH5001, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := struct {
		Title     string
		Reason    string
		Timestamp string
	}{
		Title:     "Требуется повторная авторизация",
		Reason:    reason,
		Timestamp: time.Now().Format("02.01.2006 15:04:05"),
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, msg.EH5002, http.StatusInternalServerError)

		return
	}
}

// addHistoryManagementScript - добавление скрипта управления историей.
func (m *AuthMiddleware) addHistoryManagementScript(w http.ResponseWriter) {
	w.Header().Add("X-History-Control", "no-cache")
}
