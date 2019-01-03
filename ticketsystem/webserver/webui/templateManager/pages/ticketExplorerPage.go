package pages

/*
	Html template for the Ticket Explorer Page.
*/
var TicketExplorerPage = `	

	{{ define "Title" }} Tickets {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/table">
	{{ end }}
	
	{{ define "Content" }}
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
							<th>
								State
							</th>
							<th></th>
						</tr>
                        {{range .Tickets}}
                        <tr>
                            <td>
                                {{.Title}}
                            </td>
                            <td>
                                {{if .HasEditor}}
                                    {{.Editor.LastName}},&nbsp;
                                    {{.Editor.FirstName}}&nbsp;
                                    {{.Editor.Mail}}
                                {{else}}
                                    Ticket has no editor
                                {{end}}
                            </td>
                            <td>
                                {{.Creator.LastName}},&nbsp;
                                {{.Creator.FirstName}}&nbsp;
                                {{.Creator.Mail}}
                            </td>
                            <td id="creationTime_{{.Id}}">
                            </td>
                            <td id="lastModificationTime_{{.Id}}">
                            </td>
							<script>
								var creationTime = new Date({{.CreationTime}});
								var lastModificationTime = new Date({{.LastModificationTime}});
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
								lastModificationTime = lastModificationTime.toLocaleDateString("en-GB", options);
								document.getElementById("creationTime_{{.Id}}").innerHTML = creationTime;
								document.getElementById("lastModificationTime_{{.Id}}").innerHTML = lastModificationTime;
							</script>
							<td>
								{{.State}}
							</td>
							<td>
								<button class="view-button" onclick="location.href='ticket/{{.Id}}';">
									View
								</button>
							</td>
                        </tr>
						{{end}}
					</table>
				</div>
			</div>
		</div>
	</body>

	</html>
	{{ end }}`
