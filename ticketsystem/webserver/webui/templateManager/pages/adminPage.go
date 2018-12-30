package pages

var AdminPage = `
	{{ define "Title" }} Admin {{ end }}

	{{ define "StylesAndScripts" }}
		<script src="/files/script/apiKey"></script>
		<link rel="stylesheet" href="/files/style/table">
		<link rel="stylesheet" href="/files/style/center_main">
	{{ end }}
	
	{{ define "Content" }}
		<div class="topnav">
			<a href="/">Home</a>

			<a href="/tickets">Tickets</a>
	
			<span>OP-Ticket-System</span>

			<a href="/login">Login</a>
			<a href="/register">Register</a>

			<a href="/admin">Admin</a>
		</div>
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Mail Api Keys</h2>
					<form id="form_id" method="post" name="myform" action="/set_api_keys">
						<label>Incoming Mail Api Key:</label>
						<input type="text" name="incomingMailApiKey" id="incomingMailApiKey" />
						<label>Outgoing Mail Api Key:</label>
						<input type="text" name="outgoingMailApiKey" id="outgoingMailApiKey" />
						<button type="submit" id="submitKeys" class="submit-button" disabled>Set Api Keys</button>
					</form>
				</div>
			</div>
				<div class="container">
					<h2>Locked Users</h2>
					<table>
						<tr>
							<th>
								ID
							</th>
							<th>
								Mail
                            </th>
                            <th>
                                First Name
                            </th>
                            <th>
                                Last Name
                            </th>
                            <th>
                            </th>
						</tr>
						{{range .Users}}
                        <tr>
                            <td>
                                {{.UserId}}
                            </td>
                            <td>
                                {{.Mail}}
                            </td>
                            <td>
                                {{.FirstName}}
                            </td>
                            <td>
								{{.LastName}}
                            </td>
							<td>
								<form id="formUser" method="post" name="formUser" action="/unlock_user">
									<input type = "hidden" name = "userId" value = "{{.UserId}}" />
									<button type="submit" class="view-button">Unlock</button>
								</form>
							</td>
                        </tr>
						{{end}}
					</table>
				</div>
		</div>
	{{ end }}`
