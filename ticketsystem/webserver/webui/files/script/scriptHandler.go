package script

import (
	"net/http"
	"strings"
)

/*
	The Javascript file for the login system.
*/
var loginScript = `

	function validate() {
    	var emailIsValid = false;
		var userName = document.getElementById("userName").value;
		var password = document.getElementById("password").value;
		
		console.log(userName + ":" + password)
		
		if (password == "") {
		    emailIsValid = false;
		}
		else {
		    if (validateEmail(userName)) {
		        document.getElementById("emailNotice").innerHTML = "";
		        emailIsValid = true;
		    } 
		    else {
		        document.getElementById("emailNotice").innerHTML = "Email is not Valid!";
		    }		    
		}
		
		document.getElementById("submitLogin").disabled = !emailIsValid;
	}
	
	function validateEmail(email) {
		//Source: https://stackoverflow.com/questions/46155/how-to-validate-an-email-address-in-javascript
  		var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  		return re.test(email);
	}
	
	window.onload = function(){
    	var inputs = document.getElementsByTagName('input');
    	for(var i=0; i<inputs.length; i++){
    	    inputs[i].onkeyup = validate;
    	    inputs[i].onblur = validate;
    	}
	};`

/*
	The handler for the script(Javascript) files.
*/
func HandelScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")

	s := strings.Split(r.URL.Path, "/")

	switch s[3] {
	case "login":
		w.Write([]byte(loginScript))
	}
}
