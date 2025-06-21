package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asggo/wasp/handler"
	"github.com/asggo/wasp/store"
)

// Authorizer determines if the request has proper session cookie. If so, it
// loads the user tied to the session cookie in the request. Otherwise it
// returns an invalid session error.
func Authorizer(s *store.Store) func(next http.Handler) http.Handler {
	handlerFn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sess, err := store.NewSessionFromRequest(r, s)
			if err != nil {
				e := fmt.Errorf("could not Authorizer: %v", err)
				handler.NewBadRequestError(e).Handle(w, r)
				return
			}

			if sess.IsExpired() {
				err := s.DeleteSession(sess.SessionId)
				if err != nil {
					e := fmt.Errorf("could not Authorizer: %v", err)
					handler.NewServerError(e).Handle(w, r)
				}

				e := fmt.Errorf("could not Authorizer: session expired")
				handler.NewUnauthorizedError(e).Handle(w, r)
			}

			user, err := s.GetUser(sess.UserId)
			if err != nil {
				e := fmt.Errorf("could not Authorizer: %v", err)
				handler.NewServerError(e).Handle(w, r)
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}

	return handlerFn
}

// AdminAuthorizer determines if the User object in the request context is an
// admin user. If so, they are allowed to pass through to the /admin path
// otherwise they are redirected to the /site path.
func AdminAuthorizer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(store.User)

		if !user.Admin {
			e := fmt.Errorf("could not AdminAuthorizer: %s is not admin user", user.Alias)
			handler.NewForbiddenError(e).Handle(w, r)
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
