package pages

/*
	Html template for the Register Page.
*/
var RegisterPage = `
	{{ define "Title" }} Register {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/center_main"> 
		<script src="/files/script/login"></script>
		<a href="/admin">Admin</a>
	{{ end }}
	
	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Register</h2>
					<form id="form_id" method="post" name="myform" action="/user_register">
						<label>First Name:</label>
						<input type="text" name="first_name" id="first_name" />
						<label>Last Name:</label>
						<input type="text" name="last_name" id="last_name" />
						<label>Username:</label>
						<input type="text" name="userName" id="userName" />
						<label>Password:</label>
						<input type="password" name="password" id="password" />
						<button type="submit" id="submitLogin" class="submit-button" disabled>Register</button>
					</form>
					{{if .IsRegisteringFailed }}
					<span class="error-message">
							Registering Failed!
					</span>
					</br>
        			{{end}}
					<span id="emailNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	</body>
	
	</html>
	{{ end }}`
