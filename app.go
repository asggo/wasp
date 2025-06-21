package webapp

import (
	"fmt"

	"github.com/asggo/wasp/config"
	"github.com/asggo/wasp/handler"
	"github.com/asggo/wasp/middleware"
	"github.com/asggo/wasp/store"
	"github.com/go-chi/chi/v5"
)

// Application holds our web application.
type Application struct {
	r chi.Router
}

// Router returns the Chi router for the web application.
func (a *Application) Router() chi.Router {
	return a.r
}

// NewApplication creates a new Application object using the given Config
// object.
func NewApplication() *Application {
	var app Application

	// Get our configuration
	cfg := config.NewConfiguration()

	// Setup our Store
	store, err := store.NewStore(cfg.StorePath)
	if err != nil {
		panic(fmt.Errorf("could not NewApplication: %v", err))
	}

	// Setup our router
	r := chi.NewRouter()
	r.NotFound(handler.NotFoundHandler)

	// Setup our middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Timeout(cfg.RequestTimeout))
	r.Use(middleware.SecurityHeaders)

	// Mount our sub routers
	r.Mount("/", indexRouter(&cfg, &store))
	r.Mount("/account", accountRouter(&cfg, &store))
	r.Mount("/site", siteRouter(&cfg, &store))

	app.r = r

	return &app
}
