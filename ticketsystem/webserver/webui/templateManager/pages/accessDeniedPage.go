package pages

/*
	Html template for the Access Denied Page.
*/
var AccessDeniedPage = `
	{{ define "Title" }} Access Denied {{ end }}

	{{ define "StylesAndScripts" }} {{ end }}
	
	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<h1>Access is denied: User is not logged in</h1>
			</div>
		</div>
	{{ end }}`
