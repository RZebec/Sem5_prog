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
						<label>Mail Address:</label>
						{{if .IsUserLoggedIn}}
							<input type="text" name="mail" id="mail" value="{{.UserName}}" readonly/>
						{{else}}
							<input type="text" name="mail" id="mail"/>
						{{end}}
						<label>Title:</label>
						<input type="text" name="title" id="title"/>
						<label>Initial Message:</label>
						<input type="text" name="message" id="message"/>
						<label>Internal Only:</label>
						<input type="radio" name="internal" id="internal" value="true" checked/> Yes
						<input type="radio" name="internal" id="internal" value="false"/> No
						<button type="submit" id="submitTicket" class="submit-button" disabled>Create Ticket</button>
					</form>
					<span id="mailNotice" class="error-message"></span>
					<span id="titleNotice" class="error-message"></span>
					<span id="messageNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	{{ end }}`

