package webui

import (
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

			{{if .IsUserLoggedIn }}
				<a href="/logout">Logout</a>
			{{else}}
				<a href="/login">Login</a>
			{{end}}
		</div>
		<div class="content">
			This is the Index Page
		</div>
	</body>

	</html>`

type IndexPageHandler struct {
	IsUserLoggedIn bool
}

type IndexPageData struct {
	IsUserLoggedIn bool
}

func (i IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("index").Parse(indexPageTemplate)

	data := IndexPageData{
		IsUserLoggedIn: i.IsUserLoggedIn,
	}

	t.Execute(w, data)
}
