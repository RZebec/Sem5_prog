package style

import (
	"net/http"
	"strings"
)

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

    .content {
        background-color: #333;
        width: 100%;
        color: white;
		font-size: 17px;
    }`

var loginStyle = `
    
    h2 {
        text-align: center;
    }
    
    div.container {
        height: 40em;
        margin: auto;
        display: inline-block;
    	width: 50%;
    }
    
    div.main {
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

func HandelStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	s := strings.Split(r.URL.Path, "/")

	switch s[3] {
	case "main":
		w.Write([]byte(mainStyle))
	case "login":
		w.Write([]byte(loginStyle))
	}
}
