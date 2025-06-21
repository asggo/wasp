package main

import (
	"net/http"

	"github.com/asggo/wasp"
)

func main() {
	app := webapp.NewApplication()

	http.ListenAndServe(":8000", app.Router())

}
