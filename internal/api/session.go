package api

import (
	"crypto/rand"
	"encoding/base64"
	"fgw_web_admin_panel/pkg/msg"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const (
	sessionName              = "admin_session"
	sessionPerformerKey      = "performer_tab_num"
	sessionPerformerFIOKey   = "performer_fio"
	sessionRoleAFormsKey     = "role_id_aforms"
	sessionRoleAFormsNameKey = "role_name_aforms"
	sessionAuthPerformer     = "authenticated"
	sessionCreateAtKey       = "create_at"
	sessionTokenKey          = "session_token"
	sessionMaxAgeKey         = "max_age"
	sessionLastActivityKey   = "last_activity"

	cookieActivityCheckKey = "activity_check"
	cookieActivityKey      = "activity"

	pathToDefault = "/"
	maxAge        = 1800 // maxAge время жизни сессии 3 дня.
)

func GetSessionName() string            { return sessionName }
func GetPerformerKey() string           { return sessionPerformerKey }
func GetPerformerFIOKey() string        { return sessionPerformerFIOKey }
func GetRoleAFormsKey() string          { return sessionRoleAFormsKey }
func GetRoleAFormsNameKey() string      { return sessionRoleAFormsNameKey }
func GetAuthPerformer() string          { return sessionAuthPerformer }
func GetCreateAtKey() string            { return sessionCreateAtKey }
func GetTokenKey() string               { return sessionTokenKey }
func GetMaxAgeKey() string              { return sessionMaxAgeKey }
func GetLastActivityKey() string        { return sessionLastActivityKey }
func GetPathToDefault() string          { return pathToDefault }
func GetCookieActivityCheckKey() string { return cookieActivityCheckKey }
func GetCookieActivityKey() string      { return cookieActivityKey }

var Store *sessions.CookieStore

func InitSessionStore() {
	secretKey := getSecretKey()

	Store = sessions.NewCookieStore([]byte(secretKey))

	Store.Options = &sessions.Options{
		Path:     GetPathToDefault(),
		MaxAge:   maxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	log.Println(msg.ISS200)
}

// getSecretKey получаем сгенерированный секретный ключ из .env, иначе генерируем новый.
func getSecretKey() string {
	if secretKey := os.Getenv("SESSION_SECRET"); secretKey != "" {
		return secretKey
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(msg.ESS500 + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}

func GenerateSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}
