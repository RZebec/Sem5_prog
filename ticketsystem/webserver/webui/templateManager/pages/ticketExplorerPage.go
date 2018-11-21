package pages

/*
	Html template for the Index Page.
*/
var TicketExplorerPage = `	

	{{ define "Title" }} Tickets {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/tickets"> 
	{{ end }}
	
	{{ define "Content" }}
		<div class="topnav">
			<a class="active" href="/">Home</a>
	
			<span>OP-Ticket-System</span>

			<a href="/tickets">Tickets</a>

			{{if .IsUserLoggedIn}}
				<a href="/user_logout">Logout</a>
			{{else}}
				<a href="/login">Login</a>
				<a href="/register">Register</a>
			{{end}}
		</div>
		<div class="content">
			<div class="container">
				<div class="main">
			{{range .TicketInfo}}
    				<label>Title:</label> 
    				{{.Title}}
					<label>Editor:</label> 
					{{.Editor.FirstName}}
					<label>Creator:</label> 
					{{.Creator.FirstName}}
			{{end}}
				</div>
			</div>
		</div>
	</body>

	</html>
	{{ end }}`
