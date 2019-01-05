package pages

var TicketCreatePage = `
	{{ define "Title" }} Admin {{ end }}

	{{ define "StylesAndScripts" }}
		<script src="/files/script/ticket_create"></script>
		<link rel="stylesheet" href="/files/style/table">
		<link rel="stylesheet" href="/files/style/center_main">
		<link rel="stylesheet" href="/files/style/largeMainStyle">
	{{ end }}
	
	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Create a Ticket</h2>
					<form id="form_id" method="POST" name="myform" action="/create_ticket">
						{{if .IsUserLoggedIn}}
							<label>Mail Address:</label>
							<input type="text" name="mail" id="mail" value="{{.UserName}}" readonly/>
							<label>First Name:</label>
							<input type="text" name="first_name" id="first_name" value="{{.FirstName}}" readonly/>
							<label>Last Name:</label>
							<input type="text" name="last_name" id="last_name" value="{{.LastName}}" readonly/>
						{{else}}
							<label>Mail Address:</label>
							<input type="text" name="mail" id="mail"/>
							<label>First Name:</label>
							<input type="text" name="first_name" id="first_name"/>
							<label>Last Name:</label>
							<input type="text" name="last_name" id="last_name"/>
						{{end}}
						<label>Ticket Title:</label>
						<input type="text" name="title" id="title"/>
						<label>Initial Message:</label>
						<input type="text" name="message" id="message"/>
						{{if .IsUserLoggedIn}}
							<label>Internal Only:</label>
							<input type="radio" name="internal" id="internal" value="true"/> Yes
							<input type="radio" name="internal" id="internal" value="false" checked/> No
						{{end}}
						<button type="submit" id="submitTicket" class="submit-button" disabled>Create Ticket</button>
					</form>
					<span id="mailNotice" class="error-message"></span>
					<span id="firstNameNotice" class="error-message"></span>
					<span id="lastNameNotice" class="error-message"></span>
					<span id="titleNotice" class="error-message"></span>
					<span id="messageNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	{{ end }}`
