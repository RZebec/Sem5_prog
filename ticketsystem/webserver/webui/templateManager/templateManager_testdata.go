package templateManager

/*
	Data for the test page
*/
type testPageData struct {
	Title string
}

/*
	Expected result for the access denied page test.
*/
var testPage = `<!DOCTYPE html>
	<html>

	<head>
		<title>
			 {{.Title}}
		</title>
		<link rel="stylesheet" href="/files/style/main">
		 
	</head>
	
	<body>
		
		<div class="topnav">
			<a href="/">Home</a>

			<a href="/tickets">Tickets</a>
	
			<span>OP-Ticket-System</span>

			<a href="/login">Login</a>
			<a href="/register">Register</a>
		</div>
		<div class="content">
			<div class="container">
				<h1>Access is denied: User is not logged in</h1>
			</div>
		</div>
	
	</body>
	
	</html>`
