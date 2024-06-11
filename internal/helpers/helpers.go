package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
	app.ErrorLog.Println("Client error with status of", status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated checks if a user is authenticated
func IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_token")
}
