package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"html/template"
	"net/http"
)

/*
	Html template for the Index Page.
*/
var indexPageTemplate = `<html>

	<head>
    	<link rel="stylesheet" href="/files/style/main">
	</head>

	<body>
		<div class="topnav">
			<a class="active" href="/">Home</a>
	
			<span>OP-Ticket-System</span>

			{{if .IsUserLoggedIn}}
				<a href="/user_logout">Logout</a>
			{{else}}
				<a href="/login">Login</a>
				<a href="/register">Register</a>
			{{end}}
		</div>
		<div class="content">
			<div class="container">
				<h1>This is the Index Page</h1>
			</div>
		</div>
	</body>

	</html>`

/*
	Structure for the Index Page Handler.
*/
type IndexPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
}

/*
	Structure for the Index Page Data.
*/
type indexPageData struct {
	IsUserLoggedIn bool
}

/*
	The Index Page handler.
*/
func (i IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("index").Parse(indexPageTemplate)

	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, i.UserContext, i.Config.AccessTokenCookieName)

	// Todo: HANDLE Template parsing error
	data := indexPageData{
		IsUserLoggedIn: isUserLoggedIn,
	}

	t.Execute(w, data)
}
