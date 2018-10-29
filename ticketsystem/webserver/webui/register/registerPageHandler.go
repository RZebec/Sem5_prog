package register

import (
"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
"html/template"
"net/http"
"strconv"
)

/*
	Html template for the Register Page.
*/
var registerPageTemplate = `
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

			<a href="/login">Login</a>
			<a class="active" href="/register">Register</a>
		</div>
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Register</h2>
					<form id="form_id" method="post" name="myform" action="/user_register">
						<label>First Name:</label>
						<input type="text" name="first_name" id="first_name" />
						<label>Last Name:</label>
						<input type="text" name="last_name" id="last_name" />
						<label>Username:</label>
						<input type="text" name="userName" id="userName" />
						<label>Password:</label>
						<input type="password" name="password" id="password" />
						<button type="submit" id="submitLogin" class="submit-button" disabled>Register</button>
					</form>
					{{if .IsRegisteringFailed }}
					<span class="error-message">
							Registering Failed!
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
	Structure for the Register Page Handler.
*/
type RegisterPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
}

/*
	Structure for the Register Page Data.
*/
type registerPageData struct {
	IsRegisteringFailed bool
}

/*
	The Register Page handler.
*/
func (l RegisterPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.UserContext, l.Config.AccessTokenCookieName)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	isRegisteringFailed, err := strconv.ParseBool(queryValues.Get("IsRegisteringFailed"))

	if err != nil {
		// TODO: Handle error
	}

	// Todo: HANDLE Template parsing error
	t, _ := template.New("register").Parse(registerPageTemplate)

	data := registerPageData{
		IsRegisteringFailed: isRegisteringFailed,
	}

	t.Execute(w, data)
}

