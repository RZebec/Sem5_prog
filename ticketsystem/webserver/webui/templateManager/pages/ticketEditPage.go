package pages

var TicketEditPage = `

	{{ define "Title" }} Ticket Edit {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/center_main">
		<link rel="stylesheet" href="/files/style/dropdown">
	{{ end }}

	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Ticket Edit</h2>
					<form id="merge_form" method="POST" name="merge_form" action="/merge_tickets">
						<label>Merge ticket with:</label>
						<input type="hidden" name="firstTicketId" value="{{.TicketInfo.Id}}"/>
						<select name="secondTicketId" id="secondTicketId">
						{{range .OtherTickets}}	
  							<option value="{{.Id}}">{{.Id}}-{{.Title}}</option>
						{{end}}
						</select>
						<button type="submit" id="submitMerge" class="submit-button">Merge Tickets</button>
					</form>
					<form id="change_state_form" method="POST" name="change_state_form" action="/ticket_state_change">
						<label>State:</label>
						<select name="state" id="state">
							<option value="{{.TicketInfo.State}}" selected>{{.TicketInfo.State}}</option>
  							<option value="{{.OtherState1}}">{{.OtherState1}}</option>
							<option value="{{.OtherState2}}">{{.OtherState2}}</option>
						</select>
						<button type="submit" id="submitChangeState" class="submit-button">Change State</button>
					</form>
					<form id="change_editor_form" method="POST" name="change_editor_form" action="/ticket_editor_change">
						<label>Editor:</label>
						<select name="userId" id="userId">
						{{range .Users}}	
  							<option value="{{.UserId}}">{{.Mail}}</option>
						{{end}}
						</select>
						<button type="submit" id="submitChangeEditor" class="submit-button">Change Editor</button>
					</form>
				</div>
			</div>
		</div>
	</body>

	</html>
	{{ end }}`
