package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

const (
	badRequestError     = "The request could not be processed as sent."
	unauthorizedError   = "You are not authorized to access this content."
	forbiddenError      = "Access to this content is forbidden."
	notFoundError       = "The page you are looking for was not found."
	internalServerError = "Server error. Please try your request again later."
)

var (
	errTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/nav.html", "templates/error.html"))
)

type errorHandler struct {
	status  int
	message string
	err     error
}

func (eh errorHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if eh.err != nil {
		oplog := httplog.LogEntry(r.Context())
		oplog.Error(fmt.Sprintf("%v", eh.err))
	}

	w.WriteHeader(eh.status)
	errTmpl.ExecuteTemplate(w, "layout", NewResponse(r.Context(), eh.message))
}

func NewBadRequestError(err error) errorHandler {
	return errorHandler{
		status:  http.StatusBadRequest,
		message: badRequestError,
		err:     err,
	}
}

func NewUnauthorizedError(err error) errorHandler {
	return errorHandler{
		status:  http.StatusUnauthorized,
		message: unauthorizedError,
		err:     err,
	}
}

func NewForbiddenError(err error) errorHandler {
	return errorHandler{
		status:  http.StatusForbidden,
		message: forbiddenError,
		err:     err,
	}
}

func NewNotFoundError(err error) errorHandler {
	return errorHandler{
		status:  http.StatusNotFound,
		message: notFoundError,
		err:     err,
	}
}

func NewServerError(err error) errorHandler {
	return errorHandler{
		status:  http.StatusInternalServerError,
		message: internalServerError,
		err:     err,
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	NewNotFoundError(nil).Handle(w, r)
}
