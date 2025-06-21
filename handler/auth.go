package handler

// ----------------------------------------------------------------------------
// Authentication Handler
// ----------------------------------------------------------------------------

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/asggo/wasp/config"
	"github.com/asggo/wasp/store"
)

var (
	loginTmpl          = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/login.html"))
	logoutTmpl         = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/logout.html"))
	invalidCredentials = "Invalid credentials."
)

// authHandler provides handlers for each of the endpoints within /auth
type authHandler struct {
	cfg *config.Config
	db  *store.Store
}

// Index renders the login page.
func (ah *authHandler) Index(w http.ResponseWriter, r *http.Request) {
	loginTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), nil))
}

// Logout removes the session from the store and renders the logout page.
func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := store.NewSessionFromRequest(r, ah.db)
	if err != nil {
		e := fmt.Errorf("could not AuthHandler.Logout: %v", err)
		NewServerError(e).Handle(w, r)
	}

	err = ah.db.DeleteSession(sess.SessionId)
	if err != nil {
		e := fmt.Errorf("could not AuthHandler.Logout: %v", err)
		NewServerError(e).Handle(w, r)
	}

	cookie := http.Cookie{
		Name:     "sess",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	logoutTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), nil))
}

// Login authenticates the user and sets a session cookie upon success.
func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	un := r.Form.Get("username")
	pw := r.Form.Get("password")

	user, err := ah.db.GetUserByAlias(un)
	if err != nil {
		loginTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), invalidCredentials))
		return
	}

	if !ah.db.AuthenticateUser(user.UserId, pw) {
		ah.db.IncrementFailedAuthCount(user.UserId)

		count, err := ah.db.GetFailedAuthCount(user.UserId)
		if err != nil {
			e := fmt.Errorf("could not AuthHandler.Login: %v", err)
			NewServerError(e).Handle(w, r)
			return
		}

		// Sleep based on the failed auth count.
		time.Sleep(time.Duration(25*(1<<count)) * time.Millisecond)

		loginTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), invalidCredentials))
		return
	}

	ah.db.ResetFailedAuthCount(user.UserId)

	sess, err := store.NewSession(user.UserId, ah.cfg.SessionLength)
	if err != nil {
		e := fmt.Errorf("could not AuthHandler.Login: %v", err)
		NewServerError(e).Handle(w, r)
		return
	}

	ah.db.CreateSession(sess)

	cookie := http.Cookie{
		Name:     "sess",
		Value:    sess.SessionId.String(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/site", http.StatusFound)
}

// NewAuthHandler creates a new authHandler object.
func NewAuthHandler(c *config.Config, s *store.Store) *authHandler {
	return &authHandler{cfg: c, db: s}
}
