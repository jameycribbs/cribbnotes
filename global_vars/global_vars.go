package global_vars

import (
	"github.com/gorilla/sessions"
)

type GlobalVars struct {
	SessionStore *sessions.CookieStore
}
