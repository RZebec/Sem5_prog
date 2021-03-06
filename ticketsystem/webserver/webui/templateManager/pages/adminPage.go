// 5894619, 6720876, 9793350
package pages

var AdminPage = `
	{{ define "Title" }} Admin {{ end }}

	{{ define "StylesAndScripts" }}
		<script src="/files/script/apiKey"></script>
		<link rel="stylesheet" href="/files/style/table">
		<link rel="stylesheet" href="/files/style/center_main">
		<link rel="stylesheet" href="/files/style/largeMainStyle"> 
	{{ end }}
	
	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Mail Api Keys</h2>
					<form id="form_id" method="POST" name="myform" action="/set_api_keys">
						<label>Incoming Mail Api Key:</label>
						<input type="text" name="incomingMailApiKey" id="incomingMailApiKey" value="{{.IncomingMailApiKey}}"/>
						<label>Outgoing Mail Api Key:</label>
						<input type="text" name="outgoingMailApiKey" id="outgoingMailApiKey" value="{{.OutgoingMailApiKey}}"/>
						<button type="submit" id="submitKeys" class="submit-button" disabled>Set Api Keys</button>
					</form>
					{{if eq .IsChangeFailed "yes" }}
					<span class="error-message">
							Api Key Change Failed!
					</span>
					</br>
					{{end}}
					{{if eq .IsChangeFailed "no" }}
					<span>
							Api Key Change Successful!
					</span>
					</br>
					{{end}}
				</div>
			</div>
				<div class="container">
					<div class="main">
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
		</div>
	{{ end }}`
