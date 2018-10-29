package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"html/template"
	"net/http"
	"strconv"
)

var loginPageTemplate = `
	<!DOCTYPE html>
	<html>

	<head>
		<link rel="stylesheet" href="/files/style/main">
		<link rel="stylesheet" href="/files/style/login"> 
		<script src="/files/script/login"></script>
	</head>
	
	<body>
		<div class="topnav">
			<a href="/">Home</a>
	
			<span>OP-Ticket-System</span>

			<a class="active" href="/login">Login</a>

		</div>
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Login</h2>
					<form id="form_id" method="post" name="myform" action="/user_login">
						<label>Username:</label>
						<input type="text" name="userName" id="userName" />
						<label>Password:</label>
						<input type="password" name="password" id="password" />
						<button type="submit" id="submitLogin" class="submit-button" disabled>Login</button>
					</form>
					{{if .IsLoginFailed }}
					<span class="error-message">
							Login Failed!
					</span>
					</br>
        			{{end}}
					<span id="emailNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	</body>
	
	</html>`

type LoginPageHandler struct {
	Config			config.Configuration
	SessionManager session.SessionManager
}

type loginPageData struct {
	IsLoginFailed bool
}

func (l LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.SessionManager, l.Config.AccessTokenCookieName)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	IsLoginFailed, err := strconv.ParseBool(queryValues.Get("IsLoginFailed"))

	if err != nil {
		// TODO: Handle error
	}

	// Todo: HANDLE Template parsing error
	t, _ := template.New("login").Parse(loginPageTemplate)

	data := loginPageData{
		IsLoginFailed: IsLoginFailed,
	}

	t.Execute(w, data)
}
