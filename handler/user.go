package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/asggo/wasp/config"
	"github.com/asggo/wasp/store"
)

var (
	userTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/user.html"))
	pwdTmpl  = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/changepw.html"))
)

// userHandler provides handlers for each endpoint in the /site path.
type userHandler struct {
	db  *store.Store
	cfg *config.Config
}

// Index renders the index page of the /user path.
func (uh *userHandler) Index(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(store.User)

	userTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), user))
}

// ShowChangePassword renders the change password page.
func (uh *userHandler) ShowChangePassword(w http.ResponseWriter, r *http.Request) {
	pwdTmpl.ExecuteTemplate(w, "layout", nil)
}

// ExecChangePassword resets the users password.
func (uh *userHandler) ExecChangePassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opw := r.Form.Get("old-password")
	npw := r.Form.Get("new-password")
	cpw := r.Form.Get("confirm")

	u := r.Context().Value("user").(store.User)

	if !uh.db.AuthenticateUser(u.UserId, opw) {
		pwdTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), invalidCredentials))
		return
	}

	if len(npw) < uh.cfg.MinPassphraseLength {
		e := fmt.Errorf(passwordTooShort, uh.cfg.MinPassphraseLength)
		pwdTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), e))
		return
	}

	if npw != cpw {
		pwdTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), passwordNotMatch))
		return
	}

	err := uh.db.ChangeUserPassword(u.UserId, npw)
	if err != nil {
		e := fmt.Errorf("could not UserHandler.ExecChangePassword: %v", err)
		NewServerError(e).Handle(w, r)
	}

	http.Redirect(w, r, "/account/logout", http.StatusFound)
}

// NewUserHandler creates a new userHandler with the given Store.
func NewUserHandler(c *config.Config, s *store.Store) *userHandler {
	return &userHandler{cfg: c, db: s}
}
