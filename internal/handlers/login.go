package handlers

import (
	"fmt"
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Handle posting to login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient("127.0.0.1:3002", opts...)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	client := models.NewDatabaseClient(conn)

	msg, err := client.Ping(r.Context(), &models.SimpleMessage{Msg: "Ping"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(msg.Msg)

	/*
		// Renew session token
		_ = m.App.Session.RenewToken(r.Context())

		// Parse form
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// Get form and validate fields
		form := forms.New(r.PostForm)
		form.Required("email", "password")
		if !form.Valid() {

			// Push form to session
			m.AddForms(r, map[string]*forms.Form{
				"login": form,
			})

			// Redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		// Get user data
		uEmail := form.Get("email")
		uPassword := form.Get("password")

		// Check if user DB exists
		_, err = os.Stat(dbrepo.GetUserDBPath(m.App.DBPath, uEmail, true))
		if errors.Is(err, os.ErrNotExist) {

			// Write to error log
			m.App.ErrorLog.Println(err)

			// Reset password in form
			form.Set("password", "")

			// Push form to session
			m.AddForms(r, map[string]*forms.Form{
				"login": form,
			})

			// Add error message
			m.AddErrorMsg(r, "Invalid login credentials")

			// Redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		// Create user connection
		dbconn, err := driver.ConnectSQL(dbrepo.GetUserDBPath(m.App.DBPath, uEmail, false))
		if err != nil {

			// Write to error log
			m.App.ErrorLog.Println(err)

			// Reset password in form
			form.Set("password", "")

			// Push form to session
			m.AddForms(r, map[string]*forms.Form{
				"login": form,
			})

			// Add error message
			m.AddErrorMsg(r, "Server error")

			// Redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		// Get db repo
		repo := dbrepo.NewSqliteRepo(m.App, uEmail, dbconn.SQL)

		// Authenticate user
		_, _, dbVersion, err := repo.Authenticate(uPassword)
		if err != nil {

			// Write to error log
			m.App.ErrorLog.Println(err)

			// Reset password in form
			form.Set("password", "")

			// Push form to session
			m.AddForms(r, map[string]*forms.Form{
				"login": form,
			})

			// Add error message
			m.AddErrorMsg(r, "Invalid login credentials")

			// Redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		// Get user key
		key := dbrepo.GetUserKey(uEmail)

		// Add connection to repo
		m.DB[key] = repo

		// Store user key in session
		m.App.Session.Put(r.Context(), "user_key", key)

		// Store dbVersion in session
		m.App.Session.Put(r.Context(), "db_version", dbVersion)

		// Flash message to user
		m.AddFlashMsg(r, "Logged in successfully")

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	*/
}
