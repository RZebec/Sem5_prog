package pages

/*
	Html template for the Login Page.
*/
var LoginPage = `
	{{ define "Title" }} Login {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/center_main"> 
		<script src="/files/script/login"></script>
	{{ end }}
	
	{{ define "Content" }}
		<div class="topnav">
			<a href="/">Home</a>

			<a href="/tickets">Tickets</a>
	
			<span>OP-Ticket-System</span>

			<a class="active" href="/login">Login</a>
			<a href="/register">Register</a>
			<a href="/admin">Admin</a>
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
	{{ end }}`
