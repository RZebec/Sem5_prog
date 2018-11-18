package templateManager

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var accessDeniedResultPage = `<!DOCTYPE html>
	<html>

	<head>
		<title>
			 Access Denied 
		</title>
		<link rel="stylesheet" href="/files/style/main">
		 
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
	
	</html>
	
	</body>
	
	</html>`

func TestAccessDeniedPageHandler(t *testing.T) {
	rr := httptest.NewRecorder()

	LoadTemplates()
	err := RenderTemplate(rr, "AccessDeniedPage", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := strings.TrimSpace(accessDeniedResultPage)
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			strings.TrimSpace(rr.Body.String()), expected)
	}
}
