package templateManager

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
	Expected result for the access denied page test.
 */
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
	
	</html>`

/*
	Expected result for the login page test.
 */
var loginResultPage = `<!DOCTYPE html>
	<html>

	<head>
		<title>
			 Login 
		</title>
		<link rel="stylesheet" href="/files/style/main">
		
		<link rel="stylesheet" href="/files/style/login"> 
		<script src="/files/script/login"></script>
	
	</head>
	
	<body>
		
		<div class="topnav">
			<a href="/">Home</a>
	
			<span>OP-Ticket-System</span>

			<a class="active" href="/login">Login</a>
			<a href="/register">Register</a>
		</div>
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Login</h2>
					<form id="form_id" method="post" name="myform" action="/user_login">
						<label>Username:</label>
						<input type="text" name="userName" id="userName" />
						<label>Password:</label>
						<input type="password" name="password" id="password" />
						<button type="submit" id="submitLogin" class="submit-button" disabled>Login</button>
					</form>
					
					<span id="emailNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	
	</body>
	
	</html>`


/*
	Tests if the Page is correctly rendered.
 */
func TestAccessDeniedPageRendering(t *testing.T) {
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
	result := strings.TrimSpace(rr.Body.String())
	if result != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			result, expected)
	}
}

/*
	Mock structure for the Login Page Data.
*/
type mockLoginPageData struct {
	IsLoginFailed bool
}

/*
	Tests if the Page is correctly rendered for the given data.
 */
func TestLoginPageRendering(t *testing.T) {
	rr := httptest.NewRecorder()

	testData := mockLoginPageData{
		IsLoginFailed: false,
	}

	LoadTemplates()
	err := RenderTemplate(rr, "LoginPage", testData)

	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := strings.TrimSpace(loginResultPage)
	result := strings.TrimSpace(rr.Body.String())
	if result != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			result, expected)
	}
}

/*
	Tests if the page returns an error if an page wasn´t able to be rendered.
 */
func TestErrorHandling(t *testing.T) {
	rr := httptest.NewRecorder()

	expectedError := NewError("Template doesn't exist")

	expectedText := "The template TestPage does not exist."

	err := RenderTemplate(rr, "TestPage", nil)

	if err.Error() != expectedError.Error() {
		t.Errorf("error returned unexpected text: got %v want %v",
			err, expectedError)
	}

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	result := strings.TrimSpace(rr.Body.String())
	if result != expectedText {
		t.Errorf("handler returned unexpected body: got %v want %v",
			result, expectedText)
	}
}
