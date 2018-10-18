package main

import (
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	var message = `<html>

	<head>
    	<link rel="stylesheet" href="/files/styles">
	</head>

	<body>
    	<div class="topnav">
        	<a class="active" href="#home">Home</a>
        	<a href="#news">News</a>
        	<a href="#contact">Contact</a>
        	<a href="#about">About</a>

        	<span>OP-Ticket-System</span>
    	</div>
	</body>

	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(message))
}

func cssHandler(w http.ResponseWriter, r *http.Request) {

	var message = `        

		/* Add a black background color to the top navigation */
        .topnav {
            background-color: #333;
            overflow: hidden;
        }

        /* Style the links inside the navigation bar */
        .topnav a {
            float: left;
            color: #f2f2f2;
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
            color: #f2f2f2;
            text-align: center;
            padding: 14px 16px;
            text-decoration: none;
            font-size: 20px;
        }`

	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(message))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/files/styles", cssHandler)
	if err := http.ListenAndServe(":5050", nil); err != nil {
		panic(err)
	}
}