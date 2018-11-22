package files

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var expectedStyle = `
    
    /* Add a black background color to the top navigation */
    .topnav {
        background-color: #333;
        overflow: hidden;
    }

    /* Style the links inside the navigation bar */
    .topnav a {
        float: left;
        color: white;
        text-align: center;
        padding: 14px 16px;
        text-decoration: none;
        font-size: 17px;
    }

    /* Change the color of links on hover */
    .topnav a:hover {
        background-color: #ddd;
        color: black;
    }

    /* Add a color to the active/current link */
    .topnav a.active {
        background-color: #4CAF50;
        color: white;
    }

    /* Page description Text field */
    .topnav span {
        float: right;
        color: white;
        text-align: center;
        padding: 14px 16px;
        text-decoration: none;
        font-size: 20px;
    }

	body {
		font-family: Calibri,Candara,Segoe,Segoe UI,Optima,Arial,sans-serif; 
	}

    div.content {
        background-color: #333;
        width: 100%;
        color: white;
		font-size: 17px;
     	text-align: center;
		min-height: 45em;
    }
	
	div.container {
		text-align: center;
        margin: 3em 1em auto 1em;
        display: inline-block;
    	width: 100em;
    }`

var expectedScript = `

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
	Files Handler returns valid Script File.
	Test template source: https://blog.questionable.services/article/testing-http-handlers-go/
*/
func TestFilesHandler_Script(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/files/script/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	filesHandler := FilesHandler{}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	filesHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := expectedScript
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

/*
	Files Handler returns valid Style File.
	Test template source: https://blog.questionable.services/article/testing-http-handlers-go/
*/
func TestFilesHandler_Style(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/files/style/main", nil)
	if err != nil {
		t.Fatal(err)
	}

	filesHandler := FilesHandler{}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	filesHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := expectedStyle
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
