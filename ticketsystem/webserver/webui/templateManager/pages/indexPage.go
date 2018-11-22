package pages

/*
	Html template for the Index Page.
*/
var IndexPage = `	

	{{ define "Title" }} Home {{ end }}

	{{ define "StylesAndScripts" }} {{ end }}
	
	{{ define "Content" }}
		<div class="topnav">
			<a class="active" href="/">Home</a>

			<a href="/tickets">Tickets</a>
	
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

	</html>
	{{ end }}`
