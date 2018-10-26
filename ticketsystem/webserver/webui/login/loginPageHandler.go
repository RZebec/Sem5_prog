package login

import (
	"html/template"
	"net/http"
	"strconv"
)

var loginTemplate = `
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
	IsUserLoggedIn bool
	IsLoginFailed  bool
}

type LoginPageData struct {
	IsLoginFailed bool
}

func (l LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Checks if the User is already logged in and if so redirects him to the start page
	if l.IsUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	IsLoginFailed, err := strconv.ParseBool(queryValues.Get("IsLoginFailed"))

	if err == nil {
		l.IsLoginFailed = IsLoginFailed
	}

	t, _ := template.New("login").Parse(loginTemplate)

	data := LoginPageData{
		IsLoginFailed: l.IsLoginFailed,
	}

	t.Execute(w, data)
}
