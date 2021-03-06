// 5894619, 6720876, 9793350
package pages

/*
	Html template for the User Settings Page.
*/
var UserSettingsPage = `
	{{ define "Title" }} User Settings {{ end }}

	{{ define "StylesAndScripts" }}
		<link rel="stylesheet" href="/files/style/center_main"> 
		<script src="/files/script/user_settings"></script>
	{{ end }}
	
	{{ define "Content" }}
		<div class="content">
			<div class="container">
				<div class="main">
					<h2>Vacation Mode</h2>
					<form id="toggleVacationMode" method="post" name="toggleVacationMode" action="/user_toggle_vacation">
					{{if .UserIsOnVacation }}
						<label>Your user is currently in Vacation Mode</label>
						<input type="hidden" name="vacationMode" id="vacationMode" value="false"/>
        			{{else}}
						<label>Your user is currently not in Vacation Mode</label>
						<input type="hidden" name="vacationMode" id="vacationMode" value="true"/>
					{{end}}
						<button type="submit" id="submitToggle" class="submit-button">Change Vacation Mode</button>
					</form>
					<h2>Change Password</h2>
					<form id="form_id" method="post" name="myform" action="/user_change_password">
						<label>Old Password:</label>
						<input type="password" name="old_password" id="old_password" />
						<label>New Password:</label>
						<input type="password" name="new_password" id="new_password" />
						<label>Repeat new Password:</label>
						<input type="password" name="new_repeat_password" id="new_repeat_password" />
						<button type="submit" id="submitChange" class="submit-button" disabled>Change Password</button>
					</form>
					{{if eq .IsChangeFailed "yes" }}
					<span class="error-message">
							Password Change Failed!
					</span>
					</br>
					{{end}}
					{{if eq .IsChangeFailed "no" }}
					<span>
							Password Change Successful!
					</span>
					</br>
					{{end}}
        			
					<span id="passwordNotice" class="error-message"></span>
					<span id="passwordNotTheSameNotice" class="error-message"></span>
					<span id="oldPasswordNotice" class="error-message"></span>
				</div>
			</div>
		</div>
	{{ end }}`
