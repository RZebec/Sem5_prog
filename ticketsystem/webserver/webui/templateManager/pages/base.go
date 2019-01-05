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

			<a {{if eq .Active "all_tickets" }}class="active"{{end}} href="/all_tickets">All Tickets</a>

			<a {{if eq .Active "open_tickets" }}class="active"{{end}} href="/open_tickets">Open Tickets</a>

			<a {{if eq .Active "closed_tickets" }}class="active"{{end}} href="/closed_tickets">Closed Tickets</a>

			<a {{if eq .Active "active_tickets" }}class="active"{{end}} href="/active_tickets">Active Tickets</a>

			<a {{if eq .Active "ticket_create" }}class="active"{{end}} href="/ticket_create">Create Ticket</a>
	
			<span>OP Ticket System</span>

			{{if .UserIsAuthenticated}}
				<a {{if eq .Active "user_tickets" }}class="active"{{end}} href="/user_tickets">My Tickets</a>
				{{if .UserIsAdmin}}
					<a {{if eq .Active "admin" }}class="active"{{end}} href="/admin">Server Settings</a>
				{{end}}
				<a {{if eq .Active "settings" }}class="active"{{end}} href="/user_settings">My Settings</a>
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
