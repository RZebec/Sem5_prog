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

			<a href="#login">Login</a>
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
