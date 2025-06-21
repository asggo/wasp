package handler

import (
	"html/template"
	"net/http"

	"github.com/asggo/wasp/config"
	"github.com/asggo/wasp/store"
)

var (
	idxTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/index.html"))
)

// indexHandler provides handlers for each of the endpoints in the / path.
type indexHandler struct {
	cfg *config.Config
	db  *store.Store
}

// Index renders the index page or the config page depending on whether the
// admin user account has been created.
func (ih *indexHandler) Index(w http.ResponseWriter, r *http.Request) {
	// If the admin user does not exist we need to configure it.
	if !ih.db.UserExists("admin") {
		regAdminTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), nil))
		return
	}

	idxTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), nil))
}

// NewIndexHandler returns an indexHandler using the given Config and Store.
func NewIndexHandler(c *config.Config, s *store.Store) *indexHandler {
	return &indexHandler{cfg: c, db: s}
}
