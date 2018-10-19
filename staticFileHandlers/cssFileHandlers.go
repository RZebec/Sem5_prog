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
        }

		.topnav span + a {
			float: right;
		}

		.content {
            background-color: #333;
			color: #f2f2f2;
        }`

	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(message))
}
