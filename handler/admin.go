package handler

import (
	"html/template"
	"net/http"

	"github.com/asggo/wasp/store"
)

var (
	adminTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/admin.html"))
)

// adminHandler provides handlers for all of the endpoints in the /admin path.
type adminHandler struct {
	db *store.Store
}

// Index renders the index page of the /admin path.
func (ah *adminHandler) Index(w http.ResponseWriter, r *http.Request) {
	adminTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), nil))
}

// NewAdminHandler creates a new adminHandler with the given Store.
func NewAdminHandler(s *store.Store) *adminHandler {
	return &adminHandler{db: s}
}
