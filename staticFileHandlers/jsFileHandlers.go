package staticFileHandlers

import "net/http"

func JsHandler(w http.ResponseWriter, r *http.Request) {

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
