package pages

/*
	Html template for the Ticket View Page.
*/
var TicketViewPage = `

	{{ define "Title" }} Ticket {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/table"> 
		<link rel="stylesheet" href="/files/style/message">
		<script src="/files/script/message"></script>
	{{ end }}
	
	{{ define "Content" }}
		<div class="content">
			<div class="container">
					<table>
						<tr>
							<th>
								ID
                            </th>
							<th>
								Title
                            </th>
                            <th>
                                Editor
                            </th>
                            <th>
                                Creator
                            </th>
                            <th>
                                Created on
                            </th>
                            <th>
                                Updated on
                            </th>
							<th>
								State
							</th>
							<th></th>
						</tr>
                        <tr>
							<td>
                                {{.TicketInfo.Id}}
                            </td>
                            <td>
                                {{.TicketInfo.Title}}
                            </td>
                            <td>
                                {{if .TicketInfo.HasEditor}}
                                    {{.TicketInfo.Editor.LastName}},&nbsp;
                                    {{.TicketInfo.Editor.FirstName}},&nbsp;
                                    {{.TicketInfo.Editor.Mail}}
                                {{else}}
                                    Ticket has no editor
                                {{end}}
                            </td>
                            <td>
                                {{.TicketInfo.Creator.LastName}},&nbsp;
                                {{.TicketInfo.Creator.FirstName}},&nbsp;
                                {{.TicketInfo.Creator.Mail}}
                            </td>
                            <td id="creationTime">
                            </td>
                            <td id="lastModificationTime">
                            </td>
							<script>
								var creationTime = new Date({{.TicketInfo.CreationTime}});
								var lastModificationTime = new Date({{.TicketInfo.LastModificationTime}});
								var options = {
        							weekday: "short",
        							year: "2-digit",
        							month: "2-digit",
        							day: "2-digit",
        							hour: "2-digit",
        							minute: "2-digit",
        							second: "2-digit"
    							};
								creationTime = creationTime.toLocaleDateString("en-GB", options);
								lastModificationTime = lastModificationTime.toLocaleDateString("de-DE", options);
								document.getElementById("creationTime").innerHTML = creationTime;
								document.getElementById("lastModificationTime").innerHTML = lastModificationTime;
							</script>
							<td>
								{{.TicketInfo.State}}
							</td>
							<td>
								<button class="view-button" onclick="location.href='ticket_edit/{{.TicketInfo.Id}}';">
									Edit
								</button>
							</td>
                        </tr>
					</table>
					<table>
						<tr>
							<th>
								Creator
                            </th>
                            <th>
                                Message
                            </th>
                            <th>
                                Created on
                            </th>
						</tr>
                        {{range .Messages}}
                        <tr>
							<td>
                                {{.CreatorMail}}
                            </td>
                            <td>
                                {{.Content}}
                            </td>
                            <td id="creationTime_{{.Id}}">
                            </td>
							<script>
								var creationTime = new Date({{.CreationTime}});
								var options = {
        							weekday: "short",
        							year: "2-digit",
        							month: "2-digit",
        							day: "2-digit",
        							hour: "2-digit",
        							minute: "2-digit",
        							second: "2-digit"
    							};
								creationTime = creationTime.toLocaleDateString("en-GB", options);
								document.getElementById("creationTime_{{.Id}}").innerHTML = creationTime;
							</script>
                        </tr>
						{{end}}
					</table>
					<table>
						<tr>
							<th>
								Creator
                            </th>
                            <th>
                                Message
                            </th>
							{{if .UserIsAuthenticated}}
							<th>
								<label>Internal Only</label>
							</th>
							{{end}}
                            <th>
                            </th>
						</tr>
						<tr>
							<form id="appendMessageForm" method="POST" name="appendMessageForm" action="/append_message">
								<td>
								{{if .UserIsAuthenticated}}
                                	<input type="text" name="mail" id="mail" value="{{.UserName}}" readonly/>
								{{else}}
									<input type="text" name="mail" id="mail" value=""/>
								{{end}}
                            	</td>
                            	<td>
                                	<input type="text" name="messageContent" id="messageContent" value=""/>
                            	</td>
								{{if .UserIsAuthenticated}}
								<td>
									<input type="radio" name="onlyInternal" id="onlyInternal" value="true"/> Yes
									<input type="radio" name="onlyInternal" id="onlyInternal" value="false" checked/> No
								</td>
								{{end}}
                            	<td>
									<input type="hidden" name="ticketId" value="{{.TicketInfo.Id}}"/>
									<button type="submit" id="submitAppendMessage" class="submit-button" disabled>Submit Append Message</button>
                            	</td>
							</form>
                        </tr>
						<tr>
							<td>
								<span id="mailNotice" class="error-message"></span>
							</td>
							<td>
								<span id="messageNotice" class="error-message"></span>
							</td>
							{{if .UserIsAuthenticated}}
								<td></td>
							{{end}}
							<td></td>
						</tr>
					</table>
			</div>
		</div>
	</body>

	</html>
	{{ end }}`
