package accessdenied

import (
	"html/template"
	"net/http"
)

/*
	Html template for the Access Denied Page.
 */
var accessDeniedTemplate = `
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

		</div>
		<div class="content">
			<div class="container">
				<h1>Access is denied: User is not logged in</h1>
			</div>
		</div>
	</body>
	
	</html>`

/*
	Structure for the Access Denied Page Handler.
 */
type AccessDeniedPageHandler struct {
}

/*
	Structure for the Access Denied Page Data.
 */
type accessDeniedPageData struct {
}

func (l AccessDeniedPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := accessDeniedPageData{}

	// Todo: HANDLE Template parsing error
	t, _ := template.New("accessDenied").Parse(accessDeniedTemplate)

	t.Execute(w, data)
}
