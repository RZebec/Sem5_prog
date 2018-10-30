package accessdenied

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var accessDeniedPage = `
	<!DOCTYPE html>
	<html>

	<head>
		<link rel="stylesheet" href="/files/style/main">
		<link rel="stylesheet" href="/files/style/login"> 
		<script src="/files/script/login"></script>
	</head>
	
	<body>
		<div class="topnav">
			<a href="/">Home</a>
	
			<span>OP-Ticket-System</span>

			<a href="/login">Login</a>
			<a href="/register">Register</a>
		</div>
		<div class="content">
			<div class="container">
				<h1>Access is denied: User is not logged in</h1>
			</div>
		</div>
	</body>
	
	</html>`

/*
	Access Denied Page returns valid html.
	Test template source: https://blog.questionable.services/article/testing-http-handlers-go/
 */
func TestAccessDeniedPageHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/test/access_denied", nil)
	if err != nil {
		t.Fatal(err)
	}

	accessDeniedPageHandler := AccessDeniedPageHandler{}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	accessDeniedPageHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := accessDeniedPage
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
