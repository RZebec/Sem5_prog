package pages

var ticketEditPage = `

	{{ define "Title" }} Ticket Edit {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/center_main">
	{{ end }}

	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Merge Ticket</h2>
					<form id="merge_form" method="POST" name="merge_form" action="/ticket_merge">
						<label>Ticket to merge with:</label>
						<select name="ticketId" id="ticketId">
						{{range .OtherTickets}}	
  							<option value="{{.Id}}">[{{.Id}}]-{{.Title}}</option>
						{{end}}
						</select>
						<button type="submit" id="submitMerge" class="submit-button" disabled>Merge</button>
					</form>
				</div>
				<div class="main">
					<h2>Change State</h2>
					<form id="change_state_form" method="POST" name="change_state_form" action="/ticket_state_change">
						<label>State:</label>
						<select name="ticketId" id="ticketId">
							<option value="{{.TicketInfo.State}}" selected>{{.TicketInfo.State}}</option>
  							<option value="{{.OtherStates[0]}}">{{.OtherStates[0]}}</option>
							<option value="{{.OtherStates[1]}}">{{.OtherStates[1]}}</option>
						</select>
						<button type="submit" id="submitChangeState" class="submit-button">Change</button>
					</form>
				</div>
			</div>
		</div>
	</body>

	</html>
	{{ end }}`
