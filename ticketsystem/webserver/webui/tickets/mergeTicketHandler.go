package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	A ticket merge handler.
*/
type TicketMergeHandler struct {
	Logger        logging.Logger
	TicketContext ticketData.TicketContext
	MailContext   mailData.MailContext
	UserContext   userData.UserContext
}

/*
	Merge two tickets.
*/
func (t TicketMergeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		rawFirstTickedId := r.FormValue("firstTicketId")
		rawSecondTickedId := r.FormValue("secondTicketId")

		firstTicketId, err := strconv.Atoi(rawFirstTickedId)
		if err != nil {
			t.Logger.LogError("TicketMergeHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		secondTicketId, err := strconv.Atoi(rawSecondTickedId)
		if err != nil {
			t.Logger.LogError("TicketMergeHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		firstTicketExists, firstTicket := t.TicketContext.GetTicketById(firstTicketId)
		if !firstTicketExists {
			t.Logger.LogError("TicketMergeHandler", errors.New("ticket does not exist. id: "+rawFirstTickedId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		secondTicketExists, secondTicket := t.TicketContext.GetTicketById(secondTicketId)
		if !secondTicketExists {
			t.Logger.LogError("TicketMergeHandler", errors.New("ticket does not exist. id: "+rawSecondTickedId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return

		}

		success, err := t.TicketContext.MergeTickets(firstTicketId, secondTicketId)
		if err != nil {
			t.Logger.LogError("TicketMergeHandler", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		if success {
			olderTicket := firstTicketId
			newerTicket := secondTicketId
			if newerTicket < olderTicket {
				olderTicket = secondTicketId
			}

			// Notify creator of first ticket:
			firstMailSubject := mailData.BuildTicketMergeNotificationMailSubject(firstTicket.Info().Id, olderTicket)
			firstMailContent := mailData.BuildTicketMergeNotificationMailContent(firstTicket.Info().Creator.Mail,
				firstTicketId, secondTicketId)
			err = t.MailContext.CreateNewOutgoingMail(firstTicket.Info().Creator.Mail, firstMailSubject, firstMailContent)
			if err != nil {
				t.Logger.LogError("TicketMergeHandler", err)
				http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusInternalServerError)
				return
			}
			// Notify creator of second ticket:
			secondMailSubject := mailData.BuildTicketMergeNotificationMailSubject(secondTicket.Info().Id, olderTicket)
			secondMailContent := mailData.BuildTicketMergeNotificationMailContent(secondTicket.Info().Creator.Mail,
				firstTicketId, secondTicketId)
			err = t.MailContext.CreateNewOutgoingMail(secondTicket.Info().Creator.Mail, secondMailSubject, secondMailContent)
			if err != nil {
				t.Logger.LogError("TicketMergeHandler", err)
				http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusInternalServerError)
				return
			}

			loggedInUserId := wrappers.GetUserId(r.Context())
			userExists, authenticatedUser := t.UserContext.GetUserById(loggedInUserId)
			if !userExists {
				t.Logger.LogError("TicketMergeHandler", err)
				http.Redirect(w, r, "/", http.StatusBadRequest)
				return
			}

			// Build message for history:
			messageEntry := ticketData.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
				Content: "Tickets merged: " + rawFirstTickedId + " with " + rawSecondTickedId, CreationTime: time.Now()}
			_, err = t.TicketContext.AppendMessageToTicket(olderTicket, messageEntry)
			if err != nil {
				t.Logger.LogError("TicketMergeHandler", err)
				http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusInternalServerError)
				return
			}

			// Redirect to the older ticket:
			http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusFound)

			return
		}
	}
}
