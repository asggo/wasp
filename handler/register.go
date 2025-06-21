package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/asggo/wasp/config"
	"github.com/asggo/wasp/store"
)

var (
	regTmpl          = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/register.html"))
	regAdminTmpl     = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/register_admin.html"))
	passwordNotMatch = "The passwords do not match."
	passwordTooShort = "The password must be at least %d characters."
	usernameTooShort = "The username must be at least %d characters."
	usernameTaken    = "Username is already taken."
)

// registerHandler provides handlers for all of the endpoints in the /register
// path.
type registerHandler struct {
	cfg *config.Config
	db  *store.Store
}

// Index renders the index page of the /register path.
func (rh *registerHandler) Index(w http.ResponseWriter, r *http.Request) {
	regTmpl.ExecuteTemplate(w, "layout", nil)
}

// Register creates a new User in the Store using the information provided in
// the registration form.
func (rh *registerHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	un := r.Form.Get("username")
	pw := r.Form.Get("password")
	cn := r.Form.Get("confirm")

	if len(un) < rh.cfg.MinUsernameLength {
		e := fmt.Errorf(usernameTooShort, rh.cfg.MinUsernameLength)
		regTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), e))
		return
	}

	if len(pw) < rh.cfg.MinPassphraseLength {
		e := fmt.Errorf(passwordTooShort, rh.cfg.MinPassphraseLength)
		regTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), e))
		return
	}

	if pw != cn {
		regTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), passwordNotMatch))
		return
	}

	if rh.db.UserExists(un) {
		regTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), usernameTaken))
		return
	}

	user := store.NewUser(un)

	err := rh.db.CreateUser(user, pw)
	if err != nil {
		e := fmt.Errorf("could not RegisterHandler.Register: %v", err)
		NewServerError(e).Handle(w, r)
	}

	http.Redirect(w, r, "/account", http.StatusFound)
}

// RegisterAdmin creates the first admin account.
func (rh *registerHandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pw := r.Form.Get("password")
	cn := r.Form.Get("confirm")

	if len(pw) < rh.cfg.MinPassphraseLength {
		e := fmt.Errorf(passwordTooShort, rh.cfg.MinPassphraseLength)
		regAdminTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), e))
		return
	}

	if pw != cn {
		regAdminTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), passwordNotMatch))
		return
	}

	if rh.db.UserExists("admin") {
		e := fmt.Errorf("could not RegisterHandler.RegisterAdmin: admin user exists")
		NewBadRequestError(e).Handle(w, r)
	}

	user := store.NewUser("admin")
	user.Admin = true

	err := rh.db.CreateUser(user, pw)
	if err != nil {
		e := fmt.Errorf("could not RegisterHandler.RegisterAdmin: %v", err)
		NewServerError(e).Handle(w, r)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// NewRegisterHandler creates a new registerHandler using the given Config and
// Store.
func NewRegisterHandler(c *config.Config, s *store.Store) *registerHandler {
	return &registerHandler{cfg: c, db: s}
}
