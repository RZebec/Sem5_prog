package staticFileHandlers

import "net/http"

func CssHandler(w http.ResponseWriter, r *http.Request) {

    var message = `
    
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
    }`

	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(message))
}

func LoginStyleHandler(w http.ResponseWriter, r *http.Request) {

    var message = `
    
    h2 {
        color: #ffffff;
        text-align: center;
    }
    
    div.container {
        width: 50%;
        height: 40em;
        margin: auto;
        display: inline-block;
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
        color: #4f4f4f;
        font-size: 20px;
        border-radius: 5px;
    }
    
    label {
        color: white;
        font-weight: bold;
    }
    
    input[type=button] {
        background-color: #4CAF50;;
        color: white;
        font-weight: bold;
        cursor: pointer;
        width: 100%;
        outline: none;
    }
    
    input[type=button]:hover {
        background-color: rgba(76, 175, 79, 0.466);
    }`

	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(message))
}
