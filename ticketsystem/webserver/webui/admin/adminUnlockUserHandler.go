package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
	"strconv"
)

/*
	Structure for the Login handler.
*/
type AdminUnlockUserHandler struct {
	UserContext user.UserContext
	Logger      logging.Logger
	MailContext	mail.MailContext
}

/*
	The Unlock user handler.
*/
func (a AdminUnlockUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		formId := r.FormValue("userId")

		userId, idConversionError := strconv.Atoi(formId)

		if idConversionError != nil {
			a.Logger.LogError("AdminUnlockUserHandler", idConversionError)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		accessToken := wrappers.GetUserToken(r.Context())

		unlocked, err := a.UserContext.UnlockAccount(accessToken, userId)

		if err != nil {
			a.Logger.LogError("AdminUnlockUserHandler", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		exist, existingUser := a.UserContext.GetUserById(userId)

		if unlocked && exist{
			mailSubject := mail.BuildUnlockUserNotificationMailSubject()
			mailContent := mail.BuildUnlockUserNotificationMailContent(existingUser.FirstName + " " + existingUser.LastName)

			err = a.MailContext.CreateNewOutgoingMail(existingUser.Mail, mailSubject, mailContent)

			if err != nil {
				a.Logger.LogError("AdminUnlockUserHandler", err)
				http.Redirect(w, r, "/admin", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
	}
}
