package handler

import (
	"html/template"
	"net/http"

	"github.com/asggo/wasp/store"
)

var (
	siteTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/site.html"))
)

// siteHandler provides handlers for each endpoint in the /site path.
type siteHandler struct {
	db *store.Store
}

// Index renders the index page of the /site path.
func (sh *siteHandler) Index(w http.ResponseWriter, r *http.Request) {
	siteTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), nil))
}

// NewSiteHandler creates a new siteHandler with the given Store.
func NewSiteHandler(s *store.Store) *siteHandler {
	return &siteHandler{db: s}
}
