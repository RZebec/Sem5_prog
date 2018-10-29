package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"html/template"
	"net/http"
	"strconv"
)

/*
	Html template for the Login Page.
*/
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

/*
	Structure for the Login Page Handler.
*/
type LoginPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
}

/*
	Structure for the Login Page Data.
*/
type loginPageData struct {
	IsLoginFailed bool
}

/*
	The Login Page handler.
*/
func (l LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.UserContext, l.Config.AccessTokenCookieName)

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
