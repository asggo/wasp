package webapp

import (
	"testing"

	"github.com/asggo/webtest"
)

func TestApplication(t *testing.T) {
	// Create an application router and run our tests
	router := NewApplication().Router()

	webtest.TestHandler(t, "tests/index_test.txt", router)
	webtest.TestHandler(t, "tests/admin_test.txt", router)
	webtest.TestHandler(t, "tests/account_test.txt", router)
	webtest.TestHandler(t, "tests/user_test.txt", router)
}
