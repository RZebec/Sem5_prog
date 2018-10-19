package staticFileHandlers

import "net/http"

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	var message = `<html>

	<head>
    	<link rel="stylesheet" href="/files/styles">
	</head>

	<body>
    	<div class="topnav">
        	<a class="active" href="#home">Home</a>
        	<a href="#news">News</a>
        	<a href="#contact">Contact</a>

        	<span>OP-Ticket-System</span>

			<a href="/login">Login</a>
    	</div>
		<div class="content">
			<p>
				Lorem ipsum dolor sit amet, te ius scaevola maiestatis, pro te munere ullamcorper, per te erat novum civibus. 
			</p>
		</div>
	</body>

	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(message))
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {

	var message = `
	
	<html>

	<head>
		<link rel="stylesheet" href="/files/styles">
		<link rel="stylesheet" href="/files/login-styles"> 
		<script src="/files/login"></script>
	</head>
	
	<body>
		<div class="topnav">
			<a class="active" href="#home">Home</a>
			<a href="#news">News</a>
			<a href="#contact">Contact</a>
	
			<span>OP-Ticket-System</span>
	
			<a href="/login">Login</a>
		</div>
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Login</h2>
					<form id="form_id" method="post" name="myform">
						<label>Username:</label>
						<input type="text" name="username" id="username" />
						<label>Password:</label>
						<input type="password" name="password" id="password" />
						<input type="button" value="Login" id="submit" onclick="validate()" />
					</form>
					<span>
						<b class="note">
							Note :
						</b>
						For this demo use following username and password. <br />
						<b class="valid">
							Username: admin<br />Password: admin
						</b>
					</span>
				</div>
			</div>
		</div>
		</div>
	</body>
	
	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(message))
}
