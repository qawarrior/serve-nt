package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/qawarrior/serve-nt/configuration"

	"github.com/qawarrior/secrets"

	"github.com/qawarrior/serve-nt/models"
)

func loginGet(w http.ResponseWriter, r *http.Request) {
	serveTemplate("./assets/templates/login.html", tempdata{Timestamp: time.Now(), AppName: "SERVE-NT"}, w)
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	// Get the User from the form data
	r.ParseForm()
	user := models.NewUser()
	err := fDecoder.Decode(user, r.PostForm)
	if err != nil {
		loginError(w, r, err)
		return
	}

	// Get submitted password for later comparison
	txtPwd := user.Password

	// Try to retrieve existing user
	err = user.FindByEmail()
	if err != nil {
		configuration.Lerror.Println("User does not Exist")
		http.Redirect(w, r, "/registration", 200)
		return
	}

	// Compare the password submitted against stored hash
	if secrets.ComparePassword(txtPwd, user.Password) == false {
		err = errors.New("Email and Password dont match")
		loginError(w, r, err)
		return
	}

	// Create a authenticated session
	session, err := sessionStore.Get(r, "SNT-SESSION")
	if err != nil {
		loginError(w, r, err)
		return
	}

	// If session is new, set values
	if session.IsNew {
		session.Values["UserId"] = user.ID.String()
		session.Values["LoggedIn"] = true
	}

	// Save session back to client
	err = session.Save(r, w)
	if err != nil {
		loginError(w, r, err)
		return
	}

	rurl := `/profile/` + user.ID.Hex()
	http.Redirect(w, r, rurl, 200)
}

func loginError(w http.ResponseWriter, r *http.Request, err error) {
	configuration.Lerror.Println("Login Failed:", err)
	http.Redirect(w, r, "/login", 200)
}
