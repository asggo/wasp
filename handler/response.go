package handler

import (
	"context"

	"github.com/asggo/wasp/store"
)

type Response struct {
	Auth  bool
	Admin bool
	Data  interface{}
}

// NewResponse returns an appropriate Response object based on the context
// provided.
func NewResponse(ctx context.Context, data interface{}) Response {
	val := ctx.Value("user")
	if val == nil {
		return Response{Auth: false, Admin: false, Data: data}
	}

	user := val.(store.User)
	if user.Admin {
		return Response{Auth: true, Admin: true, Data: data}
	} else {
		return Response{Auth: true, Admin: false, Data: data}
	}
}
