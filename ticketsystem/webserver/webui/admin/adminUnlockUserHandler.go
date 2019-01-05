package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
	"strconv"
)

/*
	Structure for the Login handler.
*/
type UnlockUserHandler struct {
	UserContext userData.UserContext
	Logger      logging.Logger
	MailContext mailData.MailContext
}

/*
	The Unlock userData handler.
*/
func (a UnlockUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		formId := r.FormValue("userId")

		userId, idConversionError := strconv.Atoi(formId)

		if idConversionError != nil {
			a.Logger.LogError("UnlockUserHandler", idConversionError)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		accessToken := wrappers.GetUserToken(r.Context())

		unlocked, err := a.UserContext.UnlockAccount(accessToken, userId)

		if err != nil {
			a.Logger.LogError("UnlockUserHandler", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		exist, existingUser := a.UserContext.GetUserById(userId)

		if unlocked && exist {
			mailSubject := mailData.BuildUnlockUserNotificationMailSubject()
			mailContent := mailData.BuildUnlockUserNotificationMailContent(existingUser.FirstName + " " + existingUser.LastName)

			err = a.MailContext.CreateNewOutgoingMail(existingUser.Mail, mailSubject, mailContent)

			if err != nil {
				a.Logger.LogError("UnlockUserHandler", err)
				http.Redirect(w, r, "/admin", http.StatusInternalServerError)
				return
			}

			a.Logger.LogInfo("UnlockUserHandler", "User unlocked. UserId: "+formId)
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
	}
}
