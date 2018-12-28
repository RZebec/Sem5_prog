package templateManager

/*
	A logger which sends the messages to the console.
*/
type MockConsoleLogger struct {
}

/*
	Log a debug message.
*/
func (l MockConsoleLogger) LogDebug(prefix string, message string) {

}

/*
	Log a info message.
*/
func (l MockConsoleLogger) LogInfo(prefix string, message string) {

}

/*
	Log a warning.
*/
func (l MockConsoleLogger) LogWarning(prefix string, message string) {

}

/*
	Log an error.
*/
func (l MockConsoleLogger) LogError(prefix string, err error) {

}

/*
	Expected result for the access denied page test.
*/
var accessDeniedResultPage = `<!DOCTYPE html>
	<html>

	<head>
		<title>
			 Access Denied 
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

/*
	Expected result for the login page test.
*/
var loginResultPage = `<!DOCTYPE html>
	<html>

	<head>
		<title>
			 Login 
		</title>
		<link rel="stylesheet" href="/files/style/main">
		
		<link rel="stylesheet" href="/files/style/login"> 
		<script src="/files/script/login"></script>
	
	</head>
	
	<body>
		
		<div class="topnav">
			<a href="/">Home</a>

			<a href="/tickets">Tickets</a>
	
			<span>OP-Ticket-System</span>

			<a class="active" href="/login">Login</a>
			<a href="/register">Register</a>
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
					
					<span id="emailNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	
	</body>
	
	</html>`
