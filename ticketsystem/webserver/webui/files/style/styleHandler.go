package style

import (
	"net/http"
	"strings"
)

/*
	The main CSS file.
*/
var mainStyle = `
    
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
		height: 100%;
	}

    div.content {
        background-color: #333;
        width: 100%;
		min-height: 92vh;
        color: white;
		font-size: 17px;
     	text-align: center;
    }
	
	div.container {
		text-align: center;
        margin: 3em 1em auto 1em;
        display: inline-block;
    	width: 100em;
    }`

/*
	The CSS file for the login page.
*/
var loginStyle = `
    
    h2 {
        text-align: center;
    }
    
    div.main {
		text-align: left;
        width: 10em;
        padding: 0.5em 1.5em 0.75em;
        border: 2px solid gray;
        border-radius: 10px;
        float: left;
        margin-top: 1.5em;
    }
    
    input[type=text], input[type=password] {
        width: 100%;
        margin-bottom: 0.75em;
        margin-top: 0.2em;
        border: 2px solid white;
        color: black;
        border-radius: 5px;
    }
    
	.submit-button {
    	background-color: #4caf50;
    	color: white;
    	cursor: pointer;
    	width: 100%;
    	outline: none;
    	border: none;
    	height: 2em;
    	text-decoration: none;
		margin-bottom: 0.75em;
	}
    
    .submit-button:hover {
        background-color: rgba(76, 175, 79, 0.466);
    }

	.error-message{
		color: red;
	}`

var ticketStyle = `
    
    h2 {
        text-align: center;
    }

	table {
		text-align: left;
    	border-collapse: collapse;
    	width: 95em;
		margin: 0.5em 1.5em 0.75em 1.5em;
	}

	table td, table th {
   		border: 1px solid #595959;
    	padding: 0.3em;
	}

	table tr:nth-child(even){background-color: #595959;}

	table tr:hover {background-color: #737373;}

	table th {
    	padding-top: 0.3em;
    	padding-bottom: 0.3em;
    	text-align: left;
    	background-color: #4CAF50;
    	color: white;
	}

	.view-button {
    	background-color: #4caf50;
    	color: white;
    	cursor: pointer;
    	width: 100%;
    	outline: none;
    	border: none;
    	height: 2em;
    	text-decoration: none;
	}
    
    .view-button:hover {
        background-color: rgba(76, 175, 79, 0.466);
    }`

/*
	The handler for the style(css) files.
*/
func HandelStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	s := strings.Split(r.URL.Path, "/")

	switch s[3] {
	case "main":
		w.Write([]byte(mainStyle))
	case "login":
		w.Write([]byte(loginStyle))
	case "tickets":
		w.Write([]byte(ticketStyle))
	}
}
