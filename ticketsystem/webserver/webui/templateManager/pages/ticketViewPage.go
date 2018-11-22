package pages

/*
	Html template for the Ticket View Page.
*/
var TicketViewPage = `

	{{ define "Title" }} Ticket {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/tickets"> 
	{{ end }}
	
	{{ define "Content" }}
		<div class="topnav">
			<a href="/">Home</a>
	
			<span>OP-Ticket-System</span>

			<a href="/tickets" class="active">Tickets</a>

			{{if .IsUserLoggedIn}}
				<a href="/user_logout">Logout</a>
			{{else}}
				<a href="/login">Login</a>
				<a href="/register">Register</a>
			{{end}}
		</div>
		<div class="content">
			<div class="container">
				<div class="main">
					<table>
						<tr>
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
							<th></th>
						</tr>
                        <tr>
                            <td>
                                {{.TicketInfo.Title}}
                            </td>
                            <td>
                                {{if .TicketInfo.HasEditor}}
                                    {{.TicketInfo.Editor.LastName}},&nbsp;
                                    {{.TicketInfo.Editor.FirstName}}&nbsp;
                                    {{.TicketInfo.Editor.Mail}}
                                {{else}}
                                    Ticket has no editor
                                {{end}}
                            </td>
                            <td>
                                {{.TicketInfo.Creator.LastName}},&nbsp;
                                {{.TicketInfo.Creator.FirstName}}&nbsp;
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
        							year: "numeric",
        							month: "2-digit",
        							day: "2-digit",
        							hour: "2-digit",
        							minute: "2-digit",
        							second: "2-digit"
    							};
								creationTime = creationTime.toLocaleDateString("de-DE", options);
								lastModificationTime = lastModificationTime.toLocaleDateString("de-DE", options);
								document.getElementById("creationTime").innerHTML = creationTime;
								document.getElementById("lastModificationTime").innerHTML = lastModificationTime;
							</script>
							<td>
								<button class="view-button" onclick="location.href='ticket_edit/{{.TicketInfo.Id}}';">
									Edit
								</button>
							</td>
                        </tr>
					</table>
				</div>
			</div>
		</div>
	</body>

	</html>
	{{ end }}`
