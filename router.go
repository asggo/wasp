package webapp

import (
	"net/http"

	"github.com/asggo/wasp/config"
	"github.com/asggo/wasp/handler"
	"github.com/asggo/wasp/middleware"
	"github.com/asggo/wasp/store"
	"github.com/go-chi/chi/v5"
)

// indexRouter defines all of the routes needed for the / of the site.
func indexRouter(c *config.Config, s *store.Store) http.Handler {
	r := chi.NewRouter()
	h := handler.NewIndexHandler(c, s)

	r.Get("/", h.Index)
	r.Mount("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return r
}

// accountRouter defines all of the routes needed for account creation and
// authentication.
func accountRouter(c *config.Config, s *store.Store) http.Handler {
	r := chi.NewRouter()
	ah := handler.NewAuthHandler(c, s)
	rh := handler.NewRegisterHandler(c, s)

	r.Get("/", ah.Index)
	r.Post("/login", ah.Login)
	r.Get("/logout", ah.Logout)
	r.Get("/register", rh.Index)
	r.Post("/register", rh.Register)
	r.Post("/admin", rh.RegisterAdmin)

	return r
}

// siteRouter defines all of the routes needed for the authenticated portion
// of the site.
func siteRouter(c *config.Config, s *store.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Authorizer(s))

	h := handler.NewSiteHandler(s)

	r.Get("/", h.Index)
	r.Mount("/admin", adminRouter(s))
	r.Mount("/user", userRouter(c, s))

	return r
}

// The routers below are mounted to the siteRouter and there for automatically
// use the Authorizer middleware.

// adminRouter defines all of the routes needed for the administrative portion
// of the site. Includes middleware to confirm a user is an admin.
func adminRouter(s *store.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.AdminAuthorizer)

	h := handler.NewAdminHandler(s)

	r.Get("/", h.Index)

	return r
}

// userRouter defines all of the routes needed to manage the user account.
func userRouter(c *config.Config, s *store.Store) http.Handler {
	r := chi.NewRouter()

	h := handler.NewUserHandler(c, s)

	r.Get("/", h.Index)
	r.Get("/changepw", h.ShowChangePassword)
	r.Post("/changepw", h.ExecChangePassword)

	return r
}
