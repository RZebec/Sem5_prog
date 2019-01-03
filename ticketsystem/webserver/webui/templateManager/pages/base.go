package pages

type BasePageData struct {
	UserIsAuthenticated bool
	UserIsAdmin         bool
	Active              string
}

/*
	Base Html template.
*/
var Base = `
	{{ define "Base" }}
	<!DOCTYPE html>
	<html>

	<head>
		<title>
			{{ block "Title" .}} {{ end }}
		</title>
		<link rel="stylesheet" href="/files/style/main">
		{{ block "StylesAndScripts" .}} {{ end }}
	</head>
	
	<body>
		<div class="topnav">
			<a {{if eq .Active "index" }}class="active"{{end}} href="/">Home</a>

			<a {{if eq .Active "tickets" }}class="active"{{end}} href="/tickets">Tickets</a>
	
			<span>OP-Ticket-System</span>

			{{if .UserIsAuthenticated}}
				{{if .UserIsAdmin}}
					<a {{if eq .Active "admin" }}class="active"{{end}} href="/admin">Admin</a>
				{{end}}
				<a href="/user_logout">Logout</a>				
			{{else}}
				<a {{if eq .Active "login" }}class="active"{{end}} href="/login">Login</a>
				<a {{if eq .Active "register" }}class="active"{{end}} href="/register">Register</a>				
			{{end}}
		</div>
		{{ template "Content" .}}
	</body>
	
	</html>
	{{ end }}`
