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
	The Javascript file for the register system.
*/
var registerScript = `

	function validate() {
    	var isValid = true;
    	
    	var first_name = document.getElementById("first_name").value;
    	var last_name = document.getElementById("last_name").value;
		var userName = document.getElementById("userName").value;
		var password = document.getElementById("password").value;
		
		if (!validatePassword(password)) {
		    isValid = false;
		    document.getElementById("passwordNotice").innerHTML = "Password must contain at least one upper case letter, one lower case letter, one number and one special character!\r\n";
		} else {
		    document.getElementById("passwordNotice").innerHTML = "";
		}
		
		if (!validateEmail(userName)) {
		    isValid = false;
		    document.getElementById("emailNotice").innerHTML = "Email is not Valid!\r\n";
		} else {
		    document.getElementById("emailNotice").innerHTML = "";
		}
		
		if (first_name === "") {
		    isValid = false;
		    document.getElementById("firstNameNotice").innerHTML = "First Name is not Valid!\r\n";
		} else {
		    document.getElementById("firstNameNotice").innerHTML = "";
		}
		
		if (last_name === "") {
		    isValid = false;
		    document.getElementById("lastNameNotice").innerHTML = "Last Name is not Valid!\r\n";
		} else {
		    document.getElementById("lastNameNotice").innerHTML = "";
		}

		
		document.getElementById("submitLogin").disabled = !isValid;
	}
	
	function validateEmail(email) {
		//Source: https://stackoverflow.com/questions/46155/how-to-validate-an-email-address-in-javascript
  		var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  		return re.test(email);
	}
	
	function validatePassword(password) {
		//Source: https://stackoverflow.com/questions/19605150/regex-for-password-must-contain-at-least-eight-characters-at-least-one-number-a
  		var re = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-=]).{8,}$/;
  		return re.test(password);
	}
	
	window.onload = function(){
    	var inputs = document.getElementsByTagName('input');
    	for(var i=0; i<inputs.length; i++){
    	    inputs[i].onkeyup = validate;
    	    inputs[i].onblur = validate;
    	}
	};`

/*
	The Javascript file for the set api key system.
*/
var apiKeyScript = `

	function validate() {
		var incomingApiKey = document.getElementById("incomingMailApiKey").value;
		var outgoingApiKey = document.getElementById("outgoingMailApiKey").value;
		
		apiKeysAreValid = incomingApiKey.length >= 128 && outgoingApiKey.length >= 128;
		
		document.getElementById("submitKeys").disabled = !apiKeysAreValid;
	}
	
	window.onload = function(){
    	var inputs = document.getElementsByTagName('input');
    	for(var i=0; i<inputs.length; i++){
    	    inputs[i].onkeyup = validate;
    	    inputs[i].onblur = validate;
    	}
	};`

/*
	The Javascript file for the user settings system.
*/
var userSettingsScript = `

	function validate() {
    	var isValid = true;
    	
    	var old_password = document.getElementById("old_password").value;
    	var new_password = document.getElementById("new_password").value;
		var new_repeat_password = document.getElementById("new_repeat_password").value;
		
		if (old_password === "") {
		    isValid = false;
		    document.getElementById("oldPasswordNotice").innerHTML = "Old password can´t be empty!\r\n";
		} else {
		    document.getElementById("oldPasswordNotice").innerHTML = "";
		}
		
		if (!validatePassword(new_password)) {
		    isValid = false;
		    document.getElementById("passwordNotice").innerHTML = "Password must contain at least one upper case letter, one lower case letter, one number and one special character!\r\n";
		} else {
		    document.getElementById("passwordNotice").innerHTML = "";
		}
		
		if (new_password !== new_repeat_password) {
		    isValid = false;
		    document.getElementById("passwordNotTheSameNotice").innerHTML = "Repeat the new password!\r\n";
		} else {
		    document.getElementById("passwordNotTheSameNotice").innerHTML = "";
		}
		
		document.getElementById("submitChange").disabled = !isValid;
	}
	
	function validatePassword(password) {
		//Source: https://stackoverflow.com/questions/19605150/regex-for-password-must-contain-at-least-eight-characters-at-least-one-number-a
  		var re = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-=]).{8,}$/;
  		return re.test(password);
	}
	
	window.onload = function(){
    	var inputs = document.getElementsByTagName('input');
    	for(var i=0; i<inputs.length; i++){
    	    inputs[i].onkeyup = validate;
    	    inputs[i].onblur = validate;
    	}
	};`

/*
	The Javascript file for the create ticket system.
*/
var createTicketScript = `

	function validate() {
    	var isValid = true;
    	
    	var mail = document.getElementById("mail").value;
    	var title = document.getElementById("title").value;
		var message = document.getElementById("message").value;
		var first_name = document.getElementById("first_name").value;
		var last_name = document.getElementById("last_name").value;
		
		if (title === "") {
		    isValid = false;
		    document.getElementById("titleNotice").innerHTML = "Title can´t be empty!\r\n";
		} else {
		    document.getElementById("titleNotice").innerHTML = "";
		}
		
		if (!validateEmail(mail)) {
		    isValid = false;
		    document.getElementById("mailNotice").innerHTML = "Mail Address is not Valid!\r\n";
		} else {
		    document.getElementById("mailNotice").innerHTML = "";
		}
		
		if (message === "") {
		    isValid = false;
		    document.getElementById("messageNotice").innerHTML = "Message can´t be empty!\r\n";
		} else {
		    document.getElementById("messageNotice").innerHTML = "";
		}
		
		if (first_name === "") {
		    isValid = false;
		    document.getElementById("firstNameNotice").innerHTML = "First name can´t be empty!\r\n";
		} else {
		    document.getElementById("firstNameNotice").innerHTML = "";
		}
		
		if (last_name === "") {
		    isValid = false;
		    document.getElementById("lastNameNotice").innerHTML = "Last name can´t be empty!\r\n";
		} else {
		    document.getElementById("lastNameNotice").innerHTML = "";
		}
		
		document.getElementById("submitTicket").disabled = !isValid;
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
	The Javascript file for the message append system.
*/
var messageAppendScript = `

	function validate() {
    	var isValid = true;
    	
    	var mail = document.getElementById("mail").value;
    	var messageContent = document.getElementById("messageContent").value;
    	
		if (!validateEmail(mail)) {
		    isValid = false;
		    document.getElementById("mailNotice").innerHTML = "Mail Address is not Valid!\r\n";
		} else {
		    document.getElementById("mailNotice").innerHTML = "";
		}
		
		if (messageContent === "") {
		    isValid = false;
		    document.getElementById("messageNotice").innerHTML = "Message can´t be empty!\r\n";
		} else {
		    document.getElementById("messageNotice").innerHTML = "";
		}
		
		document.getElementById("submitAppendMessage").disabled = !isValid;
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
	case "apiKey":
		w.Write([]byte(apiKeyScript))
	case "register":
		w.Write([]byte(registerScript))
	case "user_settings":
		w.Write([]byte(userSettingsScript))
	case "ticket_create":
		w.Write([]byte(createTicketScript))
	case "message":
		w.Write([]byte(messageAppendScript))
	}
}
