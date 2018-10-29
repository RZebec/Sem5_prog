package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"html/template"
	"net/http"
)

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
			{{end}}
		</div>
		<div class="content">
			<div class="container">
				<h1>This is the Index Page</h1>
			</div>
		</div>
	</body>

	</html>`

type IndexPageHandler struct {
	AccessTokenCookie helpers.Cookie
	SessionManager session.SessionManager
}

type indexPageData struct {
	IsUserLoggedIn bool
}

func (i IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("index").Parse(indexPageTemplate)

	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, i.SessionManager, i.AccessTokenCookie)

	// Todo: HANDLE Template parsing error
	data := indexPageData{
		IsUserLoggedIn: isUserLoggedIn,
	}

	t.Execute(w, data)
}