package staticFileHandlers

import "net/http"

func JSFileHandler() {
	http.HandleFunc("/files/js/index", IndexJSHandler)
	http.HandleFunc("/files/js/login", LoginJSHandler)
}

func IndexJSHandler(w http.ResponseWriter, r *http.Request) {

	var message = `        

	const Http = new XMLHttpRequest();
	const url='https://jsonplaceholder.typicode.com/posts';
	Http.open("GET", url);
	Http.send();
	Http.onreadystatechange=(e)=>{
		console.log(Http.responseText)
	}`

	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(message))
}

func LoginJSHandler(w http.ResponseWriter, r *http.Request) {

	var message = `
	
	var attempt = 3; // Variable to count number of attempts.
	// Below function Executes on click of login button.
	function validate() {
		var username = document.getElementById("username").value;
		var password = document.getElementById("password").value;
		if (username == "admin" && password == "admin") {
			alert("Login successfully");
			window.location = "/"; // Redirecting to other page.
			return false;
		}
		else {
			attempt--;// Decrementing by one.
			alert("You have left " + attempt + " attempt;");
			// Disabling fields after 3 attempts.
			if (attempt == 0) {
				document.getElementById("username").disabled = true;
				document.getElementById("password").disabled = true;
				document.getElementById("submit").disabled = true;
				return false;
			}
		}
	}`

	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(message))
}
